package groupRepository

import (
	"dental-clinic-system/models"
	"gorm.io/gorm"
)

type GroupRepository interface {
	GetGroups() ([]models.Group, error)
	GetGroup(id uint) (models.Group, error)
	CreateGroup(group models.Group) (models.Group, error)
	UpdateGroup(group models.Group) (models.Group, error)
	DeleteGroup(id uint) error
	GetClinicsByGroup(id uint) ([]models.Clinic, error)
}

func NewGroupRepository(db *gorm.DB) *groupRepository {
	return &groupRepository{DB: db}
}

type groupRepository struct {
	DB *gorm.DB
}

func (r *groupRepository) GetGroups() ([]models.Group, error) {
	var groups []models.Group
	err := r.DB.Find(&groups).Error
	return groups, err
}

func (r *groupRepository) GetGroup(id uint) (models.Group, error) {
	var group models.Group
	err := r.DB.First(&group, id).Error
	return group, err
}

func (r *groupRepository) CreateGroup(group models.Group) (models.Group, error) {
	err := r.DB.Create(&group).Error
	return group, err
}

func (r *groupRepository) UpdateGroup(group models.Group) (models.Group, error) {
	err := r.DB.Save(&group).Error
	return group, err
}

func (r *groupRepository) DeleteGroup(id uint) error {
	err := r.DB.Delete(&models.Group{}, id).Error
	return err
}

func (r *groupRepository) GetClinicsByGroup(id uint) ([]models.Clinic, error) {
	var clinics []models.Clinic
	err := r.DB.Where("group_id = ?", id).Find(&clinics).Error
	return clinics, err
}
