package main

import (
	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/suckerfish/interface/rest"
	"github.com/thebluefowl/suckerfish/service"
)

func main() {
	StartHTTPServer()
}

func StartHTTPServer() error {
	authService := service.NewAuthService()
	authHandler := rest.NewAuthHandler(authService)

	router := rest.NewRouter(authHandler)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}
