package settings

import (
	"log"
	"log/slog"
	"os"
)

var logFile *os.File
var Logger *slog.Logger

func InitLogs() {
	logFile, err := os.Create(CLIArgs.GetLogFile())
	if err != nil {
		log.Fatalf("Failed to create logFile: %v\n", err)
	}
	Logger = slog.New(slog.NewJSONHandler(logFile, nil))
}
