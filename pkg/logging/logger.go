package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/itocode21/backup-tool/pkg/config"
)

type Logger struct {
	infoLog        *log.Logger
	warnLog        *log.Logger
	errorLog       *log.Logger
	debugLog       *log.Logger
	fatalLog       *log.Logger
	output         io.Writer
	isDebugEnabled bool
}

func NewLogger(cfg *config.Config) *Logger {
	var outputs []io.Writer

	if cfg.Logging.File != "" {
		logDir := filepath.Dir(cfg.Logging.File)
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			log.Fatalf("Ошибка создания директории для логов: %v", err)
		}
		file, err := os.OpenFile(cfg.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Ошибка при открытии файла логов: %v", err)
		}
		outputs = append(outputs, file)
	}

	outputs = append(outputs, os.Stdout)

	multiWriter := io.MultiWriter(outputs...)

	logger := &Logger{
		infoLog:        log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warnLog:        log.New(multiWriter, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLog:       log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLog:       log.New(multiWriter, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLog:       log.New(multiWriter, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
		output:         multiWriter,
		isDebugEnabled: cfg.Logging.Level == "debug",
	}
	return logger
}

func (l *Logger) Info(msg string) {
	l.infoLog.Println(msg)
}

func (l *Logger) Warn(msg string) {
	l.warnLog.Println(msg)
}

func (l *Logger) Error(msg string) {
	l.errorLog.Println(msg)
}

func (l *Logger) Debug(msg string) {
	if l.isDebugEnabled {
		l.debugLog.Println(msg)
	}
}

func (l *Logger) Fatal(msg string) {
	l.fatalLog.Println(msg)
	os.Exit(1)
}

func (l *Logger) SetOutput(output io.Writer) {
	l.infoLog.SetOutput(output)
	l.warnLog.SetOutput(output)
	l.errorLog.SetOutput(output)
	l.debugLog.SetOutput(output)
	l.fatalLog.SetOutput(output)
	l.output = output
}
