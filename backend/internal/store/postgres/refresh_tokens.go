package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type RefreshToken struct {
	ID         string
	UserID     string
	TokenHash  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *string
}

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) RefreshTokenRepository {
	return RefreshTokenRepository{db: db}
}

func (r RefreshTokenRepository) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (RefreshToken, error) {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by`

	var t RefreshToken
	err := r.db.QueryRowContext(ctx, query, userID, tokenHash, expiresAt).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.RevokedAt, &t.ReplacedBy,
	)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("create refresh token: %w", err)
	}

	return t, nil
}

func (r RefreshTokenRepository) FindByHash(ctx context.Context, tokenHash string) (RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at, replaced_by
		FROM refresh_tokens
		WHERE token_hash = $1`

	var t RefreshToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.RevokedAt, &t.ReplacedBy,
	)
	if err != nil {
		return RefreshToken{}, err
	}

	return t, nil
}

func (r RefreshTokenRepository) Revoke(ctx context.Context, id string, replacedBy *string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW(), replaced_by = $2 WHERE id = $1`,
		id, replacedBy,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}

	return nil
}

func (r RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL`,
		userID,
	)
	if err != nil {
		return fmt.Errorf("revoke all refresh tokens: %w", err)
	}

	return nil
}

func (r RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM refresh_tokens WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("delete expired refresh tokens: %w", err)
	}

	return nil
}
