package logger

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger/constant"
)

func get(ctx context.Context, key string) (value string, ok bool) {
	value, ok = ctx.Value(key).(string)
	return
}

func (g Logger) ApplyContext(ctx context.Context) *Logger {
	log := &g
	if value, ok := get(ctx, constant.ContextRequestId); ok {
		log = log.WithRequestID(value)
	}

	if value, ok := get(ctx, constant.ContextUserId); ok {
		log = log.WithUserID(value)
	}
	return log
}
