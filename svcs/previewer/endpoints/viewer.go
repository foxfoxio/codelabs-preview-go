package endpoints

import (
	"encoding/json"
	"fmt"
	cp "github.com/foxfoxio/codelabs-preview-go/internal"
	"github.com/foxfoxio/codelabs-preview-go/internal/ctx_helper"
	requests2 "github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Viewer interface {
	Preview(w http.ResponseWriter, r *http.Request)
	PreviewWithQuery(w http.ResponseWriter, r *http.Request)
	Draft(w http.ResponseWriter, r *http.Request)
	Publish(w http.ResponseWriter, r *http.Request)
	View(w http.ResponseWriter, r *http.Request)
	Media(w http.ResponseWriter, r *http.Request)
	Meta(w http.ResponseWriter, r *http.Request)
	Copy(w http.ResponseWriter, r *http.Request)
	PermissionRead(w http.ResponseWriter, r *http.Request)
}

func NewViewer(viewerUsecase usecases.Viewer) Viewer {
	return &viewerEndpoint{
		viewerUsecase: viewerUsecase,
	}
}

type viewerEndpoint struct {
	viewerUsecase usecases.Viewer
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
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-store")
	log := cp.Log(ctx, "ViewerEndpoint.Draft")

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
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-store")

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

func (ep *viewerEndpoint) Meta(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-store")
	log := cp.Log(ctx, "ViewerEndpoint.Meta")

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

func (ep *viewerEndpoint) Copy(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	url := r.URL.Query().Get("url")
	name := r.URL.Query().Get("name")
	prefix := r.URL.Query().Get("prefix")
	suffix := r.URL.Query().Get("suffix")

	if url == "" {
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, "bad request")
		return
	}

	req := &requests.CopyGoogleDocRequest{
		GoogleDocPath: url,
	}

	if name != "" {
		req.FileName = &name
	}

	if prefix != "" {
		req.Prefix = &prefix
	}

	if suffix != "" {
		req.Suffix = &suffix
	}

	res, err := ep.viewerUsecase.Copy(ctx, req)

	w.Header().Set("Cache-Control", "no-store")
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err.Error() == "unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err.Error() == "invalid file path" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err.Error())
	} else {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Location", res.GoogleDocPath)
		w.WriteHeader(http.StatusMovedPermanently)
		_, _ = fmt.Fprint(w, "redirecting...")
	}
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

func (ep *viewerEndpoint) PermissionRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Cache-Control", "no-store")

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

	res, err := ep.viewerUsecase.PermissionRead(ctx, &requests.FilePermissionRequest{FileId: fileId})

	if err != nil {
		response = newResponse(1, err.Error(), nil)
		return
	}
	response = successResponse(&requests2.HttpFilePermissionResponse{
		Success: res.Success,
	})
}
