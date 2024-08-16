package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"net/http"
)

type RoleHandler struct {
	DB *gorm.DB
}

func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Role
	if result := h.DB.Find(&roles); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(roles)
	if err != nil {
		return
	}
}

func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var role models.Role
	if result := h.DB.First(&role, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&role); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var role models.Role
	if result := h.DB.First(&role, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&role)
	err := json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models.Role{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
