package config

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Auth struct {
	Github *oauth2.Config
}

func LoadAuthConfig() (*Auth, error) {
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
		return nil, err
	}
	c := &config{}
	err = viper.Unmarshal(c)
	if err != nil {
		return nil, err
	}

	authConfig := &Auth{}

	if c.Github != nil {
		authConfig.Github = &oauth2.Config{
			ClientID:     c.Github.ClientID,
			ClientSecret: c.Github.ClientSecret,
			Scopes:       c.Github.Scopes,
			Endpoint:     github.Endpoint,
			RedirectURL:  "http://localhost:7272/auth/github",
		}
	}

	return authConfig, nil
}
