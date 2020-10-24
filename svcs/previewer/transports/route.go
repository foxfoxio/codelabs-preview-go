package transports

import (
	"github.com/gorilla/mux"
	"net/http"
)

type route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Methods []string
}

type routes []route

func (r *routes) Build(router *mux.Router) {
	for _, route := range *r {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Methods...)
	}
}

func r(path string,
	handler func(w http.ResponseWriter, r *http.Request),
	methods ...string) route {
	return route{
		Path:    path,
		Handler: handler,
		Methods: methods,
	}
}
