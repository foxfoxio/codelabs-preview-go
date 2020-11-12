package requests

import (
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"time"
)

type AuthProcessSessionRequest struct {
	UserSession *entities.UserSession
}

type AuthProcessSessionResponse struct {
	IsValid     bool
	State       string
	RedirectUrl string
}

type AuthProcessOauth2CallbackRequest struct {
	State       string
	Code        string
	UserSession *entities.UserSession
}

type AuthProcessOauth2CallbackResponse struct {
	Name   string
	UserId string
	Token  string
}

type AuthProcessFirebaseAuthorizationRequest struct {
	AuthorizationToken string
}

type AuthProcessFirebaseAuthorizationResponse struct {
	UserId    string
	Email     string
	ExpiresAt time.Time
}
