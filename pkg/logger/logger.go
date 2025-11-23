package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func Init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile)
}

func Info(msg string, args ...any) {
	if len(args) > 0 {
		infoLogger.Printf(msg, args...)
	} else {
		infoLogger.Println(msg)
	}
}

func Error(msg string, args ...any) {
	if len(args) > 0 {
		errorLogger.Printf(msg, args...)
	} else {
		errorLogger.Println(msg)
	}
}
