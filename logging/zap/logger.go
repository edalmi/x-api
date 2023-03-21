package zap

import (
	"fmt"

	"github.com/edalmi/x-api/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(l *zap.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

type Logger struct {
	logger *zap.Logger
}

func (l Logger) Info(v ...interface{}) {
	l.logger.Log(zapcore.InfoLevel, fmt.Sprint(v...))
}

func (l Logger) Infof(f string, v ...interface{}) {
	l.logger.Log(zapcore.InfoLevel, fmt.Sprintf(f, v...))
}

func (l Logger) Debug(v ...interface{}) {
	l.logger.Log(zapcore.DebugLevel, fmt.Sprint(v...))
}

func (l Logger) Debugf(f string, v ...interface{}) {
	l.logger.Log(zapcore.DebugLevel, fmt.Sprintf(f, v...))
}

func (l Logger) Warn(v ...interface{}) {
	l.logger.Log(zapcore.WarnLevel, fmt.Sprint(v...))
}

func (l Logger) Warnf(f string, v ...interface{}) {
	l.logger.Log(zapcore.WarnLevel, fmt.Sprintf(f, v...))
}

func (l Logger) Error(v ...interface{}) {
	l.logger.Log(zapcore.ErrorLevel, fmt.Sprint(v...))
}

func (l Logger) Errorf(f string, v ...interface{}) {
	l.logger.Log(zapcore.ErrorLevel, fmt.Sprintf(f, v...))
}

func (l Logger) WithFields(f logging.Fields) logging.Logger {
	var fields []zap.Field

	for k, v := range f {
		fields = append(fields, zap.Field{
			Type:   zapcore.StringType,
			Key:    k,
			String: v,
		})
	}

	return &Logger{
		logger: l.logger.With(fields...),
	}
}
