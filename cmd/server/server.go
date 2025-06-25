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

	s.router.HandleFunc("/", handler.ListV2).Methods("GET")
	s.router.HandleFunc("/add", handler.AddV2).Methods("POST")
	s.router.HandleFunc("/update", handler.UpdateV2).Methods("POST") // Mantener formato original por compatibilidad
	s.router.HandleFunc("/delete", handler.DeleteV2).Methods("GET")  // Mantener formato original por compatibilidad
}

func (s *Server) Start() error {
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
