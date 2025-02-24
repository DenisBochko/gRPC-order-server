package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

const (
	key = "logger"
	RequestID = "requestID"
)

type Logger struct {
	lg *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("unable to init logger: %w", err)
	}

	ctx = context.WithValue(ctx, key, &Logger{logger})

	return ctx, nil
}

func GetLoggerFromCtx(ctx context.Context) (*Logger) {
	return ctx.Value(key).(*Logger) // приведение типа
}

func (l *Logger) Info(ctx context.Context, massage string, fields ...zap.Field) {
	if ctx.Value(RequestID) != nil {
		fields = append(fields, zap.String(RequestID, ctx.Value(RequestID).(string)))
	}

	l.lg.Info(massage, fields...)
}