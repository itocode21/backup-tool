package backup

type BackupManagerInterface interface {
	PerformFullBackup(config map[string]string) error
	RestoreBackup(config map[string]string) error
}
