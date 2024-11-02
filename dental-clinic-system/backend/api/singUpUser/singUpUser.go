package singUpUser

import (
	"dental-clinic-system/application/singUpUserService"
	"dental-clinic-system/models"
	"dental-clinic-system/redisService"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type SignUpUserHandler struct {
	singUpUserService singUpUserService.SignUpUserService
}

func NewSignUpUserHandler(singUpUserService singUpUserService.SignUpUserService) *SignUpUserHandler {
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

	user, err = s.singUpUserService.SignUpUserService(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Kullanıcıyı Redis'e kaydet
	ctx := r.Context()
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	cacheKey := uuid.New().String()

	err = redisService.Rdb.Set(ctx, cacheKey, userJSON, 10*time.Minute).Err()
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
