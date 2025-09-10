package verifyEmail

import (
	"context"
	"dental-clinic-system/models/claims"
	"net/http"
)

type EmailService interface {
	VerifyUserEmail(ctx context.Context, token string, email string) bool
}

type JwtService interface {
	ParseToken(tokenStr string) (*claims.Claims, error)
}

type VerifyEmailHandler struct {
	EmailService EmailService
	jwtService   JwtService
}

func NewVerifyEmailController(service EmailService, jwtService JwtService) *VerifyEmailHandler {
	return &VerifyEmailHandler{EmailService: service, jwtService: jwtService}
}

func (h *VerifyEmailHandler) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	claims, err := h.jwtService.ParseToken(token)
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
