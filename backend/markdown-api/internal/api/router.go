package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(loggingMiddleware(logger))

	r.Get("/health", healthHandler)

	r.Route("/documents", func(r chi.Router) {
		r.Get("/", listDocuments)
		r.Get("/{id}", getDocument)
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

func getDocument(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	jsonResponse(w, http.StatusOK, map[string]string{
		"id": id,
	})
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
