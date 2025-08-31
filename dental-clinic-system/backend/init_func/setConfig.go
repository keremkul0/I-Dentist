package init_func

import (
	"dental-clinic-system/config"

	"github.com/rs/zerolog/log"
)

func SetConfig(configPath string) *config.ConfigModel {
	configModel := config.LoadConfig(configPath)
	if configModel == nil {
		log.Fatal().Msg("Error loading config")
		panic("Error loading config")
	}
	log.Info().Msg("Config loaded successfully")

	return configModel
}
