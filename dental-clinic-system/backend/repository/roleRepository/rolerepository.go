package roleRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{DB: db}
}

type roleRepository struct {
	DB *gorm.DB
}

func (r *roleRepository) GetRoles() ([]models.Role, error) {
	var roles []models.Role
	if result := r.DB.Find(&roles); result.Error != nil {
		return nil, result.Error
	}
	return roles, nil
}

func (r *roleRepository) GetRole(id uint) (models.Role, error) {
	var role models.Role
	if result := r.DB.First(&role, id); result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (r *roleRepository) CreateRole(role models.Role) (models.Role, error) {
	if result := r.DB.Create(&role); result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (r *roleRepository) UpdateRole(role models.Role) (models.Role, error) {
	if result := r.DB.Save(&role); result.Error != nil {
		return models.Role{}, result.Error
	}
	return role, nil
}

func (r *roleRepository) DeleteRole(id uint) error {
	if result := r.DB.Delete(&models.Role{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
