//go:build windows

package logger

import (
	"os"
	"path/filepath"
	"runtime/pprof"
)

// Windows fallback: keep stdout/stderr unchanged and only ensure directory exists.
func SetLogToFile(path, alias string) {
	if alias == "" {
		return
	}
	_ = os.MkdirAll(filepath.Join(path, alias), 0755)
}

// DumpPprofToFile writes heap profile data to path/alias/alias.dat on Windows.
func DumpPprofToFile(path, alias string) error {
	if alias == "" {
		return os.ErrInvalid
	}
	dir := filepath.Join(path, alias)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	datFileName := filepath.Join(dir, alias+".dat")
	datFile, err := os.Create(datFileName)
	if err != nil {
		return err
	}
	defer datFile.Close()
	return pprof.WriteHeapProfile(datFile)
}
