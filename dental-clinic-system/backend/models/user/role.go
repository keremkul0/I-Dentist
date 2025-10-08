package user

import (
	"gorm.io/gorm"
)

type RoleName string

const (
	RoleDoctor                  RoleName = "doctor"
	RoleAssistant               RoleName = "assistant"
	RoleIntern                  RoleName = "intern"
	RoleSecretary               RoleName = "secretary"
	RoleSecurity                RoleName = "security"
	RoleManager                 RoleName = "manager"
	RoleCleaner                 RoleName = "cleaner"
	RoleRadiologyTechnician     RoleName = "radiology_technician"
	RoleAccountant              RoleName = "accountant"
	RolePatientConsultant       RoleName = "patient_consultant"
	RoleItSupportSpecialist     RoleName = "it_support_specialist"
	RoleSupplyChainManager      RoleName = "supply_chain_manager"
	RoleSterilizationTechnician RoleName = "sterilization_technician"
	RoleHrManager               RoleName = "hr_manager"
	RoleOther                   RoleName = "other"
	RoleOrthodontist            RoleName = "orthodontist"
	RoleClinicAdmin             RoleName = "clinic_admin"
	RoleSuperAdmin              RoleName = "super_admin"
)

type Role struct {
	gorm.Model
	Name  RoleName `json:"name"`
	Users []*User  `gorm:"many2many:user_roles;"`
}
