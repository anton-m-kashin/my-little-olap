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

func Iterate[Item any](a []Item) func() *Item {
	i := 0
	max := len(a)
	return func() *Item {
		if i >= max {
			return nil
		}
		defer func() { i += 1 }()
		item := a[i]
		return &item
	}
}
