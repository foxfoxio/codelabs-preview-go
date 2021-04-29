package transports

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer/endpoints"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func createCodelabsRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		// REST model
		r("/copy", viewerEp.Copy, "GET"),
		r("/copy", viewerEp.Copy, "POST"),
		r("/{fileId}/", viewerEp.Publish, "POST"),
		r("/{fileId}/permission/read", viewerEp.PermissionRead, "PUT"),
		r("/{fileId}/publish/", viewerEp.Publish, "POST"),
		r("/{fileId}/", viewerEp.View, "GET"),
		r("/{fileId}/img/{filename}", viewerEp.Media, "GET"),
		//r("/{fileId}/meta/latest", viewerEp.Meta, "GET"),
		r("/{fileId}/meta/{revision}", viewerEp.Meta, "GET"),
		r("/{fileId}/meta", viewerEp.Meta, "GET"),
		//r("/{fileId}/latest/", viewerEp.View, "GET"),
		//r("/{fileId}/latest/img/{filename}", viewerEp.Media, "GET"),
		r("/{fileId}/preview/", viewerEp.Preview, "GET"),
		r("/{fileId}/{revision}/", viewerEp.View, "GET"),
		r("/{fileId}/{revision}/img/{filename}", viewerEp.Media, "GET"),
		r("/", viewerEp.Draft, "POST"),
	}
}

func createRootRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		// for backward compat.
		r("/draft", viewerEp.Draft, "POST"),
		r("/copy", viewerEp.Copy, "GET"),
		r("/", viewerEp.PreviewWithQuery, "GET"),
	}
}

func createDraftRoutes(viewerEp endpoints.Viewer) routes {
	return routes{
		r("/", viewerEp.Draft, "POST"),
	}
}

func RegisterHttpRouter(router *mux.Router, viewerEp endpoints.Viewer) {
	rootRoutes := createRootRoutes(viewerEp)
	draftRoutes := createDraftRoutes(viewerEp)
	codeLabsRoutes := createCodelabsRoutes(viewerEp)

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
