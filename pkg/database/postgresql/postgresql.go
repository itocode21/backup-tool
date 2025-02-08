package postgresql

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itocode21/backup-tool/pkg/logging"
	"github.com/itocode21/backup-tool/pkg/storage"
)

type PostgreSQLBackup struct {
	Logger *logging.Logger
}

func (p *PostgreSQLBackup) PerformFullBackup(config map[string]string) error {
	p.Logger.Info("Starting full PostgreSQL backup...")

	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	backupFilePath := config["backup-file"]
	backupDir := strings.TrimSuffix(backupFilePath, filepath.Base(backupFilePath))
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		p.Logger.Error("Failed to create backup directory: " + err.Error())
		return nil
	}

	cmd := exec.Command("pg_dump",
		"-U", config["username"],
		"-h", config["host"],
		"-p", config["port"],
		"-d", config["dbname"],
		"-f", backupFilePath,
	)

	os.Setenv("PGPASSWORD", config["password"])
	defer os.Unsetenv("PGPASSWORD")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	p.Logger.Debug("Executing pg_dump command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		p.Logger.Error("PostgreSQL backup failed: " + err.Error() + ". Detaild: " + stderr.String())
		return nil
	}

	p.Logger.Info("PostgreSQL backup completed successfully")
	return nil
}

func (p *PostgreSQLBackup) RestoreBackup(config map[string]string) error {
	p.Logger.Info("Starting PostgreSQL restore...")

	// Проверяем обязательные параметры
	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	// Формируем команду psql
	cmd := exec.Command("psql",
		"-U", config["username"],
		"-h", config["host"],
		"-p", config["port"],
		"-d", config["dbname"],
		"-f", config["backup-file"],
	)

	// Перенаправляем stderr для логирования ошибок
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Выполняем команду
	p.Logger.Debug("Executing psql restore command with arguments: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		p.Logger.Error("PostgreSQL restore failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	p.Logger.Info("PostgreSQL restore completed successfully.")
	return nil
}

func (m *PostgreSQLBackup) UploadBackupToStorage(storage storage.Storage, bucket, key, filePath string) error {
	m.Logger.Info("Uploading PostgreSQL backup to storage...")
	err := storage.UploadFile(bucket, key, filePath)
	if err != nil {
		m.Logger.Error("Failed to upload PostgreSQL backup: " + err.Error())
		return err
	}
	m.Logger.Info("PostgreSQL backup uploaded successfully.")
	return nil
}
