package resetPassword

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Email       string `json:"email" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordResetService interface {
	ResetPassword(ctx context.Context, tokenStr string, email string, newHashedPassword string) error
}

type ResetPasswordHandler struct {
	passwordResetService PasswordResetService
}

func NewResetPasswordController(service PasswordResetService) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		passwordResetService: service,
	}
}

func (h *ResetPasswordHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest

	// JSON decode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid JSON format")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validation
	if req.Token == "" || req.Email == "" || req.NewPassword == "" {
		http.Error(w, "Token, email and new_password are required", http.StatusBadRequest)
		return
	}

	// Şifre hash'le
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("Failed to hash password")
		http.Error(w, "Password processing failed", http.StatusInternalServerError)
		return
	}

	// Service çağrısı
	err = h.passwordResetService.ResetPassword(r.Context(), req.Token, req.Email, string(hashedPassword))
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Password reset failed")
		http.Error(w, "Password reset failed", http.StatusBadRequest)
		return
	}

	log.Info().Str("email", req.Email).Msg("Password reset successful")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Password reset successful"}`))
}
