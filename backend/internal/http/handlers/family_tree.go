package handlers

import (
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type FamilyTreeHandler struct {
	memberRepo  postgres.FamilyMemberRepository
	relRepo     postgres.RelationshipRepository
}

type familyTreeResponse struct {
	Members       []domain.FamilyMember   `json:"members"`
	Relationships []domain.Relationship   `json:"relationships"`
	MembersCount  int                     `json:"members_count"`
}

func NewFamilyTreeHandler(memberRepo postgres.FamilyMemberRepository, relRepo postgres.RelationshipRepository) FamilyTreeHandler {
	return FamilyTreeHandler{memberRepo: memberRepo, relRepo: relRepo}
}

func (h FamilyTreeHandler) Tree(w http.ResponseWriter, r *http.Request) {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	members, err := h.memberRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list family members")
		return
	}

	relationships, err := h.relRepo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list relationships")
		return
	}

	if members == nil {
		members = []domain.FamilyMember{}
	}
	if relationships == nil {
		relationships = []domain.Relationship{}
	}

	writeJSON(w, http.StatusOK, familyTreeResponse{
		Members:       members,
		Relationships: relationships,
		MembersCount:  len(members),
	})
}
