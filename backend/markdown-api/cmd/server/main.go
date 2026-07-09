package main

import (
	"log/slog"
	"net/http"
	"os"

	"markdown-api/internal/api"
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

	router := api.NewRouter(logger)

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
