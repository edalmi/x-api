package slog

import (
	"fmt"

	"github.com/edalmi/x-api/logging"
	"golang.org/x/exp/slog"
)

func New(l *slog.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

type Logger struct {
	logger *slog.Logger
}

func (l Logger) Info(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

func (l Logger) Infof(f string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(f, v...))
}

func (l Logger) Debug(v ...interface{}) {
	l.logger.Debug(fmt.Sprint(v...))
}

func (l Logger) Debugf(f string, v ...interface{}) {
	l.logger.Debug(fmt.Sprintf(f, v...))
}

func (l Logger) Warn(v ...interface{}) {
	l.logger.Warn(fmt.Sprint(v...))
}

func (l Logger) Warnf(f string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(f, v...))
}

func (l Logger) Error(v ...interface{}) {
	l.logger.Error(fmt.Sprint(v...))
}

func (l Logger) Errorf(f string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(f, v...))
}

func (l *Logger) WithFields(f logging.Fields) logging.Logger {
	return l
}
