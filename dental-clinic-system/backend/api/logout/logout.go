package logout

import (
	"dental-clinic-system/application/tokenService"
	"net/http"
	"time"
)

type LogoutController struct {
	tokenService tokenService.TokenService
}

func NewLogoutController(tokenService tokenService.TokenService) *LogoutController {
	return &LogoutController{
		tokenService: tokenService,
	}
}

func (h *LogoutController) Logout(w http.ResponseWriter, r *http.Request) {
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
