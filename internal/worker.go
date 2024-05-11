package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"webhook_api/config"
)

type CallbackHandler func(ctx context.Context, request *JobRequest) error

type BackgroundWorker struct {
	sqsCfg       config.SQSConfig
	client       *sqs.Client
	handleRecord CallbackHandler
}

func NewWorker(cfg config.ServerConfig, handler CallbackHandler) *BackgroundWorker {
	return &BackgroundWorker{
		sqsCfg:       cfg.SQS,
		client:       configureSQSClient(cfg.SQS),
		handleRecord: handler,
	}
}

func (w *BackgroundWorker) Run() error {
	ctx := context.Background()
	for {
		slog.Debug("Polling SQS for messages")
		output, err := w.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:          &w.sqsCfg.QueueURL,
			WaitTimeSeconds:   w.sqsCfg.WaitTimeSeconds,
			VisibilityTimeout: w.sqsCfg.VisibilityTimeout,
		})
		if err != nil {
			slog.Error("unable to receive message from SQS", "error", err)
		}
		slog.Debug("Received message from SQS", "output", output)

		for _, rec := range output.Messages {
			var request JobRequest

			if err := json.Unmarshal([]byte(*rec.Body), &request); err != nil {
				slog.Error("unable to unmarshal request from SQS", "error", err)
			} else {
				if err := w.handleRecord(ctx, &request); err != nil {
					slog.Error("unable to process record from SQS", "error", err)
					continue
				}
			}

			if _, err := w.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &w.sqsCfg.QueueURL,
				ReceiptHandle: rec.ReceiptHandle,
			}); err != nil {
				return err
			}
			slog.Debug("Deleted message from SQS", "ID", rec.MessageId)
		}
	}
}

func NewCallbackHandler(cfg config.ServerConfig) CallbackHandler {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	return func(ctx context.Context, request *JobRequest) error {
		slog.Info("Received request from webhook", "jobID", request.ID)

		// simulate processing
		time.Sleep(300 * time.Millisecond)

		// make request
		body := strings.NewReader(fmt.Sprintf(`{"message": "success", "id": %s}`, request.ID))
		resp, err := httpClient.Post(request.CallbackURL, "application/json", body)

		if err != nil {
			return err
		}

		defer func(r io.ReadCloser) {
			if err := r.Close(); err != nil {
				slog.Warn("failed to close response body: %v", err)
			}
		}(resp.Body)

		if resp.StatusCode != http.StatusOK {
			responseBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			slog.Error("unexpected status code", "status_code", resp.StatusCode, "body", string(responseBody))
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		slog.Info("Callback request completed", "status_code", resp.StatusCode, "jobID", request.ID)
		return nil
	}
}
