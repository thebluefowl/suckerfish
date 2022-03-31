package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/suckerfish/service"
)

type AuthHandler interface {
	GetAuthURLs() echo.HandlerFunc
}

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

type authHandler struct {
	authService service.AuthService
}

func (handler *authHandler) GetAuthURLs() echo.HandlerFunc {
	return func(c echo.Context) error {
		urls := handler.authService.GetAuthURLs()
		return c.JSON(http.StatusOK, urls)
	}
}
