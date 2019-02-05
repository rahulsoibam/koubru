package logger

import (
	"io"
	"log"
)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
}

func NewLogger(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) *Logger {
	logger := Logger{}
	logger.Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Llongfile)

	logger.Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Llongfile)

	logger.Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Llongfile)

	logger.Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Llongfile)
	return &logger
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Error.Println(v...)
}

func (l *Logger) Infoln(v ...interface{}) {
	l.Info.Println(v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Info.Println(v...)
}

func (l *Logger) Traceln(v ...interface{}) {
	l.Info.Println(v...)
}
