package db

import (
	"testing"
)

func TestConnection(t *testing.T) {
	// Cadena de conexión a PostgreSQL usando los valores del docker-compose.yml
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=api_db sslmode=disable"

	err := Initialize(connStr)
	if err != nil {
		t.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer Close()

	// Si llegamos aquí sin errores, la conexión fue exitosa
	t.Log("Conexión a la base de datos exitosa")
}
