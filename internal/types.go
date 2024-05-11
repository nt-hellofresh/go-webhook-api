package internal

import (
	"net/http"
	"time"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request) error

func ToHandlerFunc(rh RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := rh(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type JobID string

type JobRequest struct {
	ID          JobID  `json:"id"`
	CallbackURL string `json:"callback_url"`
}

type NewJobResponse struct {
	CreatedAt   time.Time `json:"created_at"`
	JobID       string    `json:"job_id"`
	CallbackURL string    `json:"callback_url"`
}
