package mysql

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type MySQLBackup struct {
	Logger *logging.Logger
}

func (m *MySQLBackup) PerformFullBackup(config map[string]string) error {
	m.Logger.Info("Starting full MySQL backup...")

	defaultBackupPath := filepath.Join("backups", "mysql", config["dbname"]+".sql")
	backupFilePath := config["backup-file"]
	if backupFilePath == "" {
		backupFilePath = defaultBackupPath
	}

	backupDir := filepath.Dir(backupFilePath)
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		m.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	outputFile, err := os.Create(backupFilePath)
	if err != nil {
		m.Logger.Error("Failed to create backup file: " + err.Error())
		return err
	}
	defer outputFile.Close()

	cmd := exec.Command("mysqldump",
		"--user="+config["username"],
		"--password="+config["password"],
		"--host="+config["host"],
		"--port="+config["port"],
		config["dbname"],
	)
	cmd.Stdout = outputFile
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	m.Logger.Debug("Executing mysqldump command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MySQL backup failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MySQL backup completed successfully. File saved to: " + backupFilePath)
	return nil
}

func (m *MySQLBackup) RestoreBackup(config map[string]string) error {
	m.Logger.Info("Starting MySQL restore...")

	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	cmd := exec.Command("mysql",
		"--user="+config["username"],
		"--password="+config["password"],
		"--host="+config["host"],
		"--port="+config["port"],
		config["dbname"],
	)

	backupFile, err := os.Open(config["backup-file"])
	if err != nil {
		m.Logger.Error("Failed to open backup file: " + err.Error())
		return err
	}
	defer backupFile.Close()

	cmd.Stdin = backupFile

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	m.Logger.Debug("Executing mysql restore command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MySQL restore failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MySQL restore completed successfully.")
	return nil
}
