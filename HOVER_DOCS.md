# Logging Package - Quick Reference 📚

## Hover Documentation Testing

When you hover over any of these function calls in your IDE (VS Code, GoLand, etc.), you should see detailed documentation including:

- Function purpose and behavior
- Usage examples
- Parameter descriptions
- Context information
- Best practices

## Function Reference with Hover Tips

### 🚀 Initialization
```go
// Hover shows: one-line setup information
stacklog.Init("ServiceName")
```

### 🏗️ Infrastructure Logging (Basic)
```go
// Hover shows: when to use for startup messages
stacklog.Startup("Application starting v1.0")

// Hover shows: configuration logging purpose
stacklog.Config("Loaded configuration from file")

// Hover shows: database operation logging
stacklog.Database("Connected to PostgreSQL")

// Hover shows: general info logging with tags
stacklog.Info("CACHE", "Redis connection established")
```

### ❌ Error Logging
```go
// Hover shows: configuration error handling
stacklog.ConfigError("Failed to load config", err)

// Hover shows: database error patterns
stacklog.DBError("Connection pool exhausted", err)

// Hover shows: system/infrastructure errors
stacklog.SystemError("HTTP server failed to start", err)

// Hover shows: general error logging with tags
stacklog.Error("VALIDATION", "Input validation failed", err)
```

### 🌐 API/Request Logging (Context-Aware)
```go
// Hover shows: context-aware logging and request grouping
stacklog.API(ctx, "Processing user signup")

// Hover shows: automatic HTTP request grouping
stacklog.APIError(ctx, "User creation failed", err)
```

### 🔧 Service/Repository Pattern
```go
// Hover shows: LocalAPILogger usage for services
localAPI := stacklog.Local("UserService")

// Hover shows: defer pattern for automatic error tracking
defer localAPI.CheckAndLogAPI(ctx, &err, "Database operation failed")

// Hover shows: immediate logging without defer
localAPI.Log(ctx, "Starting data processing")
localAPI.LogError(ctx, "Validation failed", err)
```

### 🌍 HTTP Middleware
```go
// Hover shows: automatic request grouping setup
app.Use(stacklog.HTTP())

// Hover shows: direct access to basic logger
validator := NewValidator(stacklog.Logger())
```

### 🛠️ Utilities
```go
// Hover shows: simple error checking pattern
if stacklog.CheckError(err, "Operation failed") {
    return
}

// Hover shows: panic on critical errors
stacklog.Must(err, "Critical initialization failed")

// Hover shows: get current service name
serviceName := stacklog.Service()
```

## Test Your Hover Documentation

Try hovering over functions in these files:
- `/examples/hover_example.go` - Complete usage example
- Your own code after importing the package

## Expected Hover Information

Each function hover should show:
1. **Purpose**: What the function does
2. **Usage Context**: When and where to use it  
3. **Examples**: Code examples with realistic use cases
4. **Parameters**: What each parameter represents
5. **Behavior**: Special behavior (like request grouping)
6. **Notes**: Whether to use directly or prefer alternatives

## IDE Setup for Best Documentation

### VS Code
1. Install Go extension
2. Enable `go.useLanguageServer` 
3. Hover over function calls to see docs

### GoLand
1. Documentation appears automatically on hover
2. Use Ctrl+Q (Windows/Linux) or F1 (Mac) for detailed docs

### Vim/Neovim with LSP
1. Setup gopls language server
2. Use hover command for documentation

## Quality Checks ✅

Good documentation should include:
- ✅ Clear purpose statement
- ✅ Usage examples with realistic scenarios  
- ✅ When to use this vs alternatives
- ✅ Parameter explanations
- ✅ Return value descriptions
- ✅ Notes about internal vs public usage
- ✅ Links to related functions when relevant

## Package Design Philosophy

The documentation follows these principles:
1. **Beginner Friendly**: Clear explanations for newcomers
2. **Context Rich**: Shows not just what, but when and why
3. **Example Heavy**: Real usage examples in documentation
4. **Progressive Disclosure**: Basic usage first, advanced details after
5. **Consistent Style**: Same format and tone across all functions

---

**Test your hover documentation and enjoy the improved developer experience!** 🎉