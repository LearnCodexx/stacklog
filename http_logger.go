package stacklog

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	requestLogs  = make(map[string]*RequestLogEntry)
	requestMutex = sync.RWMutex{}
)

// RequestLogEntry holds information about an HTTP request and its associated
// error logs for grouped logging output.
type RequestLogEntry struct {
	Timestamp    string
	Method       string
	Path         string
	IP           string
	Status       int
	Duration     string
	ErrorLogs    []string
	HasCompleted bool
}

// HTTPLogger returns a Fiber middleware that groups API error logs with HTTP request logs.
// This middleware captures request information, collects any API errors that occur during
// request processing, and outputs them together in a visually grouped format.
//
// The middleware automatically:
//   - Assigns unique request IDs for tracking
//   - Captures request details (method, path, IP, timing)
//   - Collects API error logs that occur during request processing
//   - Groups and displays all logs together when the request completes
//
// Usage:
//
//	app.Use(logging.HTTP())  // or logging.HTTPLogger() for direct access
//
// Example output:
//
//	▌ APPLICATION LOG
//	▌ [2026-03-31 14:30:25] [ERROR] [API - UserService] Failed to create user
//	  ↳ [ user_handler.go:45 ]
//	  ↳ [ user_service.go:67 ] -> ERROR: validation failed
//	▌ [2026-03-31 14:30:25] POST 400 42.3ms from 127.0.0.1 -> /users
func HTTPLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := fmt.Sprintf("%d-%s", time.Now().UnixNano(), c.IP())

		c.Locals("requestID", requestID)

		ctx := context.WithValue(context.Background(), "fiber", c)
		ctx = context.WithValue(ctx, "requestID", requestID)
		c.SetUserContext(ctx)

		start := time.Now()

		requestMutex.Lock()
		requestLogs[requestID] = &RequestLogEntry{
			Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
			Method:       c.Method(),
			Path:         c.Path(),
			IP:           c.IP(),
			ErrorLogs:    make([]string, 0),
			HasCompleted: false,
		}
		requestMutex.Unlock()

		err := c.Next()
		duration := time.Since(start)

		requestMutex.Lock()
		if entry, exists := requestLogs[requestID]; exists {
			entry.Status = c.Response().StatusCode()
			entry.Duration = formatDuration(duration)
			entry.HasCompleted = true

			printGroupedLogs(entry)
			delete(requestLogs, requestID)
		}
		requestMutex.Unlock()

		return err
	}
}

// AddErrorToRequest adds an error message to the current HTTP request's log collection.
// This function is used internally by the API logging system to group API errors
// with their corresponding HTTP request logs.
//
// The error message will be displayed along with the HTTP request information
// when the request completes, providing a grouped view of all errors that
// occurred during request processing.
//
// Note: This function is primarily used internally by the logging system.
// Application code should use logging.APIError() which automatically calls this function.
func AddErrorToRequest(c *fiber.Ctx, errorMsg string) {
	if requestID, ok := c.Locals("requestID").(string); ok {
		requestMutex.Lock()
		if entry, exists := requestLogs[requestID]; exists && !entry.HasCompleted {
			entry.ErrorLogs = append(entry.ErrorLogs, errorMsg)
		}
		requestMutex.Unlock()
	}
}

// AddErrorToRequestFromContext adds an error message to the HTTP request log collection
// using a context that contains a Fiber context. This enables error grouping even when
// the Fiber context is not directly available.
//
// Note: This function is primarily used internally by the API logging system.
func AddErrorToRequestFromContext(ctx context.Context, errorMsg string) {
	if fiberCtx, ok := ctx.Value("fiber").(*fiber.Ctx); ok {
		AddErrorToRequest(fiberCtx, errorMsg)
	}
}

func printGroupedLogs(entry *RequestLogEntry) {
	statusColor := getEnhancedStatusColor(entry.Status)
	methodColor := getEnhancedMethodColor(entry.Method)
	resetColor := "\033[0m"
	dimColor := "\033[90m"

	hasErrors := len(entry.ErrorLogs) > 0 || entry.Status >= 400
	borderColor := "\033[32m"
	if hasErrors {
		borderColor = "\033[31m"
	}

	fmt.Printf("\n%s▌ APPLICATION LOG\033[0m\n", borderColor)

	for _, errorLog := range entry.ErrorLogs {
		fmt.Printf("%s▌\033[0m %s\n", borderColor, errorLog)
	}

	fmt.Printf("%s▌\033[0m [%s] %s%s%s %s%d%s %s %sfrom%s %s%s%s -> %s%s%s\n",
		borderColor,
		entry.Timestamp,
		methodColor, entry.Method, resetColor,
		statusColor, entry.Status, resetColor,
		entry.Duration,
		dimColor, resetColor,
		dimColor, entry.IP, resetColor,
		dimColor, entry.Path, resetColor,
	)
}

func getEnhancedMethodColor(method string) string {
	switch method {
	case "GET":
		return "\033[1;34m"
	case "POST":
		return "\033[1;36m"
	case "PUT":
		return "\033[1;33m"
	case "DELETE":
		return "\033[1;31m"
	case "PATCH":
		return "\033[1;35m"
	case "HEAD":
		return "\033[1;32m"
	case "OPTIONS":
		return "\033[1;37m"
	default:
		return "\033[1;37m"
	}
}

func getEnhancedStatusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "\033[1;92m"
	case status >= 300 && status < 400:
		return "\033[1;94m"
	case status >= 400 && status < 500:
		return "\033[1;91m"
	case status >= 500:
		return "\033[1;95m"
	default:
		return "\033[1;97m"
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%.3fns", float64(d.Nanoseconds()))
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.3fus", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.3fms", float64(d.Nanoseconds())/1000000)
	}
	return fmt.Sprintf("%.3fs", d.Seconds())
}
