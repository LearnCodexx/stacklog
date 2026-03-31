# Stacklog Package 📝

A simple and powerful logging package for Go applications with automatic HTTP request grouping and context-aware API logging.

## Features ✨

- **One-line initialization** - Setup entire logging system with single function call
- **Context-aware logging** - Request-scoped logs automatically group with HTTP responses
- **Intuitive API** - Short, clear function names that are easy to remember
- **Error stack tracing** - Automatic error path tracking with file and line information
- **HTTP request grouping** - API errors and HTTP logs appear together
- **Backwards compatible** - Maintains existing patterns while providing cleaner alternatives

## Installation 📦

```bash
go get github.com/learncodexx/stacklog
```

## Quick Start 🚀

### 1. Initialize Once in main()

```go
import "github.com/learncodexx/stacklog"

func main() {
    // One line setup for entire application
    stacklog.Init("YourServiceName")
    
    // Start using immediately
    stacklog.Startup("Application is starting up")
}
```

### 2. Basic Infrastructure Logging

```go
// Configuration
stacklog.Config("Loading application config")
if err := loadConfig(); err != nil {
    stacklog.ConfigError("Failed to load config", err)
    return
}

// Database
stacklog.Database("Connecting to PostgreSQL")
if err := connectDB(); err != nil {
    stacklog.DBError("Database connection failed", err)
    return
}
```

### 3. API Request Logging

```go
func handleUserSignup(ctx *fiber.Ctx) error {
    c := ctx.UserContext()
    
    // Log API operations with automatic request grouping
    stacklog.API(c, "Processing user signup request")
    
    user, err := createUser(c, userData)
    if err != nil {
        stacklog.APIError(c, "User creation failed", err)
        return err
    }
    
    stacklog.API(c, "User created successfully")
    return response.Success(user)
}
```

### 4. Service/Repository Pattern

```go
func (s *UserService) CreateUser(ctx context.Context, data UserData) error {
    // Get local logger for automatic error tracking
    localAPI := stacklog.Local("UserService")
    
    var err error
    defer localAPI.CheckAndLogAPI(ctx, &err, "Failed to create user")
    
    // Business logic here...
    err = s.repository.Save(ctx, user)
    return err  // Automatically logged with full stack trace if error
}
```

### 5. HTTP Middleware Setup

```go
func main() {
    stacklog.Init("MyService")
    
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.HTTPErrorHandler(stacklog.Logger()),
    })
    
    // Add automatic HTTP request logging
    app.Use(stacklog.HTTP())
    
    app.Listen(":8080")
}
```

## Complete Example 💡

```go
package main

import (
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/learncodexx/stacklog"
)

func main() {
    // Initialize logging system
    stacklog.Init("UserService")
    
    // Infrastructure logging
    stacklog.Startup("Starting User Service v1.0")
    
    // Setup database
    db, err := setupDatabase()
    if err != nil {
        stacklog.DBError("Failed to connect to database", err)
        return
    }
    stacklog.Database("Connected to PostgreSQL successfully")
    
    // Setup HTTP server
    app := fiber.New(fiber.Config{
        ErrorHandler: func(ctx *fiber.Ctx, err error) error {
            stacklog.SystemError("Unhandled error", err)
            return ctx.Status(500).JSON(fiber.Map{"error": "Internal server error"})
        },
    })
    
    // Add automatic HTTP logging with request grouping
    app.Use(stacklog.HTTP())
    
    // Setup routes
    app.Post("/users", createUserHandler)
    
    stacklog.Startup("Server listening on :8080")
    app.Listen(":8080")
}

func createUserHandler(ctx *fiber.Ctx) error {
    c := ctx.UserContext()
    
    stacklog.API(c, "Creating new user")
    
    // Simulate user creation
    localAPI := stacklog.Local("UserHandler") 
    var err error
    defer localAPI.CheckAndLogAPI(c, &err, "User creation failed")
    
    // Business logic would go here...
    // If error occurs, it will be automatically logged with stack trace
    
    stacklog.API(c, "User created successfully")
    return ctx.JSON(fiber.Map{"message": "User created"})
}
```

## API Reference 📚

### Initialization
- `Init(serviceName string)` - Initialize logging system (call once in main)

### Infrastructure Logging
- `Startup(message, args...)` - Log application startup messages
- `Config(message, args...)` - Log configuration messages  
- `Database(message, args...)` - Log database messages
- `Info(tag, message, args...)` - General info logging

### Error Logging
- `ConfigError(message, err, args...)` - Configuration errors
- `DBError(message, err, args...)` - Database errors
- `SystemError(message, err, args...)` - System/infrastructure errors
- `Error(tag, message, err, args...)` - General error logging

### API Logging (Context-Aware)
- `API(ctx, message, args...)` - API info with request grouping
- `APIError(ctx, message, err, args...)` - API errors with request grouping

### Service/Repository Helpers
- `Local(serviceName...)` - Get LocalAPILogger for services
- `localAPI.CheckAndLogAPI(ctx, &err, message)` - Auto error tracking pattern

### Direct Access
- `Logger()` - Get basic logger for special cases
- `HTTP()` - Get HTTP middleware for Fiber

### Utilities
- `CheckError(err, message, args...)` - Check and log error if not nil
- `Must(err, message)` - Panic on error (for critical initialization)

## Output Examples 📊

### Grouped Request Logs
```
▌ APPLICATION LOG
▌ [2026-03-31 14:30:25] [ERROR] [API - UserService] [ user_repository.go:123 ] User creation failed
  ↳ [ user_handler.go:45 ] 
  ↳ [ user_service.go:67 ] 
  ↳ [ user_repository.go:123 ] -> ERROR: duplicate key violates unique constraint "users_email_key"
▌ [2026-03-31 14:30:25] POST 400 42.3ms from 127.0.0.1 -> /users
```

### Infrastructure Logs
```
[2026-03-31 14:30:20] [INFO ] [STARTUP] Starting User Service v1.0
[2026-03-31 14:30:21] [INFO ] [CONFIG ] Configuration loaded successfully  
[2026-03-31 14:30:22] [INFO ] [DATABASE] Connected to PostgreSQL successfully
```

## Migration Guide 📈

### From Old Stacklog Usage:
```go
// OLD: Complex setup
print := stacklog.NewBasicPrint()
printAPI := stacklog.NewAPIPrint("")
stacklog.SetFiberErrorHook(stacklog.AddErrorToRequest)

// NEW: One line
stacklog.Init("ServiceName")
```

### From Verbose Function Names:
```go
// OLD: Verbose
stacklog.GetBasicLogger()
stacklog.GetHTTPMiddleware() 
stacklog.InfoStartup()

// NEW: Clean
stacklog.Logger()
stacklog.HTTP()
stacklog.Startup()
```

## Best Practices 💡

1. **Initialize once** in main() with your service name
2. **Use shortcuts** like `Startup()`, `ConfigError()` for common cases
3. **Use context-aware logging** (`API()`, `APIError()`) for request handling
4. **Use LocalAPILogger** with `CheckAndLogAPI()` pattern in services
5. **Consistent error handling** - always log errors at the source
6. **Meaningful messages** - include context about what operation failed

## Contributing 🤝

This package is designed to be simple and focused. When adding features:
- Keep function names short and intuitive
- Maintain backwards compatibility
- Add comprehensive documentation
- Include usage examples
- Test with real applications

---

**Happy Logging!** 🎉
defer cancel()

if err := svc.SignIn(ctx, req); err != nil {
	return stacklog.Trace(err)
}

api.Info(ctx, "signin success")
```

# HTTP Middleware Usage

```go
import (
	"github.com/gofiber/fiber/v2"
	"github.com/learncodexx/stacklog"
)

func main() {
	app := fiber.New()

	// Enable grouped request logging.
	app.Use(stacklog.HTTPLogger())

	// Let APIPrint errors join the current HTTP log group.
	stacklog.SetFiberErrorHook(stacklog.AddErrorToRequest)
}
```

## Main functions

- `Trace(err error) error`
- `SetError(message string) error`
- `ErrorPattern(err error) string`
- `TranslateError(raw string) string`
- `NewBasicPrint() *BasicPrint`
- `NewAPIPrint(service string) *APIPrint`
- `SetServiceName(ctx, service)`
- `WithTimeout(...)`, `WithDefaultTimeout(...)`
- `SetFiberErrorHook(fn)`
