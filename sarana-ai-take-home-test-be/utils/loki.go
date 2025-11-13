package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/grafana/loki-client-go/loki"
	"github.com/prometheus/common/model"
)

var LokiClient *loki.Client

// InitLoki initializes the Grafana Loki client
func InitLoki() error {
	lokiHost := os.Getenv("LOKI_HOST")
	if lokiHost == "" {
		lokiHost = "http://localhost:3100"
	}

	cfg, err := loki.NewDefaultConfig(lokiHost + "/loki/api/v1/push")
	if err != nil {
		return fmt.Errorf("failed to create Loki config: %w", err)
	}

	// Configure the Loki client
	cfg.TenantID = "" // Set if using multi-tenancy
	cfg.BatchWait = 1 * time.Second
	cfg.BatchSize = 100 * 1024 // 100KB

	client, err := loki.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create Loki client: %w", err)
	}

	LokiClient = client
	log.Println("Loki client initialized successfully")
	return nil
}

// SendToLoki sends a log entry to Grafana Loki
func SendToLoki(labels map[string]string, message string) error {
	if LokiClient == nil {
		return fmt.Errorf("Loki client not initialized")
	}

	// Convert labels to model.LabelSet
	labelSet := make(model.LabelSet)
	for k, v := range labels {
		labelSet[model.LabelName(k)] = model.LabelValue(v)
	}

	// Create log entry
	err := LokiClient.Handle(labelSet, time.Now(), message)
	if err != nil {
		return fmt.Errorf("failed to send log to Loki: %w", err)
	}

	return nil
}

// StopLoki gracefully stops the Loki client
func StopLoki() {
	if LokiClient != nil {
		LokiClient.Stop()
		log.Println("Loki client stopped")
	}
}

// SendLogToLoki is a helper function to send structured logs
func SendLogToLoki(ctx context.Context, level, method, endpoint, message string, statusCode int) error {
	labels := map[string]string{
		"job":         "notes-api",
		"level":       level,
		"method":      method,
		"endpoint":    endpoint,
		"status_code": fmt.Sprintf("%d", statusCode),
	}

	logMessage := fmt.Sprintf("[%s] %s %s - Status: %d - %s",
		level, method, endpoint, statusCode, message)

	return SendToLoki(labels, logMessage)
}
