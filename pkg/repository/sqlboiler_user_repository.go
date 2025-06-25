package repository

import (
	"context"
	"time"

	"github.com/JorgeePG/prueba-api-http-postgresql-/infraestructure/db"
	dbmodels "github.com/JorgeePG/prueba-api-http-postgresql-/infraestructure/db/models"
	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// SQLBoilerUserRepository implementa UserRepository usando SQLBoiler
type SQLBoilerUserRepository struct{}

// NewSQLBoilerUserRepository crea una nueva instancia del repositorio
func NewSQLBoilerUserRepository() UserRepository {
	return &SQLBoilerUserRepository{}
}

// Create crea un nuevo usuario en la base de datos
func (r *SQLBoilerUserRepository) Create(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	// Crear modelo SQLBoiler
	dbUser := &dbmodels.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FullName:  req.FullName,
		IsActive:  req.IsActive,
		CreatedAt: null.TimeFrom(time.Now()),
		UpdatedAt: null.TimeFrom(time.Now()),
	}

	// Insertar en base de datos
	if err := dbUser.Insert(ctx, db.DB, boil.Infer()); err != nil {
		return nil, err
	}

	// Convertir a modelo de dominio
	return r.dbUserToModel(dbUser), nil
}

// GetByID obtiene un usuario por su ID
func (r *SQLBoilerUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	dbUser, err := dbmodels.FindUser(ctx, db.DB, id)
	if err != nil {
		return nil, err
	}

	return r.dbUserToModel(dbUser), nil
}

// Update actualiza un usuario existente
func (r *SQLBoilerUserRepository) Update(ctx context.Context, id int, req *models.UpdateUserRequest) (*models.User, error) {
	// Buscar usuario existente
	dbUser, err := dbmodels.FindUser(ctx, db.DB, id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Username != "" {
		dbUser.Username = req.Username
	}
	if req.Email != "" {
		dbUser.Email = req.Email
	}
	if req.Password != "" {
		dbUser.Password = req.Password
	}
	if req.FullName.Valid {
		dbUser.FullName = req.FullName
	}
	if req.IsActive.Valid {
		dbUser.IsActive = req.IsActive
	}
	dbUser.UpdatedAt = null.TimeFrom(time.Now())

	// Guardar cambios
	if _, err := dbUser.Update(ctx, db.DB, boil.Infer()); err != nil {
		return nil, err
	}

	return r.dbUserToModel(dbUser), nil
}

// Delete elimina un usuario
func (r *SQLBoilerUserRepository) Delete(ctx context.Context, id int) error {
	dbUser, err := dbmodels.FindUser(ctx, db.DB, id)
	if err != nil {
		return err
	}

	_, err = dbUser.Delete(ctx, db.DB)
	return err
}

// List obtiene todos los usuarios
func (r *SQLBoilerUserRepository) List(ctx context.Context) ([]*models.User, error) {
	dbUsers, err := dbmodels.Users().All(ctx, db.DB)
	if err != nil {
		return nil, err
	}

	users := make([]*models.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		users[i] = r.dbUserToModel(dbUser)
	}

	return users, nil
}

// dbUserToModel convierte un modelo SQLBoiler a modelo de dominio
func (r *SQLBoilerUserRepository) dbUserToModel(dbUser *dbmodels.User) *models.User {
	return &models.User{
		ID:        dbUser.ID,
		Username:  dbUser.Username,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		FullName:  dbUser.FullName,
		IsActive:  dbUser.IsActive,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
