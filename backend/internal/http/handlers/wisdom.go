package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type WisdomHandler struct {
	storyRepo   postgres.StoryRepository
	wisdomRepo  postgres.WisdomExtractRepository
}

func NewWisdomHandler(storyRepo postgres.StoryRepository, wisdomRepo postgres.WisdomExtractRepository) WisdomHandler {
	return WisdomHandler{
		storyRepo:  storyRepo,
		wisdomRepo: wisdomRepo,
	}
}

func (h WisdomHandler) ProcessWisdom(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	storyID := r.PathValue("id")

	story, err := h.storyRepo.GetByID(r.Context(), storyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "story not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch story")
		return
	}

	excerpt := story.Content
	if utf8.RuneCountInString(excerpt) > 200 {
		excerpt = string([]rune(excerpt)[:200])
	}

	extract, err := h.wisdomRepo.Create(r.Context(), postgres.CreateWisdomExtractParams{
		StoryID:    story.ID,
		Excerpt:    excerpt,
		WisdomType: "proverb",
		Language:   story.SourceLanguage,
		Meaning:    "AI wisdom extraction will be configured with an LLM API key",
		Confidence: 0.5,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process wisdom")
		return
	}

	writeJSON(w, http.StatusCreated, extract)
}

func (h WisdomHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	storyID := r.URL.Query().Get("story_id")

	var extracts []domain.WisdomExtract
	var err error

	if storyID != "" {
		extracts, err = h.wisdomRepo.ListByStoryID(r.Context(), storyID)
	} else {
		extracts, err = h.wisdomRepo.ListByUserID(r.Context(), userID)
	}

	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list wisdom extracts")
		return
	}

	writeJSON(w, http.StatusOK, extracts)
}
