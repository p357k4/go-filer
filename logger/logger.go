package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

// Init initializes loggers. (You can later add log level filtering.)
func Init(level string) {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
