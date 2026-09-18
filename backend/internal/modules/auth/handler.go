package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	SecureCookie   bool
	AllowedOrigins []string
}

type Handler struct {
	store          store
	cookieName     string
	secureCookie   bool
	allowedOrigins map[string]struct{}
	dummyHash      string
	passwordSlots  chan struct{}
}

func New(pool *pgxpool.Pool, opts Options) (*Handler, error) {
	return newHandler(postgresStore{pool: pool}, opts)
}

func newHandler(repository store, opts Options) (*Handler, error) {
	if len(opts.AllowedOrigins) == 0 {
		return nil, errors.New("at least one authentication origin is required")
	}
	dummyHash, err := argon2id.CreateHash("dummy-password-for-constant-work", passwordParams)
	if err != nil {
		return nil, errors.New("password hashing initialization failed")
	}
	h := &Handler{
		store: repository, cookieName: "ta_session", secureCookie: opts.SecureCookie,
		allowedOrigins: make(map[string]struct{}), dummyHash: dummyHash,
		// Bound memory used by password hashing under concurrent untrusted input.
		passwordSlots: make(chan struct{}, 2),
	}
	if opts.SecureCookie {
		h.cookieName = "__Host-ta_session"
	}
	for _, origin := range opts.AllowedOrigins {
		h.allowedOrigins[origin] = struct{}{}
	}
	return h, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("POST /auth/register", h.mutation(h.passwordWork(http.HandlerFunc(h.register))))
	mux.Handle("POST /auth/login", h.mutation(h.passwordWork(http.HandlerFunc(h.login))))
	mux.Handle("GET /auth/me", h.RequireUser(http.HandlerFunc(h.me)))
	mux.Handle("POST /auth/logout", h.mutation(http.HandlerFunc(h.logout)))
}

// RequireUser is reusable by future protected endpoints. Authentication is
// resolved from PostgreSQL on every request so revocation takes effect immediately.
func (h *Handler) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		cookie, err := r.Cookie(h.cookieName)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "unauthenticated", "Sign in to continue.")
			return
		}
		digest, valid := sessionTokenHash(cookie.Value)
		if !valid {
			writeError(w, http.StatusUnauthorized, "unauthenticated", "Sign in to continue.")
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		user, err := h.store.sessionUser(ctx, digest)
		cancel()
		if errors.Is(err, errNotFound) {
			writeError(w, http.StatusUnauthorized, "unauthenticated", "Sign in to continue.")
			return
		}
		if err != nil {
			writeError(w, http.StatusServiceUnavailable, "unavailable", "Authentication is temporarily unavailable.")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userContextKey{}, user)))
	})
}

type userContextKey struct{}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userContextKey{}).(User)
	return user, ok
}

func (h *Handler) mutation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		origins := r.Header.Values("Origin")
		_, allowed := h.allowedOrigins[r.Header.Get("Origin")]
		if len(origins) != 1 || !allowed {
			writeError(w, http.StatusForbidden, "forbidden_origin", "Request origin is not allowed.")
			return
		}
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/json" {
			writeError(w, http.StatusUnsupportedMediaType, "invalid_content_type", "Send an application/json request.")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 8*1024)
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) passwordWork(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case h.passwordSlots <- struct{}{}:
			defer func() { <-h.passwordSlots }()
			next.ServeHTTP(w, r)
		default:
			w.Header().Set("Retry-After", "1")
			writeError(w, http.StatusTooManyRequests, "busy", "Please wait a moment and try again.")
		}
	})
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) bool {
	decoder := json.NewDecoder(r.Body)
	var raw json.RawMessage
	if err := decoder.Decode(&raw); err != nil || len(raw) == 0 || raw[0] != '{' {
		writeError(w, http.StatusBadRequest, "invalid_request", "Send a valid JSON object with the required fields.")
		return false
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_request", "Send a single JSON object.")
		return false
	}
	object := json.NewDecoder(bytes.NewReader(raw))
	object.DisallowUnknownFields()
	if err := object.Decode(destination); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "Send a valid JSON object with the required fields.")
		return false
	}
	return true
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var input credentialsRequest
	if !decodeJSON(w, r, &input) {
		return
	}
	email, valid := normalizeEmail(input.Email)
	if !valid {
		writeError(w, http.StatusUnprocessableEntity, "invalid_email", "Enter a valid email address.")
		return
	}
	if !validPassword(input.Password) {
		writeError(w, http.StatusUnprocessableEntity, "invalid_password", "Password must contain 8 to 128 characters.")
		return
	}
	passwordHash, err := argon2id.CreateHash(input.Password, passwordParams)
	if err != nil {
		writeUnavailable(w)
		return
	}
	token, digest, err := newSessionToken()
	if err != nil {
		writeUnavailable(w)
		return
	}
	expires := time.Now().UTC().Add(SessionLifetime)
	user, err := h.store.register(r.Context(), email, passwordHash, digest, expires)
	if errors.Is(err, errDuplicate) {
		writeError(w, http.StatusConflict, "registration_unavailable", "Unable to register with these details. Try signing in or use another email.")
		return
	}
	if err != nil {
		writeUnavailable(w)
		return
	}
	http.SetCookie(w, h.cookie(token, expires, int(SessionLifetime.Seconds())))
	writeUser(w, http.StatusCreated, user)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var input credentialsRequest
	if !decodeJSON(w, r, &input) {
		return
	}
	email, valid := normalizeEmail(input.Email)
	if !valid || input.Password == "" || !utf8.ValidString(input.Password) || utf8.RuneCountInString(input.Password) > 128 {
		writeInvalidCredentials(w)
		return
	}
	user, hash, err := h.store.credentials(r.Context(), email)
	missing := errors.Is(err, errNotFound)
	if err != nil && !missing {
		writeUnavailable(w)
		return
	}
	if missing {
		hash = h.dummyHash
	}
	match, err := argon2id.ComparePasswordAndHash(input.Password, hash)
	if err != nil {
		writeUnavailable(w)
		return
	}
	if !match || missing {
		writeInvalidCredentials(w)
		return
	}
	token, digest, err := newSessionToken()
	if err != nil {
		writeUnavailable(w)
		return
	}
	expires := time.Now().UTC().Add(SessionLifetime)
	if err := h.store.createSession(r.Context(), user.ID, digest, expires); err != nil {
		writeUnavailable(w)
		return
	}
	http.SetCookie(w, h.cookie(token, expires, int(SessionLifetime.Seconds())))
	writeUser(w, http.StatusOK, user)
}

func (h *Handler) me(w http.ResponseWriter, r *http.Request) {
	user, _ := UserFromContext(r.Context())
	writeUser(w, http.StatusOK, user)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	var input struct{}
	if !decodeJSON(w, r, &input) {
		return
	}
	if cookie, err := r.Cookie(h.cookieName); err == nil {
		if digest, valid := sessionTokenHash(cookie.Value); valid {
			if err := h.store.deleteSession(r.Context(), digest); err != nil {
				// Do not claim success or discard the cookie before revocation succeeds.
				writeUnavailable(w)
				return
			}
		}
	}
	http.SetCookie(w, h.cookie("", time.Unix(1, 0).UTC(), -1))
	w.WriteHeader(http.StatusNoContent)
}

func writeUser(w http.ResponseWriter, status int, user User) {
	writeJSON(w, status, struct {
		User User `json:"user"`
	}{User: user})
}

func writeInvalidCredentials(w http.ResponseWriter) {
	writeError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid email or password.")
}

func writeUnavailable(w http.ResponseWriter) {
	writeError(w, http.StatusServiceUnavailable, "unavailable", "Authentication is temporarily unavailable.")
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
