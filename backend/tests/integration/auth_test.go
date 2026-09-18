//go:build integration

package integration

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgconn"

	"trading-assistant/backend/db"
	"trading-assistant/backend/internal/app"
	"trading-assistant/backend/internal/modules/auth"
	"trading-assistant/backend/internal/platform/database"
)

func TestAuthenticationLifecycle(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Fatal("TEST_DATABASE_URL is required; use compose.test.yml")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	provider, sqlDB, err := db.NewProvider(pool)
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	if _, err := provider.Up(ctx); err != nil {
		t.Fatal(err)
	}
	if applied, err := provider.Up(ctx); err != nil || len(applied) != 0 {
		t.Fatalf("repeat migration must be a no-op: applied=%d error=%v", len(applied), err)
	}
	statuses, err := provider.Status(ctx)
	if err != nil || len(statuses) != 1 {
		t.Fatalf("migration status: entries=%d error=%v", len(statuses), err)
	}
	authHandler, err := auth.New(pool, auth.Options{AllowedOrigins: []string{"http://localhost:3000"}})
	if err != nil {
		t.Fatal(err)
	}
	handler := app.NewHandler(pool.Ping, authHandler)
	suffix := fmt.Sprint(time.Now().UnixNano())
	emailA, emailB := "auth-a-"+suffix+"@example.com", "auth-b-"+suffix+"@example.com"
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		if _, err := pool.Exec(cleanupCtx, `DELETE FROM users WHERE email = ANY($1::text[])`, []string{emailA, emailB}); err != nil {
			t.Error("test user cleanup failed")
		}
	}()
	request := func(method, path, body string, cookie *http.Cookie) *httptest.ResponseRecorder {
		r := httptest.NewRequest(method, path, strings.NewReader(body)).WithContext(ctx)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Origin", "http://localhost:3000")
		if cookie != nil {
			r.AddCookie(cookie)
		}
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)
		return w
	}
	status := func(response *httptest.ResponseRecorder, want int) {
		t.Helper()
		if response.Code != want {
			t.Fatalf("HTTP status=%d want=%d: %s", response.Code, want, response.Body.String())
		}
	}
	credentials := func(email, password string) string {
		body, err := json.Marshal(map[string]string{"email": email, "password": password})
		if err != nil {
			t.Fatal(err)
		}
		return string(body)
	}
	userFrom := func(response *httptest.ResponseRecorder) auth.User {
		t.Helper()
		var result struct {
			User auth.User `json:"user"`
		}
		if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		var shape map[string]map[string]any
		if err := json.Unmarshal(response.Body.Bytes(), &shape); err != nil {
			t.Fatal(err)
		}
		if len(shape) != 1 || len(shape["user"]) != 3 || result.User.ID == "" || result.User.Email == "" || result.User.CreatedAt.IsZero() {
			t.Fatal("unsafe or incomplete user response")
		}
		return result.User
	}

	status(request("POST", "/auth/register", credentials(emailA, "short"), nil), 422)
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM users WHERE email=$1`, emailA).Scan(&count); err != nil || count != 0 {
		t.Fatal("invalid registration wrote a user")
	}
	registered := request("POST", "/auth/register", credentials("  "+strings.ToUpper(emailA)+"  ", "abcdefgh"), nil)
	status(registered, 201)
	userA := userFrom(registered)
	if userA.Email != emailA {
		t.Fatal("email normalization mismatch")
	}
	cookies := registered.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatal("registration must set one session cookie")
	}
	cookieA := cookies[0]
	if !cookieA.HttpOnly || cookieA.MaxAge != 604800 || cookieA.Path != "/" || cookieA.SameSite != http.SameSiteLaxMode {
		t.Fatal("incorrect browser session cookie")
	}
	var passwordHash string
	if err := pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE id=$1`, userA.ID).Scan(&passwordHash); err != nil {
		t.Fatal(err)
	}
	if passwordHash == "abcdefgh" || !strings.HasPrefix(passwordHash, "$argon2id$") {
		t.Fatal("plaintext or unsupported password persistence")
	}
	if match, err := argon2id.ComparePasswordAndHash("abcdefgh", passwordHash); err != nil || !match {
		t.Fatal("stored password verification failed")
	}
	var storedDigest []byte
	var createdAt, expiresAt time.Time
	if err := pool.QueryRow(ctx, `SELECT token_hash,created_at,expires_at FROM sessions WHERE user_id=$1`, userA.ID).Scan(&storedDigest, &createdAt, &expiresAt); err != nil {
		t.Fatal(err)
	}
	expectedDigest := sha256.Sum256([]byte(cookieA.Value))
	if string(storedDigest) != string(expectedDigest[:]) || string(storedDigest) == cookieA.Value {
		t.Fatal("session token was not stored as SHA256 digest")
	}
	if lifetime := expiresAt.Sub(createdAt); lifetime < auth.SessionLifetime-3*time.Second || lifetime > auth.SessionLifetime+3*time.Second {
		t.Fatalf("server session lifetime=%s", lifetime)
	}
	status(request("POST", "/auth/register", credentials(emailA, "abcdefgh"), nil), 409)
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM sessions WHERE user_id=$1`, userA.ID).Scan(&count); err != nil || count != 1 {
		t.Fatal("duplicate registration created another session")
	}
	_, err = pool.Exec(ctx, `INSERT INTO users (email,password_hash) VALUES ($1,$2)`, emailA, passwordHash)
	if pgErr, ok := err.(*pgconn.PgError); !ok || pgErr.Code != "23505" {
		t.Fatal("database must enforce email uniqueness")
	}

	currentA := request("GET", "/auth/me", "", cookieA)
	status(currentA, 200)
	if userFrom(currentA).ID != userA.ID {
		t.Fatal("registration did not immediately authenticate user A")
	}
	registeredB := request("POST", "/auth/register", credentials(emailB, "ijklmnop"), nil)
	status(registeredB, 201)
	userB := userFrom(registeredB)
	cookieB := registeredB.Result().Cookies()[0]
	if userA.ID == userB.ID {
		t.Fatal("users must have distinct identities")
	}
	currentB := request("GET", "/auth/me", "", cookieB)
	status(currentB, 200)
	if userFrom(currentB).ID != userB.ID || userFrom(request("GET", "/auth/me", "", cookieA)).ID != userA.ID {
		t.Fatal("sessions are not isolated between users")
	}
	status(request("GET", "/auth/me", "", nil), 401)
	status(request("GET", "/auth/me", "", &http.Cookie{Name: cookieA.Name, Value: "invalid"}), 401)
	wrong := request("POST", "/auth/login", credentials(emailA, "bad-password"), nil)
	missing := request("POST", "/auth/login", credentials("missing-"+suffix+"@example.com", "bad-password"), nil)
	status(wrong, 401)
	status(missing, 401)
	if wrong.Body.String() != missing.Body.String() {
		t.Fatal("login discloses whether email exists")
	}
	login := request("POST", "/auth/login", credentials(emailA, "abcdefgh"), nil)
	status(login, 200)
	loginCookie := login.Result().Cookies()[0]
	if loginCookie.Value == cookieA.Value || userFrom(login).ID != userA.ID {
		t.Fatal("login did not create a fresh session for the correct user")
	}
	status(request("GET", "/auth/me", "", loginCookie), 200)

	// The reusable middleware passes only the server-resolved identity to protected work.
	protectedCalled := false
	protected := authHandler.RequireUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resolved, ok := auth.UserFromContext(r.Context())
		if !ok || resolved.ID != userA.ID {
			t.Error("protected handler received wrong user")
		}
		protectedCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))
	w := httptest.NewRecorder()
	protected.ServeHTTP(w, httptest.NewRequest("GET", "/protected", nil))
	status(w, 401)
	if protectedCalled {
		t.Fatal("unauthenticated request reached protected handler")
	}
	r := httptest.NewRequest("GET", "/protected", nil)
	r.AddCookie(cookieA)
	w = httptest.NewRecorder()
	protected.ServeHTTP(w, r)
	status(w, 204)
	if !protectedCalled {
		t.Fatal("authenticated request did not reach protected handler")
	}

	logout := request("POST", "/auth/logout", `{}`, loginCookie)
	status(logout, 204)
	clear := logout.Result().Cookies()[0]
	if clear.Name != loginCookie.Name || clear.MaxAge != -1 || clear.Value != "" || !clear.Expires.Before(time.Now()) {
		t.Fatal("logout cookie was not cleared")
	}
	status(request("GET", "/auth/me", "", loginCookie), 401)
	loginDigest := sha256.Sum256([]byte(loginCookie.Value))
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM sessions WHERE token_hash=$1`, loginDigest[:]).Scan(&count); err != nil || count != 0 {
		t.Fatal("logout did not revoke the persisted session")
	}
	status(request("GET", "/auth/me", "", cookieB), 200)
	status(request("GET", "/auth/me", "", cookieA), 200)
	if _, err := pool.Exec(ctx, `UPDATE sessions SET expires_at=CURRENT_TIMESTAMP-INTERVAL '1 second' WHERE token_hash=$1`, expectedDigest[:]); err != nil {
		t.Fatal(err)
	}
	status(request("GET", "/auth/me", "", cookieA), 401)
	status(request("POST", "/auth/logout", `{}`, cookieA), 204)
	status(request("POST", "/auth/logout", `{}`, nil), 204)
	status(request("GET", "/health", "", nil), 200)
	status(request("GET", "/ready", "", nil), 200)
}
