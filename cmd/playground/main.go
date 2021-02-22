package main

import (
	"fmt"
	"github.com/foxfoxio/codelabs-preview-go/svcs/previewer"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func main() {
	server := previewer.New()
	handler := handlers.LoggingHandler(os.Stdout, server.HttpHandler)
	srv := &http.Server{
		Addr:    "0.0.0.0:3000",
		Handler: handler,
	}

	fmt.Printf("start serving HTTP on %s\n", "localhost:3000")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("http server error: %s\n", err.Error())
	}
}
