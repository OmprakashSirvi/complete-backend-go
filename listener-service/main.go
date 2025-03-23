package main

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)


func main() {
	// try to connect to rabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		logrus.Fatalf("error while connecting to rabbit MQ: %v", err)
	}
	defer rabbitConn.Close()

	// start listening for messages


	// create consumer


	// watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	counts := 0
	backoff := 1 * time.Second
	var connection *amqp.Connection
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			logrus.Errorf("error while connecting to rabbit MQ: %v", err)
			// If counts exceeds the maximum retries threshold, then return error
			if counts > 10 {
				return nil, err
			}
			// Retry after sleeping for sometime
			counts++
			backoffPeriod := backoff + (time.Duration(counts) * time.Second)
			logrus.Infof("sleeping for: %v seconds", backoffPeriod)
			time.Sleep(backoffPeriod)
			continue
		}

		logrus.Info("connected to rabbit MQ" )
		connection = c
		break
	}

	return connection, nil
}