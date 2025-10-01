package forgotPassword

import (
	"context"
	"encoding/json"
	"net/http"
)

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetService interface {
	RequestPasswordReset(ctx context.Context, email string) error
}

type ForgotPasswordHandler struct {
	passwordResetService PasswordResetService
}

func NewForgotPasswordController(passwordResetService PasswordResetService) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{
		passwordResetService: passwordResetService,
	}
}

func (h *ForgotPasswordHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	_ = h.passwordResetService.RequestPasswordReset(ctx, req.Email)

	// Güvenlik nedeniyle her durumda aynı response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
