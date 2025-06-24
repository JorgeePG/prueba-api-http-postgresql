package main

import (
	"log"

	"github.com/JorgeePG/prueba-api-http-postgresql-/cmd/app"
)

func main() {
	application := app.New()

	if err := application.Initialize(); err != nil {
		log.Fatalf("Error al inicializar la aplicación: %v", err)
	}
	defer application.Shutdown()

	if err := application.Run(); err != nil {
		log.Fatalf("Error al ejecutar la aplicación: %v", err)
	}
}
