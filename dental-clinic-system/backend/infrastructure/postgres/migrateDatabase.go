package postgres

import (
	"dental-clinic-system/models/appointment"
	"dental-clinic-system/models/clinic"
	"dental-clinic-system/models/patient"
	"dental-clinic-system/models/procedure"
	"dental-clinic-system/models/token"
	"dental-clinic-system/models/user"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&appointment.Appointment{},
		&clinic.Clinic{},
		&patient.Patient{},
		&procedure.Procedure{},
		&user.Role{},
		&user.User{},
		&token.ExpiredTokens{},
		&token.PasswordResetToken{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating models")
		panic(err)
	}

	// Migration'dan sonra rolleri seed et
	seedRoles(db)
}

// seedRoles veritabanına tüm rolleri ekler (eğer yoksa)
func seedRoles(db *gorm.DB) {
	roles := []user.Role{
		{Name: user.RoleDoctor},
		{Name: user.RoleAssistant},
		{Name: user.RoleIntern},
		{Name: user.RoleSecretary},
		{Name: user.RoleSecurity},
		{Name: user.RoleManager},
		{Name: user.RoleCleaner},
		{Name: user.RoleRadiologyTechnician},
		{Name: user.RoleAccountant},
		{Name: user.RolePatientConsultant},
		{Name: user.RoleItSupportSpecialist},
		{Name: user.RoleSupplyChainManager},
		{Name: user.RoleSterilizationTechnician},
		{Name: user.RoleHrManager},
		{Name: user.RoleOther},
		{Name: user.RoleOrthodontist},
		{Name: user.RoleClinicAdmin},
		{Name: user.RoleSuperAdmin},
	}

	for _, role := range roles {
		// Rol zaten varsa oluşturma (FirstOrCreate kullan)
		var existingRole user.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)

		if result.Error == gorm.ErrRecordNotFound {
			// Rol yok, oluştur
			if err := db.Create(&role).Error; err != nil {
				log.Error().
					Err(err).
					Str("role_name", string(role.Name)).
					Msg("Failed to seed role")
			} else {
				log.Info().
					Str("role_name", string(role.Name)).
					Msg("Role seeded successfully")
			}
		}
	}

	log.Info().Msg("Role seeding completed")
}
