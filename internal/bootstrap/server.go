package bootstrap

import "net/http"

type Server struct {
	HttpHandler http.Handler
}
