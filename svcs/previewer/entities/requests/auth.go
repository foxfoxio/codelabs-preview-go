package requests

import "github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"

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
