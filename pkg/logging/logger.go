package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/config"
)

// Logger представляет структуру для управления логированием.
type Logger struct {
	infoLog        *log.Logger // Логгер для информационных сообщений
	warnLog        *log.Logger // Логгер для предупреждений
	errorLog       *log.Logger // Логгер для ошибок
	debugLog       *log.Logger // Логгер для отладочных сообщений
	fatalLog       *log.Logger // Логгер для критических ошибок
	output         io.Writer   // Выходной поток для всех логгеров
	isDebugEnabled bool        // Флаг для проверки, включен ли режим отладки
}

// NewLogger создает новый экземпляр логгера на основе конфигурации.
func NewLogger(cfg *config.Config) *Logger {
	var output io.Writer

	// Если указан файл для логов
	if cfg.Logging.File != "" {
		// Извлекаем директорию из пути файла
		logDir := filepath.Dir(cfg.Logging.File)

		// Проверяем, существует ли директория, и создаем её, если нет
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Ошибка создания директории для логов: %v", err)
		}

		// Создаем или открываем файл для записи логов
		file, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Ошибка при открытии файла логов: %v", err)
		}
		output = file
	} else {
		// Если файл не указан, выводим логи в стандартный вывод
		output = os.Stdout
	}

	// Создание нового экземпляра логгера
	logger := &Logger{
		infoLog:        log.New(output, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:        log.New(output, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog:       log.New(output, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog:       log.New(output, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLog:       log.New(output, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
		output:         output,
		isDebugEnabled: cfg.Logging.Level == "debug", // Установка флага для отладочного режима
	}

	return logger
}

// Info записывает информационное сообщение.
func (l *Logger) Info(msg string) {
	l.infoLog.Println(msg)
}

// Warn записывает предупреждение.
func (l *Logger) Warn(msg string) {
	l.warnLog.Println(msg)
}

// Error записывает ошибку.
func (l *Logger) Error(msg string) {
	l.errorLog.Println(msg)
}

// Debug записывает отладочное сообщение (если включен уровень debug).
func (l *Logger) Debug(msg string) {
	if l.isDebugEnabled { // Проверка флага isDebugEnabled
		l.debugLog.Println(msg)
	}
}

// Fatal записывает критическую ошибку и завершает работу программы.
func (l *Logger) Fatal(msg string) {
	l.fatalLog.Println(msg)
	os.Exit(1)
}

// SetOutput устанавливает новый выходной поток для всех логгеров.
func (l *Logger) SetOutput(output io.Writer) {
	l.infoLog.SetOutput(output)
	l.warnLog.SetOutput(output)
	l.errorLog.SetOutput(output)
	l.debugLog.SetOutput(output)
	l.fatalLog.SetOutput(output)
	l.output = output
}
