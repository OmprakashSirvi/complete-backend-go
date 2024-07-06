package main

import (
	"auth-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/tsawler/toolbox"
)

const webPort = "80"
var counts int64
var tools toolbox.Tools

type Config struct {
	DB *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")

	// CONNECT TO DB
	conn:= connectToDB()
	if (conn == nil) {
		log.Panic("Cannot conenct to postgres")
	} 

	// SET UP CONFIG
	app := Config{DB: conn, Models: data.New(conn)}

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Running auth-service at port : %s", webPort)
	err := server.ListenAndServe()

	if (err != nil) {
		log.Panic(err)
	}

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if (err != nil) {
		return nil, err
	}

	err = db.Ping()
	if (err != nil) {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	log.Printf("DataSource name : %s", dsn)

	for {
		connection, err := openDB(dsn)
		if (err != nil) {
			log.Println("DB is not yet ready...")
			counts++
		} else {
			log.Println("DB is connected")
			return connection
		}

		if (counts > 10) {
			log.Println(err)
			return nil
		}

		log.Println("Waiting for 2 seconds")
		time.Sleep(time.Second * 2)
	}
}
