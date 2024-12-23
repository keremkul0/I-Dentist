package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

//validation for int values has error in validator library
//so we need to use the min=0 to validate the int values instead of required

type ConfigModel struct {
	Server   ServerConfig   `yaml:"server" validate:"required"`
	Database DatabaseConfig `yaml:"database" validate:"required"`
	Vault    VaultConfig    `yaml:"vault" validate:"required"`
	Email    EmailConfig    `yaml:"email" validate:"required"`
	Redis    RedisConfig    `yaml:"redis" validate:"required"`
	Log      LogConfig      `yaml:"log" validate:"required"`
	JWT      JWTConfig      `validate:"required"`
}

type ServerConfig struct {
	Port int `yaml:"port" validate:"min=0"`
}

type LogConfig struct {
	Level zerolog.Level `yaml:"level" validate:"required"`
}

type DatabaseConfig struct {
	DNS string `yaml:"dns" validate:"required"`
}

type VaultConfig struct {
	UnsealKeys       []string `yaml:"unsealKeys" validate:"required,min=3,dive,required"`
	Addr             string   `yaml:"addr" validate:"required,url"`
	InitialRootToken string   `yaml:"initialRootToken" validate:"required"`
}

type EmailConfig struct {
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required,gt=0"`
	User     string `yaml:"user" validate:"required,email"`
	Password string `yaml:"password" validate:"required"`
}

type RedisConfig struct {
	Addr string `yaml:"addr" validate:"required"`
	//for now redis password is not required
	//todo add password validation
	Password string `yaml:"password" validate:""`
	DB       int    `yaml:"db" validate:"min=0"`
}

type JWTConfig struct {
	SecretKey string `validate:"required"`
}

// ValidateConfig validates the configuration using the validator
func (c *ConfigModel) ValidateConfig() error {
	validate := validator.New()
	err := validate.Struct(c)
	if err != nil {
		return err
	}
	return nil
}
