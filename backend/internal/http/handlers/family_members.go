package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type FamilyMemberHandler struct {
	repo postgres.FamilyMemberRepository
}

type familyMemberRequest struct {
	UserID          string  `json:"user_id"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	DisplayName     string  `json:"display_name"`
	Gender          string  `json:"gender"`
	BirthDate       *string `json:"birth_date"`
	DeathDate       *string `json:"death_date"`
	BirthPlace      string  `json:"birth_place"`
	Biography       string  `json:"biography"`
	IsLiving        *bool   `json:"is_living"`
	PrimaryLanguage string  `json:"primary_language"`
}

func NewFamilyMemberHandler(repo postgres.FamilyMemberRepository) FamilyMemberHandler {
	return FamilyMemberHandler{repo: repo}
}

func (h FamilyMemberHandler) Create(w http.ResponseWriter, r *http.Request) {
	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	familyMember, err := h.repo.Create(r.Context(), params)
	if err != nil {
		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "user_id does not reference an existing user")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create family member")
		return
	}

	writeJSON(w, http.StatusCreated, familyMember)
}

func (h FamilyMemberHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "user_id query parameter is required")
		return
	}

	familyMembers, err := h.repo.ListByUserID(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list family members")
		return
	}

	writeJSON(w, http.StatusOK, familyMembers)
}

func (h FamilyMemberHandler) Get(w http.ResponseWriter, r *http.Request) {
	familyMember, err := h.repo.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "family member not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch family member")
		return
	}

	writeJSON(w, http.StatusOK, familyMember)
}

func (h FamilyMemberHandler) Update(w http.ResponseWriter, r *http.Request) {
	params, ok := h.decodeRequest(w, r)
	if !ok {
		return
	}

	params.ID = r.PathValue("id")
	familyMember, err := h.repo.Update(r.Context(), params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "family member not found")
			return
		}

		if postgres.IsForeignKeyViolation(err) {
			writeError(w, http.StatusBadRequest, "user_id does not reference an existing user")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update family member")
		return
	}

	writeJSON(w, http.StatusOK, familyMember)
}

func (h FamilyMemberHandler) Delete(w http.ResponseWriter, r *http.Request) {
	err := h.repo.Delete(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "family member not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete family member")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h FamilyMemberHandler) decodeRequest(w http.ResponseWriter, r *http.Request) (postgres.UpsertFamilyMemberParams, bool) {
	var req familyMemberRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	if err := required(req.UserID, "user_id"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	if err := required(req.FirstName, "first_name"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	if err := required(req.DisplayName, "display_name"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	birthDate, err := parseDate(req.BirthDate, "birth_date")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	deathDate, err := parseDate(req.DeathDate, "death_date")
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return postgres.UpsertFamilyMemberParams{}, false
	}

	return postgres.UpsertFamilyMemberParams{
		UserID:          req.UserID,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DisplayName:     req.DisplayName,
		Gender:          req.Gender,
		BirthDate:       birthDate,
		DeathDate:       deathDate,
		BirthPlace:      req.BirthPlace,
		Biography:       req.Biography,
		IsLiving:        boolValue(req.IsLiving, true),
		PrimaryLanguage: req.PrimaryLanguage,
	}, true
}

func boolValue(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}

	return *value
}

func parseDate(value *string, field string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", *value)
	if err != nil {
		return nil, errors.New(field + " must use YYYY-MM-DD format")
	}

	return &parsed, nil
}
