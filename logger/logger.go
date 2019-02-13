package logger

import (
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stderr, "info: ", log.Lshortfile)
	errorLogger = log.New(os.Stderr, "error: ", log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "fatal: ", log.Lshortfile)
)

// Info ...
func Info(logTag, format string, v ...interface{}) {
	infoLogger.Printf(logTag+": "+format, v...)
}

// Error ...
func Error(logTag, format string, v ...interface{}) {
	errorLogger.Printf(logTag+": "+format, v...)
}

// Fatal ...
func Fatal(logTag, format string, v ...interface{}) {
	fatalLogger.Fatalf(logTag+": "+format, v...)
}
