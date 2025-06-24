package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/JorgeePG/prueba-api-http-postgresql-/db"
	"github.com/JorgeePG/prueba-api-http-postgresql-/db/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Add handles creation of new users
func Add(w http.ResponseWriter, r *http.Request) {
	// Configurar Content-Type para la respuesta
	w.Header().Set("Content-Type", "application/json")

	// Decodificar el cuerpo de la solicitud JSON
	var userData struct {
		Username string      `json:"username"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		FullName null.String `json:"full_name"`
		IsActive null.Bool   `json:"is_active"`
	}

	// Decodificar el JSON recibido
	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid request data: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar datos requeridos
	if userData.Username == "" || userData.Email == "" || userData.Password == "" {
		response := Response{
			Status:  "error",
			Message: "Username, email and password are required fields",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Crear un nuevo objeto de usuario
	user := models.User{
		Username:  userData.Username,
		Email:     userData.Email,
		Password:  userData.Password, // Nota: en producción deberías hashear la contraseña
		FullName:  userData.FullName,
		IsActive:  userData.IsActive,
		CreatedAt: null.TimeFrom(time.Now()),
		UpdatedAt: null.TimeFrom(time.Now()),
	}

	// Insertar usuario en la base de datos
	ctx := context.Background()
	err = user.Insert(ctx, db.DB, boil.Infer())

	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to create user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Devolver respuesta exitosa con el usuario creado
	response := Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    user,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Update handles updating existing users
func Update(w http.ResponseWriter, r *http.Request) {
	// Configurar Content-Type para la respuesta
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del usuario a actualizar de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := Response{
			Status:  "error",
			Message: "User ID is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convertir el ID a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid user ID: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Decodificar el cuerpo de la solicitud JSON
	var userData struct {
		Username string      `json:"username"`
		Email    string      `json:"email"`
		Password string      `json:"password,omitempty"`
		FullName null.String `json:"full_name"`
		IsActive null.Bool   `json:"is_active"`
	}

	// Decodificar el JSON recibido
	err = json.NewDecoder(r.Body).Decode(&userData)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid request data: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Crear contexto y buscar el usuario existente
	ctx := context.Background()
	user, err := models.FindUser(ctx, db.DB, id)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "User not found: " + err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Actualizar los campos del usuario
	if userData.Username != "" {
		user.Username = userData.Username
	}
	if userData.Email != "" {
		user.Email = userData.Email
	}
	if userData.Password != "" {
		user.Password = userData.Password // Nota: en producción deberías hashear la contraseña
	}
	if userData.FullName.Valid {
		user.FullName = userData.FullName
	}
	if userData.IsActive.Valid {
		user.IsActive = userData.IsActive
	}
	user.UpdatedAt = null.TimeFrom(time.Now())

	// Actualizar el usuario en la base de datos
	rowsAff, err := user.Update(ctx, db.DB, boil.Infer())
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to update user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if rowsAff == 0 {
		response := Response{
			Status:  "warning",
			Message: "No changes were made to the user",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete handles deletion of users
func Delete(w http.ResponseWriter, r *http.Request) {
	// Configurar Content-Type para la respuesta
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del usuario a eliminar de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := Response{
			Status:  "error",
			Message: "User ID is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convertir el ID a entero
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Invalid user ID: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Crear contexto y buscar el usuario existente
	ctx := context.Background()
	user, err := models.FindUser(ctx, db.DB, id)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "User not found: " + err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Eliminar el usuario de la base de datos
	rowsAff, err := user.Delete(ctx, db.DB)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Failed to delete user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	if rowsAff == 0 {
		response := Response{
			Status:  "warning",
			Message: "No user was deleted",
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data: map[string]int{
			"id": id,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// List handles listing resources
func List(w http.ResponseWriter, r *http.Request) {
	// Crear un contexto para la consulta
	ctx := context.Background()

	// Obtener la lista de usuarios usando SQLBoiler
	users, err := models.Users().All(ctx, db.DB)
	if err != nil {
		response := Response{
			Status:  "error",
			Message: "Error retrieving users: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    users,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
