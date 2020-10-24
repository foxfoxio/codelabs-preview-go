package usecases

import (
	"context"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"net/http"
)
import "github.com/gorilla/sessions"

type Session interface {
	GetSession(request *http.Request) *entities.UserSession
	GetContext(request *http.Request) context.Context
	GetContextAndSession(request *http.Request) (context.Context, *entities.UserSession)
}

func NewSession(store sessions.Store,
	sessionName string) Session {
	return &sessionUsecase{
		store:       store,
		sessionName: sessionName,
	}
}

type sessionUsecase struct {
	store       sessions.Store
	sessionName string
}

func (uc *sessionUsecase) GetSession(request *http.Request) *entities.UserSession {
	_, session := uc.GetContextAndSession(request)
	return session
}

func (uc *sessionUsecase) GetContext(request *http.Request) context.Context {
	ctx, _ := uc.GetContextAndSession(request)
	return ctx
}

func (uc *sessionUsecase) GetContextAndSession(request *http.Request) (context.Context, *entities.UserSession) {
	requestId := utils.NewID()
	ctx := ctx_helper.AppendRequestId(context.Background(), requestId)

	if request == nil {
		return ctx, nil
	}

	if r := request.Header.Get("x-foxfox-reqid"); r != "" {
		ctx = ctx_helper.AppendRequestId(ctx, r)
	}

	session, _ := uc.store.Get(request, uc.sessionName)
	userSession := entities.NewUserSession(session)

	ctx = ctx_helper.AppendUserId(ctx, userSession.UserId)
	ctx = ctx_helper.AppendSessionId(ctx, userSession.Id)
	ctx = ctx_helper.AppendSession(ctx, userSession)

	return ctx, userSession
}
