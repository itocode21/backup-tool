package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type StorageConfig struct {
	LocalPath string `mapstructure:"local_path"`
	CloudType string `mapstructure:"cloud_type"`
	Bucket    string `mapstructure:"bucket"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	File   string `mapstructure:"file"`
	Format string `mapstructure:"format"`
}

type NotificationConfig struct {
	SlackWebhookURL string `mapstructure:"slack_webhook_url"`
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
	viper.AutomaticEnv()
	viper.SetEnvPrefix("BACKUP_TOOL")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Failed to read config file: %v", err)
		return nil, err
	}

	log.Println("All settings:", viper.AllSettings())

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Failed to unmarshal config: %v", err)
		return nil, err
	}

	log.Printf("Loaded config: %+v", cfg)

	if cfg.Database.Host == "" {
		return nil, fmt.Errorf("database host is required")
	}
	if cfg.Database.DBName == "" {
		return nil, fmt.Errorf("database name is required")
	}

	validDatabaseTypes := map[string]bool{
		"mysql":      true,
		"postgresql": true,
		"mongodb":    true,
	}
	if !validDatabaseTypes[cfg.Database.Type] {
		return nil, fmt.Errorf("invalid database type: %s", cfg.Database.Type)
	}
	validLoggingLevels := map[string]bool{
		"info":  true,
		"debug": true,
		"warn":  true,
		"error": true,
	}
	if !validLoggingLevels[cfg.Logging.Level] {
		return nil, fmt.Errorf("invalid logging level: %s", cfg.Logging.Level)
	}

	if cfg.Storage.CloudType != "s3" && cfg.Storage.CloudType != "gcs" {
		return nil, fmt.Errorf("invalid cloud type: %s", cfg.Storage.CloudType)
	}

	return &cfg, nil
}
