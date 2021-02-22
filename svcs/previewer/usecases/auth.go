package usecases

import (
	"context"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/stopwatch"
	tokenUtils "github.com/foxfoxio/codelabs-preview-go/internal/token"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	"github.com/foxfoxio/codelabs-preview-go/internal/xfirebase"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"net/http"
	"strings"
	"time"
)

type Auth interface {
	Handler(h http.Handler) http.Handler
	ProcessFirebaseAuthorization(ctx context.Context, request *requests.AuthProcessFirebaseAuthorizationRequest) (*requests.AuthProcessFirebaseAuthorizationResponse, error)
}

func NewAuth(firebaseClient xfirebase.Client) Auth {
	return &authUsecase{
		firebase: firebaseClient,
	}
}

type authUsecase struct {
	firebase xfirebase.Client
}

func (uc *authUsecase) ProcessFirebaseAuthorization(ctx context.Context, request *requests.AuthProcessFirebaseAuthorizationRequest) (*requests.AuthProcessFirebaseAuthorizationResponse, error) {
	log := cp.Log(ctx, "AuthUsecase.ProcessFirebaseAuthorization")
	defer stopwatch.StartWithLogger(log).Stop()
	// verify token with firebase
	_, err := uc.firebase.VerifyIDToken(ctx, request.AuthorizationToken)

	if err != nil {
		log.WithError(err).Error("token verification failed")
		return nil, err
	}

	claim, err := tokenUtils.ExtractJwtClaims(request.AuthorizationToken)
	if err != nil {
		log.WithError(err).Error("extract claim failed")
		return nil, err
	}
	return &requests.AuthProcessFirebaseAuthorizationResponse{
		Name:      claim.Name,
		UserId:    claim.UserId,
		Email:     claim.Email,
		ExpiresAt: claim.ExpiresAt(),
	}, nil
}

func (uc *authUsecase) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			h.ServeHTTP(w, r)
			return
		}

		ctx := ctx_helper.NewContext(r.Context())
		log := cp.Log(ctx, "AuthUsecase.Middleware")
		defer stopwatch.StartWithLogger(log).Stop()
		authorizationToken := strings.ReplaceAll(r.Header.Get("authorization"), "Bearer ", "")
		if authorizationToken == "" {
			log.Info("mission authorization token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		authResponse, err := uc.ProcessFirebaseAuthorization(ctx, &requests.AuthProcessFirebaseAuthorizationRequest{AuthorizationToken: authorizationToken})
		if err != nil {
			log.WithError(err).Error("authentication failed")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userSession := &entities.UserSession{
			Id:        utils.NewID(),
			Name:      authResponse.Email,
			UserId:    authResponse.UserId,
			Email:     authResponse.Email,
			Token:     authorizationToken,
			CreatedAt: time.Now(),
		}

		ctx = ctx_helper.AppendUserId(ctx, userSession.UserId)
		ctx = ctx_helper.AppendSessionId(ctx, userSession.Id)
		ctx = ctx_helper.AppendSession(ctx, userSession)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
