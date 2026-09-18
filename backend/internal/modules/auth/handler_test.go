package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
)

type memoryAccount struct {
	user User
	hash string
}

type memorySession struct {
	userID  string
	expires time.Time
}

type memoryStore struct {
	users    map[string]memoryAccount
	sessions map[string]memorySession
	failure  error
}

func (s *memoryStore) register(ctx context.Context, email, hash string, digest []byte, expiry time.Time) (User, error) {
	if s.failure != nil {
		return User{}, s.failure
	}
	if _, exists := s.users[email]; exists {
		return User{}, errDuplicate
	}
	user := User{ID: fmt.Sprintf("user-%d", len(s.users)+1), Email: email, CreatedAt: time.Now().UTC()}
	s.users[email] = memoryAccount{user, hash}
	return user, s.createSession(ctx, user.ID, digest, expiry)
}

func (s *memoryStore) credentials(_ context.Context, email string) (User, string, error) {
	if s.failure != nil {
		return User{}, "", s.failure
	}
	account, exists := s.users[email]
	if !exists {
		return User{}, "", errNotFound
	}
	return account.user, account.hash, nil
}

func (s *memoryStore) createSession(_ context.Context, userID string, digest []byte, expiry time.Time) error {
	if s.failure != nil {
		return s.failure
	}
	s.sessions[string(digest)] = memorySession{userID, expiry}
	return nil
}

func (s *memoryStore) sessionUser(_ context.Context, digest []byte) (User, error) {
	if s.failure != nil {
		return User{}, s.failure
	}
	session, exists := s.sessions[string(digest)]
	if !exists || !session.expires.After(time.Now()) {
		return User{}, errNotFound
	}
	for _, account := range s.users {
		if account.user.ID == session.userID {
			return account.user, nil
		}
	}
	return User{}, errNotFound
}

func (s *memoryStore) deleteSession(_ context.Context, digest []byte) error {
	if s.failure != nil {
		return s.failure
	}
	delete(s.sessions, string(digest))
	return nil
}

func testAuth(t *testing.T, secure bool) (*Handler, *memoryStore, http.Handler) {
	t.Helper()
	repository := &memoryStore{users: make(map[string]memoryAccount), sessions: make(map[string]memorySession)}
	h, err := newHandler(repository, Options{SecureCookie: secure, AllowedOrigins: []string{"http://localhost:3000"}})
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return h, repository, mux
}

func authRequest(handler http.Handler, method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	r.Header.Set("Origin", "http://localhost:3000")
	r.Header.Set("Content-Type", "application/json")
	if cookie != nil {
		r.AddCookie(cookie)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	return w
}

func requireStatus(t *testing.T, response *httptest.ResponseRecorder, status int) {
	t.Helper()
	if response.Code != status {
		t.Fatalf("status = %d, want %d: %s", response.Code, status, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("auth must not be cached")
	}
}

func TestRegistrationLoginSessionLogout(t *testing.T) {
	h, repository, mux := testAuth(t, false)
	payload := `{"email":"  Sandi@Example.COM  ","password":"abcdefgh"}`
	registered := authRequest(mux, "POST", "/auth/register", payload, nil)
	requireStatus(t, registered, 201)
	cookie := registered.Result().Cookies()[0]
	if cookie.Name != "ta_session" || !cookie.HttpOnly || cookie.Secure || cookie.Path != "/" || cookie.SameSite != http.SameSiteLaxMode || cookie.MaxAge != 604800 || cookie.Domain != "" {
		t.Fatalf("unexpected cookie attributes: name=%s HttpOnly=%v Secure=%v Path=%s SameSite=%d MaxAge=%d", cookie.Name, cookie.HttpOnly, cookie.Secure, cookie.Path, cookie.SameSite, cookie.MaxAge)
	}
	if until := time.Until(cookie.Expires); until < SessionLifetime-3*time.Second || until > SessionLifetime {
		t.Fatal("incorrect cookie expiration")
	}
	account := repository.users["sandi@example.com"]
	if account.hash == "abcdefgh" || !strings.HasPrefix(account.hash, "$argon2id$") {
		t.Fatal("password not hashed")
	}
	match, err := argon2id.ComparePasswordAndHash("abcdefgh", account.hash)
	if err != nil || !match {
		t.Fatal("stored password hash cannot verify password")
	}
	if strings.Contains(registered.Body.String(), "password") || strings.Contains(registered.Body.String(), cookie.Value) {
		t.Fatal("unsafe response fields")
	}
	digest := sha256.Sum256([]byte(cookie.Value))
	if _, exists := repository.sessions[string(digest[:])]; !exists {
		t.Fatal("session digest missing")
	}
	if _, exists := repository.sessions[cookie.Value]; exists {
		t.Fatal("raw session persisted")
	}
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", cookie), 200)
	requireStatus(t, authRequest(mux, "POST", "/auth/register", payload, nil), 409)
	wrong := authRequest(mux, "POST", "/auth/login", `{"email":"sandi@example.com","password":"wrong-password"}`, nil)
	unknown := authRequest(mux, "POST", "/auth/login", `{"email":"unknown@example.com","password":"wrong-password"}`, nil)
	requireStatus(t, wrong, 401)
	requireStatus(t, unknown, 401)
	if wrong.Body.String() != unknown.Body.String() {
		t.Fatal("login discloses whether account exists")
	}
	loggedIn := authRequest(mux, "POST", "/auth/login", payload, nil)
	requireStatus(t, loggedIn, 200)
	loginCookie := loggedIn.Result().Cookies()[0]
	if loginCookie.Value == cookie.Value {
		t.Fatal("login must create a fresh session")
	}
	var current struct {
		User User `json:"user"`
	}
	if err := json.Unmarshal(loggedIn.Body.Bytes(), &current); err != nil {
		t.Fatal(err)
	}
	if current.User.ID != account.user.ID || current.User.Email != "sandi@example.com" {
		t.Fatal("incorrect current user")
	}
	protected := h.RequireUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || user.ID != current.User.ID {
			t.Error("middleware did not provide authenticated identity")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	requireStatus(t, authRequest(protected, "GET", "/future", "", cookie), 204)
	loggedOut := authRequest(mux, "POST", "/auth/logout", `{}`, loginCookie)
	requireStatus(t, loggedOut, 204)
	cleared := loggedOut.Result().Cookies()[0]
	if cleared.Value != "" || cleared.MaxAge != -1 || cleared.Path != cookie.Path || !cleared.Expires.Before(time.Now()) {
		t.Fatal("logout did not clear cookie")
	}
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", loginCookie), 401)
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", cookie), 200)
	requireStatus(t, authRequest(mux, "POST", "/auth/logout", `{}`, nil), 204)
	requireStatus(t, authRequest(mux, "POST", "/auth/logout", `{}`, loginCookie), 204)
}

func TestValidationAndRequestShape(t *testing.T) {
	_, repository, mux := testAuth(t, false)
	for _, tc := range []struct {
		name, body string
		status     int
	}{
		{"short password", `{"email":"a@example.com","password":"1234567"}`, 422},
		{"unicode length", `{"email":"a@example.com","password":"😊😊😊😊😊😊😊"}`, 422},
		{"invalid email", `{"email":"broken","password":"abcdefgh"}`, 422},
		{"display email", `{"email":"Name <a@example.com>","password":"abcdefgh"}`, 422},
		{"missing password", `{"email":"a@example.com"}`, 422},
		{"unknown field", `{"email":"a@example.com","password":"abcdefgh","admin":true}`, 400},
		{"trailing json", `{"email":"a@example.com","password":"abcdefgh"}{}`, 400},
		{"trailing junk", `{"email":"a@example.com","password":"abcdefgh"}junk`, 400},
		{"array", `[]`, 400},
		{"wrong type", `{"email":12,"password":"abcdefgh"}`, 400},
		{"empty", ``, 400},
		{"null", `null`, 400},
		{"large body", `{"email":"a@example.com","password":"` + strings.Repeat("a", 9000) + `"}`, 400},
		{"long password", `{"email":"a@example.com","password":"` + strings.Repeat("a", 129) + `"}`, 422},
	} {
		t.Run(tc.name, func(t *testing.T) {
			requireStatus(t, authRequest(mux, "POST", "/auth/register", tc.body, nil), tc.status)
		})
	}
	if len(repository.users) != 0 || len(repository.sessions) != 0 {
		t.Fatal("invalid request changed persistence")
	}
	requireStatus(t, authRequest(mux, "POST", "/auth/logout", "null", nil), 400)
	requireStatus(t, authRequest(mux, "POST", "/auth/register", `{"email":"unicode@example.com","password":"😊😊😊😊😊😊😊😊"}`, nil), 201)
}

func TestMutationOriginAndContentType(t *testing.T) {
	_, repository, mux := testAuth(t, false)
	for _, path := range []string{"/auth/register", "/auth/login", "/auth/logout"} {
		for _, origin := range []string{"", "null", "https://attacker.example", "http://localhost:3000.attacker.example", "http://localhost:3000/"} {
			r := httptest.NewRequest("POST", path, strings.NewReader(`{}`))
			r.Header.Set("Content-Type", "application/json")
			if origin != "" {
				r.Header.Set("Origin", origin)
			}
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			requireStatus(t, w, 403)
			if w.Header().Get("Access-Control-Allow-Origin") != "" {
				t.Fatal("unexpected cross-origin grant")
			}
		}
		for _, contentType := range []string{"", "text/plain", "application/x-www-form-urlencoded"} {
			r := httptest.NewRequest("POST", path, strings.NewReader(`{}`))
			r.Header.Set("Origin", "http://localhost:3000")
			r.Header.Set("Content-Type", contentType)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)
			requireStatus(t, w, 415)
		}
	}
	if len(repository.users) != 0 || len(repository.sessions) != 0 {
		t.Fatal("rejected request changed persistence")
	}
}

func TestInvalidExpiredAndUnavailableSessions(t *testing.T) {
	_, repository, mux := testAuth(t, false)
	registered := authRequest(mux, "POST", "/auth/register", `{"email":"a@example.com","password":"abcdefgh"}`, nil)
	requireStatus(t, registered, 201)
	cookie := registered.Result().Cookies()[0]
	for _, value := range []string{"invalid", base64.RawURLEncoding.EncodeToString(make([]byte, 32))} {
		requireStatus(t, authRequest(mux, "GET", "/auth/me", "", &http.Cookie{Name: cookie.Name, Value: value}), 401)
	}
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", nil), 401)
	digest, _ := sessionTokenHash(cookie.Value)
	session := repository.sessions[string(digest)]
	session.expires = time.Now().Add(-time.Second)
	repository.sessions[string(digest)] = session
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", cookie), 401)
	repository.failure = errors.New("secret database password hash and session token")
	for _, path := range []string{"/auth/register", "/auth/login", "/auth/logout"} {
		response := authRequest(mux, "POST", path, `{"email":"a@example.com","password":"abcdefgh"}`, cookie)
		if path == "/auth/logout" {
			response = authRequest(mux, "POST", path, `{}`, cookie)
		}
		requireStatus(t, response, 503)
		if strings.Contains(response.Body.String(), "secret") || response.Header().Get("Set-Cookie") != "" {
			t.Fatal("failure exposed internals or mutated cookie")
		}
	}
	requireStatus(t, authRequest(mux, "GET", "/auth/me", "", cookie), 503)
}

func TestProductionCookieAndPasswordCapacity(t *testing.T) {
	h, _, mux := testAuth(t, true)
	response := authRequest(mux, "POST", "/auth/register", `{"email":"a@example.com","password":"abcdefgh"}`, nil)
	requireStatus(t, response, 201)
	cookie := response.Result().Cookies()[0]
	if cookie.Name != "__Host-ta_session" || !cookie.Secure || !cookie.HttpOnly || cookie.Path != "/" || cookie.Domain != "" {
		t.Fatal("production cookie security flags missing")
	}
	logout := authRequest(mux, "POST", "/auth/logout", `{}`, cookie)
	requireStatus(t, logout, 204)
	clear := logout.Result().Cookies()[0]
	if clear.Name != cookie.Name || !clear.Secure || clear.MaxAge != -1 {
		t.Fatal("production logout cookie inconsistent")
	}
	h.passwordSlots <- struct{}{}
	h.passwordSlots <- struct{}{}
	busy := authRequest(mux, "POST", "/auth/login", `{"email":"a@example.com","password":"abcdefgh"}`, nil)
	requireStatus(t, busy, 429)
	if busy.Header().Get("Retry-After") != "1" {
		t.Fatal("busy response should indicate retry")
	}
}

func TestEmailAndTokenValidation(t *testing.T) {
	for _, email := range []string{"a..b@example.com", ".a@example.com", "a.@example.com", "a@-example.com", "a@example-.com", "a@example", "a@exam_ple.com", "é@example.com", "a@" + strings.Repeat("b", 64) + ".com", strings.Repeat("a", 65) + "@example.com"} {
		if _, valid := normalizeEmail(email); valid {
			t.Errorf("accepted invalid email %q", email)
		}
	}
	for _, email := range []string{"  A.B+Label@Example.COM  ", "test@example.co.id", "a_b@example.com"} {
		if _, valid := normalizeEmail(email); !valid {
			t.Errorf("rejected valid email %q", email)
		}
	}
	first, hash, err := newSessionToken()
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := newSessionToken()
	if err != nil || first == second || len(first) != 43 || len(hash) != 32 {
		t.Fatal("invalid random session generation")
	}
	if _, valid := sessionTokenHash(first); !valid {
		t.Fatal("generated session is invalid")
	}
	if _, valid := sessionTokenHash(strings.Repeat("!", 43)); valid {
		t.Fatal("accepted invalid token encoding")
	}
	if validPassword(strings.Repeat("a", 129)) || validPassword("abcdefg") || !validPassword("abcdefgh") || !validPassword(strings.Repeat("😊", 8)) {
		t.Fatal("password policy mismatch")
	}
}
