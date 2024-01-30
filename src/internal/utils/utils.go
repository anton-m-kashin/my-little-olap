package utils

import (
	"log"
	"os"
)

type Logger struct {
	Error *log.Logger
	Info  *log.Logger
}

func NewLogger() Logger {
	return Logger{
		Error: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		Info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
	}
}
