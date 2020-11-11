package endpoints

import (
	"encoding/json"
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/logger"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	requests2 "github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
	"time"
)

type Viewer interface {
	Preview(w http.ResponseWriter, r *http.Request)
	Draft(w http.ResponseWriter, r *http.Request)
}

func NewViewer(sessionUsecase usecases.Session, viewerUsecase usecases.Viewer, authUsecase usecases.Auth) Viewer {
	return &viewerEndpoint{
		sessionUsecase: sessionUsecase,
		viewerUsecase:  viewerUsecase,
		authUsecase:    authUsecase,
	}
}

type viewerEndpoint struct {
	sessionUsecase usecases.Session
	viewerUsecase  usecases.Viewer
	authUsecase    usecases.Auth
}

func (ep *viewerEndpoint) Preview(w http.ResponseWriter, r *http.Request) {
	ctx, session := ep.sessionUsecase.GetContextAndSession(r)
	authResponse, err := ep.authUsecase.ProcessSession(ctx, &requests.AuthProcessSessionRequest{UserSession: session})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "500 - OH Nooooooooo!\n%s", err.Error())
		return
	}

	if !authResponse.IsValid {
		session.State = authResponse.State
		session.RedirectUrl = r.RequestURI
		e := session.Save(r, w)
		if e != nil {
			fmt.Println("save session failed", e.Error())
		}

		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Location", authResponse.RedirectUrl)
		w.WriteHeader(http.StatusFound)
		_, _ = fmt.Fprint(w, "redirecting...")
		return
	}

	keys, ok := r.URL.Query()["file_id"]
	if !ok || len(keys[0]) < 1 {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	fileId := keys[0]
	response, err := ep.viewerUsecase.Parse(ctx, &requests.ViewerParseRequest{
		FileId: fileId,
	})
	fmt.Println(response)

	w.Header().Set("Cache-Control", "no-store")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		_, _ = fmt.Fprint(w, response.Response)
	}
}

func (ep *viewerEndpoint) Draft(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)
	log := cp.Log(ctx, "ViewerEndpoint.Draft")

	authorizationToken := ""
	if r := r.Header.Get("authorization"); r == "" {
		log.Error("missing authorization")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "unauthorized")
		return
	} else {
		authorizationToken = r
	}

	authResponse, err := ep.authUsecase.ProcessFirebaseAuthorization(ctx, &requests.AuthProcessFirebaseAuthorizationRequest{AuthorizationToken: authorizationToken})

	if err != nil {
		log.WithError(err).Error("firebase authorization failed ")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "unauthorized")
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

	var response *apiResponse
	defer func() {
		sendResponse(w, response)
	}()

	httpReq := &requests2.HttpDraftRequest{}
	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		logger.WithError(err).Error("invalid request")
		response = newResponse(1, "invalid request", nil)
		return
	}

	res, err := ep.viewerUsecase.Draft(ctx, &requests.ViewerDraftRequest{Title: httpReq.Data.Title})

	if err != nil {
		response = newResponse(1, err.Error(), nil)
		return
	}

	response = successResponse(&requests2.HttpDraftResponse{FileId: res.FileId})
}
