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

func createCodelabsRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		// REST model
		r("/{fileId}", viewerEp.Publish, "POST"),
		r("/{fileId}", viewerEp.View, "GET"),
		r("/{fileId}/meta/latest", viewerEp.Meta, "GET"),
		r("/{fileId}/meta/{revision}", viewerEp.Meta, "GET"),
		r("/{fileId}/meta", viewerEp.Meta, "GET"),
		r("/{fileId}/latest/", viewerEp.View, "GET"),
		r("/{fileId}/latest/img/{filename}", viewerEp.Media, "GET"),
		r("/{fileId}/preview", viewerEp.Preview, "GET"),
		r("/{fileId}/{revision}/", viewerEp.View, "GET"),
		r("/{fileId}/{revision}/img/{filename}", viewerEp.Media, "GET"),
		r("/", viewerEp.Draft, "POST"),
	}
}

func createRootRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		// for backward compat.
		r("/draft", viewerEp.Draft, "POST"),
		r("/", viewerEp.PreviewWithQuery, "GET"),
	}
}

func createDraftRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		r("/", viewerEp.Draft, "POST"),
	}
}

func RegisterHttpRouter(router *mux.Router, authEp endpoints.AuthHttp, viewerEp endpoints.Viewer) {
	authRoutes := createAuthRoutes(authEp)
	rootRoutes := createRootRoutes(viewerEp)
	draftRoutes := createDraftRoutes(viewerEp)
	codeLabsRoutes := createCodelabsRoutes(viewerEp)

	authRouter := router.PathPrefix("/auth").Subrouter()
	authRoutes.Build(authRouter)

	pRouter := router.PathPrefix("/p").Subrouter()
	rootRoutes.Build(pRouter)

	vRouter := router.PathPrefix("/v").Subrouter()
	codeLabsRoutes.Build(vRouter)

	draftRouter := router.PathPrefix("/draft").Subrouter()
	draftRoutes.Build(draftRouter)

	rootRouter := router.PathPrefix("/").Subrouter()
	rootRoutes.Build(rootRouter)

	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		strTime := time.Now().String()
		fmt.Println("XXX", strTime)
		_, _ = fmt.Fprintf(writer, strTime)
	}).Methods("GET")
}
