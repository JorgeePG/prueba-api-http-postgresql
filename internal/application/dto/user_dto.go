package dto

// CreateUserRequest representa la solicitud para crear un usuario
type CreateUserRequest struct {
	Username string  `json:"username" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	FullName *string `json:"full_name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// UpdateUserRequest representa la solicitud para actualizar un usuario
type UpdateUserRequest struct {
	Username string  `json:"username,omitempty"`
	Email    string  `json:"email,omitempty"`
	FullName *string `json:"full_name,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// ChangePasswordRequest representa la solicitud para cambiar contraseña
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// UserResponse representa la respuesta con datos del usuario
type UserResponse struct {
	ID        int     `json:"id"`
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	FullName  *string `json:"full_name,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// UsersListResponse representa la respuesta para lista de usuarios
type UsersListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// APIResponse representa una respuesta estándar de la API
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError representa un error en la API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
