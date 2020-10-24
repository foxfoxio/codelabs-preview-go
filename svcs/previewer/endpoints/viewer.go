package endpoints

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"net/http"
)

type Viewer interface {
	Preview(w http.ResponseWriter, r *http.Request)
}

func NewViewer(sessionUsecase usecases.Session, viewerUsecase usecases.Viewer) Viewer {
	return &viewerEndpoint{
		sessionUsecase: sessionUsecase,
		viewerUsecase:  viewerUsecase,
	}
}

type viewerEndpoint struct {
	sessionUsecase usecases.Session
	viewerUsecase  usecases.Viewer
}

func (ep *viewerEndpoint) Preview(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Preview OK")
}
