package logger

import (
	"os"
)

type Logger interface {
	Log(message interface{}) error
}

func New() (Logger, error) {
	socketPath := os.Getenv("FLUENT_SOCKET")
	if socketPath == "" {
		return &IoJsonLogger{writer: os.Stdout}, nil
	} else {
		return newFluentLogger(socketPath)
	}
}
