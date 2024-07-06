package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tsawler/toolbox"
)

const webPort = "80"
var tools toolbox.Tools

// const host = "127.0.0.1"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service at port %s\n", webPort)

	// define HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server

	err := server.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}

}
