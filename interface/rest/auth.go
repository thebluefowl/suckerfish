package rest

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/suckerfish/service"
)

type AuthHandler interface {
	GetAuthURLs() echo.HandlerFunc
	HandleRedirect() echo.HandlerFunc
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

		// json.Marshal automatical escapes non utf-8 characters.  Hand writing this to ensure URLs stay unescaped.
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		c.Response().WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(c.Response())
		encoder.SetEscapeHTML(false)
		return encoder.Encode(urls)
	}
}

func (handler *authHandler) HandleRedirect() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(service.FetchTokenRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid request"})
		}
		user, err := handler.authService.Authenticate(c.Request().Context(), request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "failed to authenticate"})
		}
		_ = user
		return c.JSON(http.StatusOK, "")
	}
}
