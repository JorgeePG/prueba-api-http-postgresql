// filepath: app/app.go
package app

import (
	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/config"
	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/server"
	"github.com/JorgeePG/prueba-api-http-postgresql-/db"
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
	// Inicializar base de datos
	if err := db.Initialize(a.config.Database.ConnectionString()); err != nil {
		return err
	}

	// Ejecutar migraciones
	if err := db.RunMigrations(); err != nil {
		return err
	}

	// Configurar servidor
	a.server = server.New(a.config.Server.Port)
	a.server.SetupRoutes()

	return nil
}

func (a *App) Run() error {
	return a.server.Start()
}

func (a *App) Shutdown() {
	db.Close()
}
