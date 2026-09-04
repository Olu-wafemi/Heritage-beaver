package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type StoryRepository struct {
	db *sql.DB
}

type UpsertStoryParams struct {
	ID             string
	UserID         string
	FamilyMemberID *string
	Title          string
	Content        string
	SourceType     string
	SourceLanguage string
	Summary        string
	OccurredAt     *time.Time
}

func NewStoryRepository(db *sql.DB) StoryRepository {
	return StoryRepository{db: db}
}

func (r StoryRepository) Create(ctx context.Context, params UpsertStoryParams) (domain.Story, error) {
	query := `
		INSERT INTO stories (user_id, family_member_id, title, content, source_type, source_language, summary, occurred_at)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8)
		RETURNING id, user_id, family_member_id, title, content, source_type, source_language,
		          COALESCE(summary, ''), occurred_at, created_at, updated_at`

	story, err := r.queryOne(ctx, query,
		params.UserID,
		params.FamilyMemberID,
		params.Title,
		params.Content,
		params.SourceType,
		params.SourceLanguage,
		params.Summary,
		params.OccurredAt,
	)
	if err != nil {
		return domain.Story{}, fmt.Errorf("create story: %w", err)
	}

	return story, nil
}

func (r StoryRepository) ListByUserID(ctx context.Context, userID string) ([]domain.Story, error) {
	query := `
		SELECT id, user_id, family_member_id, title, content, source_type, source_language,
		       COALESCE(summary, ''), occurred_at, created_at, updated_at
		FROM stories
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list stories: %w", err)
	}
	defer rows.Close()

	stories := make([]domain.Story, 0)
	for rows.Next() {
		story, err := scanStory(rows)
		if err != nil {
			return nil, fmt.Errorf("scan story: %w", err)
		}

		stories = append(stories, story)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate stories: %w", err)
	}

	return stories, nil
}

func (r StoryRepository) GetByIDForUser(ctx context.Context, id, userID string) (domain.Story, error) {
	query := `
		SELECT id, user_id, family_member_id, title, content, source_type, source_language,
		       COALESCE(summary, ''), occurred_at, created_at, updated_at
		FROM stories
		WHERE id = $1 AND user_id = $2`

	story, err := r.queryOne(ctx, query, id, userID)
	if err != nil {
		return domain.Story{}, fmt.Errorf("get story by id: %w", err)
	}

	return story, nil
}

func (r StoryRepository) Update(ctx context.Context, params UpsertStoryParams) (domain.Story, error) {
	query := `
		UPDATE stories
		SET family_member_id = $3,
		    title = $4,
		    content = $5,
		    source_type = $6,
		    source_language = $7,
		    summary = NULLIF($8, ''),
		    occurred_at = $9,
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, family_member_id, title, content, source_type, source_language,
		          COALESCE(summary, ''), occurred_at, created_at, updated_at`

	story, err := r.queryOne(ctx, query,
		params.ID,
		params.UserID,
		params.FamilyMemberID,
		params.Title,
		params.Content,
		params.SourceType,
		params.SourceLanguage,
		params.Summary,
		params.OccurredAt,
	)
	if err != nil {
		return domain.Story{}, fmt.Errorf("update story: %w", err)
	}

	return story, nil
}

func (r StoryRepository) Delete(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM stories WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete story: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete story rows affected: %w", err)
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r StoryRepository) queryOne(ctx context.Context, query string, args ...any) (domain.Story, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanStory(row)
}

type storyScanner interface {
	Scan(dest ...any) error
}

func scanStory(scanner storyScanner) (domain.Story, error) {
	var s domain.Story
	err := scanner.Scan(
		&s.ID,
		&s.UserID,
		&s.FamilyMemberID,
		&s.Title,
		&s.Content,
		&s.SourceType,
		&s.SourceLanguage,
		&s.Summary,
		&s.OccurredAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		return domain.Story{}, err
	}

	return s, nil
}
