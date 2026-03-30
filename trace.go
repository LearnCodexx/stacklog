package stacklog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Trace wraps an error with the caller file:line while preserving existing stack hints.
// It is safe to call repeatedly; it skips duplicate frames and internal runtime/stacklog frames.
func Trace(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// 0: runtime.Callers, 1: Trace, 2: caller
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()

		// Skip internal frames (runtime, stacklog) and deferred wrappers
		if strings.Contains(frame.File, "runtime/") ||
			strings.Contains(frame.File, "/stacklog/") ||
			strings.Contains(frame.Function, ".func") {
			if !more {
				break
			}
			continue
		}

		// Use PC-1 to point at the actual call site instead of the return address
		file := frame.File
		line := frame.Line
		if fn := runtime.FuncForPC(frame.PC); fn != nil {
			if f, l := fn.FileLine(frame.PC - 1); f != "" && l != 0 {
				file, line = f, l
			}
		}

		fileInfo := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

		// Avoid duplicating the same frame when bubbling the error up
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

// SetError builds an error with the caller file:line and a custom message.
func SetError(message string) error {
	_, file, line, _ := runtime.Caller(1)
	fileInfo := fmt.Sprintf("[ %s:%d ]", filepath.Base(file), line)

	return fmt.Errorf("\n  ↳ %s -> %s", fileInfo, message)
}
