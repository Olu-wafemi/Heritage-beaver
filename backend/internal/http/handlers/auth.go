package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

type AuthHandler struct {
	repo postgres.UserRepository
}

type registerRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshRequest struct {
	Token string `json:"token"`
}

type authResponse struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

type refreshResponse struct {
	Token string `json:"token"`
}

func NewAuthHandler(repo postgres.UserRepository) AuthHandler {
	return AuthHandler{repo: repo}
}

func (h AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Email, "email"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Password, "password"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.DisplayName, "display_name"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to process password")
		return
	}

	user, err := h.repo.Create(r.Context(), postgres.CreateUserParams{
		Email:          req.Email,
		DisplayName:    req.DisplayName,
		PrimaryCulture: "",
		PasswordHash:   hashedPassword,
	})
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			writeError(w, http.StatusConflict, "user email already exists")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	token, err := auth.GenerateToken(user.ID, auth.GetSecret())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusCreated, authResponse{User: user, Token: token})
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Email, "email"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Password, "password"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.repo.GetByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to look up user")
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := auth.GenerateToken(user.ID, auth.GetSecret())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: user, Token: token})
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Token, "token"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	claims, err := auth.ValidateToken(req.Token, auth.GetSecret())
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	token, err := auth.GenerateToken(claims.UserID, auth.GetSecret())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, refreshResponse{Token: token})
}
