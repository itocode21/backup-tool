package mysql

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type MySQLBackup struct {
	Logger *logging.Logger
}

func (m *MySQLBackup) PerformFullBackup(config map[string]string) error {
	m.Logger.Info("Starting full MySQL backup...")

	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	cmd := exec.Command("mysqldump",
		"--user="+config["username"],
		"--password="+config["password"],
		"--host="+config["host"],
		"--port="+config["port"],
		config["dbname"],
	)

	backupFilePath := config["backup-file"]
	backupDir := filepath.Dir(backupFilePath)
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		m.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	// Создаем файл для бэкапа
	outputFile, err := os.Create(backupFilePath)
	if err != nil {
		m.Logger.Error("Failed to create backup file: " + err.Error())
		return err
	}
	defer outputFile.Close()

	cmd.Stdout = outputFile

	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MySQL backup failed: " + err.Error())
		return err
	}

	m.Logger.Info("MySQL backup completed successfully.")
	return nil
}
