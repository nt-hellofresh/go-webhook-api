package main

import (
	"log/slog"
	"net/http"
	"os"
	"webhook_api/config"
	"webhook_api/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
		slog.String("version", config.AppVersion()),
	}))
	slog.SetDefault(logger)

	server := internal.NewServer(internal.WithWebhookRoutes)

	slog.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		slog.Error("Error starting server: %s\n", err)
	}
}
