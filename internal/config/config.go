package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"dev"`
	DB         DB         `yaml:"db"`
	HTTPServer HTTPServer `yaml:"http_server"`
	SMTP       SMTP       `yaml:"smtp"`
	Emails     Emails     `yaml:"emails"`
	Auth       Auth       `yaml:"auth"`
}

type Auth struct {
	AccessTokenDuration  time.Duration `yaml:"access_token_duration" env-required:"true"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration" env-required:"true"`
	JWTPrivateKey        string
}

type HTTPServer struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
}

type DB struct {
	Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
	UserName string
	Password string
}

type SMTP struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	UserName string
	Password string
}

type Emails struct {
	SupportEmail string `yaml:"support_email" env-required:"true"`
}

var (
	once sync.Once
	cfg  Config
)

func MustNew() *Config {
	once.Do(func() {
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			log.Fatal("CONFIG_PATH is not set")
		}

		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file does not exist: %s", configPath)
		}

		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("could not read config: %s", err)
		}

		var err error
		secretManager, err := NewSecretManager(cfg.Env)
		if err != nil {
			log.Fatalf("could not create secret manager: %s", err)
		}

		cfg.Auth.JWTPrivateKey = secretManager.MustGetSecretField("JWT", "PRIVATE_KEY")
		cfg.DB.UserName = secretManager.MustGetSecretField("DB", "USER")
		cfg.DB.Password = secretManager.MustGetSecretField("DB", "PASSWORD")
		cfg.SMTP.UserName = secretManager.MustGetSecretField("SMTP", "USERNAME")
		cfg.SMTP.Password = secretManager.MustGetSecretField("SMTP", "PASSWORD")
	})

	return &cfg
}
