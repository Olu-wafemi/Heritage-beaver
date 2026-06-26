package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
)

type WisdomExtractRepository struct {
	db *sql.DB
}

type CreateWisdomExtractParams struct {
	StoryID    string
	Excerpt    string
	WisdomType string
	Language   string
	Meaning    string
	Confidence float64
}

func NewWisdomExtractRepository(db *sql.DB) WisdomExtractRepository {
	return WisdomExtractRepository{db: db}
}

func (r WisdomExtractRepository) Create(ctx context.Context, params CreateWisdomExtractParams) (domain.WisdomExtract, error) {
	query := `
		INSERT INTO wisdom_extracts (story_id, excerpt, wisdom_type, language, meaning, confidence)
		VALUES ($1, $2, $3, COALESCE(NULLIF($4, ''), 'en'), $5, $6)
		RETURNING id, story_id, excerpt, wisdom_type, COALESCE(language, 'en'), COALESCE(meaning, ''), confidence, created_at`

	extract, err := r.queryOne(ctx, query,
		params.StoryID,
		params.Excerpt,
		params.WisdomType,
		params.Language,
		params.Meaning,
		params.Confidence,
	)
	if err != nil {
		return domain.WisdomExtract{}, fmt.Errorf("create wisdom extract: %w", err)
	}

	return extract, nil
}

func (r WisdomExtractRepository) ListByStoryID(ctx context.Context, storyID string) ([]domain.WisdomExtract, error) {
	query := `
		SELECT id, story_id, excerpt, wisdom_type, COALESCE(language, 'en'), COALESCE(meaning, ''), confidence, created_at
		FROM wisdom_extracts
		WHERE story_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, storyID)
	if err != nil {
		return nil, fmt.Errorf("list wisdom extracts by story: %w", err)
	}
	defer rows.Close()

	extracts := make([]domain.WisdomExtract, 0)
	for rows.Next() {
		extract, err := scanWisdomExtract(rows)
		if err != nil {
			return nil, fmt.Errorf("scan wisdom extract: %w", err)
		}

		extracts = append(extracts, extract)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wisdom extracts: %w", err)
	}

	return extracts, nil
}

func (r WisdomExtractRepository) ListByUserID(ctx context.Context, userID string) ([]domain.WisdomExtract, error) {
	query := `
		SELECT w.id, w.story_id, w.excerpt, w.wisdom_type, COALESCE(w.language, 'en'), COALESCE(w.meaning, ''), w.confidence, w.created_at
		FROM wisdom_extracts w
		JOIN stories s ON w.story_id = s.id
		WHERE s.user_id = $1
		ORDER BY w.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list wisdom extracts by user: %w", err)
	}
	defer rows.Close()

	extracts := make([]domain.WisdomExtract, 0)
	for rows.Next() {
		extract, err := scanWisdomExtract(rows)
		if err != nil {
			return nil, fmt.Errorf("scan wisdom extract: %w", err)
		}

		extracts = append(extracts, extract)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wisdom extracts: %w", err)
	}

	return extracts, nil
}

func (r WisdomExtractRepository) queryOne(ctx context.Context, query string, args ...any) (domain.WisdomExtract, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	return scanWisdomExtract(row)
}

type wisdomExtractScanner interface {
	Scan(dest ...any) error
}

func scanWisdomExtract(scanner wisdomExtractScanner) (domain.WisdomExtract, error) {
	var e domain.WisdomExtract
	err := scanner.Scan(
		&e.ID,
		&e.StoryID,
		&e.Excerpt,
		&e.WisdomType,
		&e.Language,
		&e.Meaning,
		&e.Confidence,
		&e.CreatedAt,
	)
	if err != nil {
		return domain.WisdomExtract{}, err
	}

	return e, nil
}
