package main

import (
	"log/slog"
	"os"
	"webhook_api/config"
	"webhook_api/internal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil).WithAttrs([]slog.Attr{
		slog.String("version", config.AppVersion()),
	}))
	slog.SetDefault(logger)

	cfg := config.MustLoadFromYAML()
	handler := internal.NewCallbackHandler(cfg)
	worker := internal.NewWorker(cfg, handler)

	if err := worker.Run(); err != nil {
		slog.Error("Worker exited with error", "error", err)
	}
}
