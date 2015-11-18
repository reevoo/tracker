package logger

import (
	"encoding/json"
	"io"
)

type IoJsonLogger struct {
	writer io.Writer
}

func (l *IoJsonLogger) Log(message interface{}) error {
	jsonBytes, _ := json.Marshal(message)
	_, error := io.WriteString(l.writer, string(jsonBytes[:])+"\n")
	return error
}
