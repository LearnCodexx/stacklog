package stacklog

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// BasicPrint handles basic infrastructure logging with simple tag-based categorization.
// This logger is thread-safe and formats messages with timestamp, level, and tag information.
//
// BasicPrint is used for infrastructure, configuration, startup, and general application
// logging that doesn't require request context awareness. For API-related logging
// that should group with HTTP requests, use APIPrint instead.
//
// This logger is used internally by the global logging functions like Info(), Error(),
// Startup(), ConfigError(), etc.
type BasicPrint struct {
	fixedLength int
	mu          sync.Mutex
}

// NewBasicPrint creates a new BasicPrint logger instance with default settings.
//
// Note: This function is primarily used internally by the global logging system.
// For application code, use logging.Init() and the global logging functions instead.
//
// Example (internal usage):
//
//	logger := NewBasicPrint()
//	logger.Info("CONFIG", "Configuration loaded successfully")
func NewBasicPrint() *BasicPrint {
	return &BasicPrint{fixedLength: 6}
}

// Info logs an informational message with the specified tag.
// The message is formatted with timestamp, INFO level, tag, and message content.
//
// Note: For application code, use the global logging functions like logging.Info(),
// logging.Startup(), logging.Config(), etc. instead of calling this directly.
//
// Example (internal usage):
//
//	logger.Info("DATABASE", "Connection established to %s", host)
func (l *BasicPrint) Info(tag, format string, a ...any) {
	l.printLog(LevelInfo, tag, format, nil, a...)
}

// Error logs an error message with the specified tag and error.
// The message is formatted with timestamp, ERROR level, tag, message, and error details.
//
// Note: For application code, use the global logging functions like logging.Error(),
// logging.ConfigError(), logging.DBError(), etc. instead of calling this directly.
//
// Example (internal usage):
//
//	logger.Error("DATABASE", "Connection failed to %s", err, host)
func (l *BasicPrint) Error(tag, format string, err error, a ...any) {
	l.printLog(LevelError, tag, format, err, a...)
}

func (l *BasicPrint) printLog(level, tag, format string, err error, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	userMsg := format
	if strings.Contains(format, "%") && len(a) > 0 {
		userMsg = fmt.Sprintf(format, a...)
	}

	typeInfo := CheckType(a...)

	message := userMsg
	if err != nil {
		message = fmt.Sprintf("%s -> %v", userMsg, err)
	}

	finalMessage := message + typeInfo

	levelDisplay := level
	if level == LevelError {
		levelDisplay = "\033[31m" + level + "\033[0m"
	}

	length := strconv.Itoa(l.fixedLength)
	fmt.Printf("[%s] [%-5s] [%-"+length+"s] %s\n", timestamp, levelDisplay, tag, finalMessage)
}
