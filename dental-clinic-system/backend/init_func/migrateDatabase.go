package init_func

import (
	"dental-clinic-system/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Appointment{},
		&models.Clinic{},
		&models.Patient{},
		&models.Procedure{},
		&models.Role{},
		&models.User{},
		&models.ExpiredTokens{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Error migrating models")
		panic(err)
	}
}
