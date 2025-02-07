package mongodb

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type MongoDBBackup struct {
	Logger *logging.Logger
}

func (m *MongoDBBackup) PerformFullBackup(config map[string]string) error {
	m.Logger.Info("Starting full MongoDB backup...")

	// Проверяем обязательные параметры
	requiredParams := []string{"host", "port", "dbname", "backup-path"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	// Создаем директорию для бэкапа
	backupPath := config["backup-path"]
	err := os.MkdirAll(backupPath, os.ModePerm)
	if err != nil {
		m.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	// Формируем команду mongodump
	args := []string{
		"--host", config["host"],
		"--port", config["port"],
		"--db", config["dbname"],
		"--out", backupPath,
	}

	// Добавляем учетные данные, если они указаны
	if config["username"] != "" && config["password"] != "" {
		args = append(args, "--username", config["username"], "--password", config["password"])
		if config["auth-db"] != "" {
			args = append(args, "--authenticationDatabase", config["auth-db"])
		} else {
			args = append(args, "--authenticationDatabase", "admin")
		}
	}

	cmd := exec.Command("mongodump", args...)

	// Перенаправляем stderr для логирования ошибок
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Выполняем команду
	m.Logger.Debug("Executing mongodump command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MongoDB backup failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MongoDB backup completed successfully.")
	return nil
}
