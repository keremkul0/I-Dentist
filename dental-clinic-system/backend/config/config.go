package config

import (
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"os"
)

const (
	ConfigName = "application"
	ConfigType = "yml"
	DefaultEnv = "qa"
)

func LoadConfig(configPath string) *ConfigModel {

	env := os.Getenv("ENV")
	if env == "" {
		log.Warn().Msg("ENV is not set, using default env")
		env = DefaultEnv
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigType(ConfigType)

	data := readConfig(env)
	return data
}

func readConfig(env string) *ConfigModel {

	viper.SetConfigName(ConfigName)
	readConfigErr := viper.ReadInConfig()
	if readConfigErr != nil {
		log.Fatal().Err(readConfigErr).Msg("Error reading config file")
		return nil
	}

	config := &ConfigModel{}
	v := viper.Sub(env)

	unMarshallErr := v.Unmarshal(config)
	if unMarshallErr != nil {
		log.Fatal().Err(unMarshallErr).Msg("Error unmarshalling config file")
		return nil
	}

	return config
}

func ReadConfigFromVault(client *api.Client, model *ConfigModel) error {
	secret, err := client.Logical().Read("secret/jwt_token")
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading secret from vault")
		return err
	}

	if secret == nil || secret.Data == nil {
		log.Fatal().Msg("JWT key not found")
		return err
	}

	jwtKey, ok := secret.Data["token"].(string)
	if !ok || jwtKey == "" {
		log.Fatal().Msg("JWT key not found")
		return err
	}
	model.JWT.SecretKey = jwtKey

	secret, err = client.Logical().Read("secret/mail_password")
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading secret from vault")
		return err
	}

	if secret == nil || secret.Data == nil {
		log.Fatal().Msg("Mail password not found")
		return err
	}

	mailPassword, ok := secret.Data["password"].(string)
	if !ok || mailPassword == "" {
		log.Fatal().Msg("Mail password not found")
		return err
	}
	model.Email.Password = mailPassword

	if err := model.ValidateConfig(); err != nil {
		log.Fatal().Err(err).Msg("Error validating config file")
		return nil
	}
	return nil
}
