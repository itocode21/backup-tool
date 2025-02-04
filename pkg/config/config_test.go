package config

import (
	"os"
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

	// Проверяем, что облачное хранилище настроено корректно
	if cfg.Storage.CloudType == "s3" && cfg.Storage.Bucket == "" {
		t.Error("Bucket is required for S3 storage")
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Устанавливаем переменные окружения
	os.Setenv("BACKUP_TOOL_DATABASE_HOST", "localhost")
	os.Setenv("BACKUP_TOOL_STORAGE_LOCAL_PATH", "/backups")
	os.Setenv("BACKUP_TOOL_LOGGING_LEVEL", "info")

	// Очищаем переменные окружения после завершения теста
	t.Cleanup(func() {
		os.Unsetenv("BACKUP_TOOL_DATABASE_HOST")
		os.Unsetenv("BACKUP_TOOL_STORAGE_LOCAL_PATH")
		os.Unsetenv("BACKUP_TOOL_LOGGING_LEVEL")
	})

	// Загружаем конфигурацию
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем значения
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

func TestLoadConfigMissingRequiredFields(t *testing.T) {
	// Указываем путь к тестовому файлу конфигурации с отсутствующими обязательными полями
	_, err := LoadConfig("test_config_missing_fields.yaml")
	if err == nil {
		t.Error("Expected error for missing required fields, got nil")
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	// Указываем путь к несуществующему файлу конфигурации
	_, err := LoadConfig("non_existent_config.yaml")
	if err == nil {
		t.Error("Expected error for missing config file, got nil")
	}
}

func TestLoadConfigInvalidValues(t *testing.T) {
	// Устанавливаем недопустимые значения в переменных окружения
	os.Setenv("BACKUP_TOOL_DATABASE_TYPE", "invalid_db_type")
	os.Setenv("BACKUP_TOOL_LOGGING_LEVEL", "invalid_level")

	// Очищаем переменные окружения после завершения теста
	t.Cleanup(func() {
		os.Unsetenv("BACKUP_TOOL_DATABASE_TYPE")
		os.Unsetenv("BACKUP_TOOL_LOGGING_LEVEL")
	})

	// Загружаем конфигурацию
	_, err := LoadConfig("test_config.yaml")
	if err == nil {
		t.Error("Expected error for invalid values, got nil")
	}
}

func TestLoadConfigPartialOverride(t *testing.T) {
	// Устанавливаем переменные окружения для частичного переопределения
	os.Setenv("BACKUP_TOOL_DATABASE_HOST", "localhost")
	os.Setenv("BACKUP_TOOL_LOGGING_LEVEL", "info")

	// Очищаем переменные окружения после завершения теста
	t.Cleanup(func() {
		os.Unsetenv("BACKUP_TOOL_DATABASE_HOST")
		os.Unsetenv("BACKUP_TOOL_LOGGING_LEVEL")
	})

	// Загружаем конфигурацию
	cfg, err := LoadConfig("test_config.yaml")
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Проверяем, что переопределенные значения изменились
	if cfg.Database.Host != "localhost" {
		t.Errorf("Expected database host 'localhost', got '%s'", cfg.Database.Host)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Expected logging level 'info', got '%s'", cfg.Logging.Level)
	}

	// Убедимся, что остальные значения остались неизменными
	if cfg.Database.Type != "mysql" {
		t.Errorf("Expected database type 'mysql', got '%s'", cfg.Database.Type)
	}
	if cfg.Storage.LocalPath != "/backups" {
		t.Errorf("Expected local path '/backups', got '%s'", cfg.Storage.LocalPath)
	}
}
