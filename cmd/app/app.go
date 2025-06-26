// filepath: app/app.go
package app

import (
	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/config"
	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/server"
	"github.com/JorgeePG/prueba-api-http-postgresql-/infraestructure/db"
	"github.com/rs/zerolog/log"
)

type App struct {
	server *server.Server
	config *config.Config
}

func New() *App {
	return &App{
		config: config.Load(),
	}
}

func (a *App) Initialize() error {
	log.Info().Msg("Initializing application...")
	// Inicializar base de datos
	if err := db.Initialize(a.config.Database.ConnectionString()); err != nil {
		log.Error().Err(err).Msg("Error initializing database")
		return err
	}

	// Ejecutar migraciones
	if err := db.RunMigrations(); err != nil {
		log.Error().Err(err).Msg("Error running migrations")
		return err
	}

	log.Info().Msg("Database initialized and migrations applied successfully")
	// Configurar servidor
	a.server = server.New(a.config.Server.Port)
	a.server.SetupRoutes()

	log.Info().Msg("Server routes set up successfully")
	return nil
}

func (a *App) Run() error {
	log.Info().Msg("Starting application server...")
	return a.server.Start()
}

func (a *App) Shutdown() {
	log.Info().Msg("Shutting down application...")
	db.Close()
}
