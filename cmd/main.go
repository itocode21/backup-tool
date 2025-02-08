package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/config"
	"github.com/itocode21/backup-tool/pkg/database"
	"github.com/itocode21/backup-tool/pkg/logging"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("pkg/config/test_config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	logger := logging.NewLogger(cfg)

	// Создание экземпляра Backup для выбранной СУБД
	backup, err := database.NewBackup(cfg.Database.Type, logger)
	if err != nil {
		logger.Fatal("Unsupported database type: " + cfg.Database.Type)
	}

	// Создание экземпляра Storage для выбранной системы хранения
	storageSystem, err := database.NewStorage(cfg.Storage.CloudType)
	if err != nil {
		logger.Fatal("Unsupported storage type: " + cfg.Storage.CloudType)
	}

	// Конфигурация для бэкапа
	backupConfig := map[string]string{
		"host":     cfg.Database.Host,
		"port":     fmt.Sprintf("%d", cfg.Database.Port),
		"username": cfg.Database.Username,
		"password": cfg.Database.Password,
		"dbname":   cfg.Database.DBName,
	}

	// Добавляем специфичные параметры для каждой СУБД
	switch cfg.Database.Type {
	case "mysql":
		backupConfig["backup-file"] = "backups/mysql_backup.sql"
	case "postgresql":
		backupConfig["backup-file"] = "backups/postgresql_backup.sql"
	case "mongodb":
		backupConfig["backup-path"] = "backups/mongodb_backup"
	default:
		logger.Fatal("Invalid database type: " + cfg.Database.Type)
	}

	// Выполнение полного бэкапа
	err = backup.PerformFullBackup(backupConfig)
	if err != nil {
		logger.Fatal("Backup failed: " + err.Error())
	}

	// Если указано хранение в облаке, загружаем бэкап
	if cfg.Storage.CloudType == "s3" {
		filePath := ""
		switch cfg.Database.Type {
		case "mysql", "postgresql":
			filePath = backupConfig["backup-file"]
		case "mongodb":
			filePath = filepath.Join(backupConfig["backup-path"], cfg.Database.DBName)
		}

		bucket := cfg.Storage.Bucket
		key := fmt.Sprintf("%s/%s", cfg.Database.Type, filepath.Base(filePath))

		err = backup.UploadBackupToStorage(storageSystem, bucket, key, filePath)
		if err != nil {
			logger.Fatal("Failed to upload backup to storage: " + err.Error())
		}
	}

	logger.Info("Backup completed successfully.")
}
