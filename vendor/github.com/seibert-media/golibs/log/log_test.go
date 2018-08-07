package log_test

import (
	"context"
	"errors"
	"testing"

	"github.com/seibert-media/golibs/log"
	"go.uber.org/zap"
)

func Test_NewDebug(t *testing.T) {
	logger := log.New("", true)
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_From(t *testing.T) {
	l := log.New("", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = context.Background()
	if !log.From(ctx).IsNop() {
		t.Fatal("logger should be nop")
	}
}

func Test_WithFields(t *testing.T) {
	l := log.New("", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = log.WithFields(ctx, zap.String("test-new-field", "test"))

	log.From(ctx).Debug("test", zap.String("test", "test"))
}

func Test_WithFieldsOverwrite(t *testing.T) {
	l := log.New("", true)
	ctx := context.Background()

	ctx = log.WithLogger(ctx, l)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	log.WithFieldsOverwrite(ctx, zap.String("test-new-field", "test"))

	log.From(ctx).Debug("test", zap.String("test", "test"))
}

func Test_To(t *testing.T) {
	l := log.New("", true)
	ctx := context.Background()

	ctx = l.To(ctx)
	if log.From(ctx).IsNop() {
		t.Fatal("logger should not be nop")
	}
	log.From(ctx).Debug("test", zap.String("test", "test"))

	ctx = context.Background()
	if !log.From(ctx).IsNop() {
		t.Fatal("logger should be nop")
	}
}

func Test_NewNoDebug(t *testing.T) {
	logger := log.New("", false)
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}

func Test_NewInvalidSentryURL(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("New() should have panicked")
			}
		}()
		log.New("^", true)
	}()
}

func Test_NewNop(t *testing.T) {
	logger := log.NewNop()
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger = logger.WithFields(zap.String("test", "test"), zap.Int("num", 0))
	if logger == nil {
		t.Fatal("ctx is nil")
	}
	if logger.Logger == nil {
		t.Fatal("logger is nil")
	}
	if logger.Sentry == nil {
		t.Fatal("sentry is nil")
	}
	logger.Debug("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Info("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1))
	logger.Error("test", zap.String("test", "test"), zap.Int("num", 1), zap.Error(errors.New("test")))
}
