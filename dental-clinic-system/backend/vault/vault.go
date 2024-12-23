package vault

import (
	"dental-clinic-system/config"
	"fmt"
	"github.com/hashicorp/vault/api"

	"github.com/rs/zerolog/log"
)

func ConnectVault(vaultConfig config.VaultConfig) (*api.Client, error) {

	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = vaultConfig.Addr

	client, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating vault client")
		return nil, err
	}

	client.SetToken(vaultConfig.InitialRootToken)

	sealStatus, err := client.Sys().SealStatus()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting vault seal status")
		return nil, err
	}

	if sealStatus.Sealed {
		for _, unsealKey := range vaultConfig.UnsealKeys {

			if unsealKey == "" {
				log.Fatal().Msg("Unseal key not found")
				return nil, fmt.Errorf("unseal key not found")
			}
			if unsealResponse, err := client.Sys().Unseal(unsealKey); err != nil {
				log.Fatal().Err(err).Msg("Error unsealing vault")
				return nil, err
			} else if !unsealResponse.Sealed {
				log.Info().Msg("Vault successfully unsealed")
				break
			}
		}

		if sealStatus, err = client.Sys().SealStatus(); err != nil || sealStatus.Sealed {
			log.Fatal().Err(err).Msg("Error unsealing vault")
			return nil, err
		}
	} else {
		log.Info().Msg("Vault is already unsealed")
	}
	return client, nil
}
