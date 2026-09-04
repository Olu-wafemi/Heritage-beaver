package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/llm"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type WisdomHandler struct {
	storyRepo  postgres.StoryRepository
	wisdomRepo postgres.WisdomExtractRepository
	llm        *llm.Client
}

const wisdomSystemPrompt = `You are an expert in oral tradition, proverbs, and intergenerational wisdom.
Given a family story, identify up to 3 pieces of wisdom embedded in it: proverbs, life lessons,
values, or advice. Reply ONLY with a JSON array, no prose, no markdown fences. Each element:
{"excerpt": "short quote or paraphrase from the story", "wisdom_type": "proverb|lesson|value|advice",
"meaning": "what this wisdom teaches, one sentence", "confidence": 0.0-1.0}`

func NewWisdomHandler(storyRepo postgres.StoryRepository, wisdomRepo postgres.WisdomExtractRepository, llmClient *llm.Client) WisdomHandler {
	return WisdomHandler{
		storyRepo:  storyRepo,
		wisdomRepo: wisdomRepo,
		llm:        llmClient,
	}
}

func (h WisdomHandler) ProcessWisdom(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	storyID := r.PathValue("id")

	story, err := h.storyRepo.GetByIDForUser(r.Context(), storyID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "story not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch story")
		return
	}

	if h.llm == nil || !h.llm.Enabled() {
		h.storeFallbackExtract(w, r, &story)
		return
	}

	extracts, err := h.extractWithLLM(r, &story)
	if err != nil {
		log.Printf("wisdom extraction failed for story %s: %v", story.ID, err)
		writeError(w, http.StatusBadGateway, "failed to extract wisdom from story")
		return
	}

	res := wisdomProcessResponse{Extracts: extracts}
	if len(extracts) == 0 {
		res.Message = noWisdomMessage
	}

	writeJSON(w, http.StatusCreated, res)
}

type wisdomProcessResponse struct {
	Extracts []domain.WisdomExtract `json:"extracts"`
	Message  string                 `json:"message,omitempty"`
}

type llmExtract struct {
	Excerpt    string  `json:"excerpt"`
	WisdomType string  `json:"wisdom_type"`
	Meaning    string  `json:"meaning"`
	Confidence float64 `json:"confidence"`
}

const noWisdomMessage = "No wisdom could be drawn from this story yet. Stories work best when they include what happened, what someone said or taught, and how it changed things — proverbs, advice, and life lessons are extracted automatically."

func (h WisdomHandler) extractWithLLM(r *http.Request, story *domain.Story) ([]domain.WisdomExtract, error) {
	raw, err := h.llm.Complete(r.Context(), wisdomSystemPrompt, story.Content)
	if err != nil {
		return nil, fmt.Errorf("llm complete: %w", err)
	}

	payload, err := llm.ExtractJSON(raw)
	if err != nil {
		log.Printf("llm output unparseable for story %s (first 500 chars): %s", story.ID, truncate(raw, 500))
		return nil, fmt.Errorf("parse llm output: %w", err)
	}

	var found []llmExtract
	if err := json.Unmarshal(payload, &found); err != nil || len(found) == 0 {
		var fallback []string
		if err2 := json.Unmarshal(payload, &fallback); err2 != nil || len(fallback) == 0 {
			log.Printf("llm output unmarshal failed for story %s: %v", story.ID, err)
			return nil, fmt.Errorf("unmarshal extracts: %w", err)
		}
		for _, s := range fallback {
			found = append(found, llmExtract{Excerpt: s, WisdomType: "lesson", Meaning: s})
		}
	}

	created := make([]domain.WisdomExtract, 0, len(found))
	for _, f := range found {
		excerpt := f.Excerpt
		if utf8.RuneCountInString(excerpt) > 200 {
			excerpt = string([]rune(excerpt)[:200])
		}
		if excerpt == "" || f.Meaning == "" {
			continue
		}

		confidence := f.Confidence
		if confidence <= 0 || confidence > 1 {
			confidence = 0.5
		}

		extract, err := h.wisdomRepo.Create(r.Context(), postgres.CreateWisdomExtractParams{
			StoryID:    story.ID,
			Excerpt:    excerpt,
			WisdomType: f.WisdomType,
			Language:   story.SourceLanguage,
			Meaning:    f.Meaning,
			Confidence: confidence,
		})
		if err != nil {
			return nil, fmt.Errorf("store extract: %w", err)
		}

		created = append(created, extract)
	}

	return created, nil
}

func (h WisdomHandler) storeFallbackExtract(w http.ResponseWriter, r *http.Request, story *domain.Story) {
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

	writeJSON(w, http.StatusCreated, wisdomProcessResponse{Extracts: []domain.WisdomExtract{extract}})
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
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
		if _, err := h.storyRepo.GetByIDForUser(r.Context(), storyID, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeError(w, http.StatusNotFound, "story not found")
				return
			}

			writeError(w, http.StatusInternalServerError, "failed to fetch story")
			return
		}

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
