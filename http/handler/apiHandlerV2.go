package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/models"
	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/repository"
	"github.com/JorgeePG/prueba-api-http-postgresql-/pkg/service"
)

// userService es una instancia global del servicio de usuarios
var userService *service.UserService

// init inicializa el servicio de usuarios
func init() {
	userRepo := repository.NewSQLBoilerUserRepository()
	userService = service.NewUserService(userRepo)
}

// AddV2 handles creation of new users using the new architecture
func AddV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Invalid request data: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := userService.CreateUser(r.Context(), &req)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Failed to create user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "User created successfully",
		Data:    user.ToAPIUser(),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateV2 handles updating existing users using the new architecture
func UpdateV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el ID del usuario a actualizar de la URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := models.Response{
			Status:  "error",
			Message: "User ID is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Invalid user ID: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Invalid request data: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := userService.UpdateUser(r.Context(), id, &req)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Failed to update user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "User updated successfully",
		Data:    user.ToAPIUser(),
	}
	json.NewEncoder(w).Encode(response)
}

// DeleteV2 handles deletion of users using the new architecture
func DeleteV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		response := models.Response{
			Status:  "error",
			Message: "User ID is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Invalid user ID: " + err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = userService.DeleteUser(r.Context(), id)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Failed to delete user: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "User deleted successfully",
		Data: map[string]int{
			"id": id,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// ListV2 handles listing users using the new architecture
func ListV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := userService.ListUsers(r.Context())
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Error retrieving users: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convertir a respuesta API (sin contraseñas)
	apiUsers := make([]*models.User, len(users))
	for i, user := range users {
		apiUsers[i] = user.ToAPIUser()
	}

	response := models.Response{
		Status:  "success",
		Message: "Users retrieved successfully",
		Data:    apiUsers,
	}
	json.NewEncoder(w).Encode(response)
}
