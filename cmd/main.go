package main

import (
	"fmt"
	"log"
	"os"

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

	// Конфигурация для бэкапа/восстановления
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
		backupConfig["backup-file"] = "backups/pgsql_backup.sql"
	case "mongodb":
		backupConfig["backup-path"] = "backups/mongodb_backup"
	default:
		logger.Fatal("Invalid database type: " + cfg.Database.Type)
	}

	// Восстановление данных
	if len(os.Args) > 1 && os.Args[1] == "restore" {
		err = backup.RestoreBackup(backupConfig)
		if err != nil {
			logger.Fatal("Restore failed: " + err.Error())
		}
		logger.Info("Restore completed successfully.")
		return
	}

	// Если не указано "restore", выполняем бэкап
	err = backup.PerformFullBackup(backupConfig)
	if err != nil {
		logger.Fatal("Backup failed: " + err.Error())
	}

	logger.Info("Backup completed successfully.")
}
