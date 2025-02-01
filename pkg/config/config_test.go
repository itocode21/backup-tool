package config

import (
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Указываем путь к тестовому файлу конфигурации
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем, что поля заполнены корректно
	if cfg.Database.Type != "mysql" {
		t.Errorf("Expected database type 'mysql', got '%s'", cfg.Database.Type)
	}
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Storage.LocalPath != "/backups" {
		t.Errorf("Expected local path '/backups', got '%s'", cfg.Storage.LocalPath)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected logging level 'info', got '%s'", cfg.Logging.Level)
	}
}
