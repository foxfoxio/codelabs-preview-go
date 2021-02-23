package usecases

import (
	"context"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
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
	AccessTokenMiddleware(next http.Handler) http.Handler
	ApiKeyMiddleware(next http.Handler) http.Handler
	ProcessFirebaseAuthorization(ctx context.Context, request *requests.AuthProcessFirebaseAuthorizationRequest) (*requests.AuthProcessFirebaseAuthorizationResponse, error)
}

func NewAuth(firebaseClient xfirebase.Client, apiKey string) Auth {
	return &authUsecase{
		firebase: firebaseClient,
		apiKey:   apiKey,
	}
}

type authUsecase struct {
	firebase xfirebase.Client
	apiKey   string
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

func (uc *authUsecase) AccessTokenMiddleware(next http.Handler) http.Handler {
	getTokenFromTicketCookie := func(r *http.Request) string {
		v, err := r.Cookie("ticket")

		if err != nil || v.Value == "" {
			return ""
		}

		return tokenUtils.Unswizzle(v.Value)
	}

	getTokenFromAuthorizationHeader := func(r *http.Request) string {
		return strings.ReplaceAll(r.Header.Get("authorization"), "Bearer ", "")
	}

	doCheckAccessToken := func(log *logger.Logger, ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
		defer stopwatch.StartWithLogger(log).Stop()
		authorizationToken := getTokenFromAuthorizationHeader(r)
		if authorizationToken == "" {
			authorizationToken = getTokenFromTicketCookie(r)
		}
		if authorizationToken == "" {
			log.Info("mission authorization token")
			w.WriteHeader(http.StatusUnauthorized)
			return false
		}

		authResponse, err := uc.ProcessFirebaseAuthorization(ctx, &requests.AuthProcessFirebaseAuthorizationRequest{AuthorizationToken: authorizationToken})
		if err != nil {
			log.WithError(err).Error("authentication failed")
			w.WriteHeader(http.StatusUnauthorized)
			return false
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
		return true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		ctx := ctx_helper.NewContext(r.Context())
		log := cp.Log(ctx, "AuthUsecase.AccessTokenMiddleware")

		if !doCheckAccessToken(log, ctx, w, r) {
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (uc *authUsecase) ApiKeyMiddleware(next http.Handler) http.Handler {
	doCheckApiKey := func(log *logger.Logger, w http.ResponseWriter, r *http.Request) bool {
		defer stopwatch.StartWithLogger(log).Stop()
		apiKey := r.Header.Get("x-api-key")
		if apiKey != uc.apiKey {
			log.WithField("api-key", apiKey).Error("mismatch apikey")
			w.WriteHeader(http.StatusForbidden)
			return false
		}
		return true
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		if uc.apiKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		ctx := ctx_helper.NewContext(r.Context())
		log := cp.Log(ctx, "AuthUsecase.ApiKeyMiddleware")
		if !doCheckApiKey(log, w, r) {
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
