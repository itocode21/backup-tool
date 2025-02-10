package mongodb

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/itocode21/backup-tool/pkg/logging"
)

type MongoDBBackup struct {
	Logger *logging.Logger
}

func (m *MongoDBBackup) PerformFullBackup(config map[string]string) error {
	m.Logger.Info("Starting full MongoDB backup...")

	// Проверка обязательных параметров
	requiredParams := []string{"host", "port", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	// Установка пути к файлу резервной копии
	backupFilePath := config["backup-file"]
	backupDir := filepath.Dir(backupFilePath)

	// Создание директории для резервной копии
	err := os.MkdirAll(backupDir, os.ModePerm)
	if err != nil {
		m.Logger.Error("Failed to create backup directory: " + err.Error())
		return err
	}

	// Подготовка аргументов для mongodump
	args := []string{
		"--host", config["host"],
		"--port", config["port"],
		"--db", config["dbname"],
		"--out", backupDir, // Используем директорию
	}
	if config["username"] != "" && config["password"] != "" {
		args = append(args, "--username", config["username"], "--password", config["password"])
		if config["auth-db"] != "" {
			args = append(args, "--authenticationDatabase", config["auth-db"])
		} else {
			args = append(args, "--authenticationDatabase", "admin")
		}
	}

	// Выполнение команды mongodump
	cmd := exec.Command("mongodump", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	m.Logger.Debug("Executing mongodump command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MongoDB backup failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MongoDB backup completed successfully. Files saved to: " + backupDir)
	return nil
}

func (m *MongoDBBackup) RestoreBackup(config map[string]string) error {
	fmt.Println("debug: start mongodb restore...")
	m.Logger.Info("Starting MongoDB restore...")
	requiredParams := []string{"host", "port", "dbname", "backup-path"}
	for _, param := range requiredParams {
		if config[param] == "" {
			m.Logger.Error("Missing required parameter: " + param)
			return errors.New("missing required parameter: " + param)
		}
	}

	args := []string{
		"--host", config["host"],
		"--port", config["port"],
		filepath.Join(config["backup-path"], config["dbname"]),
	}

	if config["username"] != "" && config["password"] != "" {
		args = append(args, "--username", config["username"], "--password", config["password"])
	}

	m.Logger.Debug("Executing mongorestore command with arguments: " + strings.Join(args, " "))
	cmd := exec.Command("mongorestore", args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		m.Logger.Error("MongoDB restore failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MongoDB restore completed successfully.")
	return nil
}
