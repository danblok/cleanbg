package main

import (
	"fmt"
	"os"

	"github.com/danblok/cleanbg/internal/api"
	"github.com/danblok/cleanbg/internal/httpapi"
	"github.com/danblok/cleanbg/internal/log"
)

func main() {
	log := log.New("local")
	log.Warn("logger enabled")

	client, err := api.Connect(os.Getenv("GRPC_ADDR"))
	if err != nil {
		log.Error("grpc connection failed", "error", err)
	}

	port := os.Getenv("HTTP_ADDR")
	server, err := httpapi.NewHTTPServer(client, port, log)
	if err != nil {
		log.Error("failed to create new server", "error", err)
		return
	}

	log.Info(fmt.Sprintf("started http server on %s", port))
	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start server", "error", err)
	}
}
