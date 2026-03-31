package stacklog

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// LocalAPILogger provides a wrapper for service and repository logging with automatic
// error tracking and stack trace generation. It maintains compatibility with existing
// CheckAndLogAPI patterns while using the global logging system internally.
//
// This logger is designed for use in services, repositories, and business logic where
// you want automatic error path tracking and request grouping.
type LocalAPILogger struct {
	serviceName string
}

// NewLocalAPILogger creates a new LocalAPILogger instance with the specified service name.
// The service name will be used in log messages for better identification.
//
// Example:
//
//	logger := NewLocalAPILogger("UserService")
//	defer logger.CheckAndLogAPI(ctx, &err, "User creation failed")
func NewLocalAPILogger(serviceName string) *LocalAPILogger {
	ensureInit()
	return &LocalAPILogger{
		serviceName: serviceName,
	}
}

// CheckAndLogAPI provides automatic error checking and logging with stack trace generation.
// This is the primary pattern for service and repository error handling.
//
// Usage pattern:
//
//	var err error
//	defer localAPI.CheckAndLogAPI(ctx, &err, "Operation description")
//	// ... do work that may set err ...
//	return err  // Will be automatically logged if not nil
//
// The function will:
//   - Check if err is not nil
//   - Generate a stack trace showing the error path
//   - Log the error with context for request grouping
//   - Return true if an error was logged, false otherwise
//
// Example:
//
//	func (s *UserService) CreateUser(ctx context.Context, data UserData) error {
//		var err error
//		defer s.logger.CheckAndLogAPI(ctx, &err, "Failed to create user")
//
//		err = s.validateUser(data)
//		if err != nil {
//			return err  // Will be logged automatically
//		}
//
//		err = s.repository.Save(ctx, user)
//		return err  // Will be logged automatically if error
//	}
func (l *LocalAPILogger) CheckAndLogAPI(ctx context.Context, err *error, message string, args ...any) bool {
	if err != nil && *err != nil {
		_, file, line, _ := runtime.Caller(1)
		handlerPath := fmt.Sprintf("\n  ↳ [ %s:%d ]", filepath.Base(file), line)

		errStr := (*err).Error()
		var finalMessage string

		if strings.Contains(message, "%") && len(args) > 0 {
			finalMessage = fmt.Sprintf(message, args...)
		} else {
			finalMessage = message
		}

		if strings.Contains(errStr, "↳") {
			*err = fmt.Errorf("%s %w", handlerPath, *err)
		} else {
			*err = fmt.Errorf("%s -> %w", handlerPath, *err)
		}

		// Use global API logger with automatic HTTP request grouping
		globalAPI.Error(ctx, finalMessage, *err, args...)
		return true
	}
	return false
}

// Log logs an informational message using the global API logger with context awareness.
// These messages will automatically group with HTTP request logs when used within
// a request context from Fiber middleware.
//
// Example:
//
//	logger.Log(ctx, "Processing user data for ID %s", userID)
func (l *LocalAPILogger) Log(ctx context.Context, message string, args ...any) {
	globalAPI.Info(ctx, message, args...)
}

// LogError logs an error message using the global API logger with context awareness.
// Unlike CheckAndLogAPI, this logs the error immediately rather than deferring.
//
// Example:
//
//	if err := validateInput(data); err != nil {
//		logger.LogError(ctx, "Input validation failed", err)
//		return err
//	}
func (l *LocalAPILogger) LogError(ctx context.Context, message string, err error, args ...any) {
	globalAPI.Error(ctx, message, err, args...)
}

// Info logs an informational message - alias for Log() for backwards compatibility.
// Deprecated: Use Log() instead for cleaner code.
func (l *LocalAPILogger) Info(ctx context.Context, message string, args ...any) {
	l.Log(ctx, message, args...)
}

// Error logs an error message - alias for LogError() for backwards compatibility.
// Deprecated: Use LogError() instead for cleaner code.
func (l *LocalAPILogger) Error(ctx context.Context, message string, err error, args ...any) {
	l.LogError(ctx, message, err, args...)
}

// GetLocalAPILogger creates and returns a LocalAPILogger instance for use in services
// and repositories. Optionally specify a service name for better log identification.
//
// This is a convenience function equivalent to NewLocalAPILogger().
//
// Example:
//
//	// In a service constructor
//	func NewUserService(repo UserRepo) *UserService {
//		return &UserService{
//			repo:   repo,
//			logger: logging.GetLocalAPILogger("UserService"),
//		}
//	}
//
//	// Or get one inline
//	localAPI := logging.GetLocalAPILogger("OrderProcessor")
//	defer localAPI.CheckAndLogAPI(ctx, &err, "Order processing failed")
func GetLocalAPILogger(serviceName ...string) *LocalAPILogger {
	name := "API"
	if len(serviceName) > 0 && serviceName[0] != "" {
		name = serviceName[0]
	}
	return NewLocalAPILogger(name)
}
