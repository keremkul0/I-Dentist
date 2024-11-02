package user

import (
	"dental-clinic-system/application/userService"
	"dental-clinic-system/helpers"
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService userService.UserService
}

func NewUserController(service userService.UserService) *UserHandler {
	return &UserHandler{userService: service}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	users, err := h.userService.GetUsers(user.ClinicID)
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
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(uint(id))
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

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if mainUser.ClinicID != user.ClinicID || (!helpers.ContainsRole(mainUser, "Clinic Admin") && (!helpers.ContainsRole(mainUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tempUserGetModel := helpers.UserConvertor(user)
	if h.userService.CheckUserExist(tempUserGetModel) {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	user.Password = helpers.HashPassword(user.Password) // Hash the password if not already hashed

	createdUser, err := h.userService.CreateUser(user)

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
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if mainUser.ClinicID != user.ClinicID || (!helpers.ContainsRole(mainUser, "Clinic Admin") && (!helpers.ContainsRole(mainUser, "Superadmin"))) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	createdUser, err := h.userService.UpdateUser(user)
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
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	claims, err := helpers.TokenEmailHelper(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	mainUser, err := h.userService.GetUserByEmail(claims.Email)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	requestedUser, err := h.userService.GetUser(uint(id))
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

	err = h.userService.DeleteUser(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
