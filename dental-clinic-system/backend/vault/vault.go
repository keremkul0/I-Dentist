package vault

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
)

func ConnectVault() (*api.Client, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := api.DefaultConfig()
	config.Address = os.Getenv("VAULT_ADDR")

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(os.Getenv("Initial_Root_Token"))

	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		return nil, err
	}

	if sealStatus.Sealed {
		for i := 1; i <= 5; i++ {
			unsealKey := os.Getenv(fmt.Sprintf("Unseal_Key_%d", i))
			if unsealKey == "" {
				return nil, fmt.Errorf("unseal key %d not found", i)
			}

			if unsealResponse, err := client.Sys().Unseal(unsealKey); err != nil {
				return nil, fmt.Errorf("vault unseal failed: %v", err)
			} else if !unsealResponse.Sealed {
				log.Info().Msg("Vault successfully unsealed")
				break
			}
		}

		if sealStatus, err = client.Sys().SealStatus(); err != nil || sealStatus.Sealed {
			return nil, fmt.Errorf("vault unseal failed")
		}
	} else {
		log.Info().Msg("Vault is already unsealed")
	}
	return client, nil
}

func GetJWTKeyFromVault() ([]byte, error) {
	client, err := ConnectVault()
	if err != nil {
		return nil, err
	}

	secret, err := client.Logical().Read("secret/jwt_token")
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		return nil, errors.New("data could not be retrieved from Vault or secret/jwt_token not found")
	}

	jwtKey, ok := secret.Data["token"].(string)
	if !ok || jwtKey == "" {
		return nil, errors.New("JWT key not found or format is incorrect")
	}

	return []byte(jwtKey), nil
}
