package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

type CreateUserParams struct {
	Email          string
	DisplayName    string
	PrimaryCulture string
	PasswordHash   string
}

type UpdateUserParams struct {
	ID             string
	DisplayName    string
	PrimaryCulture string
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

func (r UserRepository) Create(ctx context.Context, params CreateUserParams) (domain.User, error) {
	query := `
		INSERT INTO users (email, display_name, primary_culture, password_hash)
		VALUES ($1, $2, COALESCE(NULLIF($3, ''), 'yoruba'), $4)
		RETURNING id, email, display_name, primary_culture, password_hash, created_at, updated_at`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, params.Email, params.DisplayName, params.PrimaryCulture, params.PasswordHash).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PrimaryCulture,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r UserRepository) List(ctx context.Context) ([]domain.User, error) {
	query := `
		SELECT id, email, display_name, primary_culture, password_hash, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Email, &user.DisplayName, &user.PrimaryCulture, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return users, nil
}

func (r UserRepository) GetByID(ctx context.Context, id string) (domain.User, error) {
	query := `
		SELECT id, email, display_name, primary_culture, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PrimaryCulture,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, display_name, primary_culture, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PrimaryCulture,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}

func (r UserRepository) Update(ctx context.Context, params UpdateUserParams) (domain.User, error) {
	query := `
		UPDATE users
		SET display_name = $2,
		    primary_culture = COALESCE(NULLIF($3, ''), primary_culture),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, email, display_name, primary_culture, password_hash, created_at, updated_at`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, params.ID, params.DisplayName, params.PrimaryCulture).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PrimaryCulture,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf("update user: %w", err)
	}

	return user, nil
}

func (r UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
