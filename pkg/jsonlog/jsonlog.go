package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"

	"sync"
	"time"
)

type Level int8

const (
	LevelInfo Level = iota
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

type Middlewares struct {
	AfterPrintError func(err error)
}

type Logger interface {
	PrintInfo(message string, properties map[string]string)
	PrintError(err error, properties map[string]string)
	PrintFatal(err error, properties map[string]string)
}

type logger struct {
	out         io.Writer
	minLevel    Level
	mu          sync.Mutex
	middlewares Middlewares
}

func New(out io.Writer, minLevel Level, middlewares Middlewares) Logger {
	return &logger{
		out:         out,
		minLevel:    minLevel,
		middlewares: middlewares,
	}
}

func (l *logger) PrintInfo(message string, properties map[string]string) {
	l.print(LevelInfo, message, properties)
}

func (l *logger) PrintError(err error, properties map[string]string) {
	l.print(LevelError, err.Error(), properties)

	if l.middlewares.AfterPrintError != nil {
		l.middlewares.AfterPrintError(err)
	}
}

func (l *logger) PrintFatal(err error, properties map[string]string) {
	l.print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

// print is an internal method for writing the log entry.
func (l *logger) print(level Level, message string, properties map[string]string) (int, error) {
	if level < l.minLevel {
		return 0, nil
	}

	aux := details{
		Level:      level.String(),
		Time:       time.Now().UTC().Format(time.RFC3339),
		Message:    message,
		Properties: properties,
	}

	if level >= LevelError {
		aux.Trace = string(debug.Stack())
	}

	var line []byte

	line, err := json.Marshal(aux)
	if err != nil {
		line = []byte(LevelError.String() + ": unable to marshal log message:" + err.Error())
	}

	// Prevent concurrent writes
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.out.Write(append(line, '\n'))
}

type details struct {
	Level      string            `json:"level"`
	Time       string            `json:"time"`
	Message    string            `json:"message"`
	Properties map[string]string `json:"properties,omitempty"`
	Trace      string            `json:"trace,omitempty"`
}

// Write satisfies the io.Writer interface
func (l *logger) Write(message []byte) (n int, err error) {
	return l.print(LevelError, string(message), nil)
}
