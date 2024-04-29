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
	zapSyncOnce  sync.Once
	zapAppLogger AppLogger
)

type zapConcrete struct {
	l *otelzap.Logger
}

func newZap() *zapConcrete {
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

	oz := otelzap.New(zapLogger, otelzap.WithErrorStatusLevel(zapcore.ErrorLevel), otelzap.WithCallerDepth(1), otelzap.WithCaller(true))

	return &zapConcrete{
		l: oz,
	}
}

func GetZap() AppLogger {
	zapSyncOnce.Do(func() {
		zapAppLogger = newZap()
	})

	return zapAppLogger
}

func (zc *zapConcrete) buildFields(ctx context.Context, fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))

	for k, f := range fields {
		var zapField zap.Field

		switch f.fieldType {
		case fieldTypeString:
			zapField = zap.String(f.Key, f.String)
		case fieldTypeInt:
			zapField = zap.Int(f.Key, f.Int)
		case fieldTypeFloat:
			zapField = zap.Float32(f.Key, f.Float)
		case fieldTypeByte:
			zapField = zap.Binary(f.Key, f.Byte)
		case fieldTypeError:
			zapField = zap.Error(f.Error)
		case fieldTypeBool:
			zapField = zap.Bool(f.Key, f.Bool)
		}

		zapFields[k] = zapField
	}

	if traceID := zc.getTraceID(ctx); traceID != "" {
		zapFields = append(zapFields, zap.String("trace_id", traceID))
	}

	return zapFields
}

func (zc *zapConcrete) Info(ctx context.Context, msg string, fields ...Field) {
	zapFields := zc.buildFields(ctx, fields)
	zc.l.InfoContext(ctx, msg, zapFields...)
}

func (zc *zapConcrete) Error(ctx context.Context, msg string, fields ...Field) {
	zapFields := zc.buildFields(ctx, fields)
	zc.l.ErrorContext(ctx, msg, zapFields...)
}

func (zc *zapConcrete) getTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}
