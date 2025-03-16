package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tsawler/toolbox"
)

const webPort = "80"

var tools toolbox.Tools
var appLogger *logrus.Logger

type Config struct{}

func main() {
	// Setting up ENV variables
	viper.AutomaticEnv()

	// Setting up logger service
	appLogger = logrus.New()
	appLogger.Formatter = new(logrus.JSONFormatter)

	env := viper.GetString("MY_ENV")
	appLogger.Level = logrus.InfoLevel
	if env == "" || env == "local" {
		appLogger.Level = logrus.TraceLevel
	}
	appLogger.Infof("log level: %v", appLogger.Level)

	app := Config{}

	appLogger.Infof("Starting broker service at port %s\n", webPort)

	// define HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server

	err := server.ListenAndServe()

	if err != nil {
		appLogger.Panic(err)
	}

}
