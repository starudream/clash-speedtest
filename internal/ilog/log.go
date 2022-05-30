package ilog

import (
	"bytes"
	"io"
	slog "log"
	"runtime/debug"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.ErrorStackMarshaler = func(err error) any { return strings.TrimSuffix(string(debug.Stack()), "\n") }
}

type levelWriter struct {
	logger zerolog.Logger
	prefix string
}

func New(logger zerolog.Logger, prefix string) *levelWriter {
	if prefix != "" {
		prefix = "[" + prefix + "] "
	}
	return &levelWriter{logger, prefix}
}

func (l levelWriter) WithLevel(level zerolog.Level) *slog.Logger {
	return slog.New(&writer{l.logger, level}, l.prefix, 0)
}

type writer struct {
	logger zerolog.Logger
	level  zerolog.Level
}

var _ io.Writer = (*writer)(nil)

func (w writer) Write(p []byte) (n int, err error) {
	w.logger.WithLevel(w.level).Msg(string(bytes.TrimSuffix(p, []byte{'\n'})))
	return len(p), nil
}

func WrapError(err error, prefixes ...string) bool {
	if err == nil {
		return true
	}
	if len(prefixes) > 0 && prefixes[0] != "" {
		log.Error().CallerSkipFrame(1).Msgf("[%s] %v", prefixes[0], err)
	} else {
		log.Error().CallerSkipFrame(1).Msgf("%v", err)
	}
	return false
}
