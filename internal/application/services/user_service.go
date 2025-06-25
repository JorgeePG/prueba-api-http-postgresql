package services

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/application/dto"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/domain/entities"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/domain/repositories"
)

// UserService encapsula la lógica de negocio para usuarios
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService crea una nueva instancia del servicio de usuarios
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario
func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// Verificar si el usuario ya existe
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// Crear entidad de usuario
	user, err := entities.NewUser(req.Username, req.Email, req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user entity: %w", err)
	}

	// Establecer campos opcionales
	if req.FullName != nil {
		user.FullName = req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = req.IsActive
	}

	// Guardar en repositorio
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// GetUser obtiene un usuario por ID
func (s *UserService) GetUser(ctx context.Context, id int) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return s.toUserResponse(user), nil
}

// UpdateUser actualiza un usuario existente
func (s *UserService) UpdateUser(ctx context.Context, id int, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Verificar unicidad de username y email si se están actualizando
	if req.Username != "" && req.Username != user.Username {
		existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
		if existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	if req.Email != "" && req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmail(ctx, req.Email)
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	// Actualizar entidad
	if err := user.Update(req.Username, req.Email, req.FullName, req.IsActive); err != nil {
		return nil, fmt.Errorf("failed to update user entity: %w", err)
	}

	// Guardar cambios
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// DeleteUser elimina un usuario
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ListUsers obtiene una lista paginada de usuarios
func (s *UserService) ListUsers(ctx context.Context, page, perPage int) (*dto.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	offset := (page - 1) * perPage

	users, err := s.userRepo.List(ctx, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	userResponses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(user)
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &dto.UsersListResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}, nil
}

// ChangePassword cambia la contraseña de un usuario
func (s *UserService) ChangePassword(ctx context.Context, id int, req dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return errors.New("user not found")
	}

	// Aquí deberías verificar la contraseña anterior con hash
	// Por simplicidad, asumimos que la verificación es correcta
	if user.Password != req.OldPassword {
		return errors.New("invalid old password")
	}

	if err := user.ChangePassword(req.NewPassword); err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user password: %w", err)
	}

	return nil
}

// toUserResponse convierte una entidad User a UserResponse
func (s *UserService) toUserResponse(user *entities.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
