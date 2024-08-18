package group

import (
	models2 "dental-clinic-system/repository/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type GroupHandlerService interface {
	GetGroups(w http.ResponseWriter, r *http.Request)
	GetGroup(w http.ResponseWriter, r *http.Request)
	CreateGroup(w http.ResponseWriter, r *http.Request)
	UpdateGroup(w http.ResponseWriter, r *http.Request)
	DeleteGroup(w http.ResponseWriter, r *http.Request)
	GetClinicsByGroup(w http.ResponseWriter, r *http.Request)
}

func NewGroupHandlerService(db *gorm.DB) *GroupHandler {
	return &GroupHandler{DB: db}
}

type GroupHandler struct {
	DB *gorm.DB
}

func (h *GroupHandler) GetGroups(w http.ResponseWriter, r *http.Request) {
	var groups []models2.Group
	if result := h.DB.Find(&groups); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(groups)
	if err != nil {
		return
	}
}

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var group models2.Group
	if result := h.DB.First(&group, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	err := json.NewEncoder(w).Encode(group)
	if err != nil {
		return
	}
}

func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group models2.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if result := h.DB.Create(&group); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(group)
	if err != nil {
		return
	}
}

func (h *GroupHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var group models2.Group
	if result := h.DB.First(&group, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.DB.Save(&group)
	err := json.NewEncoder(w).Encode(group)
	if err != nil {
		return
	}
}

func (h *GroupHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if result := h.DB.Delete(&models2.Group{}, params["id"]); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) GetClinicsByGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var clinics []models2.Clinic
	if result := h.DB.Where("group_id = ?", params["id"]).Find(&clinics); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(clinics)
	if err != nil {
		return
	}
}
