package log

import (
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"path/filepath"
)

var Logger *slog.Logger

const (
	LevelDebug = slog.Level(-4)
	LevelInfo  = slog.Level(0)
	LevelWarn  = slog.Level(4)
	LevelError = slog.Level(8)
)

func init() {
	handler := slog.NewTextHandler(os.Stdout,
		&slog.HandlerOptions{Level: LevelDebug})
	Logger = slog.New(handler)
}

// Get the file and line number of the caller
func FileLine() (string, string) {
	_, file, line, ok := runtime.Caller(1) // Adjusted to 2 to account for the additional stack frame
	// Strip all bit the last part of the file path
	lastPart :=filepath.Base(file)
	if !ok {
		Logger.Error("Unable to get caller information")
		return "", ""
	}
	return lastPart, strconv.Itoa(line)
}
func FileLine3() (string, string) {
	_, file, line, ok := runtime.Caller(2) // Adjusted to 2 to account for the additional stack frame
	if !ok {
		Logger.Error("Unable to get caller information")
		return "", ""
	}
	return file, strconv.Itoa(line)
}

func ErrorLine(msg string) {
	file, line := FileLine3()
	Logger.Error(msg, "file", file, "line", line)
}
