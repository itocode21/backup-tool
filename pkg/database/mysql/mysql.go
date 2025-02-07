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

	// Проверяем обязательные параметры
	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	// Создаем директорию для файла бэкапа
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

	// Формируем команду mysqldump
	cmd := exec.Command("mysqldump",
		"--user="+config["username"],
		"--password="+config["password"],
		"--host="+config["host"],
		"--port="+config["port"],
		config["dbname"],
	)

	// Перенаправляем stdout в файл
	cmd.Stdout = outputFile

	// Перенаправляем stderr для логирования ошибок
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Выполняем команду
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MySQL backup failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MySQL backup completed successfully.")
	return nil
}

func (m *MySQLBackup) RestoreBackup(config map[string]string) error {
	m.Logger.Info("Starting MySQL restore...")

	// Проверяем обязательные параметры
	requiredParams := []string{"host", "port", "username", "password", "dbname", "backup-file"}
	for _, param := range requiredParams {
		if config[param] == "" {
			return errors.New("missing required parameter: " + param)
		}
	}

	// Формируем команду mysql
	cmd := exec.Command("mysql",
		"--user="+config["username"],
		"--password="+config["password"],
		"--host="+config["host"],
		"--port="+config["port"],
		config["dbname"],
	)

	// Открываем файл бэкапа
	backupFile, err := os.Open(config["backup-file"])
	if err != nil {
		m.Logger.Error("Failed to open backup file: " + err.Error())
		return err
	}
	defer backupFile.Close()

	// Передаем содержимое файла в stdin команды mysql
	cmd.Stdin = backupFile

	// Перенаправляем stderr для логирования ошибок
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Выполняем команду
	m.Logger.Debug("Executing mysql restore command with arguments: " + strings.Join(cmd.Args, " "))
	err = cmd.Run()
	if err != nil {
		m.Logger.Error("MySQL restore failed: " + err.Error() + ". Details: " + stderr.String())
		return err
	}

	m.Logger.Info("MySQL restore completed successfully.")
	return nil
}
