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

	r.Use(loggingMiddleware(logger))

	r.Get("/health", healthHandler)

	handler := NewDocumentHandler(
		repository,
		storage,
	)

	r.Route("/documents", func(r chi.Router) {
		r.Get("/", listDocuments)
		r.Get("/{id}", handler.Get)
		r.Post("/", createDocument)
		r.Delete("/{id}", deleteDocument)
	})

	return r
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func listDocuments(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, []string{})
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusCreated, map[string]string{
		"message": "not implemented",
	})
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
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
