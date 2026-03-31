# Stacklog Improvements

## NEW: Instance-Based Logging (Preferred)

For better testability and concurrent usage, you can now create isolated logger instances:

```go
// Create isolated logger instance
logger := stacklog.NewStacklog("UserService")

// Use instance methods
logger.Startup("Service starting")
logger.Config("Configuration loaded")
logger.Database("Connected to PostgreSQL")
logger.API(ctx, "Processing user request")
logger.APIError(ctx, "User validation failed", err)
```

## NEW: Configurable Error Patterns

You can now customize error message translations:

```go
// Add custom error pattern
stacklog.AddErrorMapping("payment_failed", "Payment could not be processed at this time.", false)

// Or create custom registry
customRegistry := stacklog.NewErrorPatternRegistry()
customRegistry.AddMapping("custom_error", "Custom user-friendly message", false)
stacklog.SetDefaultRegistry(customRegistry)
```

## NEW: Automatic Cleanup

HTTP logger now automatically cleans up old request logs to prevent memory leaks:

```go
// Cleanup starts automatically on first HTTP request
// Manual control:
stacklog.StartCleanup()  // Start background cleanup
stacklog.StopCleanup()   // Stop on shutdown

// Monitor memory usage
count := stacklog.GetRequestLogCount()
fmt.Printf("Active request logs: %d\n", count)
```

## Backward Compatibility

All existing functions continue to work unchanged:

```go
stacklog.Init("MyService")
stacklog.Startup("App starting")
stacklog.ConfigError("Config failed", err)
stacklog.API(ctx, "Processing request")
```

## Interface Support

For dependency injection and testing:

```go
var logger stacklog.LoggerInterface = myLogger
var apiLogger stacklog.APILoggerInterface = myAPILogger
```

## Constants Update

Updated for consistency:
- `LogLevelInfo`, `LogLevelError`
- `LogTagAPI`, `LogTagConfig`, `LogTagDatabase`, `LogTagSystem`, `LogTagStartup`
- `ContextKeyService`, `ContextKeyRequestID`