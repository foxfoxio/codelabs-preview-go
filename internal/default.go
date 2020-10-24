package cp

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
)

func Log(ctx context.Context, serviceInfo string) *logger.Logger {
	return ctx_helper.GetLogger(ctx).WithServiceInfo(serviceInfo)
}
