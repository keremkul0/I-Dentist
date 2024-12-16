package logout

import (
	"context"
	"net/http"
	"time"
)

type TokenService interface {
	AddTokenToBlacklist(ctx context.Context, token string, expireTime time.Time) error
}

type LogoutController struct {
	tokenService TokenService
}

func NewLogoutController(tokenService TokenService) *LogoutController {
	return &LogoutController{
		tokenService: tokenService,
	}
}

func (h *LogoutController) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	token, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "No token found", http.StatusUnauthorized)
		return
	}

	err = h.tokenService.AddTokenToBlacklist(ctx, token.Value, time.Now().Add(24*time.Hour))
	if err != nil {
		http.Error(w, "logout failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
	w.WriteHeader(http.StatusOK)
}
