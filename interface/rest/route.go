package rest

import "github.com/labstack/echo/v4"

type EchoRouter interface {
	AddRoutes(*echo.Echo)
}

func NewRouter(authHandler AuthHandler) EchoRouter {
	return &router{
		authHandler: authHandler,
	}
}

type router struct {
	authHandler AuthHandler
}

func (r *router) AddRoutes(e *echo.Echo) {
	e.GET("/auth/urls", r.authHandler.GetAuthURLs())
	e.GET("/auth/:provider", r.authHandler.HandleRedirect())
}
