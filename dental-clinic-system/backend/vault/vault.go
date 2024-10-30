package vault

import (
	"context"
	"errors"
	"github.com/hashicorp/vault/api"
	"os"
	"time"
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

func AddTokenToVaultBlacklist(client *api.Client, token string, expiry time.Time) error {
	data := map[string]interface{}{
		"token":  token,
		"expiry": expiry.Format(time.RFC3339),
	}

	// KV store’da saklayın (örnek yol: auth/blacklist/)
	_, err := client.KVv2("auth-secrets").Put(context.Background(), "blacklist/"+token, data)
	return err
}

func IsTokenBlacklisted(client *api.Client, token string) bool {
	_, err := client.KVv2("auth-secrets").Get(context.Background(), "blacklist/"+token)
	return err == nil
}
