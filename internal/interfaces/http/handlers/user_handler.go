package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/application/dto"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/application/services"
	"github.com/gorilla/mux"
)

// UserHandler maneja las peticiones HTTP relacionadas con usuarios
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler crea una nueva instancia del handler de usuarios
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser maneja la creación de usuarios
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "CREATE_FAILED", "Failed to create user", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusCreated, "User created successfully", user)
}

// GetUser maneja la obtención de un usuario por ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		h.sendErrorResponse(w, http.StatusBadRequest, "MISSING_ID", "User ID is required", "")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", err.Error())
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", "")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "GET_FAILED", "Failed to get user", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "User retrieved successfully", user)
}

// UpdateUser maneja la actualización de usuarios
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		h.sendErrorResponse(w, http.StatusBadRequest, "MISSING_ID", "User ID is required", "")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, req)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", "")
			return
		}
		h.sendErrorResponse(w, http.StatusBadRequest, "UPDATE_FAILED", "Failed to update user", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "User updated successfully", user)
}

// DeleteUser maneja la eliminación de usuarios
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		h.sendErrorResponse(w, http.StatusBadRequest, "MISSING_ID", "User ID is required", "")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", err.Error())
		return
	}

	err = h.userService.DeleteUser(r.Context(), id)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", "")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "DELETE_FAILED", "Failed to delete user", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers maneja la obtención de lista de usuarios
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page := 1
	perPage := 10

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	users, err := h.userService.ListUsers(r.Context(), page, perPage)
	if err != nil {
		h.sendErrorResponse(w, http.StatusInternalServerError, "LIST_FAILED", "Failed to list users", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Users retrieved successfully", users)
}

// ChangePassword maneja el cambio de contraseña
func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, exists := vars["id"]
	if !exists {
		h.sendErrorResponse(w, http.StatusBadRequest, "MISSING_ID", "User ID is required", "")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID", err.Error())
		return
	}

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body", err.Error())
		return
	}

	err = h.userService.ChangePassword(r.Context(), id, req)
	if err != nil {
		if err.Error() == "user not found" {
			h.sendErrorResponse(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found", "")
			return
		}
		if err.Error() == "invalid old password" {
			h.sendErrorResponse(w, http.StatusBadRequest, "INVALID_PASSWORD", "Invalid old password", "")
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "PASSWORD_CHANGE_FAILED", "Failed to change password", err.Error())
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, "Password changed successfully", nil)
}

// sendSuccessResponse envía una respuesta exitosa
func (h *UserHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse envía una respuesta de error
func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, code, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := dto.APIResponse{
		Status:  "error",
		Message: message,
		Error: &dto.APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}

	json.NewEncoder(w).Encode(response)
}
