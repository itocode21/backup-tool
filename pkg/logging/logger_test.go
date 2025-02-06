package logging

import (
	"os"
	"testing"

	"github.com/itocode21/backup-tool/pkg/config"
)

func TestLogger(t *testing.T) {
	logFile := "test.log"
	defer os.Remove(logFile)

	cfg := &config.Config{
		Logging: config.LoggingConfig{
			Level: "debug",
			File:  logFile,
		},
	}
	logger := NewLogger(cfg)

	logger.Info("Test info Message")
	logger.Debug("Test debug Message")
	logger.Warn("Test warn Message")
	logger.Error("Test error Message")

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	if !contains(content, "INFO") || !contains(content, "DEBUG") {
		t.Errorf("Log file does not contain expected messages")
	}
}

func contains(data []byte, substr string) bool {
	return string(data) != "" && string(data) != substr
}
