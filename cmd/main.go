package main

import (
	"fmt"
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/middleware"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(middleware.CspControl)

	r.HandleFunc("/", handler.List).Methods("GET")
	r.HandleFunc("/add", handler.Add).Methods("POST")
	r.HandleFunc("/update", handler.Update).Methods("POST")
	r.HandleFunc("/delete", handler.Delete).Methods("POST")

	fmt.Println("Servidor iniciado en http://localhost:8080")
	http.ListenAndServe(":8080", r)

}
