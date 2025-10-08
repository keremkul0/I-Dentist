package roleService

import (
	"context"
	"dental-clinic-system/models/user"
	"strings"
)

type RoleRepository interface {
	GetRoles(ctx context.Context) ([]user.Role, error)
	GetRole(ctx context.Context, id uint) (user.Role, error)
	CreateRole(ctx context.Context, role user.Role) (user.Role, error)
	UpdateRole(ctx context.Context, role user.Role) (user.Role, error)
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

func (s *roleService) GetRoles(ctx context.Context) ([]user.Role, error) {
	return s.roleRepository.GetRoles(ctx)
}

func (s *roleService) GetRole(ctx context.Context, id uint) (user.Role, error) {
	return s.roleRepository.GetRole(ctx, id)
}

func (s *roleService) CreateRole(ctx context.Context, role user.Role) (user.Role, error) {
	return s.roleRepository.CreateRole(ctx, role)
}

func (s *roleService) UpdateRole(ctx context.Context, role user.Role) (user.Role, error) {
	return s.roleRepository.UpdateRole(ctx, role)
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	return s.roleRepository.DeleteRole(ctx, id)
}

func (s *roleService) UserHasRole(user user.UserGetModel, roleName string) bool {
	for _, role := range user.Roles {
		if strings.EqualFold(string(role.Name), roleName) {
			return true
		}
	}

	return false
}
