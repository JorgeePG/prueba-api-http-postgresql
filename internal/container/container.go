package container

import (
	"database/sql"

	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/application/services"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/domain/repositories"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/infrastructure/persistence/sqlboiler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/interfaces/http/handlers"
)

// Container contiene todas las dependencias de la aplicación
type Container struct {
	// Repositories
	UserRepository repositories.UserRepository

	// Services
	UserService *services.UserService

	// Handlers
	UserHandler *handlers.UserHandler
}

// NewContainer crea un nuevo contenedor con todas las dependencias
func NewContainer(db *sql.DB) *Container {
	// Repositories
	userRepo := sqlboiler.NewUserRepository(db)

	// Services
	userService := services.NewUserService(userRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)

	return &Container{
		UserRepository: userRepo,
		UserService:    userService,
		UserHandler:    userHandler,
	}
}
