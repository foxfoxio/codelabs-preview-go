package endpoints

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
)

type AuthHttp interface {
	Oauth2Callback(w http.ResponseWriter, r *http.Request)
}

func NewAuth(sessionUsecase usecases.Session) AuthHttp {
	return &authHttpEndpoint{
		sessionUsecase: sessionUsecase,
	}
}

type authHttpEndpoint struct {
	sessionUsecase usecases.Session
}

func (ep *authHttpEndpoint) Oauth2Callback(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Oauth2Callback OK")
}
