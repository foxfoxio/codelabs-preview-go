package endpoints

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/entities/requests"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
)

type Viewer interface {
	Preview(w http.ResponseWriter, r *http.Request)
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
