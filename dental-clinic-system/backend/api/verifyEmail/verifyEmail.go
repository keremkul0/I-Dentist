package verifyEmail

import (
	"context"
	"dental-clinic-system/helpers"
	"net/http"
	"time"
)

type EmailService interface {
	VerifyUserEmail(ctx context.Context, token string, email string) bool
}

type VerifyEmailHandler struct {
	EmailService EmailService
}

func NewVerifyEmailController(service EmailService) *VerifyEmailHandler {
	return &VerifyEmailHandler{EmailService: service}
}

func (h *VerifyEmailHandler) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	
	if h.EmailService.VerifyUserEmail(ctx, token, claims.Email) {
		_, err := w.Write([]byte("Email verified"))
		if err != nil {
			return
		}
		return
	}

	http.Error(w, "Invalid token", http.StatusBadRequest)
}
