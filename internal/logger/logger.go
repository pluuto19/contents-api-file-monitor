package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}

	return &Logger{
		infoLogger:  log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func Info(l *Logger, msg string) {
	if l == nil {
		return
	}
	l.infoLogger.Println(msg)
}

func Infof(l *Logger, format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.infoLogger.Printf(format, args...)
}

func Error(l *Logger, msg string) {
	if l == nil {
		return
	}
	l.errorLogger.Println(msg)
}

func Errorf(l *Logger, format string, args ...interface{}) {
	if l == nil {
		return
	}
	l.errorLogger.Printf(format, args...)
}

func ErrorWithErr(l *Logger, msg string, err error) {
	if l == nil {
		return
	}
	l.errorLogger.Printf("%s: %v", msg, err)
}
