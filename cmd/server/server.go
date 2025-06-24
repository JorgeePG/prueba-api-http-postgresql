package server

import (
	"fmt"
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/middleware"
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

	s.router.HandleFunc("/", handler.List).Methods("GET")
	s.router.HandleFunc("/add", handler.Add).Methods("POST")
	s.router.HandleFunc("/update", handler.Update).Methods("POST")
	s.router.HandleFunc("/delete", handler.Delete).Methods("GET")
}

func (s *Server) Start() error {
	fmt.Printf("Servidor iniciado en http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
