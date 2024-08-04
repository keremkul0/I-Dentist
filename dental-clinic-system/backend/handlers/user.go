package handlers

import (
    "encoding/json"
    "net/http"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "dental-clinic-system/models"
)

type UserHandler struct {
    DB *gorm.DB
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var user models.User
    json.NewDecoder(r.Body).Decode(&user)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Failed to hash password", http.StatusInternalServerError)
        return
    }
    user.Password = string(hashedPassword)
    h.DB.Create(&user)
    json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var requestUser models.User
    json.NewDecoder(r.Body).Decode(&requestUser)
    var user models.User
    h.DB.Where("email = ?", requestUser.Email).First(&user)
    if user.ID == 0 {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }
    err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestUser.Password))
    if err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }
    // Token generation (JWT or similar) should be done here
    json.NewEncoder(w).Encode(user)
}
