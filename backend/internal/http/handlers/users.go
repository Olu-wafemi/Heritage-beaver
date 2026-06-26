package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type UserHandler struct {
	repo postgres.UserRepository
}

type createUserRequest struct {
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	PrimaryCulture string `json:"primary_culture"`
}

type updateUserRequest struct {
	DisplayName    string `json:"display_name"`
	PrimaryCulture string `json:"primary_culture"`
}

func NewUserHandler(repo postgres.UserRepository) UserHandler {
	return UserHandler{repo: repo}
}

func (h UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Email, "email"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.DisplayName, "display_name"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.Create(r.Context(), postgres.CreateUserParams{
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		PrimaryCulture: req.PrimaryCulture,
	})
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			writeError(w, http.StatusConflict, "user email already exists")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func (h UserHandler) Get(w http.ResponseWriter, r *http.Request) {
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

var _ = domain.User{}
