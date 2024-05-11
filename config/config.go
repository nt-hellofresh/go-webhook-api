package config

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"log/slog"
	"os"
)

type SQSConfig struct {
	QueueURL          string `yaml:"queue_url"`
	Endpoint          string `yaml:"endpoint_url"`
	Region            string `yaml:"region"`
	WaitTimeSeconds   int32  `yaml:"wait_time_seconds"`
	VisibilityTimeout int32  `yaml:"visibility_timeout"`
}

type ServerConfig struct {
	Environment string    `yaml:"environment"`
	SQS         SQSConfig `yaml:"sqs"`
}

func MustLoadFromYAML() ServerConfig {
	var mapConfig map[string]ServerConfig

	fd, err := os.Open("config/config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	defer func(r io.Closer) {
		if err := r.Close(); err != nil {
			slog.Warn("failed to close file descriptor: %v", err)
		}
	}(fd)

	if err := yaml.NewDecoder(fd).Decode(&mapConfig); err != nil {
		log.Fatal(err)
	}

	env := os.Getenv("ENV")
	if env != "" {
		return mapConfig[env]

	}

	slog.Info("falling back to default environment")
	return mapConfig["default"]
}

func AppVersion() string {
	version := os.Getenv("VERSION")
	if version != "" {
		return version
	}
	return "0.1.0"
}
