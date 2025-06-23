package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/db"
	"github.com/JorgeePG/prueba-api-http-postgresql-/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Configurar la conexión a la base de datos
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=api_db sslmode=disable"
	err := db.Initialize(connStr)
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer db.Close()

	// Ejecutar migraciones para crear o actualizar el esquema
	if err = db.RunMigrations(); err != nil {
		log.Fatalf("Error al ejecutar migraciones: %v", err)
	}

	// Configurar el router
	r := mux.NewRouter()
	r.Use(middleware.CspControl)

	r.HandleFunc("/", handler.List).Methods("GET")
	r.HandleFunc("/add", handler.Add).Methods("POST")
	r.HandleFunc("/update", handler.Update).Methods("POST")
	r.HandleFunc("/delete", handler.Delete).Methods("POST")

	// Iniciar el servidor
	fmt.Println("Servidor iniciado en http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
