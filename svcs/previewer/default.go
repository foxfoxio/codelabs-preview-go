package previewer

import (
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/transports"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/usecases"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func New(rootRouter *mux.Router) {

	store := sessions.NewCookieStore([]byte("t0p-secret"))

	sessionUsecase := usecases.NewSession(store, "user-session")
	viewerUsecase := usecases.NewViewer()

	authEp := endpoints.NewAuth(sessionUsecase)
	viewerEp := endpoints.NewViewer(sessionUsecase, viewerUsecase)

	transports.RegisterHttpRouter(rootRouter, authEp, viewerEp)
}
