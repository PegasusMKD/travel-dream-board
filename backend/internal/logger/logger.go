package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// ParseLevel converts a string level to LogLevel
func ParseLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn", "warning":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

// globalConfig holds the logger configuration set during initialization
type globalConfig struct {
	stdoutLevel LogLevel
	fileLevel   LogLevel
	filePath    string
	filePrefix  string
	initialized bool
	mu          sync.RWMutex
}

var globalCfg = &globalConfig{
	stdoutLevel: INFO,
	fileLevel:   DEBUG,
	filePath:    "./logs",
	filePrefix:  "depgraph",
}

type output struct {
	writer   io.Writer
	minLevel LogLevel
}

type Logger struct {
	component string
	fields    map[string]any
	outputs   []output
}

// Initialize sets up the global logger configuration
func Initialize(stdoutLevel, fileLevel, filePath, filePrefix string) error {
	globalCfg.mu.Lock()
	defer globalCfg.mu.Unlock()

	globalCfg.stdoutLevel = ParseLevel(stdoutLevel)
	globalCfg.fileLevel = ParseLevel(fileLevel)
	globalCfg.filePath = filePath
	globalCfg.filePrefix = filePrefix

	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(filePath, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	globalCfg.initialized = true
	return nil
}

// New creates a new Logger instance with the configured outputs
func New(component string) *Logger {
	globalCfg.mu.RLock()
	defer globalCfg.mu.RUnlock()

	logger := &Logger{
		component: component,
		fields:    make(map[string]any),
		outputs:   make([]output, 0),
	}

	// Always add stdout output
	logger.outputs = append(logger.outputs, output{
		writer:   os.Stdout,
		minLevel: globalCfg.stdoutLevel,
	})

	// Add file output if initialized
	if globalCfg.initialized {
		// Generate timestamped filename
		timestamp := time.Now().Format("2006-01-02_15-04-05")
		filename := fmt.Sprintf("%s_%s.log", globalCfg.filePrefix, timestamp)
		fullPath := filepath.Join(globalCfg.filePath, filename)

		file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			logger.outputs = append(logger.outputs, output{
				writer:   file,
				minLevel: globalCfg.fileLevel,
			})
		} else {
			// Log error to stdout but don't fail
			log.Printf("[ERROR] [Logger] Failed to open log file: %v\n", err)
		}
	}

	return logger
}

func (l *Logger) format(level, msg string, fields ...any) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	base := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level, l.component, msg)

	// Combine persistent fields with new fields
	allFields := make(map[string]any)
	for k, v := range l.fields {
		allFields[k] = v
	}

	// Add new fields from this log call
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprint(fields[i])
			allFields[key] = fields[i+1]
		}
	}

	// Append fields if present
	if len(allFields) > 0 {
		var attrs []string
		for k, v := range allFields {
			attrs = append(attrs, fmt.Sprintf("%s=%v", k, v))
		}
		base += " | " + strings.Join(attrs, " ")
	}

	return base
}

func (l *Logger) log(level LogLevel, levelStr, msg string, fields ...any) {
	formatted := l.format(levelStr, msg, fields...)

	for _, out := range l.outputs {
		if level >= out.minLevel {
			fmt.Fprintln(out.writer, formatted)
		}
	}
}

func (l *Logger) Debug(msg string, fields ...any) {
	l.log(DEBUG, "DEBUG", msg, fields...)
}

func (l *Logger) Info(msg string, fields ...any) {
	l.log(INFO, "INFO", msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...any) {
	l.log(WARN, "WARN", msg, fields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	l.log(ERROR, "ERROR", msg, fields...)
}

// With creates a new logger with additional persistent fields
// These fields will be included in all subsequent log calls
func (l *Logger) With(fields ...any) *Logger {
	newLogger := &Logger{
		component: l.component,
		fields:    make(map[string]any),
		outputs:   l.outputs, // Share the same outputs
	}

	// Copy existing fields
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprint(fields[i])
			newLogger.fields[key] = fields[i+1]
		}
	}

	return newLogger
}
