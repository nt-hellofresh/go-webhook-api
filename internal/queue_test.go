package internal

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
	"webhook_api/config"
)

var (
	queueURL = "http://localhost:4566/000000000000/webhook-queue"
	endpoint = "http://localhost:4566"
	region   = "ap-southeast-2"
)

func TestIntegrationSendMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})
	awsCfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(region),
		awsConfig.WithEndpointResolverWithOptions(resolver),
		//awsConfig.WithClientLogMode(aws.LogResponseWithBody),
	)

	assert.NoError(t, err)

	sqsClient := sqs.NewFromConfig(awsCfg)
	processor := NewSQSWorkQueue(config.SQSConfig{QueueURL: queueURL}, sqsClient)

	_, err = processor.CreateJob(context.TODO(), &JobRequest{CallbackURL: "http://localhost:8080/callback"})

	assert.NoError(t, err)
}
