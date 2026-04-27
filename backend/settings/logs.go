package settings

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var logFile *os.File
var Logger *slog.Logger

func InitLogs() {
	logPath := CLIArgs.GetLogFile()

	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v\n", err)
	}

	logFile, err := os.Create(logPath)

	if err != nil {
		log.Fatalf("Failed to create logFile: %v\n", err)
	}

	Logger = slog.New(slog.NewJSONHandler(logFile, nil))
}
