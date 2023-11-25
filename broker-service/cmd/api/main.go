package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	RabbitMQ *amqp.Connection
}

func main() {

	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	// initialize the app
	app := Config{
		RabbitMQ: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			counts++
			log.Printf("Failed to connect to RabbitMQ: %s", err)

		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 8 {
			fmt.Println("Failed to connect to RabbitMQ after 8 attempts")
			fmt.Printf("error: %s\n", err)
			return nil, err
		}
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Backing off...")
		log.Printf("\n Retrying in %d seconds", backOff/time.Second)
		time.Sleep(backOff)
		continue
	}

	return connection, nil

}
