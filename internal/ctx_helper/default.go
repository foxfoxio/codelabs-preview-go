package ctx_helper

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper/constant"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
)

func NewContext(ctx context.Context) context.Context {
	if requestID := GetRequestId(ctx); requestID == "" {
		ctx = AppendRequestId(ctx, utils.NewID())
	}

	var log = GetLogger(ctx).ApplyContext(ctx)
	ctx = AppendLogger(ctx, log)
	return ctx
}

func AppendLogger(ctx context.Context, log *logger.Logger) context.Context {
	return context.WithValue(ctx, constant.ContextLogger, log)
}

func AppendRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, constant.ContextRequestId, requestId)
}

func AppendUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, constant.ContextUserId, userId)
}

func AppendSessionId(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, constant.ContextSessionId, sessionId)
}

func GetLogger(ctx context.Context) *logger.Logger {
	if log, ok := ctx.Value(constant.ContextLogger).(*logger.Logger); ok && log != nil {
		return log
	} else {
		return logger.Default()
	}
}

func GetRequestId(ctx context.Context) string {
	if token, ok := ctx.Value(constant.ContextRequestId).(string); ok {
		return token
	} else {
		return ""
	}
}

func GetUserId(ctx context.Context) string {
	if token, ok := ctx.Value(constant.ContextUserId).(string); ok {
		return token
	} else {
		return ""
	}
}

func GetSessionId(ctx context.Context) string {
	if token, ok := ctx.Value(constant.ContextSessionId).(string); ok {
		return token
	} else {
		return ""
	}
}
