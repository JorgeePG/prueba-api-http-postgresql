package entities

import (
	"errors"
	"time"
)

// User representa la entidad de usuario en el dominio
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // No se expone en JSON por seguridad
	FullName  *string   `json:"full_name,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser crea una nueva instancia de User con validaciones básicas
func NewUser(username, email, password string) (*User, error) {
	if err := validateUserData(username, email, password); err != nil {
		return nil, err
	}

	now := time.Now()
	isActive := true

	return &User{
		Username:  username,
		Email:     email,
		Password:  password,
		IsActive:  &isActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Update actualiza los campos modificables del usuario
func (u *User) Update(username, email string, fullName *string, isActive *bool) error {
	if username != "" {
		u.Username = username
	}
	if email != "" {
		u.Email = email
	}
	if fullName != nil {
		u.FullName = fullName
	}
	if isActive != nil {
		u.IsActive = isActive
	}
	u.UpdatedAt = time.Now()
	return nil
}

// ChangePassword cambia la contraseña del usuario
func (u *User) ChangePassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}
	u.Password = newPassword
	u.UpdatedAt = time.Now()
	return nil
}

// validateUserData valida los datos básicos del usuario
func validateUserData(username, email, password string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if email == "" {
		return errors.New("email is required")
	}
	if password == "" {
		return errors.New("password is required")
	}
	// Aquí podrías agregar más validaciones como formato de email, etc.
	return nil
}
