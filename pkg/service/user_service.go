package service

import (
	"context"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"
	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/repository"
)

// UserService maneja la lógica de negocio para usuarios
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService crea una nueva instancia del servicio
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Validar datos
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Establecer valores por defecto
	req.SetDefaults()

	// Crear usuario
	return s.userRepo.Create(ctx, req)
}

// GetUser obtiene un usuario por ID
func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUser actualiza un usuario
func (s *UserService) UpdateUser(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	return s.userRepo.Update(ctx, id, req)
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers obtiene todos los usuarios
func (s *UserService) ListUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.List(ctx)
}
