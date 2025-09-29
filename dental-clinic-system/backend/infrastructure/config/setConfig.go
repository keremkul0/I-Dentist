package config

import (
	"github.com/rs/zerolog/log"
)

func SetConfig(configPath string) *ConfigModel {
	configModel := LoadConfig(configPath)
	if configModel == nil {
		log.Fatal().Msg("Error loading config")
		panic("Error loading config")
	}
	log.Info().Msg("Config loaded successfully")

	return configModel
}
