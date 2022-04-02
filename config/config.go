package config

import (
	"github.com/spf13/viper"
	"github.com/thebluefowl/suckerfish/db"
)

var AppConfig *Config

type Config struct {
	Port     string       `mapstructure:"port"`
	Postgres *db.PGConfig `mapstructure:"postgres"`
}

func LoadAppConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return nil
	}

	return nil
}
