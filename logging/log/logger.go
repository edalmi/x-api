package log

import (
	"fmt"
	"log"
)

func New(l *log.Logger) *Logger {
	return &Logger{
		logger: l,
	}
}

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Debug(v ...interface{}) {
	l.logger.Println("DEBUG", fmt.Sprint(v...))
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	l.logger.Println("DEBUG", fmt.Sprintf(f, v...))
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.Println("INFO", fmt.Sprint(v...))
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.logger.Println("INFO", fmt.Sprintf(f, v...))
}

func (l *Logger) Warn(v ...interface{}) {
	l.logger.Println("WARN", fmt.Sprint(v...))
}

func (l *Logger) Warnf(f string, v ...interface{}) {
	l.logger.Println("WARN", fmt.Sprintf(f, v...))
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.Println("ERROR", fmt.Sprint(v...))
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.logger.Println("ERROR", fmt.Sprintf(f, v...))
}
