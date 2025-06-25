package middleware

import (
	"log"
	"net/http"
	"time"
)

// CORS middleware para manejar las políticas de CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JSONMiddleware establece el Content-Type para JSON
func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware registra las peticiones HTTP
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Crear un ResponseWriter que capture el status code
		wrapper := &responseWrapper{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)

		log.Printf(
			"%s %s %d %v",
			r.Method,
			r.RequestURI,
			wrapper.statusCode,
			time.Since(start),
		)
	})
}

// responseWrapper es un wrapper para capturar el status code
type responseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
