package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type VerificationToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	Consumed  bool
}

type VerificationTokenRepository struct {
	db *sql.DB
}

func NewVerificationTokenRepository(db *sql.DB) VerificationTokenRepository {
	return VerificationTokenRepository{db: db}
}

func (r VerificationTokenRepository) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (VerificationToken, error) {
	query := `
		INSERT INTO verification_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, created_at, consumed_at IS NOT NULL`

	var t VerificationToken
	err := r.db.QueryRowContext(ctx, query, userID, tokenHash, expiresAt).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.Consumed,
	)
	if err != nil {
		return VerificationToken{}, fmt.Errorf("create verification token: %w", err)
	}

	return t, nil
}

func (r VerificationTokenRepository) FindByHash(ctx context.Context, tokenHash string) (VerificationToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, consumed_at IS NOT NULL
		FROM verification_tokens
		WHERE token_hash = $1`

	var t VerificationToken
	err := r.db.QueryRowContext(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.CreatedAt, &t.Consumed,
	)
	if err != nil {
		return VerificationToken{}, err
	}

	return t, nil
}

func (r VerificationTokenRepository) Consume(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE verification_tokens SET consumed_at = NOW() WHERE id = $1 AND consumed_at IS NULL AND expires_at > NOW()`,
		id,
	)
	if err != nil {
		return fmt.Errorf("consume verification token: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("consume verification token rows: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r VerificationTokenRepository) DeleteForUser(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM verification_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("delete verification tokens: %w", err)
	}

	return nil
}
