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

func TestListUsers(t *testing.T) {
	// Primero agregar algunos usuarios para asegurar que hay datos para listar
	addURL := apiURL + "/add"
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)

	// Crear primer usuario
	u1 := user{
		Username: fmt.Sprintf("listtest1_%d_%d", timestamp, random),
		Email:    fmt.Sprintf("listtest1_%d_%d@example.com", timestamp, random),
		Password: "password123",
		FullName: fmt.Sprintf("List Test User 1 %d", timestamp),
		IsActive: true,
	}

	jsonData1, err := json.Marshal(u1)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	resp1, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData1))
	if err != nil {
		t.Fatalf("Failed to create test user 1: %v", err)
	}
	resp1.Body.Close()

	// Crear segundo usuario
	u2 := user{
		Username: fmt.Sprintf("listtest2_%d_%d", timestamp, random),
		Email:    fmt.Sprintf("listtest2_%d_%d@example.com", timestamp, random),
		Password: "password456",
		FullName: fmt.Sprintf("List Test User 2 %d", timestamp),
		IsActive: false,
	}

	jsonData2, err := json.Marshal(u2)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	resp2, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData2))
	if err != nil {
		t.Fatalf("Failed to create test user 2: %v", err)
	}
	resp2.Body.Close()

	// Ahora hacer la petición GET para listar usuarios
	listURL := apiURL + "/"

	respList, err := http.Get(listURL)
	if err != nil {
		t.Fatalf("Failed to make HTTP request to list users: %v", err)
	}
	defer respList.Body.Close()

	// Verificar el código de estado
	if respList.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, respList.StatusCode)
	}

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(respList.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode list response: %v", err)
	}

	if response.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", response.Status)
	}

	// Verificar que la respuesta contiene datos
	if response.Data == nil {
		t.Error("Expected data in response, got nil")
	}

	t.Logf("List users test completed successfully: %s", response.Message)
}

func TestUpdateUserOk(t *testing.T) {
	// Primero crear un usuario para poder actualizarlo
	addURL := apiURL + "/add"
	timestamp := time.Now().UnixNano()
	random := rand.Intn(10000)

	originalUser := user{
		Username: fmt.Sprintf("updatetest_%d_%d", timestamp, random),
		Email:    fmt.Sprintf("updatetest_%d_%d@example.com", timestamp, random),
		Password: "originalpassword",
		FullName: fmt.Sprintf("Original User %d", timestamp),
		IsActive: true,
	}

	// Crear el usuario
	jsonData, err := json.Marshal(originalUser)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request to add user: %v", err)
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
	b, _ := json.Marshal(addResp.Data)
	json.Unmarshal(b, &uid)
	if uid.ID == 0 {
		t.Fatalf("No se pudo obtener el id del usuario creado")
	}

	// Ahora actualizar el usuario
	updateURL := apiURL + "/update"
	updatedUser := user{
		Username: fmt.Sprintf("updated_%d_%d", timestamp, random),
		Email:    fmt.Sprintf("updated_%d_%d@example.com", timestamp, random),
		Password: "newpassword123",
		FullName: fmt.Sprintf("Updated User %d", timestamp),
		IsActive: false,
	}

	// Convertir a JSON
	updateJsonData, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("Failed to marshal updated user to JSON: %v", err)
	}

	// Hacer la petición HTTP POST al endpoint de update
	updateResp, err := http.Post(updateURL+fmt.Sprintf("?id=%d", uid.ID), "application/json", bytes.NewBuffer(updateJsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request to update: %v", err)
	}
	defer updateResp.Body.Close()

	// Verificar el código de estado
	if updateResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, updateResp.StatusCode)
	}

	// Verificar la respuesta JSON
	var updateResponse Response
	err = json.NewDecoder(updateResp.Body).Decode(&updateResponse)
	if err != nil {
		t.Fatalf("Failed to decode update response: %v", err)
	}

	if updateResponse.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", updateResponse.Status)
	}

	t.Logf("Update user test completed successfully: %s", updateResponse.Message)
}

func TestUpdateUserBad(t *testing.T) {
	// Intentar actualizar un usuario con ID inválido
	updateURL := apiURL + "/update"

	updateUser := user{
		Username: "shouldnotwork",
		Email:    "shouldnotwork@example.com",
		Password: "password123",
		FullName: "Should Not Work",
		IsActive: true,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(updateUser)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Hacer la petición HTTP POST con un ID inválido
	updateResp, err := http.Post(updateURL+"?id=-1", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request to update: %v", err)
	}
	defer updateResp.Body.Close()

	// Verificar el código de estado
	if updateResp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, updateResp.StatusCode)
	}

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(updateResp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode update response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", response.Status)
	}

	t.Logf("Update user bad test completed successfully: %s", response.Message)
}

func TestUpdateUserBadMissingID(t *testing.T) {
	// Intentar actualizar un usuario sin proporcionar ID
	updateURL := apiURL + "/update"

	updateUser := user{
		Username: "shouldnotwork",
		Email:    "shouldnotwork@example.com",
		Password: "password123",
		FullName: "Should Not Work",
		IsActive: true,
	}

	// Convertir a JSON
	jsonData, err := json.Marshal(updateUser)
	if err != nil {
		t.Fatalf("Failed to marshal user to JSON: %v", err)
	}

	// Hacer la petición HTTP POST sin ID
	updateResp, err := http.Post(updateURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make HTTP request to update: %v", err)
	}
	defer updateResp.Body.Close()

	// Verificar el código de estado
	if updateResp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, updateResp.StatusCode)
	}

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(updateResp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode update response: %v", err)
	}

	if response.Status != "error" {
		t.Errorf("Expected status 'error', got '%s'", response.Status)
	}

	t.Logf("Update user missing ID test completed successfully: %s", response.Message)
}
