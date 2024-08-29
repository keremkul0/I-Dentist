package groupService

import (
	"dental-clinic-system/models"
	"dental-clinic-system/repository/groupRepository"
)

type GroupService interface {
	GetGroups() ([]models.Group, error)
	GetGroup(id uint) (models.Group, error)
	CreateGroup(group models.Group) (models.Group, error)
	UpdateGroup(group models.Group) (models.Group, error)
	DeleteGroup(id uint) error
	GetClinicsByGroup(id uint) ([]models.Clinic, error)
}

type groupService struct {
	groupRepository groupRepository.GroupRepository
}

func NewGroupService(groupRepository groupRepository.GroupRepository) *groupService {
	return &groupService{
		groupRepository: groupRepository,
	}
}

func (s *groupService) GetGroups() ([]models.Group, error) {
	return s.groupRepository.GetGroups()
}

func (s *groupService) GetGroup(id uint) (models.Group, error) {
	return s.groupRepository.GetGroup(id)
}

func (s *groupService) CreateGroup(group models.Group) (models.Group, error) {
	return s.groupRepository.CreateGroup(group)
}

func (s *groupService) UpdateGroup(group models.Group) (models.Group, error) {
	return s.groupRepository.UpdateGroup(group)
}

func (s *groupService) DeleteGroup(id uint) error {
	return s.groupRepository.DeleteGroup(id)
}

func (s *groupService) GetClinicsByGroup(id uint) ([]models.Clinic, error) {
	return s.groupRepository.GetClinicsByGroup(id)
}
