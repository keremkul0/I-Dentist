package login

import (
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
	"time"
)

type LoginService interface {
	Login(email string, password string) (models.Login, error)
}

type LoginHandler struct {
	loginService LoginService
}

func NewLoginController(service LoginService) *LoginHandler {
	return &LoginHandler{loginService: service}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Login
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := h.loginService.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Hour * 24)
	tokenString, err := helpers.GenerateJWTToken(user.Email, expirationTime)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}
