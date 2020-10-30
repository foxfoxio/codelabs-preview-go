package usecases

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
)

func getSession(ctx context.Context) *entities.UserSession {
	if session, ok := ctx_helper.GetSession(ctx).(*entities.UserSession); ok {
		return session
	}

	return nil
}
