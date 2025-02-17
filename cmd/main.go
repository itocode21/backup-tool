package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/config"
	"github.com/itocode21/backup-tool/pkg/database"
	"github.com/itocode21/backup-tool/pkg/logging"
)

func main() {
	configPath := flag.String("config", "", "Path to the configuration file (required)")
	dbType := flag.String("type", "", "Database type (mysql|postgresql|mongodb) (required)")
	command := flag.String("command", "", "Command to execute (backup|restore) (required)")
	backupFile := flag.String("backup-file", "", "Path to the backup file (optional for restore/backup)")
	flag.Parse()

	if *configPath == "" || *dbType == "" || *command == "" {
		log.Fatal("Missing required parameters: --config, --type, and --command are required.")
	}

	fullConfigPath := *configPath
	if !filepath.IsAbs(*configPath) {	
		workingDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		fullConfigPath = filepath.Join(workingDir, *configPath)
	}

	cfg, err := config.LoadConfig(fullConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logging.NewLogger(cfg)

	backup, err := database.NewBackup(*dbType, logger)
	if err != nil {
		log.Fatalf("Failed to create backup instance: %v", err)
	}

	switch *command {
	case "backup":
		config := map[string]string{
			"host":        cfg.Database.Host,
			"port":        fmt.Sprintf("%d", cfg.Database.Port),
			"username":    cfg.Database.Username,
			"password":    cfg.Database.Password,
			"dbname":      cfg.Database.DBName,
			"backup-file": *backupFile,
		}
		if err := backup.PerformFullBackup(config); err != nil {
			log.Fatalf("Backup failed: %v", err)
		}
		fmt.Println("Backup completed successfully.")
	case "restore":
		config := map[string]string{
			"host":        cfg.Database.Host,
			"port":        fmt.Sprintf("%d", cfg.Database.Port),
			"username":    cfg.Database.Username,
			"password":    cfg.Database.Password,
			"dbname":      cfg.Database.DBName,
			"backup-file": *backupFile,
		}
		if err := backup.RestoreBackup(config); err != nil {
			log.Fatalf("Restore failed: %v", err)
		}
		fmt.Println("Restore completed successfully.")
	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}
