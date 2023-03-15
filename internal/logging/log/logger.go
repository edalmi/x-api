package stdlog

import "log"

type Logger struct {
	*log.Logger
}

func (l *Logger) Debug(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Debugf(f string, v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Infof(f string, v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Warnf(f string, v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.Logger.Println(v...)
}

func (l *Logger) Errorf(f string, v ...interface{}) {
	l.Logger.Println(v...)
}
