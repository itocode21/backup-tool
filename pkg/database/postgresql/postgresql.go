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
