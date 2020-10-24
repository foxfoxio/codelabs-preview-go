package transports

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func createAuthRoutes(authEp endpoints.AuthHttp) routes {
	return routes{
		r("/oauth2/callback", authEp.Oauth2Callback, "GET"),
	}
}

func createRootRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		r("/", viewerEp.Preview, "GET"),
	}
}

func RegisterHttpRouter(router *mux.Router, authEp endpoints.AuthHttp, viewerEp endpoints.Viewer) {
	authRoutes := createAuthRoutes(authEp)
	rootRoutes := createRootRoutes(viewerEp)

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRoutes.Build(authRouter)

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRoutes.Build(rootRouter)

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		strTime := time.Now().String()
		fmt.Println("XXX", strTime)
		_, _ = fmt.Fprintf(writer, strTime)
	}).Methods("GET")
}
