// A simple logger service
package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"

	"github.com/tsawler/toolbox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	// TODO: Get from env variable
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001"
)

var client *mongo.Client
var tools toolbox.Tools

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()

	if err != nil {
		log.Panic(err)
	}

	log.Println("connected to mongo")

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	defer func() {
		log.Println("disconnecting")
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	app := Config{Models: data.New(client)}

	// start web server
	srv := &http.Server{Addr: fmt.Sprintf(":%s", webPort), Handler: app.route()}

	log.Println("created server")
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
	log.Printf("server is up and running on port %s\n", webPort)
}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("connecting to mongo")
	clientOptions := options.Client().ApplyURI(mongoURL)

	// TODO: Get from env variable
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	return mongo.Connect(ctx, clientOptions)
}
