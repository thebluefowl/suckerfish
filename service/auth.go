package service

import "golang.org/x/oauth2/github"

type AuthService interface {
	GetAuthURLs() []AuthURL
}

func NewAuthService() AuthService {
	return new(authService)
}

type authService struct {
}

type AuthURL struct {
	Provider string `json:"provider"`
	AuthURL  string `json:"url"`
}

func (service *authService) GetAuthURLs() []AuthURL {
	urls := []AuthURL{
		{
			Provider: "github",
			AuthURL:  github.Endpoint.AuthURL,
		},
	}
	return urls
}
