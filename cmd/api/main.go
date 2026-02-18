package main

import (
	"log"

	"github.com/Matheus-Lima-Moreira/financial-pocket/internal/bootstrap"
)

func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info().
		Str("port", app.Config.Port).
		Msg("starting server")

	if err := app.Router.Run(":" + app.Config.Port); err != nil {
		log.Fatal(err)
	}
}
