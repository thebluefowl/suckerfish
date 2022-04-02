package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var AuthConfig *Auth

type Auth struct {
	Github *oauth2.Config
}

func LoadAuthConfig() error {
	type config struct {
		Github *struct {
			ClientID     string   `mapstructure:"client_id"`
			ClientSecret string   `mapstructure:"client_secret"`
			Scopes       []string `mapstructure:"scopes"`
		} `yaml:"github"`
	}
	viper.AddConfigPath("./")
	viper.SetConfigName("auth")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}

	AuthConfig = &Auth{}

	if c.Github != nil {
		AuthConfig.Github = &oauth2.Config{
			ClientID:     c.Github.ClientID,
			ClientSecret: c.Github.ClientSecret,
			Scopes:       c.Github.Scopes,
			Endpoint:     github.Endpoint,
			RedirectURL:  "http://localhost:7272/auth/github",
		}
	}

	return nil
}
