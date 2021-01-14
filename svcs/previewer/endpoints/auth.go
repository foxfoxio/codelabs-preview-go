package endpoints

import (
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
	"strings"
)

type AuthHttp interface {
	Oauth2Callback(w http.ResponseWriter, r *http.Request)
	LoginWithToken(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
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

func (ep *authHttpEndpoint) LoginWithToken(w http.ResponseWriter, r *http.Request) {
	ctx, session := ep.sessionUsecase.GetContextAndSession(r)
	log := cp.Log(ctx, "AuthHttp.LoginWithToken")

	authorizationToken := ""
	if r := r.Header.Get("authorization"); r == "" {
		log.Warn("authorization token is empty")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, http.StatusText(http.StatusUnauthorized))
		return
	} else {
		authorizationToken = r
	}

	authResponse, err := ep.authUsecase.ProcessFirebaseAuthorization(ctx, &requests.AuthProcessFirebaseAuthorizationRequest{AuthorizationToken: authorizationToken})

	if err != nil {
		log.WithError(err).Error("process firebase authen failed")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "invalid token")
		return
	}

	session.UserId = authResponse.UserId
	session.Name = authResponse.Email
	session.ExpiresAt = authResponse.ExpiresAt
	session.Token = authorizationToken
	session.State = ""
	session.RedirectUrl = ""

	if !session.IsValid() {
		log.WithField("ExpiresAt", session.ExpiresAt).Error("token is not valid")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprintf(w, "invalid token")
		return

	}

	err = session.Save(r, w)
	if err != nil {
		fmt.Println("xxx save session failed", err)
	}
}

func (ep *authHttpEndpoint) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, session := ep.sessionUsecase.GetContextAndSession(r)
	log := cp.Log(ctx, "AuthHttp.Logout")

	err := session.Invalidate(r, w)
	if err != nil {
		log.WithError(err).Println("xxx save session failed", err)
	}
}
