package user

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type UserService interface {
	GetUsers(ctx context.Context, ClinicID uint) ([]models.UserGetModel, error)
	GetUser(ctx context.Context, id uint) (models.UserGetModel, error)
	GetUserByEmail(ctx context.Context, email string) (models.UserGetModel, error)
	CreateUser(ctx context.Context, user models.User) (models.UserGetModel, error)
	UpdateUser(ctx context.Context, user models.User) (models.UserGetModel, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, user models.UserGetModel) bool
}

type UserHandler struct {
	userService UserService
}

func NewUserController(service UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	users, err := h.userService.GetUsers(ctx, user.ClinicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if user.ClinicID != requestedUser.ClinicID {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if mainUser.ClinicID != user.ClinicID || (!helpers.ContainsRole(mainUser, "Clinic Admin") && (!helpers.ContainsRole(mainUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tempUserGetModel := mapper.UserMapper(user)
	if h.userService.CheckUserExist(ctx, tempUserGetModel) {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	user.Password = helpers.HashPassword(user.Password) // Hash the password if not already hashed

	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		return
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if mainUser.ClinicID != user.ClinicID || (!helpers.ContainsRole(mainUser, "Clinic Admin") && (!helpers.ContainsRole(mainUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createdUser, err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(createdUser)
	if err != nil {
		return
	}
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if mainUser.ClinicID != requestedUser.ClinicID || (!helpers.ContainsRole(mainUser, "Clinic Admin") && (!helpers.ContainsRole(mainUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if helpers.ContainsRole(requestedUser, "Superadmin") {
		http.Error(w, "Cannot delete Superadmin", http.StatusBadRequest)
		return
	}

	err = h.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
