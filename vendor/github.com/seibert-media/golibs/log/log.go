package log

import (
	"context"
	"os"

	"github.com/blendle/zapdriver"
	"github.com/getsentry/raven-go"
	"github.com/tchap/zapext/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger implements Context
type Logger struct {
	*zap.Logger
	Sentry *raven.Client
	Level  zap.AtomicLevel

	dsn   string
	nop   bool
	local bool
}

// CtxLoggerKey defines the key under which the logger is being stored
type CtxLoggerKey string

// DefaultCtxLoggerKey defines the key under which the logger is being stored
var DefaultCtxLoggerKey = CtxLoggerKey("sm-ctx-logger")

// WithLogger returns context containing Logger
func WithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, DefaultCtxLoggerKey, l)
}

// From retrieves the logger stored in context if existing or returns NopLogger otherwise
func From(ctx context.Context) *Logger {
	l, ok := ctx.Value(DefaultCtxLoggerKey).(*Logger)
	if !ok {
		return NewNop()
	}
	return l
}

// WithFields adds all passed in zap fields to the Logger stored in ctx and overwrites it for further use
func WithFields(ctx context.Context, fields ...zapcore.Field) context.Context {
	l := From(ctx).WithFields(fields...)
	return WithLogger(ctx, l)
}

// SetLevel of the logger stored in ctx
func SetLevel(ctx context.Context, to zapcore.Level) {
	From(ctx).Level.SetLevel(to)
}

// WithFieldsOverwrite adds all passed in zap fields to the Logger stored in ctx and overwrites it for further use
// WARNING: This might kill thread safety - Experimental and bad practice - DO NOT USE!
func WithFieldsOverwrite(ctx context.Context, fields ...zapcore.Field) *Logger {
	l := From(ctx)
	n := l.WithFields(fields...)
	*l = *n
	return l
}

// New Logger instance with an optional sentry key.
// If no sentry dsn is provided, the sentry encoding is disabled
// If local is true, logs will be provided in a human readable format, false will print stackdriver conformant logs as json
func New(dsn string, local bool) (*Logger, error) {
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	return NewWithLevel(dsn, local, level)
}

// NewWithLevel builds a Logger instance with an optional sentry key and the predefined level.
// If no sentry dsn is provided, the sentry encoding is disabled
// If local is true, logs will be provided in a human readable format, false will print stackdriver conformant logs as json
func NewWithLevel(dsn string, local bool, level zap.AtomicLevel) (*Logger, error) {
	var (
		cores  []zapcore.Core
		sentry *raven.Client
		err    error
	)

	if len(dsn) > 0 {
		sentry, err = raven.New(dsn)
		if err != nil {
			return nil, err
		}
		cores = append(cores, zapsentry.NewCore(zapcore.ErrorLevel, sentry))
	}

	if local {
		cores = append(cores, buildConsoleLogger(level))
	} else {
		stackdriver := zapdriver.NewProductionConfig()
		stackdriver.Level = level

		l, err := stackdriver.Build()
		if err != nil {
			return nil, err
		}
		cores = append(cores, l.Core())
	}

	logger := zap.New(zapcore.NewTee(cores...)).WithOptions(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	return &Logger{
		Logger: logger,
		Sentry: sentry,
		Level:  level,

		dsn:   dsn,
		nop:   false,
		local: local,
	}, nil
}

// NewNop returns Logger with empty logging, tracing and ErrorReporting
func NewNop() *Logger {
	sentry, _ := raven.New("")
	logger := zap.NewNop()

	log := &Logger{
		Logger: logger,
		Sentry: sentry,
		nop:    true,
	}

	return log
}

// WithFields wrapper around zap.With
func (l *Logger) WithFields(fields ...zapcore.Field) *Logger {
	if l.nop {
		return l
	}
	log, err := NewWithLevel(l.dsn, l.local, l.Level)
	if err != nil {
		l.Error("creating new logger", zap.Error(err))
		return l
	}

	if l.Sentry != nil {
		log.Sentry.SetRelease(l.Sentry.Release())
	}
	log.Logger = l.Logger.With(fields...)
	return log
}

// IsNop returns the nop status of Logger (mainly for testing)
func (l *Logger) IsNop() bool {
	return l.nop
}

// To stores the current logger in the passed in context
func (l *Logger) To(ctx context.Context) context.Context {
	return WithLogger(ctx, l)
}

// WithRelease returns a new logger updating the internal sentry client with release info
// This should be the first change to the logger (before adding fields) as otherwise the change
// might not be persisted
func (l *Logger) WithRelease(info string) *Logger {
	if l.nop {
		return l
	}
	l.Sentry.SetRelease(info)
	return l
}

// SetLevel of the underlying zap.Logger
func (l *Logger) SetLevel(to zapcore.Level) {
	l.Level.SetLevel(to)
}

func buildConsoleLogger(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.Lock(os.Stdout)

	config := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewConsoleEncoder(config)

	return zapcore.NewCore(encoder, stdout, level)
}
