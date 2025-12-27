package logger

import (
	"golang.org/x/sys/unix"
	"log"
	"os"
	"path/filepath"
)

// SetLogToFile redirects stderr and stdout to log files in /data/logs/{alias}/.
// Creates separate .err and .log files using the provided alias.
// Uses O_APPEND for multi-process safety.
func SetLogToFile(alias string) {
	if alias == "" {
		panic("SetLogToFile called with empty alias")
	}

	dir := filepath.Join(PATH_TO_LOGS, alias)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	errorFileName := filepath.Join(dir, alias+".err")
	logFileName := filepath.Join(dir, alias+".log")

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY

	errorFile, err := os.OpenFile(errorFileName, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to create error file: %v", err)
	}

	logFile, err := os.OpenFile(logFileName, flags, 0644)
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	if err := unix.Dup2(int(errorFile.Fd()), int(os.Stderr.Fd())); err != nil {
		log.Fatalf("Failed to redirect stderr: %v", err)
	}
	errorFile.Close()

	if err := unix.Dup2(int(logFile.Fd()), int(os.Stdout.Fd())); err != nil {
		log.Fatalf("Failed to redirect stdout: %v", err)
	}
	logFile.Close()
}
