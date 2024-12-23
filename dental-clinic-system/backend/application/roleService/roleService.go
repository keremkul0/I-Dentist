package roleService

import (
	"context"
	"dental-clinic-system/models"
)

type RoleRepository interface {
	GetRoles(ctx context.Context) ([]models.Role, error)
	GetRole(ctx context.Context, id uint) (models.Role, error)
	CreateRole(ctx context.Context, role models.Role) (models.Role, error)
	UpdateRole(ctx context.Context, role models.Role) (models.Role, error)
	DeleteRole(ctx context.Context, id uint) error
}

type roleService struct {
	roleRepository RoleRepository
}

func NewRoleService(roleRepository RoleRepository) *roleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}

func (s *roleService) GetRoles(ctx context.Context) ([]models.Role, error) {
	return s.roleRepository.GetRoles(ctx)
}

func (s *roleService) GetRole(ctx context.Context, id uint) (models.Role, error) {
	return s.roleRepository.GetRole(ctx, id)
}

func (s *roleService) CreateRole(ctx context.Context, role models.Role) (models.Role, error) {
	return s.roleRepository.CreateRole(ctx, role)
}

func (s *roleService) UpdateRole(ctx context.Context, role models.Role) (models.Role, error) {
	return s.roleRepository.UpdateRole(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepository.DeleteRole(ctx, id)
}
