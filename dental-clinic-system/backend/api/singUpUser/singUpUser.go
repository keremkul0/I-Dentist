package singUpUser

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
)

type SignUpUserService interface {
	SignUpUserService(user models.User) (string, error)
}

type SignUpUserHandler struct {
	singUpUserService SignUpUserService
}

func NewSignUpUserHandler(singUpUserService SignUpUserService) *SignUpUserHandler {
	return &SignUpUserHandler{
		singUpUserService: singUpUserService,
	}
}

func (s *SignUpUserHandler) SignUpUser(w http.ResponseWriter, r *http.Request) {

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cacheKey, err := s.singUpUserService.SignUpUserService(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Write response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(cacheKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
