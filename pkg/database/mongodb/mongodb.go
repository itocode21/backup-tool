package mongodb

import (
	"os"
	"os/exec"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type MongoDBBackup struct {
	Logger *logging.Logger
}

func (m *MongoDBBackup) PerformFullBackup(config map[string]string) error {
	// Проверяем существование директории для бэкапа
	backupPath := config["backup-path"]
	err := os.MkdirAll(backupPath, os.ModePerm)
	if err != nil {
		m.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	// Формируем команду mongodump
	cmd := exec.Command("mongodump",
		"--host", config["host"],
		"--port", config["port"],
		"--username", config["username"],
		"--password", config["password"],
		"--db", config["dbname"],
		"--out", backupPath,
	)

	// Выполняем команду
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MongoDB backup failed: " + err.Error())
		return err
	}

	m.Logger.Info("MongoDB backup completed successfully.")
	return nil
}
