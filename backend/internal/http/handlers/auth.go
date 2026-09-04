package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/auth"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/domain"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/mail"
	"github.com/oluwafemiomotoso/heritage-beaver/backend/internal/store/postgres"
)

const verificationTTL = 24 * time.Hour

type AuthHandler struct {
	users   postgres.UserRepository
	refresh postgres.RefreshTokenRepository
	verify  postgres.VerificationTokenRepository
	mailer  *mail.Mailer
	baseURL string
	secret  string
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
	RefreshToken string `json:"refresh_token"`
}

type verifyEmailRequest struct {
	Token string `json:"token"`
}

type resendVerificationRequest struct {
	Email string `json:"email"`
}

type authResponse struct {
	User         domain.User `json:"user"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}

type registerResponse struct {
	User    domain.User `json:"user"`
	Message string      `json:"message"`
}

func NewAuthHandler(repo postgres.UserRepository, refreshRepo postgres.RefreshTokenRepository, verifyRepo postgres.VerificationTokenRepository, mailer *mail.Mailer, baseURL, secret string) AuthHandler {
	return AuthHandler{users: repo, refresh: refreshRepo, verify: verifyRepo, mailer: mailer, baseURL: strings.TrimSuffix(baseURL, "/"), secret: secret}
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

	if err := checkPassword(req.Password); err != nil {
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

	user, err := h.users.Create(r.Context(), postgres.CreateUserParams{
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

		log.Printf("register user %s: %v", req.Email, err)
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	if err := h.sendVerificationEmail(r, user); err != nil {
		log.Printf("send verification email for %s: %v", user.ID, err)
	}

	writeJSON(w, http.StatusCreated, registerResponse{
		User:    user,
		Message: "Account created. Check your inbox for a confirmation link to begin.",
	})
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

	user, err := h.users.GetByEmail(r.Context(), req.Email)
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

	if !user.EmailVerified {
		writeError(w, http.StatusForbidden, "email not verified — check your inbox for the confirmation link")
		return
	}

	h.writeTokenPair(w, r, user)
}

func (h AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req verifyEmailRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Token, "token"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	stored, err := h.verify.FindByHash(r.Context(), auth.HashRefreshToken(req.Token))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusBadRequest, "invalid or expired confirmation link")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to confirm email")
		return
	}

	if stored.Consumed || time.Now().After(stored.ExpiresAt) {
		writeError(w, http.StatusBadRequest, "invalid or expired confirmation link")
		return
	}

	if err := h.verify.Consume(r.Context(), stored.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusBadRequest, "invalid or expired confirmation link")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to confirm email")
		return
	}

	if err := h.users.SetEmailVerified(r.Context(), stored.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to confirm email")
		return
	}

	user, err := h.users.GetByID(r.Context(), stored.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to confirm email")
		return
	}

	h.writeTokenPair(w, r, user)
}

func (h AuthHandler) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req resendVerificationRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.Email, "email"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Always 200: never reveal whether an address is registered.
	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err == nil && !user.EmailVerified {
		if err := h.sendVerificationEmail(r, user); err != nil {
			log.Printf("resend verification email for %s: %v", user.ID, err)
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "If the address is registered and unconfirmed, a new link is on its way.",
	})
}

func (h AuthHandler) sendVerificationEmail(r *http.Request, user domain.User) error {
	token, hash, err := auth.GenerateRefreshToken()
	if err != nil {
		return err
	}

	_ = h.verify.DeleteForUser(r.Context(), user.ID)

	if _, err := h.verify.Create(r.Context(), user.ID, hash, time.Now().Add(verificationTTL)); err != nil {
		return err
	}

	link := h.baseURL + "/verify-email?token=" + token
	return h.mailer.SendVerificationEmail(r.Context(), user.Email, link)
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.RefreshToken, "refresh_token"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	stored, err := h.refresh.FindByHash(r.Context(), auth.HashRefreshToken(req.RefreshToken))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to look up session")
		return
	}

	// Reuse of a rotated (revoked) token signals theft: kill the whole family.
	if stored.RevokedAt != nil {
		_ = h.refresh.RevokeAllForUser(r.Context(), stored.UserID)
		writeError(w, http.StatusUnauthorized, "session reused; signed out everywhere")
		return
	}

	if time.Now().After(stored.ExpiresAt) {
		_ = h.refresh.Revoke(r.Context(), stored.ID, nil)
		writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	user, err := h.users.GetByID(r.Context(), stored.UserID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	// Rotate: mint the replacement first, then revoke the presented token.
	token, hash, err := auth.GenerateRefreshToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to rotate session")
		return
	}

	replacement, err := h.refresh.Create(r.Context(), stored.UserID, hash, time.Now().Add(auth.RefreshTokenTTL))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to rotate session")
		return
	}

	if err := h.refresh.Revoke(r.Context(), stored.ID, &replacement.ID); err != nil {
		_ = h.refresh.Revoke(r.Context(), replacement.ID, nil)
		writeError(w, http.StatusInternalServerError, "failed to rotate session")
		return
	}

	access, err := auth.GenerateAccessToken(user.ID, h.secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: user, Token: access, RefreshToken: token})
}

func (h AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := required(req.RefreshToken, "refresh_token"); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	stored, err := h.refresh.FindByHash(r.Context(), auth.HashRefreshToken(req.RefreshToken))
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_ = h.refresh.Revoke(r.Context(), stored.ID, nil)
	w.WriteHeader(http.StatusNoContent)
}

func (h AuthHandler) writeTokenPair(w http.ResponseWriter, r *http.Request, user domain.User) {
	access, err := auth.GenerateAccessToken(user.ID, h.secret)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	token, hash, err := auth.GenerateRefreshToken()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	if _, err := h.refresh.Create(r.Context(), user.ID, hash, time.Now().Add(auth.RefreshTokenTTL)); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusOK, authResponse{User: user, Token: access, RefreshToken: token})
}

func checkPassword(password string) error {
	if utf8.RuneCountInString(password) < auth.MinPasswordChars {
		return errors.New("password must be at least 8 characters")
	}

	return nil
}
