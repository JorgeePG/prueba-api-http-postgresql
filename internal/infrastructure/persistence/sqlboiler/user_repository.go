package sqlboiler

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/JorgeePG/prueba-api-http-postgresql-/infraestructure/db/models"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/domain/entities"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// UserRepository implementa el repositorio de usuarios usando SQLBoiler
type UserRepository struct {
	db boil.ContextExecutor
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db boil.ContextExecutor) *UserRepository {
	return &UserRepository{db: db}
}

// Create crea un nuevo usuario en la base de datos
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
	dbUser := r.toSQLBoilerUser(user)

	if err := dbUser.Insert(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Actualizar el ID de la entidad con el ID generado por la DB
	user.ID = dbUser.ID
	return nil
}

// GetByID obtiene un usuario por su ID
func (r *UserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	dbUser, err := models.FindUser(ctx, r.db, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No encontrado
		}
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return r.toEntity(dbUser), nil
}

// GetByUsername obtiene un usuario por su username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	dbUser, err := models.Users(
		models.UserWhere.Username.EQ(username),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No encontrado
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return r.toEntity(dbUser), nil
}

// GetByEmail obtiene un usuario por su email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	dbUser, err := models.Users(
		models.UserWhere.Email.EQ(email),
	).One(ctx, r.db)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No encontrado
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return r.toEntity(dbUser), nil
}

// Update actualiza un usuario existente
func (r *UserRepository) Update(ctx context.Context, user *entities.User) error {
	dbUser := r.toSQLBoilerUser(user)

	if _, err := dbUser.Update(ctx, r.db, boil.Infer()); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete elimina un usuario por su ID
func (r *UserRepository) Delete(ctx context.Context, id int) error {
	dbUser := &models.User{ID: id}

	if _, err := dbUser.Delete(ctx, r.db); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// List obtiene una lista paginada de usuarios
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	dbUsers, err := models.Users(
		qm.OrderBy(models.UserColumns.CreatedAt+" DESC"),
		qm.Limit(limit),
		qm.Offset(offset),
	).All(ctx, r.db)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	users := make([]*entities.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.toEntity(dbUser)
	}

	return users, nil
}

// Count obtiene el total de usuarios
func (r *UserRepository) Count(ctx context.Context) (int, error) {
	count, err := models.Users().Count(ctx, r.db)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return int(count), nil
}

// toSQLBoilerUser convierte una entidad User a modelo SQLBoiler
func (r *UserRepository) toSQLBoilerUser(user *entities.User) *models.User {
	dbUser := &models.User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: null.TimeFrom(user.CreatedAt),
		UpdatedAt: null.TimeFrom(user.UpdatedAt),
	}

	if user.FullName != nil {
		dbUser.FullName = null.StringFrom(*user.FullName)
	}

	if user.IsActive != nil {
		dbUser.IsActive = null.BoolFrom(*user.IsActive)
	}

	return dbUser
}

// toEntity convierte un modelo SQLBoiler a entidad User
func (r *UserRepository) toEntity(dbUser *models.User) *entities.User {
	user := &entities.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		CreatedAt: dbUser.CreatedAt.Time,
		UpdatedAt: dbUser.UpdatedAt.Time,
	}

	if dbUser.FullName.Valid {
		fullName := dbUser.FullName.String
		user.FullName = &fullName
	}

	if dbUser.IsActive.Valid {
		isActive := dbUser.IsActive.Bool
		user.IsActive = &isActive
	}

	return user
}
