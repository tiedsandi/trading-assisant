package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/alexedwards/argon2id"
)

const SessionLifetime = 7 * 24 * time.Hour

// Explicit parameters also apply to the dummy hash used for unknown accounts.
var passwordParams = &argon2id.Params{
	Memory: 64 * 1024, Iterations: 3, Parallelism: 2, SaltLength: 16, KeyLength: 32,
}

// This intentionally uses a practical ASCII email shape, matching the frontend.
var emailPattern = regexp.MustCompile("^[a-z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(?:\\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)+$")

func normalizeEmail(email string) (string, bool) {
	email = strings.ToLower(strings.TrimSpace(email))
	parts := strings.Split(email, "@")
	if len(email) > 254 || len(parts) != 2 || len(parts[0]) > 64 || !emailPattern.MatchString(email) {
		return "", false
	}
	if strings.HasPrefix(parts[0], ".") || strings.HasSuffix(parts[0], ".") || strings.Contains(parts[0], "..") {
		return "", false
	}
	for _, label := range strings.Split(parts[1], ".") {
		if len(label) > 63 {
			return "", false
		}
	}
	return email, true
}

func validPassword(password string) bool {
	length := utf8.RuneCountInString(password)
	return utf8.ValidString(password) && length >= 8 && length <= 128
}

func newSessionToken() (string, []byte, error) {
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return "", nil, err
	}
	token := base64.RawURLEncoding.EncodeToString(random[:])
	digest := sha256.Sum256([]byte(token))
	return token, digest[:], nil
}

func sessionTokenHash(token string) ([]byte, bool) {
	if len(token) != 43 {
		return nil, false
	}
	raw, err := base64.RawURLEncoding.Strict().DecodeString(token)
	if err != nil || len(raw) != 32 {
		return nil, false
	}
	digest := sha256.Sum256([]byte(token))
	return digest[:], true
}

func (h *Handler) cookie(value string, expires time.Time, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name: h.cookieName, Value: value, Path: "/", Expires: expires,
		MaxAge: maxAge, HttpOnly: true, Secure: h.secureCookie, SameSite: http.SameSiteLaxMode,
	}
}
