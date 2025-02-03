package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`     // Тип базы данных (mysql, postgresql, mongodb и т.д.)
	Host     string `mapstructure:"host"`     // Хост базы данных
	Port     int    `mapstructure:"port"`     // Порт базы данных
	Username string `mapstructure:"username"` // Имя пользователя
	Password string `mapstructure:"password"` // Пароль
	DBName   string `mapstructure:"dbname"`   // Имя базы данных
}

type StorageConfig struct {
	LocalPath string `mapstructure:"local_path"` // Путь для локального хранения
	CloudType string `mapstructure:"cloud_type"` // Тип облачного хранилища (s3, gcs, azure)
	Bucket    string `mapstructure:"bucket"`     // Имя бакета в облаке
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // Уровень логирования (info, warn, error)
	File   string `mapstructure:"file"`   // Файл для логирования
	Format string `mapstructure:"format"` // Формат логов (json, text)
}

type NotificationConfig struct {
	SlackWebhookURL string `mapstructure:"slack_webhook_url"` // URL для Slack-уведомлений
}

type Config struct {
	Database     DatabaseConfig     `mapstructure:"database"`
	Storage      StorageConfig      `mapstructure:"storage"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	Notification NotificationConfig `mapstructure:"notification"`
}

func LoadConfig(path string) (*Config, error) {
	viper.Reset() // Сбросить кэш

	// Указываем путь к файлу конфигурации
	viper.SetConfigFile(path)

	// Включаем чтение переменных окружения
	viper.AutomaticEnv()

	// Устанавливаем префикс для переменных окружения (опционально)
	viper.SetEnvPrefix("BACKUP_TOOL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Чтение конфигурации из файла
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}

	// Вывод всех загруженных ключей и значений для отладки
	log.Println("All settings:", viper.AllSettings())

	// Загрузка конфигурации в структуру
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
		return nil, err
	}

	// Вывод загруженной конфигурации для отладки
	log.Printf("Loaded config: %+v", cfg)

	// Проверка обязательных полей
	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("database host is required")
	}
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("database name is required")
	}

	// Проверка допустимых значений для database.type
	validDatabaseTypes := map[string]bool{
		"mysql":    true,
		"postgres": true,
	}
	if !validDatabaseTypes[cfg.Database.Type] {
		return nil, fmt.Errorf("invalid database type: %s", cfg.Database.Type)
	}

	// Проверка допустимых значений для logging.level
	validLoggingLevels := map[string]bool{
		"info":  true,
		"debug": true,
		"warn":  true,
		"error": true,
	}
	if !validLoggingLevels[cfg.Logging.Level] {
		return nil, fmt.Errorf("invalid logging level: %s", cfg.Logging.Level)
	}

	// Проверка других полей (при необходимости)
	if cfg.Storage.CloudType != "s3" && cfg.Storage.CloudType != "gcs" {
		return nil, fmt.Errorf("invalid cloud type: %s", cfg.Storage.CloudType)
	}

	return &cfg, nil
}
