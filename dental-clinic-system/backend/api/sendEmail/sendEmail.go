package sendEmail

import (
	"context"
	"dental-clinic-system/models/claims"
	"net/http"
	"time"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email string, token string) error
}

type JwtService interface {
	GenerateJWTToken(email string, time time.Time) (string, error)
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
}

type SendEmailHandler struct {
	EmailService EmailService
	jwtService   JwtService
}

func NewSendEmailController(service EmailService, jwtService JwtService) *SendEmailHandler {
	return &SendEmailHandler{EmailService: service, jwtService: jwtService}
}

func (h *SendEmailHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := h.jwtService.ParseTokenFromCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := h.jwtService.GenerateJWTToken(claims.Email, time.Now().Add(time.Minute*5))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.EmailService.SendVerificationEmail(ctx, claims.Email, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
