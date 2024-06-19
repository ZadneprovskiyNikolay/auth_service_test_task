package config

import (
	"context"
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

type SecretManager interface {
	MustGetSecretField(path, field string) string
}

type secretManager struct {
	env         string
	vaultClient *vault.Client
}

func NewSecretManager(env string) (SecretManager, error) {
	var vaultClient *vault.Client
	if env == EnvDev || env == EnvProd {
		var err error
		vaultCfg := vault.DefaultConfig()
		vaultClient, err = vault.NewClient(vaultCfg)
		if err != nil {
			return nil, err
		}
	}

	return &secretManager{
		env:         strings.ToLower(env),
		vaultClient: vaultClient,
	}, nil
}

func (s *secretManager) MustGetSecretField(path, field string) string {
	switch s.env {
	case EnvLocal:
		return os.Getenv(strings.Join([]string{path, field}, "_"))
	case EnvDev, EnvProd:
		secret, err := s.vaultClient.KVv2("secret").Get(context.Background(), path)
		if err != nil {
			log.Fatalf("could not get secret %s: %v", path, err)
		}

		value, ok := secret.Data[field].(string)
		if !ok {
			log.Fatalf("secret field type assertion failed: %T %#v", secret.Data[field], secret.Data[field])
		}

		return value
	default:
		return ""
	}
}
