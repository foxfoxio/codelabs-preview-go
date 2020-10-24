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
	GetSession(request *http.Request) *sessions.Session
	GetContext(request *http.Request) context.Context
	GetContextAndSession(request *http.Request) (context.Context, *sessions.Session)
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

func (uc *sessionUsecase) GetSession(request *http.Request) *sessions.Session {
	_, session := uc.GetContextAndSession(request)
	return session
}

func (uc *sessionUsecase) GetContext(request *http.Request) context.Context {
	ctx, _ := uc.GetContextAndSession(request)
	return ctx
}

func (uc *sessionUsecase) GetContextAndSession(request *http.Request) (context.Context, *sessions.Session) {
	ctx := ctx_helper.AppendRequestId(context.Background(), utils.NewID())

	if request == nil {
		return ctx, nil
	}

	session, _ := uc.store.Get(request, uc.sessionName)
	if value, ok := session.Values[entities.SessionKeyUserSession].(string); ok {
		userSession := entities.UnmarshalSession(value)
		ctx = ctx_helper.AppendUserId(ctx, userSession.UserId)
		ctx = ctx_helper.AppendSessionId(ctx, userSession.Id)
	}

	return ctx, session
}
