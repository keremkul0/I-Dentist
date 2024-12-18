package roleRepository

import (
	"context"
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

func NewRoleRepository(db *gorm.DB) *roleRepository {
	return &roleRepository{DB: db}
}

type roleRepository struct {
	DB *gorm.DB
}

func (r *roleRepository) GetRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	result := r.DB.WithContext(ctx).Find(&roles)
	return roles, result.Error
}

func (r *roleRepository) GetRole(ctx context.Context, id uint) (models.Role, error) {
	var role models.Role
	result := r.DB.WithContext(ctx).First(&role, id)
	return role, result.Error
}

func (r *roleRepository) CreateRole(ctx context.Context, role models.Role) (models.Role, error) {
	result := r.DB.WithContext(ctx).Create(&role)
	return role, result.Error
}

func (r *roleRepository) UpdateRole(ctx context.Context, role models.Role) (models.Role, error) {
	result := r.DB.WithContext(ctx).Save(&role)
	return role, result.Error
}

func (r *roleRepository) DeleteRole(ctx context.Context, id uint) error {
	result := r.DB.WithContext(ctx).Delete(&models.Role{}, id)
	return result.Error
}
