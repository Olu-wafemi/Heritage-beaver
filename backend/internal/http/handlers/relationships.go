package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type RelationshipHandler struct {
	repo    postgres.RelationshipRepository
	members postgres.FamilyMemberRepository
}

type relationshipRequest struct {
	SourceMemberID   string `json:"source_member_id"`
	TargetMemberID   string `json:"target_member_id"`
	RelationshipType string `json:"relationship_type"`
}

func NewRelationshipHandler(repo postgres.RelationshipRepository, members postgres.FamilyMemberRepository) RelationshipHandler {
	return RelationshipHandler{repo: repo, members: members}
}

func (h RelationshipHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	if err := h.verifyMembersOwned(r, userID, params.SourceMemberID, params.TargetMemberID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
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
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
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
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rel, err := h.repo.GetByIDForUser(r.Context(), r.PathValue("id"), userID)
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
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	if err := h.verifyMembersOwned(r, userID, params.SourceMemberID, params.TargetMemberID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	params.ID = r.PathValue("id")
	params.UserID = userID
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
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	err := h.repo.Delete(r.Context(), r.PathValue("id"), userID)
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

// verifyMembersOwned ensures both linked members belong to the requester and
// that a person is not linked to themselves.
func (h RelationshipHandler) verifyMembersOwned(r *http.Request, userID, sourceID, targetID string) error {
	if sourceID == targetID {
		return errors.New("a person can't be related to themselves")
	}

	for _, id := range []string{sourceID, targetID} {
		if _, err := h.members.GetByIDForUser(r.Context(), id, userID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return errors.New("source_member_id or target_member_id does not reference one of your family members")
			}

			return errors.New("failed to verify family members")
		}
	}

	return nil
}

func (h RelationshipHandler) decodeRequest(w http.ResponseWriter, r *http.Request) (postgres.UpsertRelationshipParams, bool) {
	var req relationshipRequest
	if err := decodeJSON(r, &req); err != nil {
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
		SourceMemberID:   req.SourceMemberID,
		TargetMemberID:   req.TargetMemberID,
		RelationshipType: req.RelationshipType,
	}, true
}
