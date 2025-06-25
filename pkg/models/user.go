package models

import (
	"errors"

	"github.com/volatiletech/null/v8"
)

// User representa el modelo de usuario para la API (sin SQLBoiler)
type User struct {
	ID        int         `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Password  string      `json:"-"` // No se expone en JSON
	FullName  null.String `json:"full_name,omitempty"`
	IsActive  null.Bool   `json:"is_active,omitempty"`
	CreatedAt null.Time   `json:"created_at,omitempty"`
	UpdatedAt null.Time   `json:"updated_at,omitempty"`
}

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	FullName null.String `json:"full_name"`
	IsActive null.Bool   `json:"is_active"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Password string      `json:"password,omitempty"`
	FullName null.String `json:"full_name"`
	IsActive null.Bool   `json:"is_active"`
}

// Response representa una respuesta estándar de la API
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToAPIUser convierte el modelo de dominio a respuesta API
func (u *User) ToAPIUser() *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		// Password no se incluye por seguridad
	}
}

// Validate valida los campos requeridos del usuario
func (req *CreateUserRequest) Validate() error {
	if req.Username == "" {
		return errors.New("username is required")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// SetDefaults establece valores por defecto
func (req *CreateUserRequest) SetDefaults() {
	if !req.IsActive.Valid {
		req.IsActive = null.BoolFrom(true)
	}
}
