package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	authhttp "navigate-business-api/internal/auth/http"
)

func NewRouter(authHandler *authhttp.Handler, authenticator *authhttp.Authenticator) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(8 * time.Second))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:4321", "http://127.0.0.1:4321"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Content-Type", "Authorization"},
	}))
	router.Get("/api/v1/hello-world", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Hello, world!"})
	})

	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Use(authenticator.RequireGoogle)
		r.Post("/sync", authHandler.Sync)
		r.Get("/me", authHandler.Current)
	})

	return router
}
