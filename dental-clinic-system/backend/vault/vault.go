package vault

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"os"
)

func ConnectVault() (*api.Client, error) {
	config := api.DefaultConfig()
	err := config.ReadEnvironment()
	if err != nil {
		return nil, err
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(os.Getenv("Root_Token"))

	err = client.SetAddress("http://127.0.0.1:8200")
	if err != nil {
		return nil, err
	}

	return client, nil
}

func GetJWTKeyFromVault() ([]byte, error) {
	client, err := ConnectVault()

	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Read("jwt-secrets/jwt-key")
	if err != nil {
		return nil, err
	}

	jwtKey := secret.Data["jwt_key"].(string)
	if jwtKey == "" {
		return nil, errors.New("JWT key not found")
	}
	return []byte(jwtKey), nil
}
