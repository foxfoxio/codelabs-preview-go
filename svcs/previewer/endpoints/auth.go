package endpoints

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
	"strings"
)

type AuthHttp interface {
	Oauth2Callback(w http.ResponseWriter, r *http.Request)
}

func NewAuth(sessionUsecase usecases.Session, authUsecase usecases.Auth) AuthHttp {
	return &authHttpEndpoint{
		sessionUsecase: sessionUsecase,
		authUsecase:    authUsecase,
	}
}

type authHttpEndpoint struct {
	sessionUsecase usecases.Session
	authUsecase    usecases.Auth
}

func (ep *authHttpEndpoint) Oauth2Callback(w http.ResponseWriter, r *http.Request) {
	ctx, session := ep.sessionUsecase.GetContextAndSession(r)

	state := r.FormValue("state")
	code := r.FormValue("code")

	authResponse, err := ep.authUsecase.ProcessOauth2Callback(ctx, &requests.AuthProcessOauth2CallbackRequest{
		State:       state,
		Code:        code,
		UserSession: session,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "500 - OH Nooooooooo!\n%s", err.Error())
		return
	}

	redirectingTo := session.RedirectUrl
	if strings.Contains(redirectingTo, "callback") {
		redirectingTo = "/"
	}

	session.UserId = authResponse.UserId
	session.Name = authResponse.Name
	session.Token = authResponse.Token
	session.State = ""
	session.RedirectUrl = ""
	err = session.Save(r, w)
	if err != nil {
		fmt.Println("xxx save session failed", err)
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Location", redirectingTo)
	w.WriteHeader(http.StatusFound)
	_, _ = fmt.Fprint(w, "redirecting...")
}
