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
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating models")
		panic(err)
	}
}
