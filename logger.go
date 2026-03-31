// Package stacklog provides a simple and powerful logging system with automatic HTTP request grouping
// and context-aware API logging for Go applications.
//
// Features:
//   - One-line initialization with stacklog.Init()
//   - Context-aware logging that groups API errors with HTTP requests
//   - Intuitive function names for common use cases
//   - Automatic error stack tracing
//   - Backwards compatible with existing patterns
//
// Basic usage:
//
//	stacklog.Init("MyService")
//	stacklog.Startup("Application starting")
//	stacklog.ConfigError("Failed to load config", err)
package stacklog

import (
	"context"
	"fmt"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// Global loggers for easy access
var (
	globalBasic *BasicPrint
	globalAPI   *APIPrint
	once        sync.Once
	serviceName string
)

// Init initializes the global logging system with the specified service name.
// This should be called once at application startup before any other logging functions.
// It sets up automatic HTTP request grouping and prepares all loggers.
//
// Example:
//
//	logging.Init("UserService")
//	logging.Startup("Application is starting")
func Init(service string) {
	once.Do(func() {
		serviceName = service
		globalBasic = NewBasicPrint()
		globalAPI = NewAPIPrint(service)
		
		// Automatically setup HTTP logger integration
		SetFiberErrorHook(AddErrorToRequest)
	})
}

// ========================================
// BASIC LOGGING (untuk infrastructure, config, startup)
// ========================================

// Info logs general informational messages with a tag and message.
// Use this for infrastructure, configuration, and general application logging.
//
// Example:
//
//	logging.Info("DATABASE", "Connection pool initialized")
//	logging.Info("CONFIG", "Loaded %d configuration keys", keyCount)
func Info(tag, message string, args ...any) {
	ensureInit()
	globalBasic.Info(tag, message, args...)
}

// Error logs general error messages with a tag, message, and error.
// Use this for infrastructure, configuration, and system-level errors.
//
// Example:
//
//	logging.Error("DATABASE", "Connection failed", err)
//	logging.Error("CONFIG", "Invalid config key %s", err, keyName)
func Error(tag, message string, err error, args ...any) {
	ensureInit()
	globalBasic.Error(tag, message, err, args...)
}

// ========================================
// SHORTCUTS untuk use cases umum
// ========================================

// Startup logs application startup and initialization messages.
// This is a shortcut for Info("STARTUP", message, args...).
//
// Example:
//
//	logging.Startup("User Service v1.2.0 starting")
//	logging.Startup("Loaded %d plugins", pluginCount)
func Startup(message string, args ...any) {
	Info("STARTUP", message, args...)
}

// Config logs configuration-related informational messages.
// This is a shortcut for Info("CONFIG", message, args...).
//
// Example:
//
//	logging.Config("Configuration loaded from %s", configPath)
//	logging.Config("Using %s environment", env)
func Config(message string, args ...any) {
	Info("CONFIG", message, args...)
}

// Database logs database-related informational messages.
// This is a shortcut for Info("DATABASE", message, args...).
//
// Example:
//
//	logging.Database("Connected to PostgreSQL on %s", host)
//	logging.Database("Migration completed successfully")
func Database(message string, args ...any) {
	Info("DATABASE", message, args...)
}

// ConfigError logs configuration-related errors.
// This is a shortcut for Error("CONFIG", message, err, args...).
//
// Example:
//
//	logging.ConfigError("Failed to load config file", err)
//	logging.ConfigError("Invalid config value for %s", err, keyName)
func ConfigError(message string, err error, args ...any) {
	Error("CONFIG", message, err, args...)
}

// DBError logs database-related errors.
// This is a shortcut for Error("DATABASE", message, err, args...).
//
// Example:
//
//	logging.DBError("Connection pool exhausted", err)
//	logging.DBError("Query failed for table %s", err, tableName)
func DBError(message string, err error, args ...any) {
	Error("DATABASE", message, err, args...)
}

// SystemError logs system and infrastructure errors.
// This is a shortcut for Error("SYSTEM", message, err, args...).
//
// Example:
//
//	logging.SystemError("Failed to start HTTP server", err)
//	logging.SystemError("Memory allocation failed", err)
func SystemError(message string, err error, args ...any) {
	Error("SYSTEM", message, err, args...)
}

// ========================================
// API LOGGING (untuk request/response, business logic)
// ========================================

// API logs API-related informational messages with context awareness.
// These messages will automatically group with HTTP request logs when used
// within a request context from Fiber middleware.
//
// Example:
//
//	logging.API(ctx, "Processing user signup request")
//	logging.API(ctx, "User %s created successfully", userEmail)
func API(ctx context.Context, message string, args ...any) {
	ensureInit()
	globalAPI.Info(ctx, message, args...)
}

// APIError logs API-related errors with context awareness and automatic grouping.
// These error messages will automatically group with HTTP request logs and
// include full error stack traces when used within a request context.
//
// Example:
//
//	logging.APIError(ctx, "User creation failed", err)
//	logging.APIError(ctx, "Validation failed for field %s", err, fieldName)
func APIError(ctx context.Context, message string, err error, args ...any) {
	ensureInit()
	globalAPI.Error(ctx, message, err, args...)
}

// ========================================
// SIMPLE ACCESS (tanpa perlu getter functions)
// ========================================

// Logger returns the global basic logger for special cases where you need
// direct access to the logger instance (e.g., passing to third-party libraries).
//
// Example:
//
//	validator := NewValidator(ctx, logging.Logger())
//	server := NewServer(logging.Logger())
func Logger() *BasicPrint {
	ensureInit()
	return globalBasic
}

// Local returns a LocalAPILogger instance for use in services and repositories.
// This logger provides the CheckAndLogAPI pattern for automatic error tracking
// with stack traces. Optionally specify a service name for better log identification.
//
// Example:
//
//	localAPI := logging.Local("UserService")
//	defer localAPI.CheckAndLogAPI(ctx, &err, "Failed to create user")
func Local(serviceName ...string) *LocalAPILogger {
	name := "API"
	if len(serviceName) > 0 && serviceName[0] != "" {
		name = serviceName[0]
	}
	return GetLocalAPILogger(name)
}

// HTTP returns the configured HTTP middleware for Fiber applications.
// This middleware automatically logs HTTP requests and groups them with
// any API errors that occur during request processing.
//
// Example:
//
//	app.Use(logging.HTTP())
func HTTP() fiber.Handler {
	ensureInit()
	return HTTPLogger()
}

// ========================================
// UTILITY FUNCTIONS
// ========================================

// Must panics with a formatted message if err is not nil.
// This is useful for critical initialization failures where the application
// cannot continue running.
//
// Example:
//
//	db, err := sql.Open("postgres", dsn)
//	logging.Must(err, "Failed to connect to database")
func Must(err error, message string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %v", message, err))
	}
}

// CheckError logs the error if it is not nil and returns true if an error occurred.
// This provides a simple way to check and log errors in a single line.
//
// Example:
//
//	if logging.CheckError(err, "Database operation failed") {
//		return // or handle error
//	}
func CheckError(err error, message string, args ...any) bool {
	if err != nil {
		Error("ERROR", message, err, args...)
		return true
	}
	return false
}

// Service returns the current service name that was set during Init().
//
// Example:
//
//	serviceName := logging.Service() // returns "UserService" if Init("UserService") was called
func Service() string {
	ensureInit()
	return serviceName
}

// ensureInit ensures the logging system is initialized with defaults
// if Init() was never called explicitly
func ensureInit() {
	if globalBasic == nil || globalAPI == nil {
		Init("DefaultService")
	}
}

// ========================================
// BACKWARDS COMPATIBILITY - DEPRECATED
// These functions are kept for compatibility but prefer the newer shorter names
// ========================================

// GetBasicLogger returns basic logger. 
// Deprecated: Use Logger() instead for cleaner code.
func GetBasicLogger() *BasicPrint {
	return Logger()
}

// GetAPILogger returns API logger.
// Deprecated: Use the API() and APIError() functions directly, or access via Logger() if needed.
func GetAPILogger() *APIPrint {
	ensureInit()
	return globalAPI
}

// GetHTTPMiddleware returns HTTP middleware.
// Deprecated: Use HTTP() instead for cleaner code.
func GetHTTPMiddleware() fiber.Handler {
	return HTTP()
}

// InfoStartup logs startup message.
// Deprecated: Use Startup() instead for cleaner code.
func InfoStartup(message string, args ...any) {
	Startup(message, args...)
}

// InfoConfig logs config message.
// Deprecated: Use Config() instead for cleaner code.
func InfoConfig(message string, args ...any) {
	Config(message, args...)
}

// InfoDatabase logs database message.
// Deprecated: Use Database() instead for cleaner code.  
func InfoDatabase(message string, args ...any) {
	Database(message, args...)
}

// ErrorConfig logs config error.
// Deprecated: Use ConfigError() instead for cleaner code.
func ErrorConfig(message string, err error, args ...any) {
	ConfigError(message, err, args...)
}

// ErrorDatabase logs database error.
// Deprecated: Use DBError() instead for cleaner code.
func ErrorDatabase(message string, err error, args ...any) {
	DBError(message, err, args...)
}

// ErrorAPI logs API error.
// Deprecated: Use APIError() instead for cleaner code.
func ErrorAPI(ctx context.Context, message string, err error, args ...any) {
	APIError(ctx, message, err, args...)
}

// APIInfo logs API info.
// Deprecated: Use API() instead for cleaner code.
