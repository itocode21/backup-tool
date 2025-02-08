package database

import (
	"errors"

	"github.com/itocode21/backup-tool/pkg/database/mongodb"
	"github.com/itocode21/backup-tool/pkg/database/mysql"
	"github.com/itocode21/backup-tool/pkg/database/postgresql"
	"github.com/itocode21/backup-tool/pkg/logging"
	"github.com/itocode21/backup-tool/pkg/storage"
)

type Backup interface {
	PerformFullBackup(config map[string]string) error
	RestoreBackup(config map[string]string) error
	UploadBackupToStorage(storage storage.Storage, bucket, key, filepath string) error
}

// NewBackup создает экземпляр Backup для конкретной СУБД.
func NewBackup(dbType string, logger *logging.Logger) (Backup, error) {
	switch dbType {
	case "mysql":
		return &mysql.MySQLBackup{Logger: logger}, nil
	case "postgresql":
		return &postgresql.PostgreSQLBackup{Logger: logger}, nil
	case "mongodb":
		return &mongodb.MongoDBBackup{Logger: logger}, nil
	default:
		return nil, errors.New("unsupported database type: " + dbType)
	}
}

// NewStorage создает экземпляр Storage для конкретной системы хранения.
func NewStorage(storageType string) (storage.Storage, error) {
	switch storageType {
	case "s3":
		return storage.NewS3Storage(), nil
	default:
		return nil, errors.New("unsupported storage type: " + storageType)
	}
}

var ErrUnsupportedDBType = errors.New("unsupported database type")
