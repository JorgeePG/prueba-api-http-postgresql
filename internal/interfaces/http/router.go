package http

import (
	"net/http"

	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/interfaces/http/handlers"
	"github.com/JorgeePG/prueba-api-http-postgresql-/internal/interfaces/http/middleware"
	"github.com/gorilla/mux"
)

// Router configura las rutas de la aplicación
type Router struct {
	userHandler *handlers.UserHandler
}

// NewRouter crea una nueva instancia del router
func NewRouter(userHandler *handlers.UserHandler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

// Setup configura todas las rutas de la aplicación
func (router *Router) Setup() *mux.Router {
	r := mux.NewRouter()

	// Middleware global
	r.Use(middleware.CORS)
	r.Use(middleware.JSONMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// User routes
	users := api.PathPrefix("/users").Subrouter()
	users.HandleFunc("", router.userHandler.CreateUser).Methods("POST")
	users.HandleFunc("", router.userHandler.ListUsers).Methods("GET")
	users.HandleFunc("/{id:[0-9]+}", router.userHandler.GetUser).Methods("GET")
	users.HandleFunc("/{id:[0-9]+}", router.userHandler.UpdateUser).Methods("PUT")
	users.HandleFunc("/{id:[0-9]+}", router.userHandler.DeleteUser).Methods("DELETE")
	users.HandleFunc("/{id:[0-9]+}/password", router.userHandler.ChangePassword).Methods("PUT")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Server is running"}`))
	}).Methods("GET")

	return r
}
