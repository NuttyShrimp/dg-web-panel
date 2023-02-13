package log

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	With(args ...interface{}) Logger

	// Loggers on level (fmt.Sprint)
	Debug(msg string, kvpPairs ...interface{})
	Info(msg string, kvpPairs ...interface{})
	Error(msg string, kvpPairs ...interface{})
	Fatal(msg string, kvpPairs ...interface{})
	// Loggers on level (fmt.Sprintf), easier spring formatting
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

type logger struct {
	*zap.SugaredLogger
}

func New(isDevEnv bool) Logger {
	var l *zap.Logger
	hooks := zap.Hooks(func(entry zapcore.Entry) error {
		if entry.Level == zapcore.ErrorLevel {
			defer sentry.Flush(2 * time.Second)
			sentry.CaptureMessage(fmt.Sprintf("%s, Line No: %d :: %s", entry.Caller.File, entry.Caller.Line, entry.Message))
		}
		return nil
	})
	if isDevEnv {
		l, _ = zap.NewDevelopment(hooks)
	} else {
		l, _ = zap.NewProduction(hooks)
	}
	loggerInstance := &logger{l.Sugar()}
	return loggerInstance
}

func (l *logger) With(args ...interface{}) Logger {
	if len(args) > 0 {
		return &logger{l.SugaredLogger.With(args...)}
	}
	return l
}

func (l logger) Debug(msg string, kvpPair ...interface{}) {
	l.SugaredLogger.Debugw(msg, kvpPair...)
}
func (l logger) Info(msg string, kvpPair ...interface{}) {
	l.SugaredLogger.Infow(msg, kvpPair...)
}
func (l logger) Error(msg string, kvpPair ...interface{}) {
	l.SugaredLogger.Errorw(msg, kvpPair...)
}
func (l logger) Fatal(msg string, kvpPair ...interface{}) {
	l.SugaredLogger.Fatalw(msg, kvpPair...)
}
