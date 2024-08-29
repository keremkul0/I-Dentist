package signup

import (
	"dental-clinic-system/models"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SignupHandler struct {
	DB *gorm.DB
}

func (h *SignupHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var signupData struct {
		User   models.User   `json:"userRepository"`
		Group  models.Group  `json:"groupRepository"`
		Clinic models.Clinic `json:"clinicRepository"`
	}

	if err := json.NewDecoder(r.Body).Decode(&signupData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create or retrieve Group
	signupData.Group.CreatedAt = time.Now()
	signupData.Group.UpdatedAt = time.Now()
	if err := h.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&signupData.Group).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			http.Error(w, "Group name already exists", http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if signupData.Group.ID == 0 {
		var existingGroup models.Group
		if err := h.DB.Where("name = ?", signupData.Group.Name).First(&existingGroup).Error; err != nil {
			http.Error(w, "Failed to retrieve existing groupRepository", http.StatusInternalServerError)
			return
		}
		signupData.Group.ID = existingGroup.ID
	}

	// Create Clinic
	signupData.Clinic.GroupID = signupData.Group.ID
	if err := h.DB.Create(&signupData.Clinic).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create User
	signupData.User.GroupID = signupData.Group.ID
	signupData.User.ClinicID = signupData.Clinic.ID
	signupData.User.Password = hashPassword(signupData.User.Password) // Hash the password
	signupData.User.CreatedAt = time.Now()
	signupData.User.UpdatedAt = time.Now()
	signupData.User.LastLogin = time.Now()
	signupData.User.IsActive = true

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// Retrieve and assign existing roles to the userRepository
		var roles []*models.Role
		roleNames := []string{}
		for _, role := range signupData.User.Roles {
			roleNames = append(roleNames, role.Name)
		}
		if err := tx.Where("name IN ?", roleNames).Find(&roles).Error; err != nil {
			return err
		}
		signupData.User.Roles = roles

		// Create User and associate roles
		if err := tx.Create(&signupData.User).Error; err != nil {
			return err
		}
		if err := tx.Model(&signupData.User).Association("Roles").Replace(roles); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		http.Error(w, "Transaction failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful response
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User, Group, and Clinic created successfully!",
		"user_id": signupData.User.ID,
	})
	if err != nil {
		return
	}
}

func hashPassword(password string) string {
	// Hashing the password (using bcrypt)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
