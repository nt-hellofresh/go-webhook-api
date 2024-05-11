package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
	"log/slog"
	"net/http"
	"webhook_api/config"
)

type ServerOpts func(mux *http.ServeMux)

func NewServer(opts ...ServerOpts) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", ToHandlerFunc(Index))

	for _, opt := range opts {
		opt(mux)
	}

	// This really should be part of a different server
	// but included here for testing purposes
	mux.HandleFunc("POST /api-2/complete", ToHandlerFunc(OnJobComplete))
	return mux
}

func WithWebhookRoutes(mux *http.ServeMux) {
	cfg := config.MustLoadFromYAML()
	slog.Info("Loaded config from YAML", "config", cfg)

	processor := configureSQSBacklog(cfg.SQS)
	handler := NewWebhookHandler(processor)

	mux.HandleFunc("POST /api/jobs", ToHandlerFunc(handler.SubmitJob))
	mux.HandleFunc("GET /api/jobs/{job_id}", ToHandlerFunc(handler.GetJobStatus))
}

func configureSQSBacklog(cfg config.SQSConfig) WebhookProcessor {
	sqsClient := configureSQSClient(cfg)
	return NewSQSWorkQueue(cfg, sqsClient)
}

func configureSQSClient(cfg config.SQSConfig) *sqs.Client {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           cfg.Endpoint,
			SigningRegion: cfg.Region,
		}, nil
	})
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithEndpointResolverWithOptions(resolver),
	)

	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(awsCfg)
	return sqsClient
}
