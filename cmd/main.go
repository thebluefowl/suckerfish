package main

import (
	"flag"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/thebluefowl/suckerfish/config"
	"github.com/thebluefowl/suckerfish/db"
	"github.com/thebluefowl/suckerfish/interface/rest"
	"github.com/thebluefowl/suckerfish/persistence/sql"
	"github.com/thebluefowl/suckerfish/service"
)

func main() {
	migrate := flag.Bool("migrate", false, "run migration")
	flag.Parse()

	appConfig, err := config.LoadAppConfig()
	if err != nil {
		panic(err)
	}

	authConfig, err := config.LoadAuthConfig()
	if err != nil {
		panic(err)
	}

	dbClient, err := db.GetPGClient(appConfig.Postgres)
	if err != nil {
		panic(err)
	}
	if *migrate {
		Migrate(dbClient)
		return
	}
	StartHTTPServer(appConfig, authConfig, dbClient)
}

func StartHTTPServer(appConfig *config.Config, authConfig *config.Auth, dbClient *db.PGClient) error {
	userRepository := sql.NewUserRepository(dbClient)
	authService := service.NewAuthService(userRepository, authConfig, appConfig)
	authHandler := rest.NewAuthHandler(authService)

	router := rest.NewRouter(authHandler)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(fmt.Sprintf(":%s", appConfig.Port))
}

func Migrate(client *db.PGClient) {
	fmt.Println("running migrations")
	if err := sql.Run(client); err != nil {
		panic(err)
	}
}
