package main

import (
	"os"

	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/app"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	application := app.New()

	log.Info().Msg("Iniciando inicialización de la aplicación")
	if err := application.Initialize(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Error al inicializar la aplicación")
		return
	}
	log.Info().
		Msg("Aplicación inicializada correctamente")
	defer application.Shutdown()

	if err := application.Run(); err != nil {
		log.Fatal().
			Err(err).
			Msg("Error al ejecutar la aplicación")
	}
	log.Info().
		Msg("Aplicación finalizada correctamente")
}
