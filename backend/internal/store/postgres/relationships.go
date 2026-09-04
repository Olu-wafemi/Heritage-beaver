package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type RelationshipRepository struct {
	db *sql.DB
}

type UpsertRelationshipParams struct {
	ID               string
	UserID           string
	SourceMemberID   string
	TargetMemberID   string
	RelationshipType string
}

func NewRelationshipRepository(db *sql.DB) RelationshipRepository {
	return RelationshipRepository{db: db}
}

func (r RelationshipRepository) Create(ctx context.Context, params UpsertRelationshipParams) (domain.Relationship, error) {
	query := `
		INSERT INTO relationships (user_id, source_member_id, target_member_id, relationship_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, source_member_id, target_member_id, relationship_type, created_at`

	rel, err := r.queryOne(ctx, query,
		params.UserID,
		params.SourceMemberID,
		params.TargetMemberID,
		params.RelationshipType,
	)
	if err != nil {
		return domain.Relationship{}, fmt.Errorf("create relationship: %w", err)
	}

	return rel, nil
}

func (r RelationshipRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Relationship, error) {
	query := `
		SELECT id, user_id, source_member_id, target_member_id, relationship_type, created_at
		FROM relationships
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list relationships: %w", err)
	}
	defer rows.Close()

	relationships := make([]domain.Relationship, 0)
	for rows.Next() {
		rel, err := scanRelationship(rows)
		if err != nil {
			return nil, fmt.Errorf("scan relationship: %w", err)
		}

		relationships = append(relationships, rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate relationships: %w", err)
	}

	return relationships, nil
}

func (r RelationshipRepository) GetByIDForUser(ctx context.Context, id, userID string) (domain.Relationship, error) {
	query := `
		SELECT id, user_id, source_member_id, target_member_id, relationship_type, created_at
		FROM relationships
		WHERE id = $1 AND user_id = $2`

	rel, err := r.queryOne(ctx, query, id, userID)
	if err != nil {
		return domain.Relationship{}, fmt.Errorf("get relationship by id: %w", err)
	}

	return rel, nil
}

func (r RelationshipRepository) Update(ctx context.Context, params UpsertRelationshipParams) (domain.Relationship, error) {
	query := `
		UPDATE relationships
		SET source_member_id = $3,
		    target_member_id = $4,
		    relationship_type = $5
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, source_member_id, target_member_id, relationship_type, created_at`

	rel, err := r.queryOne(ctx, query,
		params.ID,
		params.UserID,
		params.SourceMemberID,
		params.TargetMemberID,
		params.RelationshipType,
	)
	if err != nil {
		return domain.Relationship{}, fmt.Errorf("update relationship: %w", err)
	}

	return rel, nil
}

func (r RelationshipRepository) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM relationships WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete relationship: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete relationship rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r RelationshipRepository) queryOne(ctx context.Context, query string, args ...any) (domain.Relationship, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanRelationship(row)
}

type relationshipScanner interface {
	Scan(dest ...any) error
}

func scanRelationship(scanner relationshipScanner) (domain.Relationship, error) {
	var rel domain.Relationship
	err := scanner.Scan(
		&rel.ID,
		&rel.UserID,
		&rel.SourceMemberID,
		&rel.TargetMemberID,
		&rel.RelationshipType,
		&rel.CreatedAt,
	)
	if err != nil {
		return domain.Relationship{}, err
	}

	return rel, nil
}
