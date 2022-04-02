package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/thebluefowl/suckerfish/config"
	"github.com/thebluefowl/suckerfish/domain"
	"github.com/thebluefowl/suckerfish/pkg/httpclient"
	"golang.org/x/oauth2"
)

type AuthService interface {
	GetAuthURLs() []AuthURL
	Authenticate(context.Context, *FetchTokenRequest) (*domain.User, error)
}

func NewAuthService(userRepository domain.UserRepository) AuthService {
	return &authService{userRepository: userRepository}
}

type authService struct {
	userRepository domain.UserRepository
}

type AuthURL struct {
	Provider string      `json:"provider"`
	AuthURL  interface{} `json:"url"`
}

func (*authService) GetAuthURLs() []AuthURL {
	c := config.AuthConfig

	urls := []AuthURL{}
	if c.Github != nil {
		state := ksuid.New().String()
		urls = append(urls, AuthURL{domain.ProviderGithub, c.Github.AuthCodeURL(state)})
	}
	return urls
}

type FetchTokenRequest struct {
	Provider string `param:"provider"`
	Code     string `query:"code"`
	State    string `query:"state"`
}

func (service *authService) Authenticate(ctx context.Context, request *FetchTokenRequest) (*domain.User, error) {
	c := config.AuthConfig
	user := &domain.User{}
	switch request.Provider {
	case domain.ProviderGithub:
		token, err := c.Github.Exchange(ctx, request.Code)
		if err != nil {
			return nil, err
		}
		user, err = service.fetchGithubUser(ctx, token)
		if err != nil {
			return nil, err
		}
	}
	isNewUser, err := service.IsNewUser(user)
	if err != nil {
		return nil, err
	}

	if !isNewUser {
		user.IsNewUser = false
		return user, nil
	}

	user.ID = ksuid.New().String()
	user.IsNewUser = true
	if err := service.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (service *authService) IsNewUser(user *domain.User) (bool, error) {
	user, err := service.userRepository.GetByEmail(user.Email)
	if err != nil {
		return true, err
	}
	if user != nil {
		return false, nil
	}
	return true, nil
}

func (*authService) fetchGithubUser(ctx context.Context, token *oauth2.Token) (*domain.User, error) {
	request, err := http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("token %s", token.AccessToken))
	client := httpclient.GetClient(&httpclient.HTTPClientOpts{ConnTimeout: 2 * time.Second, ReadTimeout: 10 * time.Second})
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	type UserResponse struct {
		AvatarURL string `json:"avatar_url"`
		Company   string `json:"company"`
		Email     string `json:"email"`
		Location  string `json:"location"`
		Name      string `json:"name"`
	}

	userResponse := &UserResponse{}
	if err := decoder.Decode(userResponse); err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:      userResponse.Name,
		Email:     userResponse.Email,
		Provider:  domain.ProviderGithub,
		AvatarURL: userResponse.AvatarURL,
		Location:  userResponse.Location,
		Company:   userResponse.Company,
		Token:     token.AccessToken,
	}

	return user, nil
}
