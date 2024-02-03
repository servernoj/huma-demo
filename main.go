package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
	"github.com/servernoj/huma-demo/api"
	v1 "github.com/servernoj/huma-demo/api/v1"
	v2 "github.com/servernoj/huma-demo/api/v2"
)

type ServiceOptions struct {
	Port int `help:"Port to listen on" short:"p" default:"8001"`
}

func main() {
	Server()
}

func Server() {
	// -- Huma CLI
	cli := huma.NewCLI(func(hooks huma.Hooks, options *ServiceOptions) {
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		// http server basedd on Gin router
		port, _ := strconv.Atoi(os.Getenv("PORT"))
		if port == 0 {
			port = options.Port
		}
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: router,
		}
		// -- V1
		api.Setup[v1.VersionedImpl](router, api.VersionConfig{
			Tag:    "v1",
			SemVer: "1.0.0",
		})
		// -- V2
		api.Setup[v2.VersionedImpl](router, api.VersionConfig{
			Tag:    "v2",
			SemVer: "2.0.0",
		})

		// Hooks
		hooks.OnStart(func() {
			log.Printf("starting server on port: %d\n", port)
			log.Fatal(server.ListenAndServe())
		})
		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		})
	})
	cli.Run()
}
