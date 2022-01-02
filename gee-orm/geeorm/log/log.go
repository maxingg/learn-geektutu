package log

import (
	"io"
	"log"
	"os"
	"sync"
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if level > InfoLevel {
		infoLog.SetOutput(io.Discard)
	}

	if level > ErrorLevel {
		errorLog.SetOutput(io.Discard)
	}
}
