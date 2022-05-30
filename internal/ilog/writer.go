package ilog

import (
	"io"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewConsoleWriter() io.Writer {
	return &zerolog.ConsoleWriter{
		Out:        colorable.NewColorableStdout(),
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}
}

func NewFileWriter(filename string) io.Writer {
	return &zerolog.ConsoleWriter{
		Out: &lumberjack.Logger{
			Filename:  filename,
			MaxSize:   100,
			MaxAge:    360,
			LocalTime: true,
		},
		NoColor:    true,
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
	}
}
