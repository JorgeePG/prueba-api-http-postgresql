package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	// Hacer la petición HTTP POST al endpoint
	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request without error")
	defer resp.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err, "Should decode response without error")
	assert.Equal(t, "success", response.Status, "Response status should be success")

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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	// Hacer la petición HTTP POST al endpoint
	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request without error")
	defer resp.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 Bad Request for invalid data")

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err, "Should decode response without error")
	assert.Equal(t, "error", response.Status, "Response status should be error for bad request")

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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	// Hacer la petición HTTP POST al endpoint
	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request without error")
	defer resp.Body.Close()

	// Decodificar la respuesta para obtener el id
	var addResp Response
	err = json.NewDecoder(resp.Body).Decode(&addResp)
	assert.NoError(t, err, "Should decode add user response without error")

	// Extraer el id del usuario creado
	type userID struct {
		ID int `json:"id"`
	}
	var uid userID
	// Si Data es un objeto con campo id
	b, _ := json.Marshal(addResp.Data)
	json.Unmarshal(b, &uid)
	assert.NotZero(t, uid.ID, "Should get valid user ID from creation response")

	// Ahora eliminar el usuario usando el id
	deleteURL := apiURL + "/delete"

	respDel, err := http.Get(deleteURL + fmt.Sprintf("?id=%d", uid.ID))
	assert.NoError(t, err, "Should make HTTP request to delete without error")
	defer respDel.Body.Close()

	assert.Equal(t, http.StatusOK, respDel.StatusCode, "Should return 200 OK for valid delete")

	var responseDel Response
	err = json.NewDecoder(respDel.Body).Decode(&responseDel)
	assert.NoError(t, err, "Should decode delete response without error")
	assert.Equal(t, "success", responseDel.Status, "Delete response status should be success")

	t.Logf("Delete test completed successfully: %s", responseDel.Message)
}

func TestDeleteBad(t *testing.T) {
	// Eliminar un usuario usando un id inválido
	deleteURL := apiURL + "/delete"

	respDel, err := http.Get(deleteURL + fmt.Sprintf("?id=%d", -1)) // Usando un id inválido
	assert.NoError(t, err, "Should make HTTP request to delete without error")
	defer respDel.Body.Close()

	assert.Equal(t, http.StatusNotFound, respDel.StatusCode, "Should return 404 Not Found for invalid ID")

	var responseDel Response
	err = json.NewDecoder(respDel.Body).Decode(&responseDel)
	assert.NoError(t, err, "Should decode delete response without error")
	assert.Equal(t, "error", responseDel.Status, "Delete response status should be error for invalid ID")

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
	assert.NoError(t, err, "Should marshal user 1 to JSON without error")

	resp1, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData1))
	assert.NoError(t, err, "Should create test user 1 without error")
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
	assert.NoError(t, err, "Should marshal user 2 to JSON without error")

	resp2, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData2))
	assert.NoError(t, err, "Should create test user 2 without error")
	resp2.Body.Close()

	// Ahora hacer la petición GET para listar usuarios
	listURL := apiURL + "/"

	respList, err := http.Get(listURL)
	assert.NoError(t, err, "Should make HTTP request to list users without error")
	defer respList.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, respList.StatusCode, "Should return 200 OK for list users")

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(respList.Body).Decode(&response)
	assert.NoError(t, err, "Should decode list response without error")
	assert.Equal(t, "success", response.Status, "List response status should be success")
	assert.NotNil(t, response.Data, "Should have data in response")

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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	resp, err := http.Post(addURL, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request to add user without error")
	defer resp.Body.Close()

	// Decodificar la respuesta para obtener el id
	var addResp Response
	err = json.NewDecoder(resp.Body).Decode(&addResp)
	assert.NoError(t, err, "Should decode add user response without error")

	// Extraer el id del usuario creado
	type userID struct {
		ID int `json:"id"`
	}
	var uid userID
	b, _ := json.Marshal(addResp.Data)
	json.Unmarshal(b, &uid)
	assert.NotZero(t, uid.ID, "Should get valid user ID from creation response")

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
	assert.NoError(t, err, "Should marshal updated user to JSON without error")

	// Hacer la petición HTTP POST al endpoint de update
	updateResp, err := http.Post(updateURL+fmt.Sprintf("?id=%d", uid.ID), "application/json", bytes.NewBuffer(updateJsonData))
	assert.NoError(t, err, "Should make HTTP request to update without error")
	defer updateResp.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusOK, updateResp.StatusCode, "Should return 200 OK for valid update")

	// Verificar la respuesta JSON
	var updateResponse Response
	err = json.NewDecoder(updateResp.Body).Decode(&updateResponse)
	assert.NoError(t, err, "Should decode update response without error")
	assert.Equal(t, "success", updateResponse.Status, "Update response status should be success")

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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	// Hacer la petición HTTP POST con un ID inválido
	updateResp, err := http.Post(updateURL+"?id=-1", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request to update without error")
	defer updateResp.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusNotFound, updateResp.StatusCode, "Should return 404 Not Found for invalid ID")

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(updateResp.Body).Decode(&response)
	assert.NoError(t, err, "Should decode update response without error")
	assert.Equal(t, "error", response.Status, "Response status should be error for invalid ID")

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
	assert.NoError(t, err, "Should marshal user to JSON without error")

	// Hacer la petición HTTP POST sin ID
	updateResp, err := http.Post(updateURL, "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err, "Should make HTTP request to update without error")
	defer updateResp.Body.Close()

	// Verificar el código de estado
	assert.Equal(t, http.StatusBadRequest, updateResp.StatusCode, "Should return 400 Bad Request for missing ID")

	// Verificar la respuesta JSON
	var response Response
	err = json.NewDecoder(updateResp.Body).Decode(&response)
	assert.NoError(t, err, "Should decode update response without error")
	assert.Equal(t, "error", response.Status, "Response status should be error for missing ID")

	t.Logf("Update user missing ID test completed successfully: %s", response.Message)
}
