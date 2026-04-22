package settings

import (
	"flag"
	"fmt"
	"time"
)

// CLI Arguments
type Args struct {
	dev     *bool
	port    *int
	logFile *string
}

var CLIArgs Args

// Parse CLI Arguments
func InitArgs() {
	CLIArgs.dev = flag.Bool("dev", false, "run in development mode")
	CLIArgs.port = flag.Int("p", 8080, "port to run server on")
	defaultLogFile := fmt.Sprintf("./data/logs/logs-%s", time.Now().Format("01-01-2025"))
	CLIArgs.logFile = flag.String("f", defaultLogFile, "name of the file to store logs in")
	flag.Parse()
}

func (args *Args) GetPort() int {
	return *args.port
}

func (args *Args) GetDev() bool {
	return *args.dev
}

func (args *Args) GetLogFile() string {
	return *args.logFile
}
