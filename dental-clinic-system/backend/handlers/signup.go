package handlers

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignupHandler struct {
    DB *gorm.DB
}

func (h *SignupHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var signupData struct {
        User  models.User  `json:"user"`
        Group models.Group `json:"group"`
        Clinic models.Clinic `json:"clinic"`
    }

    if err := json.NewDecoder(r.Body).Decode(&signupData); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Grup Oluştur
    signupData.Group.CreatedAt = time.Now()
    signupData.Group.UpdatedAt = time.Now()
    if err := h.DB.Create(&signupData.Group).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Klinik Oluştur
    signupData.Clinic.GroupID = signupData.Group.ID
    if err := h.DB.Create(&signupData.Clinic).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Kullanıcı Oluştur
    signupData.User.GroupID = signupData.Group.ID
    signupData.User.ClinicID= signupData.Clinic.ID
    signupData.User.Password = hashPassword(signupData.User.Password) // Şifreyi hashle
    signupData.User.CreatedAt = time.Now()
    signupData.User.UpdatedAt = time.Now()
    signupData.User.LastLogin = time.Now()
    signupData.User.IsActive = true
    if err := h.DB.Create(&signupData.User).Error; err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Başarılı yanıt
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "User, Group, and Clinic created successfully!",
        "user_id": signupData.User.ID,
    })
}

func hashPassword(password string) string {
    // Şifreleme işlemi (bcrypt kullanabiliriz)
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hashedPassword)
}
