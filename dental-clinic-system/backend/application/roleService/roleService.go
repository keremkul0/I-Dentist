package roleService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/roleRepository"
)

type RoleService interface {
	GetRoles() ([]models.Role, error)
	GetRole(id uint) (models.Role, error)
	CreateRole(role models.Role) (models.Role, error)
	UpdateRole(role models.Role) (models.Role, error)
	DeleteRole(id uint) error
}

type roleService struct {
	roleRepository roleRepository.RoleRepository
}

func NewRoleService(roleRepository roleRepository.RoleRepository) *roleService {
	return &roleService{
		roleRepository: roleRepository,
	}
}

func (s *roleService) GetRoles() ([]models.Role, error) {
	return s.roleRepository.GetRoles()
}

func (s *roleService) GetRole(id uint) (models.Role, error) {
	return s.roleRepository.GetRole(id)
}

func (s *roleService) CreateRole(role models.Role) (models.Role, error) {
	return s.roleRepository.CreateRole(role)
}

func (s *roleService) UpdateRole(role models.Role) (models.Role, error) {
	return s.roleRepository.UpdateRole(role)
}

func (s *roleService) DeleteRole(id uint) error {
	return s.roleRepository.DeleteRole(id)
}
