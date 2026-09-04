package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/llm"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type MythologyHandler struct {
	storyRepo postgres.StoryRepository
	mythRepo  postgres.MythChapterRepository
	llm       *llm.Client
}

const mythSystemPrompt = `You are a keeper of family mythology in the tradition of African oral storytelling.
Given a set of family stories, weave them into ONE myth chapter: a timeless episode with
moral weight, vivid but respectful of the source facts. Reply ONLY with a JSON object, no prose,
no markdown fences: {"title": "chapter title", "theme": "one-word theme e.g. resilience",
"narrative": "3-6 paragraphs of mythic narrative"}`

type generateChapterRequest struct {
	Theme       string   `json:"theme"`
	ChapterType string   `json:"chapter_type"`
	StoryIDs    []string `json:"story_ids"`
}

type llmChapter struct {
	Title     string `json:"title"`
	Theme     string `json:"theme"`
	Narrative string `json:"narrative"`
}

type updateChapterRequest struct {
	Title       string `json:"title"`
	Theme       string `json:"theme"`
	ChapterType string `json:"chapter_type"`
	Narrative   string `json:"narrative"`
}

func NewMythologyHandler(storyRepo postgres.StoryRepository, mythRepo postgres.MythChapterRepository, llmClient *llm.Client) MythologyHandler {
	return MythologyHandler{storyRepo: storyRepo, mythRepo: mythRepo, llm: llmClient}
}

func (h MythologyHandler) Generate(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req generateChapterRequest
	if err := decodeJSON(r, &req); err != nil {
		if err.Error() == "request body is required" {
			req = generateChapterRequest{}
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	stories, err := h.storiesForGeneration(r, userID, req.StoryIDs)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to load stories")
		return
	}
	if len(stories) == 0 {
		writeError(w, http.StatusBadRequest, "at least one story is required to generate a myth chapter")
		return
	}

	chapterType := strings.TrimSpace(req.ChapterType)
	if chapterType == "" {
		chapterType = "origin"
	}

	if h.llm == nil || !h.llm.Enabled() {
		h.storeFallbackChapter(w, r, userID, chapterType, req.Theme, stories)
		return
	}

	chapter, err := h.generateWithLLM(r, userID, chapterType, req.Theme, stories)
	if err != nil {
		log.Printf("myth generation failed for user %s: %v", userID, err)
		writeError(w, http.StatusBadGateway, "failed to generate myth chapter")
		return
	}

	writeJSON(w, http.StatusCreated, chapter)
}

func (h MythologyHandler) storiesForGeneration(r *http.Request, userID string, storyIDs []string) ([]domain.Story, error) {
	if len(storyIDs) == 0 {
		return h.storyRepo.ListByUserID(r.Context(), userID)
	}

	stories := make([]domain.Story, 0, len(storyIDs))
	for _, id := range storyIDs {
		story, err := h.storyRepo.GetByIDForUser(r.Context(), id, userID)
		if err != nil {
			return nil, fmt.Errorf("get story %s: %w", id, err)
		}
		stories = append(stories, story)
	}

	return stories, nil
}

func (h MythologyHandler) generateWithLLM(r *http.Request, userID, chapterType, theme string, stories []domain.Story) (domain.MythChapter, error) {
	var sb strings.Builder
	if theme != "" {
		sb.WriteString("Requested theme: " + theme + "\n\n")
	}
	for i, s := range stories {
		sb.WriteString(fmt.Sprintf("Story %d — %s:\n%s\n\n", i+1, s.Title, s.Content))
		if sb.Len() > 12000 {
			break
		}
	}

	raw, err := h.llm.Complete(r.Context(), mythSystemPrompt, sb.String())
	if err != nil {
		return domain.MythChapter{}, fmt.Errorf("llm complete: %w", err)
	}

	payload, err := llm.ExtractJSON(raw)
	if err != nil {
		log.Printf("myth llm output unparseable (first 500 chars): %s", truncate(raw, 500))
		return domain.MythChapter{}, fmt.Errorf("parse llm output: %w", err)
	}

	// Model may return an object or a single-element array.
	var found llmChapter
	if err := json.Unmarshal(payload, &found); err != nil || found.Narrative == "" {
		var arr []llmChapter
		if err2 := json.Unmarshal(payload, &arr); err2 != nil || len(arr) == 0 || arr[0].Narrative == "" {
			return domain.MythChapter{}, fmt.Errorf("unmarshal chapter: %v", err)
		}
		found = arr[0]
	}

	if strings.TrimSpace(found.Title) == "" {
		found.Title = "Untitled Chapter"
	}
	if strings.TrimSpace(found.Theme) == "" {
		found.Theme = firstNonEmpty(theme, "untold")
	}

	return h.mythRepo.Create(r.Context(), postgres.UpsertMythChapterParams{
		UserID:      userID,
		Title:       found.Title,
		Theme:       found.Theme,
		ChapterType: chapterType,
		Narrative:   found.Narrative,
	})
}

func (h MythologyHandler) storeFallbackChapter(w http.ResponseWriter, r *http.Request, userID, chapterType, theme string, stories []domain.Story) {
	titles := make([]string, 0, len(stories))
	for _, s := range stories {
		titles = append(titles, s.Title)
	}
	narrative := fmt.Sprintf("Gathered from %d stor%s (%s). Configure an LLM API key to weave these into a living myth.",
		len(stories), plural(len(stories)), strings.Join(titles, "; "))

	chapter, err := h.mythRepo.Create(r.Context(), postgres.UpsertMythChapterParams{
		UserID:      userID,
		Title:       "Unwoven Chapter",
		Theme:       firstNonEmpty(theme, "untold"),
		ChapterType: chapterType,
		Narrative:   narrative,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create myth chapter")
		return
	}

	writeJSON(w, http.StatusCreated, chapter)
}

func (h MythologyHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapters, err := h.mythRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list myth chapters")
		return
	}

	writeJSON(w, http.StatusOK, chapters)
}

func (h MythologyHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	chapter, err := h.mythRepo.GetByIDForUser(r.Context(), r.PathValue("id"), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "myth chapter not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch myth chapter")
		return
	}

	writeJSON(w, http.StatusOK, chapter)
}

func (h MythologyHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req updateChapterRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Title, "title"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Narrative, "narrative"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	chapter, err := h.mythRepo.Update(r.Context(), postgres.UpsertMythChapterParams{
		ID:          r.PathValue("id"),
		UserID:      userID,
		Title:       req.Title,
		Theme:       req.Theme,
		ChapterType: req.ChapterType,
		Narrative:   req.Narrative,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "myth chapter not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update myth chapter")
		return
	}

	writeJSON(w, http.StatusOK, chapter)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func plural(n int) string {
	if n == 1 {
		return "y"
	}
	return "ies"
}
