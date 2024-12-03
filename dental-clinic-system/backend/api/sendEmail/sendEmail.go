package sendEmail

import (
	"dental-clinic-system/helpers"
	"net/http"
	"time"
)

type EmailService interface {
	SendVerificationEmail(email string, token string) error
}

type SendEmailHandler struct {
	EmailService EmailService
}

func NewSendEmailController(service EmailService) *SendEmailHandler {
	return &SendEmailHandler{EmailService: service}
}

func (h *SendEmailHandler) SendVerificationEmail(w http.ResponseWriter, r *http.Request) {
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

	err = h.EmailService.SendVerificationEmail(claims.Email, token)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

}
