package logger

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

func (l *Logger) log(level Level, name string, data any) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	entry := struct {
		Timestamp string `json:"timestamp"`
		Level     string `json:"level"`
		Name      string `json:"name"`
		Data      any    `json:"data,omitempty"`
	}{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     level.String(),
		Name:      name,
		Data:      data,
	}

	if err, ok := data.(error); ok {
		entry.Data = err.Error()
	}

	line, err := json.Marshal(entry)
	if err != nil {
		line = []byte("ERROR: Failed to marshal log entry: " + err.Error())
	}

	l.out.Write(append(line, '\n'))
}

func (l *Logger) Debug(name string, data any) {
	l.log(LevelDebug, name, data)
}

func (l *Logger) Info(name string, data any) {
	l.log(LevelInfo, name, data)
}

func (l *Logger) Warn(name string, data any) {
	l.log(LevelWarn, name, data)
}

func (l *Logger) Error(name string, data any) {
	l.log(LevelError, name, data)
}
