package main

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	previewer.New(router)
	srv := &http.Server{
		Addr:    "localhost:3000",
		Handler: router,
	}

	fmt.Printf("start serving HTTP on %s", "localhost:3000")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("http server error")
	}
}
