package stacklog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Trace wraps err with caller file:line and preserves existing stack hints.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()

		if strings.Contains(frame.File, "runtime/") ||
			strings.Contains(frame.File, "/logging/") ||
			strings.Contains(frame.Function, ".func") {
			if !more {
				break
			}
			continue
		}

		file := frame.File
		line := frame.Line
		if fn := runtime.FuncForPC(frame.PC); fn != nil {
			if f, l := fn.FileLine(frame.PC - 1); f != "" && l != 0 {
				file, line = f, l
			}
		}

		fileInfo := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

		if strings.Contains(errStr, fileInfo) {
			if !more {
				break
			}
			continue
		}

		level := strings.Count(errStr, "↳")
		if level == 0 {
			return fmt.Errorf("\n  ↳ %s -> ERROR: %w", fileInfo, err)
		}

		indent := strings.Repeat("  ", level)
		return fmt.Errorf("\n%s↳ %s %w", indent, fileInfo, err)
	}

	return err
}

// SetError builds a new error with caller file:line.
func SetError(message string) error {
	_, file, line, _ := runtime.Caller(1)
	fileInfo := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

	return fmt.Errorf("\n  ↳ %s -> %s", fileInfo, message)
}

// GetFileName extracts the file name token from a traced error string.
func GetFileName(errStr string) string {
	start := strings.Index(errStr, "[")
	end := strings.Index(errStr, ".go")

	if start == -1 || end == -1 || end < start {
		return ""
	}

	return strings.TrimSpace(errStr[start+1 : end+3])
}
