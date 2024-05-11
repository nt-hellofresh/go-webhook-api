package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testJobID = JobID("mock-job-id")

type MockWebhookProcessor struct{}

func (m *MockWebhookProcessor) CreateJob(ctx context.Context, request *JobRequest) (JobID, error) {
	return testJobID, nil
}

func (m *MockWebhookProcessor) GetStatus(jobID JobID) (string, error) {
	return "mock-status", nil
}

func TestNewWebhookHandler(t *testing.T) {
	processor := &MockWebhookProcessor{}
	handler := NewWebhookHandler(processor)

	t.Run("should return correct job response", func(t *testing.T) {
		data := struct {
			CallbackURL string `json:"callback_url"`
		}{
			CallbackURL: "http://localhost:8080/my-callback-endpoint",
		}

		jsonBytes, err := json.Marshal(data)
		assert.NoError(t, err)

		resp := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/jobs", bytes.NewReader(jsonBytes))

		handler.SubmitJob(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		outData := struct {
			JobID       JobID  `json:"job_id"`
			CallbackURL string `json:"callback_url"`
		}{}

		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&outData))
		assert.Equal(t, testJobID, outData.JobID)
		assert.Equal(t, data.CallbackURL, outData.CallbackURL)
	})
}
