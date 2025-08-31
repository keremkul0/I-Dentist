package user

import (
	"context"
	"dental-clinic-system/helpers"
	"dental-clinic-system/mapper"
	"dental-clinic-system/models/user"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserService interface {
	GetUsers(ctx context.Context, ClinicID uint) ([]user.UserGetModel, error)
	GetUser(ctx context.Context, id uint) (user.UserGetModel, error)
	GetUserByEmail(ctx context.Context, email string) (user.UserGetModel, error)
	CreateUser(ctx context.Context, user user.User) (user.UserGetModel, error)
	UpdateUser(ctx context.Context, user user.User) (user.UserGetModel, error)
	DeleteUser(ctx context.Context, id uint) error
	CheckUserExist(ctx context.Context, user user.UserGetModel) (bool, error)
}

type UserHandler struct {
	userService UserService
}

func NewUserController(service UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	users, err := h.userService.GetUsers(ctx, authenticatedUser.ClinicID)
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
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid authenticatedUser ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != requestedUser.ClinicID {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err = json.NewEncoder(w).Encode(authenticatedUser)
	if err != nil {
		return
	}
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(authenticatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var newUser user.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != newUser.ClinicID || (!helpers.ContainsRole(authenticatedUser, "Clinic Admin") && (!helpers.ContainsRole(authenticatedUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tempUserGetModel := mapper.MapUserToUserGetModel(newUser)
	exists, err := h.userService.CheckUserExist(ctx, tempUserGetModel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	newUser.Password = helpers.HashPassword(newUser.Password) // Hash the password if not already hashed

	createdUser, err := h.userService.CreateUser(ctx, newUser)
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
	var updateUser user.User
	err := json.NewDecoder(r.Body).Decode(&updateUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.CookieTokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != updateUser.ClinicID || (!helpers.ContainsRole(authenticatedUser, "Clinic Admin") && (!helpers.ContainsRole(authenticatedUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createdUser, err := h.userService.UpdateUser(ctx, updateUser)
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
	authenticatedUser, err := h.userService.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if authenticatedUser.ClinicID != requestedUser.ClinicID || (!helpers.ContainsRole(authenticatedUser, "Clinic Admin") && (!helpers.ContainsRole(authenticatedUser, "Superadmin"))) {
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
