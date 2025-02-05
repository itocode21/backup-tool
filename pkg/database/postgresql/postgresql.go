package postgresql

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type PostgreSQLBackup struct {
	Logger *logging.Logger
}

func (p *PostgreSQLBackup) PerformFullBackup(config map[string]string) error {
	// Проверяем существование директории для файла бэкапа
	backupFilePath := config["backup-file"]
	backupDir := filepath.Dir(backupFilePath)
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		p.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	// Формируем команду pg_dump
	cmd := exec.Command("pg_dump",
		"-U", config["username"],
		"-h", config["host"],
		"-p", config["port"],
		"-d", config["dbname"],
		"-f", backupFilePath,
	)

	// Выполняем команду
	err = cmd.Run()
	if err != nil {
		p.Logger.Error("PostgreSQL backup failed: " + err.Error())
		return err
	}

	p.Logger.Info("PostgreSQL backup completed successfully.")
	return nil
}
