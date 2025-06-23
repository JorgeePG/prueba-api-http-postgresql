package handler

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Add handles creation of new resources
func Add(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Resource added successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Update handles updating existing resources
func Update(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Resource updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Delete handles deletion of resources
func Delete(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Resource deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// List handles listing resources
func List(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Status:  "success",
		Message: "Resources retrieved successfully",
		Data:    []string{}, // Replace with actual data
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
