package roleRepository

import (
	"context"
	"dental-clinic-system/models/user"
	"errors"

	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

// Repository handles role-related database operations
type Repository struct {
	DB *gorm.DB
}

// NewRepository creates a new instance of Repository
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

// GetRoles retrieves all roles from the database
func (repo *Repository) GetRoles(ctx context.Context) ([]user.Role, error) {
	var rolesList []user.Role
	result := repo.DB.WithContext(ctx).Find(&rolesList)
	if result.Error != nil {
		log.Error().
			Str("operation", "GetRoles").
			Err(result.Error).
			Msg("Failed to retrieve roles")
		return nil, result.Error
	}
	log.Info().
		Str("operation", "GetRoles").
		Int("count", len(rolesList)).
		Msg("Retrieved roles successfully")
	return rolesList, nil
}

// GetRole retrieves a single role by its ID
func (repo *Repository) GetRole(ctx context.Context, id uint) (user.Role, error) {
	var rl user.Role
	result := repo.DB.WithContext(ctx).First(&rl, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Warn().
				Str("operation", "GetRole").
				Err(result.Error).
				Uint("role_id", id).
				Msg("Role not found")
		} else {
			log.Error().
				Str("operation", "GetRole").
				Err(result.Error).
				Uint("role_id", id).
				Msg("Failed to retrieve role")
		}
		return user.Role{}, result.Error
	}
	log.Info().
		Str("operation", "GetRole").
		Uint("role_id", id).
		Msg("Retrieved role successfully")
	return rl, nil
}

// CreateRole creates a new role record in the database
func (repo *Repository) CreateRole(ctx context.Context, newRole user.Role) (user.Role, error) {
	result := repo.DB.WithContext(ctx).Create(&newRole)
	if result.Error != nil {
		log.Error().
			Str("operation", "CreateRole").
			Err(result.Error).
			Msg("Failed to create role")
		return user.Role{}, result.Error
	}
	log.Info().
		Str("operation", "CreateRole").
		Uint("role_id", newRole.ID).
		Msg("Role created successfully")
	return newRole, nil
}

// UpdateRole updates an existing role record in the database
func (repo *Repository) UpdateRole(ctx context.Context, updatedRole user.Role) (user.Role, error) {
	result := repo.DB.WithContext(ctx).Save(&updatedRole)
	if result.Error != nil {
		log.Error().
			Str("operation", "UpdateRole").
			Err(result.Error).
			Uint("role_id", updatedRole.ID).
			Msg("Failed to update role")
		return user.Role{}, result.Error
	}
	log.Info().
		Str("operation", "UpdateRole").
		Uint("role_id", updatedRole.ID).
		Msg("Role updated successfully")
	return updatedRole, nil
}

// DeleteRole deletes a role record from the database by its ID
func (repo *Repository) DeleteRole(ctx context.Context, id uint) error {
	result := repo.DB.WithContext(ctx).Delete(&user.Role{}, id)
	if result.Error != nil {
		log.Error().
			Str("operation", "DeleteRole").
			Err(result.Error).
			Uint("role_id", id).
			Msg("Failed to delete role")
		return result.Error
	}
	log.Info().
		Str("operation", "DeleteRole").
		Uint("role_id", id).
		Msg("Role deleted successfully")
	return nil
}
