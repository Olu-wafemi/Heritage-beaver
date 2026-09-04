package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type UserHandler struct {
	repo    postgres.UserRepository
	refresh postgres.RefreshTokenRepository
}

type updateUserRequest struct {
	DisplayName    string `json:"display_name"`
	PrimaryCulture string `json:"primary_culture"`
}

func NewUserHandler(repo postgres.UserRepository, refreshRepo postgres.RefreshTokenRepository) UserHandler {
	return UserHandler{repo: repo, refresh: refreshRepo}
}

// requireSelf ensures the caller only ever touches their own account.
func requireSelf(w http.ResponseWriter, r *http.Request) string {
	userID := auth.UserIDFromContext(r.Context())
	if userID == "" {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return ""
	}

	if r.PathValue("id") != userID {
		writeError(w, http.StatusForbidden, "forbidden")
		return ""
	}

	return userID
}

func (h UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	if requireSelf(w, r) == "" {
		return
	}

	user, err := h.repo.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if requireSelf(w, r) == "" {
		return
	}

	var req updateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.DisplayName, "display_name"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.Update(r.Context(), postgres.UpdateUserParams{
		ID:             r.PathValue("id"),
		DisplayName:    req.DisplayName,
		PrimaryCulture: req.PrimaryCulture,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := requireSelf(w, r)
	if userID == "" {
		return
	}

	// Deleting your account signs you out everywhere first.
	_ = h.refresh.RevokeAllForUser(r.Context(), userID)

	err := h.repo.Delete(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
