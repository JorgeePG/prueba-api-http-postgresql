package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

type user struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

var apiURL = "http://localhost:8080"

func TestAddUserOk(t *testing.T) {
	addURL := apiURL + "/add"

	// Generar datos únicos para cada ejecución
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)

	u1 := user{
		Username: fmt.Sprintf("testuser_%d_%d", timestamp, random),
		Email:    fmt.Sprintf("test_%d_%d@example.com", timestamp, random),
		Password: "password123",
		FullName: fmt.Sprintf("Test User %d", timestamp),
		IsActive: true,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(u1)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Hacer la petición HTTP POST al endpoint
	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Verificar el código de estado
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	t.Logf("Test completed successfully: %s", response.Message)
}

func TestAddUserBad(t *testing.T) {
	addURL := apiURL + "/add"

	u1 := user{
		Username: "", //Lo dejo vacío sabiedo que tiene que fallar
		Email:    "",
		Password: "password123",
		FullName: "Test User 1",
		IsActive: true,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(u1)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Hacer la petición HTTP POST al endpoint
	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()
	// Verificar el código de estado
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", response.Status)
	}

	t.Logf("Test completed successfully: %s", response.Message)
}
