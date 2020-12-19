package main

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	router := mux.NewRouter()
	previewer.New(router)
	handler := handlers.LoggingHandler(os.Stdout, router)
	srv := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: handler,
	}

	fmt.Printf("start serving HTTP on %s\n", "localhost:3000")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Printf("http server error")
	}
}
