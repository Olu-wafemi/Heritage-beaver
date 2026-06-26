package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type RelationshipHandler struct {
	repo postgres.RelationshipRepository
}

type relationshipRequest struct {
	UserID           string `json:"user_id"`
	SourceMemberID   string `json:"source_member_id"`
	TargetMemberID   string `json:"target_member_id"`
	RelationshipType string `json:"relationship_type"`
}

func NewRelationshipHandler(repo postgres.RelationshipRepository) RelationshipHandler {
	return RelationshipHandler{repo: repo}
}

func (h RelationshipHandler) Create(w http.ResponseWriter, r *http.Request) {
	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	rel, err := h.repo.Create(r.Context(), params)
	if err != nil {
		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "source_member_id or target_member_id does not reference an existing family member")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create relationship")
		return
	}

	writeJSON(w, http.StatusCreated, rel)
}

func (h RelationshipHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	relationships, err := h.repo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list relationships")
		return
	}

	writeJSON(w, http.StatusOK, relationships)
}

func (h RelationshipHandler) Get(w http.ResponseWriter, r *http.Request) {
	rel, err := h.repo.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "relationship not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch relationship")
		return
	}

	writeJSON(w, http.StatusOK, rel)
}

func (h RelationshipHandler) Update(w http.ResponseWriter, r *http.Request) {
	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	params.ID = r.PathValue("id")
	rel, err := h.repo.Update(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "relationship not found")
			return
		}

		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "source_member_id or target_member_id does not reference an existing family member")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update relationship")
		return
	}

	writeJSON(w, http.StatusOK, rel)
}

func (h RelationshipHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.repo.Delete(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "relationship not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete relationship")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h RelationshipHandler) decodeRequest(w http.ResponseWriter, r *http.Request) (postgres.UpsertRelationshipParams, bool) {
	var req relationshipRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertRelationshipParams{}, false
	}

	if err := required(req.UserID, "user_id"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertRelationshipParams{}, false
	}

	if err := required(req.SourceMemberID, "source_member_id"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertRelationshipParams{}, false
	}

	if err := required(req.TargetMemberID, "target_member_id"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertRelationshipParams{}, false
	}

	if err := required(req.RelationshipType, "relationship_type"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertRelationshipParams{}, false
	}

	return postgres.UpsertRelationshipParams{
		UserID:           req.UserID,
		SourceMemberID:   req.SourceMemberID,
		TargetMemberID:   req.TargetMemberID,
		RelationshipType: req.RelationshipType,
	}, true
}
