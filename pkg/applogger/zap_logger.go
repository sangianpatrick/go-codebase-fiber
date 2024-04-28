package applogger

import (
	"context"
	"log"
	"sync"

	"github.com/sangianpatrick/go-codebase-fiber/config"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zl       *ZapLogger
	syncOnce sync.Once
)

type ZapLogger struct {
	l *otelzap.Logger
}

func constructZapLogger() *ZapLogger {
	c := config.Get()

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.LevelKey = "severity"
	cfg.DisableStacktrace = true

	zapLogger, err := cfg.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("service_name", c.Service.Name),
		),
	)

	if err != nil {
		log.Println(err)
		zapLogger = &zap.Logger{}
	}

	return &ZapLogger{
		l: otelzap.New(zapLogger, otelzap.WithErrorStatusLevel(zapcore.ErrorLevel)),
	}
}

func GetZapLogger() *ZapLogger {
	syncOnce.Do(func() {
		zl = constructZapLogger()
	})

	return zl
}

func (z *ZapLogger) Error(ctx context.Context, message string, fields ...zap.Field) {
	z.setTraceID(ctx, &fields)
	z.l.Ctx(ctx).Error(message, fields...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	z.setTraceID(ctx, &fields)
	z.l.Ctx(ctx).Info(msg, fields...)
}

func (z *ZapLogger) setTraceID(ctx context.Context, f *[]zap.Field) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	traceID := span.SpanContext().TraceID().String()

	*f = append(*f, zap.String("trace_id", traceID))
}
