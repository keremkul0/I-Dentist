package verifyEmail

import (
	"dental-clinic-system/helpers"
	"net/http"
)

type EmailService interface {
	VerifyUserEmail(token string, email string) bool
}

type VerifyEmailHandler struct {
	EmailService EmailService
}

func NewVerifyEmailController(service EmailService) *VerifyEmailHandler {
	return &VerifyEmailHandler{EmailService: service}
}

func (h *VerifyEmailHandler) VerifyUserEmail(w http.ResponseWriter, r *http.Request) {
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

	if h.EmailService.VerifyUserEmail(token, claims.Email) {
		http.Redirect(w, r, "/verify-email-success", http.StatusSeeOther)
		return
	}

	http.Error(w, "Invalid token", http.StatusBadRequest)
}
