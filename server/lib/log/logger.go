package log

import (
	"github.com/TheZeroSlave/zapsentry"
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

func modifyToSentryLogger(log *zap.Logger, sentryDSN string) *zap.Logger {
	cfg := zapsentry.Configuration{
		Level:             zapcore.ErrorLevel, //when to send message to sentry
		EnableBreadcrumbs: true,               // enable sending breadcrumbs to Sentry
		BreadcrumbLevel:   zapcore.InfoLevel,  // at what level should we sent breadcrumbs to sentry
		Tags: map[string]string{
			"component": "system",
		},
	}
	core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromDSN(sentryDSN))

	//in case of err it will return noop core. so we can safely attach it
	if err != nil {
		log.Warn("failed to init zap", zap.Error(err))
	}

	log = zapsentry.AttachCoreToLogger(core, log)

	// to use breadcrumbs feature - create new scope explicitly
	// and attach after attaching the core
	return log.With(zapsentry.NewScope())
}

func New(isDevEnv bool, sentryDSN string) Logger {
	var l *zap.Logger
	if isDevEnv {
		l, _ = zap.NewDevelopment()
	} else {
		l, _ = zap.NewProduction()
	}
	sl := modifyToSentryLogger(l, sentryDSN)
	loggerInstance := &logger{sl.Sugar()}

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
