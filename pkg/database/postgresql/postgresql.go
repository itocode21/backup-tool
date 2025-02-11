package postgresql

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type PostgreSQLBackup struct {
	Logger *logging.Logger
}

func (p *PostgreSQLBackup) PerformFullBackup(config map[string]string) error {
	p.Logger.Info("Starting full PostgreSQL backup...")

	requiredParams := []string{"host", "port", "username", "password", "dbname"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	defaultBackupFile := filepath.Join("backups", "postgresql", config["dbname"]+".sql")
	backupFilePath := config["backup-file"]
	if backupFilePath == "" {
		backupFilePath = defaultBackupFile
	}

	backupDir := filepath.Dir(backupFilePath)
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		p.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
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
		p.Logger.Error("PostgreSQL backup failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	p.Logger.Info("PostgreSQL backup completed successfully. File saved to: " + backupFilePath)
	return nil
}

func (p *PostgreSQLBackup) RestoreBackup(config map[string]string) error {
	p.Logger.Info("Starting PostgreSQL restore...")

	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	cmd := exec.Command("psql",
		"-U", config["username"],
		"-h", config["host"],
		"-p", config["port"],
		"-d", config["dbname"],
		"-f", config["backup-file"],
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	p.Logger.Debug("Executing psql restore command with arguments: " + strings.Join(cmd.Args, " "))
	err := cmd.Run()
	if err != nil {
		p.Logger.Error("PostgreSQL restore failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	p.Logger.Info("PostgreSQL restore completed successfully.")
	return nil
}
