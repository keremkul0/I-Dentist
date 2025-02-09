package login

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models/auth"
	"encoding/json"
	"net/http"
	"time"
)

type LoginService interface {
	Login(ctx context.Context, email string, password string) (auth.Login, error)
}

type LoginHandler struct {
	loginService LoginService
}

func NewLoginController(service LoginService) *LoginHandler {
	return &LoginHandler{loginService: service}
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var creds auth.Login
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	user, err := h.loginService.Login(ctx, creds.Email, creds.Password)
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
	w.WriteHeader(http.StatusOK)
}
