package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type StoryHandler struct {
	repo postgres.StoryRepository
}

type storyRequest struct {
	FamilyMemberID *string `json:"family_member_id"`
	Title          string  `json:"title"`
	Content        string  `json:"content"`
	SourceType     string  `json:"source_type"`
	SourceLanguage string  `json:"source_language"`
	Summary        string  `json:"summary"`
	OccurredAt     *string `json:"occurred_at"`
}

func NewStoryHandler(repo postgres.StoryRepository) StoryHandler {
	return StoryHandler{repo: repo}
}

func (h StoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	params.UserID = userID

	story, err := h.repo.Create(r.Context(), params)
	if err != nil {
		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "family_member_id does not reference an existing family member")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create story")
		return
	}

	writeJSON(w, http.StatusCreated, story)
}

func (h StoryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stories, err := h.repo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list stories")
		return
	}

	writeJSON(w, http.StatusOK, stories)
}

func (h StoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	story, err := h.repo.GetByIDForUser(r.Context(), r.PathValue("id"), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "story not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch story")
		return
	}

	writeJSON(w, http.StatusOK, story)
}

func (h StoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	params.ID = r.PathValue("id")
	params.UserID = userID

	story, err := h.repo.Update(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "story not found")
			return
		}

		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "family_member_id does not reference an existing family member")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update story")
		return
	}

	writeJSON(w, http.StatusOK, story)
}

func (h StoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.repo.Delete(r.Context(), r.PathValue("id"), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "story not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete story")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h StoryHandler) decodeRequest(w http.ResponseWriter, r *http.Request) (postgres.UpsertStoryParams, bool) {
	var req storyRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertStoryParams{}, false
	}

	if err := required(req.Title, "title"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertStoryParams{}, false
	}

	if err := required(req.Content, "content"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertStoryParams{}, false
	}

	if err := required(req.SourceType, "source_type"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertStoryParams{}, false
	}

	occurredAt, err := parseDate(req.OccurredAt, "occurred_at")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertStoryParams{}, false
	}

	return postgres.UpsertStoryParams{
		FamilyMemberID: req.FamilyMemberID,
		Title:          req.Title,
		Content:        req.Content,
		SourceType:     req.SourceType,
		SourceLanguage: req.SourceLanguage,
		Summary:        req.Summary,
		OccurredAt:     occurredAt,
	}, true
}
