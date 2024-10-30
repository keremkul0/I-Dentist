package logout

import (
	"dental-clinic-system/application/tokenService"
	"net/http"
	"time"
)

type LogoutHandler struct {
	tokenService tokenService.TokenService
}

func NewLogoutHandler(tokenService tokenService.TokenService) *LogoutHandler {
	return &LogoutHandler{
		tokenService: tokenService,
	}
}

func (h *LogoutHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "No token found", http.StatusUnauthorized)
		return
	}

	h.tokenService.AddTokenToBlacklistService(token.Value, time.Now().Add(24*time.Hour))

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
