package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const TokenTTL = 12 * time.Hour

type Config struct {
	AppName      string `mapstructure:"APP_NAME"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPass       string `mapstructure:"DB_PASSWORD"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBName       string `mapstructure:"DB_NAME"`
	DBSSL        string `mapstructure:"DB_SSL"`
	DefaultCoins int    `mapstructure:"DEFAULT_COINS"`
	SigningKey   string `mapstructure:"JWT_SIGNING_KEY"`
	Salt         string `mapstructure:"HASH_SALT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	viper.SetDefault("APP_NAME", "avito_shop")
	viper.SetDefault("DB_PORT", "5432")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Env file not found")
		} else {
			return config, fmt.Errorf("failed to read config: %w", err)
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
