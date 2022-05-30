package ihttp

import (
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
)

type logger struct {
	zerolog.Logger
}

var _ resty.Logger = (*logger)(nil)

func (l logger) Errorf(format string, v ...any) {
	l.Logger.Error().CallerSkipFrame(2).Msgf("[http] "+format, v...)
}

func (l logger) Warnf(format string, v ...any) {
	l.Logger.Warn().CallerSkipFrame(2).Msgf("[http] "+format, v...)
}

func (l logger) Debugf(format string, v ...any) {
	l.Logger.Debug().CallerSkipFrame(5).Msgf("[http] "+format, v...)
}
