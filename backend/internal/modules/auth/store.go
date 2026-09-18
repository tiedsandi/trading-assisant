package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User contains the only identity fields returned by the authentication API.
// Future user-owned records should reference ID, never email.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	errNotFound  = errors.New("auth record not found")
	errDuplicate = errors.New("registration unavailable")
)

type store interface {
	register(context.Context, string, string, []byte, time.Time) (User, error)
	credentials(context.Context, string) (User, string, error)
	createSession(context.Context, string, []byte, time.Time) error
	sessionUser(context.Context, []byte) (User, error)
	deleteSession(context.Context, []byte) error
}

type postgresStore struct{ pool *pgxpool.Pool }

func (s postgresStore) register(ctx context.Context, email, passwordHash string, tokenHash []byte, expiresAt time.Time) (User, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var user User
	err = tx.QueryRow(ctx, `INSERT INTO users (email, password_hash) VALUES ($1, $2)
		RETURNING id::text, email, created_at`, email, passwordHash).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_key" {
			return User{}, errDuplicate
		}
		return User{}, err
	}
	if _, err := tx.Exec(ctx, `INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`, user.ID, tokenHash, expiresAt); err != nil {
		return User{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}
	return user, nil
}

func (s postgresStore) credentials(ctx context.Context, email string) (User, string, error) {
	var user User
	var passwordHash string
	err := s.pool.QueryRow(ctx, `SELECT id::text, email, created_at, password_hash FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.CreatedAt, &passwordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", errNotFound
	}
	return user, passwordHash, err
}

func (s postgresStore) createSession(ctx context.Context, userID string, tokenHash []byte, expiresAt time.Time) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO sessions (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`, userID, tokenHash, expiresAt)
	return err
}

func (s postgresStore) sessionUser(ctx context.Context, tokenHash []byte) (User, error) {
	var user User
	err := s.pool.QueryRow(ctx, `SELECT u.id::text, u.email, u.created_at
		FROM sessions s JOIN users u ON u.id = s.user_id
		WHERE s.token_hash = $1 AND s.expires_at > CURRENT_TIMESTAMP`, tokenHash).
		Scan(&user.ID, &user.Email, &user.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, errNotFound
	}
	return user, err
}

func (s postgresStore) deleteSession(ctx context.Context, tokenHash []byte) error {
	_, err := s.pool.Exec(ctx, `DELETE FROM sessions WHERE token_hash = $1`, tokenHash)
	return err
}
