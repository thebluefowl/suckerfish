package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"github.com/thebluefowl/suckerfish/config"
	"github.com/thebluefowl/suckerfish/domain"
	"github.com/thebluefowl/suckerfish/pkg/httpclient"
	"golang.org/x/oauth2"
)

type AuthService interface {
	GetAuthURLs() []AuthURL
	GenerateJWT(*domain.User) (string, error)
	AuthenticateFromProvider(context.Context, *FetchTokenRequest) (*domain.User, error)
}

func NewAuthService(userRepository domain.UserRepository, authConfig *config.Auth, appConfig *config.Config) AuthService {
	return &authService{
		userRepository: userRepository,
		authConfig:     authConfig,
		appConfig:      appConfig,
	}
}

type authService struct {
	userRepository domain.UserRepository
	authConfig     *config.Auth
	appConfig      *config.Config
}

type AuthURL struct {
	Provider string      `json:"provider"`
	AuthURL  interface{} `json:"url"`
}

func (service *authService) GetAuthURLs() []AuthURL {
	urls := []AuthURL{}
	if service.authConfig.Github != nil {
		state := ksuid.New().String()
		urls = append(urls, AuthURL{domain.ProviderGithub, service.authConfig.Github.AuthCodeURL(state)})
	}
	return urls
}

type FetchTokenRequest struct {
	Provider string `param:"provider"`
	Code     string `query:"code"`
	State    string `query:"state"`
}

func (service *authService) AuthenticateFromProvider(ctx context.Context, request *FetchTokenRequest) (*domain.User, error) {
	user := &domain.User{}
	switch request.Provider {
	case domain.ProviderGithub:
		token, err := service.authConfig.Github.Exchange(ctx, request.Code)
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

type CustomClaims struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
	Provider  string `json:"provider"`
	Company   string `json:"company"`
	Location  string `json:"location"`
	IsStaff   bool   `json:"is_staff"`
	jwt.StandardClaims
}

func (service *authService) GenerateJWT(user *domain.User) (string, error) {
	signingKey := []byte(service.appConfig.SigningKey)
	claims := CustomClaims{
		user.ID,
		user.Name,
		user.AvatarURL,
		user.Email,
		user.Provider,
		user.Company,
		user.Location,
		user.IsStaff,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (service *authService) Authenticate(tokenString string) (*domain.User, error) {
	signingKey := []byte(service.appConfig.SigningKey)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &domain.User{
			ID:        claims["id"].(string),
			Name:      claims["name"].(string),
			AvatarURL: claims["avatar_url"].(string),
			Email:     claims["email"].(string),
			Provider:  claims["provider"].(string),
			Company:   claims["company"].(string),
			Location:  claims["location"].(string),
			IsStaff:   claims["is_staff"].(bool),
		}, nil
	} else {
		return nil, err
	}
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

func (*authService) fetchGithubUser(_ context.Context, token *oauth2.Token) (*domain.User, error) {
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
