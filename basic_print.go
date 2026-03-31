package stacklog

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type BasicPrint struct {
	fixedLength int
	mu          sync.Mutex
}

func NewBasicPrint() *BasicPrint {
	return &BasicPrint{fixedLength: 6}
}

func (l *BasicPrint) Info(tag, format string, a ...any) {
	l.printLog(LevelInfo, tag, format, nil, a...)
}

func (l *BasicPrint) Error(tag, format string, err error, a ...any) {
	l.printLog(LevelError, tag, format, err, a...)
}

func (l *BasicPrint) printLog(level, tag, format string, err error, a ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")

	userMsg := format
	if strings.Contains(format, "%") && len(a) > 0 {
		userMsg = fmt.Sprintf(format, a...)
	}

	typeInfo := CheckType(a...)

	message := userMsg
	if err != nil {
		message = fmt.Sprintf("%s -> %v", userMsg, err)
	}

	finalMessage := message + typeInfo

	levelDisplay := level
	if level == LevelError {
		levelDisplay = "\033[31m" + level + "\033[0m"
	}

	length := strconv.Itoa(l.fixedLength)
	fmt.Printf("[%s] [%-5s] [%-"+length+"s] %s\n", timestamp, levelDisplay, tag, finalMessage)
}
