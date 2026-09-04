package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type MythChapterRepository struct {
	db *sql.DB
}

type UpsertMythChapterParams struct {
	ID          string
	UserID      string
	Title       string
	Theme       string
	ChapterType string
	Narrative   string
}

func NewMythChapterRepository(db *sql.DB) MythChapterRepository {
	return MythChapterRepository{db: db}
}

func (r MythChapterRepository) Create(ctx context.Context, params UpsertMythChapterParams) (domain.MythChapter, error) {
	query := `
		INSERT INTO myth_chapters (user_id, title, theme, chapter_type, narrative)
		VALUES ($1, NULLIF($2, ''), COALESCE(NULLIF($3, ''), 'untold'), COALESCE(NULLIF($4, ''), 'origin'), $5)
		RETURNING id, user_id, title, theme, chapter_type, narrative, created_at, updated_at`

	chapter, err := r.queryOne(ctx, query,
		params.UserID,
		params.Title,
		params.Theme,
		params.ChapterType,
		params.Narrative,
	)
	if err != nil {
		return domain.MythChapter{}, fmt.Errorf("create myth chapter: %w", err)
	}

	return chapter, nil
}

func (r MythChapterRepository) ListByUserID(ctx context.Context, userID string) ([]domain.MythChapter, error) {
	query := `
		SELECT id, user_id, title, theme, chapter_type, narrative, created_at, updated_at
		FROM myth_chapters
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list myth chapters: %w", err)
	}
	defer rows.Close()

	chapters := make([]domain.MythChapter, 0)
	for rows.Next() {
		chapter, err := scanMythChapter(rows)
		if err != nil {
			return nil, fmt.Errorf("scan myth chapter: %w", err)
		}

		chapters = append(chapters, chapter)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate myth chapters: %w", err)
	}

	return chapters, nil
}

func (r MythChapterRepository) GetByIDForUser(ctx context.Context, id, userID string) (domain.MythChapter, error) {
	query := `
		SELECT id, user_id, title, theme, chapter_type, narrative, created_at, updated_at
		FROM myth_chapters
		WHERE id = $1 AND user_id = $2`

	chapter, err := r.queryOne(ctx, query, id, userID)
	if err != nil {
		return domain.MythChapter{}, fmt.Errorf("get myth chapter by id: %w", err)
	}

	return chapter, nil
}

func (r MythChapterRepository) Update(ctx context.Context, params UpsertMythChapterParams) (domain.MythChapter, error) {
	query := `
		UPDATE myth_chapters
		SET title = $2,
		    theme = $3,
		    chapter_type = $4,
		    narrative = $5,
		    updated_at = NOW()
		WHERE id = $1 AND user_id = $6
		RETURNING id, user_id, title, theme, chapter_type, narrative, created_at, updated_at`

	chapter, err := r.queryOne(ctx, query,
		params.ID,
		params.Title,
		params.Theme,
		params.ChapterType,
		params.Narrative,
		params.UserID,
	)
	if err != nil {
		return domain.MythChapter{}, fmt.Errorf("update myth chapter: %w", err)
	}

	return chapter, nil
}

func (r MythChapterRepository) queryOne(ctx context.Context, query string, args ...any) (domain.MythChapter, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanMythChapter(row)
}

type mythChapterScanner interface {
	Scan(dest ...any) error
}

func scanMythChapter(scanner mythChapterScanner) (domain.MythChapter, error) {
	var c domain.MythChapter
	err := scanner.Scan(
		&c.ID,
		&c.UserID,
		&c.Title,
		&c.Theme,
		&c.ChapterType,
		&c.Narrative,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return domain.MythChapter{}, err
	}

	return c, nil
}
