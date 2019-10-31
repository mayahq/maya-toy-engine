package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/codegangsta/cli"
	"github.com/sibeshkar/maya-engine/server"
)

func serve(c *cli.Context) {
	restfulAPI := server.NewRESTfulAPI()
	go restfulAPI.Start(c.String("addr"), c.String("static"))

	// Ctrl+C handling
	handler := make(chan os.Signal, 1)
	signal.Notify(handler, os.Interrupt)
	for sig := range handler {
		if sig == os.Interrupt {
			fmt.Println("Caught interrupt signal...")
			break
		}

	}
}
