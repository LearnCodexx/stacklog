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

type APIPrint struct {
	defaultService string
	fixedLength    int
}

func NewAPIPrint(defaultService string) *APIPrint {
	return &APIPrint{
		defaultService: defaultService,
		fixedLength:    6,
	}
}

// Info prints a colored, tagged info log. Pass context to override service tag via KeyAPIPrint.
func (fl *APIPrint) Info(ctx context.Context, format string, a ...any) {
	_, file, line, _ := runtime.Caller(2)
	path := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

	userMsg := fmt.Sprintf(format, a...)
	typeInfo := CheckType(a...)

	finalMessage := fmt.Sprintf("%s %s%s", path, userMsg, typeInfo)
	fl.printFromContext(ctx, LevelInfo, TagAPI, finalMessage)
}

// Error prints a colored, grouped error log and keeps the stack trace tidy.
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

		// Try to group with HTTP request if in request context
		if fiberCtx != nil {
			if AddErrorToRequestFromMiddleware != nil {
				errorLog := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, levelColor, finalTag, fullMessage)
				AddErrorToRequestFromMiddleware(fiberCtx, errorLog)
				return // Successfully grouped, don't print standalone
			}
		}
	}

	// Print standalone log
	borderColor := "\033[33m" // Yellow for warnings/info
	if level == LevelError {
		borderColor = "\033[31m" // Red border for errors
	}

	fmt.Printf("\n%s▌ APPLICATION LOG\033[0m\n", borderColor)
	fmt.Printf("%s▌\033[0m [%s] [%s] [%-"+length+"s] %s\n",
		borderColor,
		timestamp,
		levelColor,
		finalTag,
		fullMessage,
	)
}

func (fl *APIPrint) calculateTagAndLength(tag, serviceName string) (string, string) {
	if tag == TagAPI && serviceName != "" {
		finalTag := concatTagService(tag, serviceName)
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

func concatTagService(tag, service string) string {
	return tag + " - " + service
}

// AddErrorToRequestFromMiddleware references the middleware function to avoid circular imports
var AddErrorToRequestFromMiddleware func(*fiber.Ctx, string)
