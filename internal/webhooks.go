package internal

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type WebhookHandler struct {
	processor WebhookProcessor
}

func NewWebhookHandler(processor WebhookProcessor) *WebhookHandler {
	return &WebhookHandler{
		processor: processor,
	}
}

func (h *WebhookHandler) SubmitJob(w http.ResponseWriter, r *http.Request) error {
	defer func(r io.Closer) {
		if err := r.Close(); err != nil {
			slog.Warn("failed to close request body: %v", err)
		}
	}(r.Body)

	var request JobRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return err
	}

	jobID, err := h.processor.CreateJob(r.Context(), &request)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(NewJobResponse{
		CreatedAt:   time.Now(),
		JobID:       string(jobID),
		CallbackURL: request.CallbackURL,
	})
}

func (h *WebhookHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) error {
	jobID := r.PathValue("job_id")

	if jobID == "" {
		http.Error(w, "missing job_id query parameter", http.StatusBadRequest)
		return nil
	}

	status, err := h.processor.GetStatus(JobID(jobID))

	if err != nil {
		return err
	}

	type statusResp struct {
		Status string
	}
	return json.NewEncoder(w).Encode(statusResp{Status: status})
}
