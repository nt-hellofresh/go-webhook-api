package internal

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"webhook_api/config"
)

type WebhookProcessor interface {
	CreateJob(ctx context.Context, request *JobRequest) (JobID, error)
	GetStatus(jobID JobID) (string, error)
}

type SQSWorkQueue struct {
	sqsConfig config.SQSConfig
	client    *sqs.Client
}

func NewSQSWorkQueue(cfg config.SQSConfig, client *sqs.Client) *SQSWorkQueue {
	return &SQSWorkQueue{
		sqsConfig: cfg,
		client:    client,
	}
}

func (q *SQSWorkQueue) queueMessage(ctx context.Context, r *JobRequest) error {
	bytes, err := json.Marshal(r)

	if err != nil {
		return err
	}

	slog.Info("args", "jobID", r.ID, "payload", string(bytes))
	output, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.sqsConfig.QueueURL),
		MessageBody: aws.String(string(bytes)),
	})

	if err != nil {
		return err
	}

	slog.Info("Successfully queued message onto SQS", "ID", output.MessageId)
	return nil
}

func (q *SQSWorkQueue) CreateJob(ctx context.Context, request *JobRequest) (JobID, error) {
	jobID := JobID(ulid.Make().String())
	request.ID = jobID
	err := q.queueMessage(ctx, request)
	return jobID, err
}

func (q *SQSWorkQueue) GetStatus(_ JobID) (string, error) {
	// faking response for now
	return "in progress", nil
}
