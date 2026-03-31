package stacklog

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// APIPrint handles context-aware logging for API operations with automatic
// HTTP request grouping and error stack trace generation.
//
// This logger is used internally by the global logging functions API() and APIError().
// For most use cases, prefer using the global functions directly rather than
// creating APIPrint instances manually.
type APIPrint struct {
	defaultService string
	fixedLength    int
}

// NewAPIPrint creates a new APIPrint instance with the specified default service name.
// The service name is used in log messages for service identification.
//
// Note: This function is primarily used internally by the global logging system.
// For application code, use logging.Init() and the global logging functions instead.
//
// Example (internal usage):
//
//	apiLogger := NewAPIPrint("UserService")
//	apiLogger.Info(ctx, "Processing user request")
func NewAPIPrint(defaultService string) *APIPrint {
	return &APIPrint{
		defaultService: defaultService,
		fixedLength:    6,
	}
}

// Info logs an informational message with context awareness and automatic file/line tracking.
// Messages logged through this function will automatically group with HTTP request logs
// when used within a Fiber request context.
//
// Note: For application code, use the global logging.API() function instead of calling this directly.
//
// Example (internal usage):
//
//	apiLogger.Info(ctx, "User %s logged in successfully", userEmail)
func (fl *APIPrint) Info(ctx context.Context, format string, a ...any) {
	_, file, line, _ := runtime.Caller(2)
	path := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

	userMsg := fmt.Sprintf(format, a...)
	typeInfo := CheckType(a...)
	finalMessage := fmt.Sprintf("%s %s%s", path, userMsg, typeInfo)

	fl.printFromContext(ctx, LevelInfo, TagAPI, finalMessage)
}

// Error logs an error message with context awareness, automatic stack trace generation,
// and HTTP request grouping. Error messages will automatically include file/line information
// and will group with HTTP request logs when used within a Fiber request context.
//
// The function handles error stack trace parsing and formatting, ensuring that
// multi-level error paths are properly displayed with visual indicators.
//
// Note: For application code, use the global logging.APIError() function instead of calling this directly.
//
// Example (internal usage):
//
//	apiLogger.Error(ctx, "Failed to create user", err)
//	apiLogger.Error(ctx, "Validation failed for field %s", err, fieldName)
func (fl *APIPrint) Error(ctx context.Context, format string, err error, a ...any) {
	_, file, line, _ := runtime.Caller(2)
	handlerPath := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

	baseMessage := format
	if err != nil {
		baseMessage = fmt.Sprintf("%s %v", format, err)
	}

	userMessage := fmt.Sprintf(baseMessage, a...)
	typeInfo := CheckType(a...)

	rootPath := handlerPath
	rootMessage := ""
	cleanStackTrace := userMessage

	lines := strings.Split(userMessage, "↳")
	if len(lines) > 1 {
		lastIndex := len(lines) - 1
		lastLine := strings.TrimSpace(lines[lastIndex])

		start := strings.Index(lastLine, "[")
		end := strings.Index(lastLine, "]")
		if start != -1 && end != -1 {
			rootPath = "\033[1;31m" + lastLine[start:end+1] + "\033[0m"

			msgAfterPath := strings.TrimSpace(lastLine[end+1:])
			fullErrorDetail := ""

			if strings.Contains(msgAfterPath, "->") {
				parts := strings.SplitN(msgAfterPath, "->", 2)
				rootMessage = strings.TrimSpace(parts[0])

				rawDetail := strings.TrimSpace(parts[1])
				cleanDetail := strings.ReplaceAll(rawDetail, "ERROR:", "")
				cleanDetail = strings.TrimSpace(cleanDetail)
				fullErrorDetail = " -> ERROR: " + cleanDetail
			} else {
				rootMessage = msgAfterPath
			}

			lines[lastIndex] = fmt.Sprintf(" %s%s", lastLine[start:end+1], fullErrorDetail)
		}

		cleanStackTrace = strings.Join(lines, "↳")
	}

	header := rootPath
	if rootMessage != "" {
		header = fmt.Sprintf("%s %s", rootPath, rootMessage)
	}

	finalCombinedMessage := header + " " + cleanStackTrace + typeInfo
	fl.printFromContext(ctx, LevelError, TagAPI, finalCombinedMessage)
}

func (fl *APIPrint) printFromContext(ctx context.Context, level, tag, fullMessage string) {
	serviceName := fl.getServiceName(ctx)
	finalTag, length := fl.calculateTagAndLength(tag, serviceName)
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	var fiberCtx *fiber.Ctx
	if fc, ok := ctx.Value("fiber").(*fiber.Ctx); ok {
		fiberCtx = fc
	}

	levelColor := level
	if level == LevelError {
		levelColor = "\033[31m" + level + "\033[0m"

		if fiberCtx != nil && AddErrorToRequestFromMiddleware != nil {
			errorLog := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, levelColor, finalTag, fullMessage)
			AddErrorToRequestFromMiddleware(fiberCtx, errorLog)
			return
		}
	}

	borderColor := "\033[33m"
	if level == LevelError {
		borderColor = "\033[31m"
	}

	fmt.Printf("\n%s▌ APPLICATION LOG\033[0m\n", borderColor)
	fmt.Printf("%s▌\033[0m [%s] [%s] [%-"+length+"s] %s\n",
		borderColor, timestamp, levelColor, finalTag, fullMessage)
}

func (fl *APIPrint) calculateTagAndLength(tag, serviceName string) (string, string) {
	if tag == TagAPI && serviceName != "" {
		finalTag := tag + " - " + serviceName
		return finalTag, strconv.Itoa(len(finalTag))
	}

	return tag, strconv.Itoa(fl.fixedLength)
}

func (fl *APIPrint) getServiceName(ctx context.Context) string {
	if serviceName, ok := ctx.Value(KeyAPIPrint).(string); ok {
		return serviceName
	}

	return fl.defaultService
}

// AddErrorToRequestFromMiddleware holds the function used to add API errors
// to HTTP request logs for grouped output. This variable is set automatically
// during logging system initialization.
var AddErrorToRequestFromMiddleware func(*fiber.Ctx, string)

// SetFiberErrorHook configures the function used to integrate API error logging
// with HTTP request logging. This enables automatic grouping of API errors
// with their corresponding HTTP request logs.
//
// This function is called automatically during logging.Init() to set up
// the integration between APIPrint and HTTPLogger.
//
// Note: This function is used internally by the logging system during initialization.
// Application code doesn't need to call this directly.
func SetFiberErrorHook(fn func(*fiber.Ctx, string)) {
	AddErrorToRequestFromMiddleware = fn
}
