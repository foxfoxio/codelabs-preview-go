package usecases

import (
	"context"
	"fmt"
	tokenUtils "github.com/foxfoxio/codelabs-preview-go/internal/token"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"golang.org/x/oauth2"
	"time"
)

type Auth interface {
	ProcessSession(ctx context.Context, request *requests.AuthProcessSessionRequest) (*requests.AuthProcessSessionResponse, error)
	ProcessOauth2Callback(ctx context.Context, request *requests.AuthProcessOauth2CallbackRequest) (*requests.AuthProcessOauth2CallbackResponse, error)
	ProcessFirebaseAuthorization(ctx context.Context, request *requests.AuthProcessFirebaseAuthorizationRequest) (*requests.AuthProcessFirebaseAuthorizationResponse, error)
}

func NewAuth(config *oauth2.Config) Auth {

	return &authUsecase{
		config: config,
	}
}

type authUsecase struct {
	config *oauth2.Config
}

func (uc *authUsecase) ProcessSession(ctx context.Context, request *requests.AuthProcessSessionRequest) (*requests.AuthProcessSessionResponse, error) {
	isValid := request.UserSession != nil && request.UserSession.IsValid()
	redirectUrl := ""
	randState := ""

	if !isValid {
		randState = fmt.Sprintf("st%d", time.Now().UnixNano())
		redirectUrl = uc.config.AuthCodeURL(randState)
	}

	return &requests.AuthProcessSessionResponse{
		IsValid:     isValid,
		State:       randState,
		RedirectUrl: redirectUrl,
	}, nil
}

func (uc *authUsecase) ProcessOauth2Callback(ctx context.Context, request *requests.AuthProcessOauth2CallbackRequest) (*requests.AuthProcessOauth2CallbackResponse, error) {
	if request.UserSession.State != request.State {
		return nil, fmt.Errorf("invalid state")
	}

	if request.Code == "" {
		return nil, fmt.Errorf("invalid code")
	}

	fmt.Println("xxx ProcessOauth2Callback request", *request)

	token, err := uc.config.Exchange(ctx, request.Code)
	if err != nil {
		return nil, fmt.Errorf("exchange code failed: %s", err.Error())
	}

	userId := ""
	name := ""
	if rawIDToken, ok := token.Extra("id_token").(string); ok {
		jwtClaim, e := tokenUtils.ExtractJwtClaims(rawIDToken)
		if e != nil {
			fmt.Println("extract jwt claim failed", e.Error())
		} else {
			userId = jwtClaim.Email
			name = jwtClaim.Name
		}
	}

	encodedToken, err := tokenUtils.EncodeBase64(token)

	if err != nil {
		fmt.Println("encode token failed", err.Error())
	}

	return &requests.AuthProcessOauth2CallbackResponse{
		Name:   name,
		UserId: userId,
		Token:  encodedToken,
	}, nil
}

func (uc *authUsecase) ProcessFirebaseAuthorization(ctx context.Context, request *requests.AuthProcessFirebaseAuthorizationRequest) (*requests.AuthProcessFirebaseAuthorizationResponse, error) {
	// TODO: verify token with firebase
	claim, err := tokenUtils.ExtractJwtClaims(request.AuthorizationToken)

	if err != nil {
		return nil, err
	}
	return &requests.AuthProcessFirebaseAuthorizationResponse{
		UserId:    claim.UserId,
		Email:     claim.Email,
		ExpiresAt: claim.ExpiresAt(),
	}, nil
}
