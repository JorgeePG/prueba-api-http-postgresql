package repositories

import (
	"context"

	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/domain/entities"
)

// UserRepository define las operaciones de persistencia para User
type UserRepository interface {
	// Create crea un nuevo usuario en la base de datos
	Create(ctx context.Context, user *entities.User) error

	// GetByID obtiene un usuario por su ID
	GetByID(ctx context.Context, id int) (*entities.User, error)

	// GetByUsername obtiene un usuario por su username
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetByEmail obtiene un usuario por su email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update actualiza un usuario existente
	Update(ctx context.Context, user *entities.User) error

	// Delete elimina un usuario por su ID
	Delete(ctx context.Context, id int) error

	// List obtiene una lista paginada de usuarios
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count obtiene el total de usuarios
	Count(ctx context.Context) (int, error)
}
