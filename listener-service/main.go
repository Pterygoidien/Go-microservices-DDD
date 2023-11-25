package main

import (
	"fmt"
	event "listener/events"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		os.Exit(1)
	}

	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create a consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}
	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
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
