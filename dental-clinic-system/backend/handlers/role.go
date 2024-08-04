package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type RoleHandler struct {
	DB *gorm.DB
}

func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Role
	h.DB.Find(&roles)
	json.NewEncoder(w).Encode(roles)
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	json.NewDecoder(r.Body).Decode(&role)
	h.DB.Create(&role)
	json.NewEncoder(w).Encode(role)
}
