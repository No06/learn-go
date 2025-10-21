package logger

import (
	"log"
	"os"
)

// Logger wraps the standard library logger to allow future replacement.
type Logger struct {
	*log.Logger
}

// New creates a configured logger instance.
func New() *Logger {
	return &Logger{Logger: log.New(os.Stdout, "[learn-go] ", log.LstdFlags|log.Lshortfile)}
}
