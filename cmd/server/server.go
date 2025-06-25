package server

import (
	"fmt"
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/http/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/http/middleware"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	port   string
}

func New(port string) *Server {
	return &Server{
		router: mux.NewRouter(),
		port:   port,
	}
}

func (s *Server) SetupRoutes() {
	s.router.Use(middleware.CspControl)

	// Rutas originales (mantener compatibilidad)
	s.router.HandleFunc("/", handler.List).Methods("GET")
	s.router.HandleFunc("/add", handler.Add).Methods("POST")
	s.router.HandleFunc("/update", handler.Update).Methods("POST")
	s.router.HandleFunc("/delete", handler.Delete).Methods("GET")

	// Nuevas rutas con arquitectura mejorada (v2)
	api := s.router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/users", handler.ListV2).Methods("GET")
	api.HandleFunc("/users", handler.AddV2).Methods("POST")
	api.HandleFunc("/users/update", handler.UpdateV2).Methods("POST") // Mantener formato original por compatibilidad
	api.HandleFunc("/users/delete", handler.DeleteV2).Methods("GET")  // Mantener formato original por compatibilidad
}

func (s *Server) Start() error {
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
