package internal

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) error {
	_, err := fmt.Fprintln(w, "This is the index page!")
	return err
}

func OnJobComplete(w http.ResponseWriter, r *http.Request) error {
	req, err := io.ReadAll(r.Body)

	if err != nil {
		return err
	}

	defer func(r io.ReadCloser) {
		if err := r.Close(); err != nil {
			slog.Error("unable to close request body", "error", err)
		}
	}(r.Body)

	slog.Info("Callback URL called", "request", req)

	w.WriteHeader(http.StatusOK)
	return nil
}
