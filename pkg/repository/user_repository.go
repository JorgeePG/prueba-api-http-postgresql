package repository

import (
	"context"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"
)

// UserRepository define las operaciones de base de datos para usuarios
type UserRepository interface {
	Create(ctx context.Context, user *models.CreateUserRequest) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	Update(ctx context.Context, id int, user *models.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*models.User, error)
}
