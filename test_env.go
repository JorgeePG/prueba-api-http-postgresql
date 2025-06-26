package main

import (
	"fmt"
	"os"

	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/config"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load("../.env"); err != nil {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found, using system environment variables")
		}
	} else {
		fmt.Println("Environment variables loaded from .env file")
	}

	// Verificar variables de entorno directamente
	fmt.Printf("DB_PASSWORD from env: '%s'\n", os.Getenv("DB_PASSWORD"))
	fmt.Printf("DB_HOST from env: '%s'\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_USER from env: '%s'\n", os.Getenv("DB_USER"))

	// Cargar configuración
	cfg := config.Load()

	// Mostrar configuración cargada
	fmt.Printf("Config Password: '%s'\n", cfg.Database.Password)
	fmt.Printf("Config Host: '%s'\n", cfg.Database.Host)
	fmt.Printf("Config User: '%s'\n", cfg.Database.User)
	fmt.Printf("Config Port: %d\n", cfg.Database.Port)

	// Mostrar connection string
	fmt.Printf("Connection String: '%s'\n", cfg.Database.ConnectionString())
}
