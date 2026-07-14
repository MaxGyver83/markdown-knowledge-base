package api

import (
	"encoding/json"
	"log/slog"
	"markdown-api/internal/documents"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(
	logger *slog.Logger,
	repository *documents.Repository,
	storage Storage,
) http.Handler {
	r := chi.NewRouter()

	r.Use(corsMiddleware)
	r.Use(loggingMiddleware(logger))

	r.Get("/health", healthHandler)

	handler := NewDocumentHandler(
		repository,
		storage,
	)

	resetHandler := NewResetHandler(logger, repository, storage)

	r.Route("/documents", func(r chi.Router) {
		r.Post("/reset", resetHandler.Reset)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.Get)
		r.Post("/", handler.Create)
		r.Put("/{id}", handler.Update)
		r.Delete("/{id}", handler.Delete)
	})

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func jsonResponse(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			logger.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
			)

			next.ServeHTTP(w, r)
		})
	}
}
