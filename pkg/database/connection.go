package database

import (
	"errors"

	"github.com/itocode21/backup-tool/pkg/database/mongodb"
	"github.com/itocode21/backup-tool/pkg/database/mysql"
	"github.com/itocode21/backup-tool/pkg/database/postgresql"
	"github.com/itocode21/backup-tool/pkg/logging"
)

type Backup interface {
	PerformFullBackup(config map[string]string) error
	RestoreBackup(config map[string]string) error
}

func NewBackup(dbType string, logger *logging.Logger) (Backup, error) {
	switch dbType {
	case "mysql":
		return &mysql.MySQLBackup{Logger: logger}, nil
	case "postgresql":
		return &postgresql.PostgreSQLBackup{Logger: logger}, nil
	case "mongodb":
		return &mongodb.MongoDBBackup{Logger: logger}, nil
	default:
		return nil, ErrUnsupportedDBType
	}
}

var ErrUnsupportedDBType = errors.New("unsupported database type")
