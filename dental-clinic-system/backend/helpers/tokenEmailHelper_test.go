package helpers

import (
	"dental-clinic-system/models"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenEmailHelper(t *testing.T) {
	// Setup a valid JWT token
	claims := &models.Claims{
		Email: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTKey())
	if err != nil {
		t.Fatalf("Failed to sign token: %v", err)
	}

	tests := []struct {
		name      string
		setup     func() *http.Request
		wantEmail string
		wantErr   bool
	}{
		{
			name: "Valid token",
			setup: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})
				return req
			},
			wantEmail: "test@example.com",
			wantErr:   false,
		},
		{
			name: "No token cookie",
			setup: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/", nil)
			},
			wantEmail: "",
			wantErr:   true,
		},
		{
			name: "Invalid token",
			setup: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.AddCookie(&http.Cookie{Name: "token", Value: "invalidtoken"})
				return req
			},
			wantEmail: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setup()
			claims, err := TokenEmailHelper(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenEmailHelper() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && claims.Email != tt.wantEmail {
				t.Errorf("TokenEmailHelper() email = %v, want %v", claims.Email, tt.wantEmail)
			}
		})
	}
}
