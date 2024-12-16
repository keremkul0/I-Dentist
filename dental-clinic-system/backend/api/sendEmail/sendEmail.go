package sendEmail

import (
	"context"
	"dental-clinic-system/helpers"
	"net/http"
	"time"
)

type EmailService interface {
	SendVerificationEmail(ctx context.Context, email string, token string) error
}

type SendEmailHandler struct {
	EmailService EmailService
}

func NewSendEmailController(service EmailService) *SendEmailHandler {
	return &SendEmailHandler{EmailService: service}
}

func (h *SendEmailHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	token, err := helpers.GenerateJWTToken(claims.Email, time.Now().Add(time.Minute*5))
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
