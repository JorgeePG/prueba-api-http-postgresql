package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/JorgeePG/prueba-api-http-postgresql-/mqtt/subscriber"
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

func AddTopicSubscriber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		response := models.Response{
			Status:  "error",
			Message: "Topic is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar que el topic no esté vacío después de trim spaces
	topic = strings.TrimSpace(topic)
	if topic == "" {
		response := models.Response{
			Status:  "error",
			Message: "Topic cannot be empty or just spaces",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validar formato del topic MQTT
	if !isValidMQTTTopic(topic) {
		response := models.Response{
			Status:  "error",
			Message: "Invalid MQTT topic format",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Intentar agregar el suscriptor
	if err := subscriber.AddTopicSubscriber(topic); err != nil {
		response := models.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Subscriber added successfully for topic: " + topic,
	}
	json.NewEncoder(w).Encode(response)
}

func DeleteTopicSubscriber(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	topic := r.URL.Query().Get("topic")
	if topic == "" {
		response := models.Response{
			Status:  "error",
			Message: "Topic is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Intentar eliminar el suscriptor
	if err := subscriber.DeleteTopicSubscriber(topic); err != nil {
		response := models.Response{
			Status:  "error",
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "Subscriber removed successfully for topic: " + topic,
	}
	json.NewEncoder(w).Encode(response)
}

func ListMqttMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Obtener el parámetro num_mensajes de la URL
	numMensajesStr := r.URL.Query().Get("num_mensajes")
	limit := 100 // Valor por defecto

	// Si se proporciona el parámetro, convertirlo a int
	if numMensajesStr != "" {
		var err error
		limit, err = strconv.Atoi(numMensajesStr)
		if err != nil {
			response := models.Response{
				Status:  "error",
				Message: "Invalid num_mensajes parameter: " + err.Error(),
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Validar que el límite sea positivo
		if limit <= 0 {
			response := models.Response{
				Status:  "error",
				Message: "num_mensajes must be a positive number",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	messages, err := subscriber.ListMqttMessages(limit)
	if err != nil {
		response := models.Response{
			Status:  "error",
			Message: "Error retrieving MQTT messages: " + err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := models.Response{
		Status:  "success",
		Message: "MQTT messages retrieved successfully",
		Data:    messages,
	}
	json.NewEncoder(w).Encode(response)
}

// isValidMQTTTopic valida si un topic MQTT es válido
func isValidMQTTTopic(topic string) bool {
	if topic == "" || len(topic) > 65535 {
		return false
	}
	// MQTT topic no debe contener caracteres null
	if strings.Contains(topic, "\x00") {
		return false
	}
	return true
}
