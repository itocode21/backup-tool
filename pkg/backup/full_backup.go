package backup

import (
	"github.com/itocode21/backup-tool/pkg/database"
	"github.com/itocode21/backup-tool/pkg/logging"
)

type BackupManager struct {
	DatabaseType string
	Backup       database.Backup
	Logger       *logging.Logger
}

func NewBackupManager(dbtype string, logger *logging.Logger) (*BackupManager, error) {
	backup, err := database.NewBackup(dbtype, logger)
	if err != nil {
		return nil, err
	}

	return &BackupManager{
		DatabaseType: dbtype,
		Backup:       backup,
		Logger:       logger,
	}, nil
}

func (b *BackupManager) PerformFullBackup(config map[string]string) error {
	b.Logger.Info("Starting full backup for " + b.DatabaseType)
	return b.Backup.RestoreBackup(config)
}

func (b *BackupManager) RestoreBackup(config map[string]string) error {
	b.Logger.Info("Starting restore for " + b.DatabaseType)
	return b.Backup.RestoreBackup(config)
}
