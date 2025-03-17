package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tsawler/toolbox"
)

type Config struct {
	Mailer Mail
}

const webPort = "80"

var tools toolbox.Tools
var appLogger *logrus.Logger

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

	app := Config{
		Mailer: createMail(),
	}

	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// Creates a new mail object with default parameters
func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
