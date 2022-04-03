package config

import (
	"github.com/spf13/viper"
	"github.com/thebluefowl/suckerfish/db"
)

type Config struct {
	// Port is the HTTP web server port
	Port string `mapstructure:"port"`
	// SigningKey is used to sign the JWT token
	SigningKey string `mapstructure:"signing_key"`

	// Postgres config
	Postgres *db.PGConfig `mapstructure:"postgres"`
}

func LoadAppConfig() (*Config, error) {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	appConfig := &Config{}
	if err := viper.Unmarshal(appConfig); err != nil {
		return nil, err
	}

	return appConfig, nil
}
