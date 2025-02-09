package role

import (
	"context"
	"dental-clinic-system/models/user"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]user.Role, error)
	GetRole(ctx context.Context, id uint) (user.Role, error)
	CreateRole(ctx context.Context, role user.Role) (user.Role, error)
	UpdateRole(ctx context.Context, role user.Role) (user.Role, error)
	DeleteRole(ctx context.Context, id uint) error
}

type RoleHandler struct {
	roleService RoleService
}

func NewRoleController(roleService RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) GetRoles(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roles, err := h.roleService.GetRoles(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(roles)
	if err != nil {
		return
	}
}

func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}
	role, err := h.roleService.GetRole(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var role user.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	role, err = h.roleService.CreateRole(ctx, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	var role user.Role
	err := json.NewDecoder(r.Body).Decode(&role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	role, err = h.roleService.UpdateRole(ctx, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(role)
	if err != nil {
		return
	}
}

func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancelFunc := context.WithTimeout(ctx, 2*time.Second)
	defer cancelFunc()
	params := mux.Vars(r)
	idStr := params["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}
	err = h.roleService.DeleteRole(ctx, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
}
