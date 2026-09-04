package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type FamilyMemberRepository struct {
	db *sql.DB
}

type UpsertFamilyMemberParams struct {
	ID              string
	UserID          string
	FirstName       string
	LastName        string
	DisplayName     string
	Gender          string
	BirthDate       *time.Time
	DeathDate       *time.Time
	BirthPlace      string
	Biography       string
	IsLiving        bool
	PrimaryLanguage string
}

func NewFamilyMemberRepository(db *sql.DB) FamilyMemberRepository {
	return FamilyMemberRepository{db: db}
}

func (r FamilyMemberRepository) Create(ctx context.Context, params UpsertFamilyMemberParams) (domain.FamilyMember, error) {
	query := `
		INSERT INTO family_members (
			user_id, first_name, last_name, display_name, gender, birth_date, death_date,
			birth_place, biography, is_living, primary_language
		)
		VALUES ($1, $2, NULLIF($3, ''), $4, NULLIF($5, ''), $6, $7, NULLIF($8, ''), NULLIF($9, ''), $10, NULLIF($11, ''))
		RETURNING id, user_id, first_name, COALESCE(last_name, ''), display_name, COALESCE(gender, ''),
		          birth_date, death_date, COALESCE(birth_place, ''), COALESCE(biography, ''), is_living,
		          COALESCE(primary_language, ''), created_at, updated_at`

	member, err := r.queryOne(ctx, query,
		params.UserID,
		params.FirstName,
		params.LastName,
		params.DisplayName,
		params.Gender,
		params.BirthDate,
		params.DeathDate,
		params.BirthPlace,
		params.Biography,
		params.IsLiving,
		params.PrimaryLanguage,
	)
	if err != nil {
		return domain.FamilyMember{}, fmt.Errorf("create family member: %w", err)
	}

	return member, nil
}

func (r FamilyMemberRepository) ListByUserID(ctx context.Context, userID string) ([]domain.FamilyMember, error) {
	query := `
		SELECT id, user_id, first_name, COALESCE(last_name, ''), display_name, COALESCE(gender, ''),
		       birth_date, death_date, COALESCE(birth_place, ''), COALESCE(biography, ''), is_living,
		       COALESCE(primary_language, ''), created_at, updated_at
		FROM family_members
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list family members: %w", err)
	}
	defer rows.Close()

	members := make([]domain.FamilyMember, 0)
	for rows.Next() {
		member, err := scanFamilyMember(rows)
		if err != nil {
			return nil, fmt.Errorf("scan family member: %w", err)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate family members: %w", err)
	}

	return members, nil
}

func (r FamilyMemberRepository) GetByIDForUser(ctx context.Context, id, userID string) (domain.FamilyMember, error) {
	query := `
		SELECT id, user_id, first_name, COALESCE(last_name, ''), display_name, COALESCE(gender, ''),
		       birth_date, death_date, COALESCE(birth_place, ''), COALESCE(biography, ''), is_living,
		       COALESCE(primary_language, ''), created_at, updated_at
		FROM family_members
		WHERE id = $1 AND user_id = $2`

	member, err := r.queryOne(ctx, query, id, userID)
	if err != nil {
		return domain.FamilyMember{}, fmt.Errorf("get family member by id: %w", err)
	}

	return member, nil
}

func (r FamilyMemberRepository) Update(ctx context.Context, params UpsertFamilyMemberParams) (domain.FamilyMember, error) {
	query := `
		UPDATE family_members
		SET first_name = $3,
		    last_name = NULLIF($4, ''),
		    display_name = $5,
		    gender = NULLIF($6, ''),
		    birth_date = $7,
		    death_date = $8,
		    birth_place = NULLIF($9, ''),
		    biography = NULLIF($10, ''),
		    is_living = $11,
		    primary_language = NULLIF($12, ''),
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, first_name, COALESCE(last_name, ''), display_name, COALESCE(gender, ''),
		          birth_date, death_date, COALESCE(birth_place, ''), COALESCE(biography, ''), is_living,
		          COALESCE(primary_language, ''), created_at, updated_at`

	member, err := r.queryOne(ctx, query,
		params.ID,
		params.UserID,
		params.FirstName,
		params.LastName,
		params.DisplayName,
		params.Gender,
		params.BirthDate,
		params.DeathDate,
		params.BirthPlace,
		params.Biography,
		params.IsLiving,
		params.PrimaryLanguage,
	)
	if err != nil {
		return domain.FamilyMember{}, fmt.Errorf("update family member: %w", err)
	}

	return member, nil
}

func (r FamilyMemberRepository) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM family_members WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete family member: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete family member rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r FamilyMemberRepository) queryOne(ctx context.Context, query string, args ...any) (domain.FamilyMember, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanFamilyMember(row)
}

type familyMemberScanner interface {
	Scan(dest ...any) error
}

func scanFamilyMember(scanner familyMemberScanner) (domain.FamilyMember, error) {
	var member domain.FamilyMember
	err := scanner.Scan(
		&member.ID,
		&member.UserID,
		&member.FirstName,
		&member.LastName,
		&member.DisplayName,
		&member.Gender,
		&member.BirthDate,
		&member.DeathDate,
		&member.BirthPlace,
		&member.Biography,
		&member.IsLiving,
		&member.PrimaryLanguage,
		&member.CreatedAt,
		&member.UpdatedAt,
	)
	if err != nil {
		return domain.FamilyMember{}, err
	}

	return member, nil
}
