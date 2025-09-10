package postgres

import (
	"dental-clinic-system/infrastructure/config"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(DBConfig config.DatabaseConfig) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DBConfig.DNS))
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
		panic(err)
	}
	var result int
	err = db.Raw("SELECT 1").Scan(&result).Error
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to database")
		panic(err)
	}
	log.Info().Msg("Database connection successful")
	db.Exec("CREATE SCHEMA IF NOT EXISTS public")
	return db
}
