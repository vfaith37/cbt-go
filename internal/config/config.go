package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Address string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Address: viper.GetString("SERVER_ADDRESS"),
		},
		Database: DatabaseConfig{
			DSN: viper.GetString("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret: viper.GetString("JWT_SECRET"),
			TTL:    viper.GetDuration("JWT_TTL"),
		},
	}

	return config, nil
}
