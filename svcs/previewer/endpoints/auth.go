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

//firebase header
//authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjJmOGI1NTdjMWNkMWUxZWM2ODBjZTkyYWFmY2U0NTIxMWUxZTRiNDEiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vZm94Zm94LWxlYXJuIiwiYXVkIjoiZm94Zm94LWxlYXJuIiwiYXV0aF90aW1lIjoxNjA1MDMyMzI5LCJ1c2VyX2lkIjoiMU9TM1FBZ1lRc2d0bjNwTHpIOXJLQ0dqY1A5MyIsInN1YiI6IjFPUzNRQWdZUXNndG4zcEx6SDlyS0NHamNQOTMiLCJpYXQiOjE2MDUwMzIzMjksImV4cCI6MTYwNTAzNTkyOSwiZW1haWwiOiJzc2RmQGFzZGZhZGYuc2RmZC5jb20iLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImZpcmViYXNlIjp7ImlkZW50aXRpZXMiOnsiZW1haWwiOlsic3NkZkBhc2RmYWRmLnNkZmQuY29tIl19LCJzaWduX2luX3Byb3ZpZGVyIjoicGFzc3dvcmQifX0.c7SZROTypgKcrEQEozAMQvJxkcD4ZnGD6XSuMYN4X6UxwSFhCaC7Gs1ns-ur1XMokM-7eO2BQK9-1R0VV5enOfIrzJV1ltqNMbj3sqRvKyC93dHTJ8lVEbSHzcI4Zo08R40I8tmqJWJDoWLdHUA1hg_rqOFIy8nvDcRvz6vzy66JJi_JdQyVFDECv08P5UNe7LhCn4cOELdoaR2uSUtilJDrX4ykEZHISQWotcY3adHlvCmA7KBXOFLnRRcSLHpo2a3aUS3Y5o7nk_BITyby2shqXdSWyKe2iyiUCArzlvlBrg8coryFRyV6qXHjtyOEAnpDqMfbeIwHN9kudgAyVw
