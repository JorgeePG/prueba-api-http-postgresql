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

func TestDeleteUserOk(t *testing.T) {

	//Primero se crea un usuario para poder eliminarlo
	addURL := apiURL + "/add"

	// Generar datos únicos para cada ejecución
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)
	username := fmt.Sprintf("testuser_%d_%d", timestamp, random)
	email := fmt.Sprintf("test_%d_%d@example.com", timestamp, random)
	password := "password123"
	fullName := fmt.Sprintf("Test User %d", timestamp)
	u1 := user{
		Username: username,
		Email:    email,
		Password: password,
		FullName: fullName,
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

	// Decodificar la respuesta para obtener el id
	var addResp Response
	err = json.NewDecoder(resp.Body).Decode(&addResp)
	if err != nil {
		t.Fatalf("Failed to decode add user response: %v", err)
	}

	// Extraer el id del usuario creado
	type userID struct {
		ID int `json:"id"`
	}
	var uid userID
	// Si Data es un objeto con campo id
	b, _ := json.Marshal(addResp.Data)
	json.Unmarshal(b, &uid)
	if uid.ID == 0 {
		t.Fatalf("No se pudo obtener el id del usuario creado")
	}

	// Ahora eliminar el usuario usando el id

	deleteURL := apiURL + "/delete"

	respDel, err := http.Get(deleteURL + fmt.Sprintf("?id=%d", uid.ID))
	if err != nil {
		t.Fatalf("Failed to make HTTP request to delete: %v", err)
	}
	defer respDel.Body.Close()

	if respDel.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, respDel.StatusCode)
	}

	var responseDel Response
	err = json.NewDecoder(respDel.Body).Decode(&responseDel)
	if err != nil {
		t.Fatalf("Failed to decode delete response: %v", err)
	}
	if responseDel.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", responseDel.Status)
	}

	t.Logf("Delete test completed successfully: %s", responseDel.Message)
}

func TestDeleteBad(t *testing.T) {

	// Ahora eliminar el usuario usando el id

	deleteURL := apiURL + "/delete"

	respDel, err := http.Get(deleteURL + fmt.Sprintf("?id=%d", -1)) // Usando un id inválido
	if err != nil {
		t.Fatalf("Failed to make HTTP request to delete: %v", err)
	}
	defer respDel.Body.Close()

	if respDel.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, respDel.StatusCode)
	}

	var responseDel Response
	err = json.NewDecoder(respDel.Body).Decode(&responseDel)
	if err != nil {
		t.Fatalf("Failed to decode delete response: %v", err)
	}
	if responseDel.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", responseDel.Status)
	}

	t.Logf("Delete test completed successfully: %s", responseDel.Message)
}
