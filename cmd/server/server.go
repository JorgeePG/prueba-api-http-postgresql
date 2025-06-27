package server

import (
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/http/handler"
	"github.com/JorgeePG/prueba-api-http-postgresql-/http/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server struct {
	router  *mux.Router
	port    string
	useSSL  bool
	sslCert string
	sslKey  string
}

func New(port string, useSSL bool, sslCert, sslKey string) *Server {
	return &Server{
		router:  mux.NewRouter(),
		port:    port,
		useSSL:  useSSL,
		sslCert: sslCert,
		sslKey:  sslKey,
	}
}

func (s *Server) SetupRoutes() {
	// Middlewares de seguridad
	s.router.Use(middleware.SecurityHeaders)
	if s.useSSL {
		s.router.Use(middleware.HTTPSRedirect)
	}
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
	address := ":" + s.port

	if s.useSSL {
		log.Info().Msgf("🔒 Iniciando servidor HTTPS en el puerto %s", s.port)
		log.Info().Msgf("📜 Usando certificado: %s", s.sslCert)
		log.Info().Msgf("🔑 Usando clave privada: %s", s.sslKey)
		return http.ListenAndServeTLS(address, s.sslCert, s.sslKey, s.router)
	} else {
		log.Info().Msgf("⚠️  Iniciando servidor HTTP (no seguro) en el puerto %s", s.port)
		return http.ListenAndServe(address, s.router)
	}
}
