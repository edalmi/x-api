package logging

type Logger interface {
	Debug(v ...interface{})
	Debugf(f string, v ...interface{})
	Info(v ...interface{})
	Infof(f string, v ...interface{})
	Warn(v ...interface{})
	Warnf(f string, v ...interface{})
	Error(v ...interface{})
	Errorf(f string, v ...interface{})
}
