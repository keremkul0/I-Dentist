package singUpUser

import (
	"context"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
)

type SignUpUserService interface {
	SignUpUser(ctx context.Context, user user.User) (string, error)
}

type SignUpUserHandler struct {
	signUpUserService SignUpUserService
}

func NewSignUpUserHandler(signUpUserService SignUpUserService) *SignUpUserHandler {
	return &SignUpUserHandler{
		signUpUserService: signUpUserService,
	}
}

func (s *SignUpUserHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user user.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cacheKey, err := s.signUpUserService.SignUpUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cacheKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
