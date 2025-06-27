package server

import (
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/http/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/http/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	s.router.HandleFunc("/update", handler.UpdateV2).Methods("POST")
	s.router.HandleFunc("/delete", handler.DeleteV2).Methods("GET")
	// MQTT routes
	s.router.HandleFunc("/mqtt/topic/add", handler.AddTopicSubscriber).Methods("GET")
	s.router.HandleFunc("/mqtt/topic/delete", handler.DeleteTopicSubscriber).Methods("GET")
	s.router.HandleFunc("/mqtt/messages", handler.ListMqttMessages).Methods("GET")
}

func (s *Server) Start() error {
	log.Info().Msgf("Iniciando servidor en el puerto %s", s.port)
	return http.ListenAndServe(":"+s.port, s.router)
}
