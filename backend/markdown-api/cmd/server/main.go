package main

import (
	"log/slog"
	"net/http"
	"os"

	"markdown-api/internal/api"
	"markdown-api/internal/database"
	"markdown-api/internal/documents"
	"markdown-api/internal/storage"
)

func main() {
	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelInfo,
		},
	))

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.Open("documents.db")
	if err != nil {
		logger.Error("database error", "error", err)
		os.Exit(1)
	}

	repository := documents.NewRepository(db)

	markdownStorage := storage.NewMarkdownStorage("./data/documents")

	router := api.NewRouter(
		logger,
		repository,
		markdownStorage,
	)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info("starting server",
		"port", port,
	)

	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped",
			"error", err,
		)
		os.Exit(1)
	}
}
