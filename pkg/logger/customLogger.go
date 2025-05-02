package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type CustomLogger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errLogger  *log.Logger
	mu         sync.RWMutex
}

func NewCustomLogger() (*CustomLogger, error) {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	err := os.MkdirAll("../logs", 0o755)
	if err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	fileInfo, err := os.OpenFile("../logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open info log file: %w", err)
	}
	fileWarn, err := os.OpenFile("../logs/warning.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open warning log file: %w", err)
	}
	fileErr, err := os.OpenFile("../logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to open error log file: %w", err)
	}

	Logger := &CustomLogger{
		infoLogger: log.New(fileInfo, "INFO: ", flags),
		warnLogger: log.New(fileWarn, "WARN: ", flags),
		errLogger:  log.New(fileErr, "ERROR: ", flags),
	}

	return Logger, nil
}

func (l *CustomLogger) Info(msg ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.infoLogger.Println(msg...)
}

func (l *CustomLogger) Warn(msg ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.warnLogger.Println(msg...)
}

func (l *CustomLogger) Error(msg ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.errLogger.Println(msg...)
}
