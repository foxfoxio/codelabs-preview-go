package endpoints

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	"github.com/foxfoxio/codelabs-preview-go/internal/utils"
	requests2 "github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type Viewer interface {
	Preview(w http.ResponseWriter, r *http.Request)
	PreviewWithQuery(w http.ResponseWriter, r *http.Request)
	Draft(w http.ResponseWriter, r *http.Request)
	Publish(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request)
	Media(w http.ResponseWriter, r *http.Request)
	Meta(w http.ResponseWriter, r *http.Request)
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

func (ep *viewerEndpoint) PreviewWithQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)

	fileId := r.URL.Query().Get("file_id")
	if fileId == "" {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	response, err := ep.viewerUsecase.Parse(ctx, &requests.ViewerParseRequest{
		FileId: fileId,
	})

	w.Header().Set("Cache-Control", "no-store")
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		_, _ = fmt.Fprint(w, response.Response)
	}
}

func (ep *viewerEndpoint) Preview(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)

	params := mux.Vars(r)
	fileId := ""

	if id, ok := params["fileId"]; ok {
		fileId = id
	}

	if fileId == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	response, err := ep.viewerUsecase.Parse(ctx, &requests.ViewerParseRequest{
		FileId: fileId,
	})

	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		responseGZip(w, []byte(response.Response), "text/html")
	}
}

func (ep *viewerEndpoint) Draft(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)
	log := cp.Log(ctx, "ViewerEndpoint.Draft")
	ctx, err := ep.authenticate(ctx, r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "unauthorized")
		return
	}

	var response *apiResponse
	defer func() {
		sendResponse(w, response)
	}()

	httpReq := &requests2.HttpDraftRequest{}
	if err := json.NewDecoder(r.Body).Decode(&httpReq); err != nil {
		log.WithError(err).Error("invalid request")
		response = newResponse(1, "invalid request", nil)
		return
	}

	res, err := ep.viewerUsecase.Draft(ctx, &requests.ViewerDraftRequest{MetaData: httpReq.Data})

	if err != nil {
		log.WithError(err).Error("process draft failed")
		response = newResponse(1, err.Error(), nil)
		return
	}

	response = successResponse(&requests2.HttpDraftResponse{FileId: res.FileId})
}

func (ep *viewerEndpoint) Publish(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)
	ctx, err := ep.authenticate(ctx, r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "unauthorized")
		return
	}

	var response *apiResponse
	defer func() {
		sendResponse(w, response)
	}()

	params := mux.Vars(r)
	fileId := ""

	if id, ok := params["fileId"]; ok {
		fileId = id
	}

	if fileId == "" {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	res, err := ep.viewerUsecase.Publish(ctx, &requests.ViewerPublishRequest{FileId: fileId})

	if err != nil {
		response = newResponse(1, err.Error(), nil)
		return
	}
	meta, _ := structToMap(res.Meta)
	response = successResponse(&requests2.HttpPublishResponse{
		Revision: res.Revision,
		Meta:     meta,
	})
}

func (ep *viewerEndpoint) View(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)

	params := mux.Vars(r)
	fileId := ""
	revision := 0

	if id, ok := params["fileId"]; ok {
		fileId = id
	}

	if rev, ok := params["revision"]; ok {
		if rev == "latest" {
			revision = 0
		} else if r, e := strconv.ParseInt(rev, 10, 32); e != nil {
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "bad request")
		} else {
			revision = int(r)
		}
	}

	if fileId == "" {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	response, err := ep.viewerUsecase.View(ctx, &requests.ViewerViewRequest{
		FileId:   fileId,
		Revision: revision,
	})

	w.Header().Set("Cache-Control", "no-store")
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		responseGZip(w, []byte(response.Response), "text/html")
	}
}

func (ep *viewerEndpoint) Media(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)

	params := mux.Vars(r)
	fileId := ""
	revision := 0
	filename := ""

	if id, ok := params["fileId"]; ok {
		fileId = id
	}
	if f, ok := params["filename"]; ok {
		filename = f
	}

	if rev, ok := params["revision"]; ok {
		if rev == "latest" {
			revision = 0
		} else if r, e := strconv.ParseInt(rev, 10, 32); e != nil {
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "bad request")
		} else {
			revision = int(r)
		}
	}

	if fileId == "" {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	response, err := ep.viewerUsecase.Media(ctx, &requests.ViewerMediaRequest{
		FileId:   fileId,
		Revision: revision,
		Filename: filename,
	})

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		responseGZip(w, response.Content, response.ContentType)
	}
}

func (ep *viewerEndpoint) authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	log := cp.Log(ctx, "ViewerEndpoint.authenticate")
	authorizationToken := ""
	if r := r.Header.Get("authorization"); r == "" {
		log.Error("missing authorization")
		return ctx, errors.New("unauthorized")

	} else {
		authorizationToken = r
	}

	authResponse, err := ep.authUsecase.ProcessFirebaseAuthorization(ctx, &requests.AuthProcessFirebaseAuthorizationRequest{AuthorizationToken: authorizationToken})

	if err != nil {
		log.WithError(err).Error("firebase authorization failed ")
		return ctx, errors.New("unauthorized")
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

	return ctx, nil
}

func (ep *viewerEndpoint) Meta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	ctx := ctx_helper.NewContextFromRequest(r)
	log := cp.Log(ctx, "ViewerEndpoint.Meta")

	ctx, err := ep.authenticate(ctx, r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = fmt.Fprint(w, "unauthorized")
		return
	}

	var response *apiResponse
	defer func() {
		sendResponse(w, response)
	}()

	params := mux.Vars(r)
	fileId := ""
	revision := 0

	if id, ok := params["fileId"]; ok {
		fileId = id
	}

	if rev, ok := params["revision"]; ok {
		if rev == "latest" {
			revision = 0
		} else if r, e := strconv.ParseInt(rev, 10, 32); e != nil {
			w.Header().Set("Cache-Control", "no-store")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, "bad request")
		} else {
			revision = int(r)
		}
	}

	if fileId == "" {
		log.Error("empty fileId")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}
	resp, err := ep.viewerUsecase.Meta(ctx, &requests.ViewerMetaRequest{
		FileId:   fileId,
		Revision: revision,
	})

	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		response = newResponse(1, err.Error(), nil)
		return
	}

	meta, _ := structToMap(resp.Meta)

	response = successResponse(&requests2.HttpMetaResponse{
		Meta: meta,
	})
}

func structToMap(data interface{}) (map[string]interface{}, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}
