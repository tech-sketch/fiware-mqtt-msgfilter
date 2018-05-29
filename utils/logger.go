/*
Package utils : utilities.

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

type iWriter interface {
	Printf(format string, v ...interface{})
}

/*
Logger : simple Logger
*/
type Logger struct {
	Name   string
	writer iWriter
	now    func() time.Time
}

/*
NewLogger : a factory method to create Logger
*/
func NewLogger(name string) *Logger {
	return &Logger{
		Name:   name,
		writer: log.New(os.Stdout, "", 0),
		now:    time.Now,
	}
}

/*
Debugf : output debug log
*/
func (l *Logger) Debugf(msg string, args ...interface{}) {
	l.logf("debug", msg, args...)
}

/*
Infof : output info log
*/
func (l *Logger) Infof(msg string, args ...interface{}) {
	l.logf("info ", msg, args...)
}

/*
Warnf : output info log
*/
func (l *Logger) Warnf(msg string, args ...interface{}) {
	l.logf("warn ", msg, args...)
}

/*
Errorf : output info log
*/
func (l *Logger) Errorf(msg string, args ...interface{}) {
	l.logf("error", msg, args...)
}

func (l *Logger) logf(level string, msg string, args ...interface{}) {
	baseMsg := fmt.Sprintf("[APP] %s |%s| [%s] %s", l.now().Format("2006/01/02 - 15:04:05"), level, l.Name, msg)
	l.writer.Printf(baseMsg, args...)
}
