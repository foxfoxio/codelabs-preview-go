package ctx_helper

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper/constant"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"net/http"
)

func NewContext(ctx context.Context) context.Context {
	if requestID := GetRequestId(ctx); requestID == "" {
		ctx = AppendRequestId(ctx, utils.NewID())
	}

	var log = GetLogger(ctx).ApplyContext(ctx)
	ctx = AppendLogger(ctx, log)
	return ctx
}

func NewContextFromRequest(r *http.Request) context.Context {
	requestId := utils.NewID()
	if t := r.Header.Get("request-id"); t != "" {
		requestId = t
	}

	ctx := AppendRequestId(context.Background(), requestId)
	var log = GetLogger(ctx).ApplyContext(ctx)
	return AppendLogger(ctx, log)
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

func AppendSession(ctx context.Context, session interface{}) context.Context {
	return context.WithValue(ctx, constant.ContextSession, session)
}

func AppendAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, constant.ContextAuthorization, authorization)
}

func GetLogger(ctx context.Context) (l *logger.Logger) {
	defer func() {
		if l.RequestID != "" {
			return
		}

		if r := GetRequestId(ctx); r != "" {
			l = l.WithRequestID(r)
		}
	}()

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

func GetAuthorization(ctx context.Context) string {
	if token, ok := ctx.Value(constant.ContextAuthorization).(string); ok {
		return token
	} else {
		return ""
	}
}

func GetSession(ctx context.Context) interface{} {
	return ctx.Value(constant.ContextSession)
}
