package config

import (
	"fmt"
	"log"

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
	viper.SetConfigFile(path)
	viper.AutomaticEnv() // Чтение переменных окружения

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}

	// Вывод всех загруженных ключей и значений
	log.Println("All settings:", viper.AllSettings())

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

	return &cfg, nil
}
