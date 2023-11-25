package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer is a struct that holds the connection to RabbitMQ
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer creates a new Consumer instance with the given RabbitMQ connection.
// It returns the Consumer instance and an error if any.
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup initializes the RabbitMQ channel and declares the exchange.
// It returns an error if any.
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload represents the data structure of the message payload.
type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// Listen starts listening for messages on the specified topics.
// It binds the queue to the topics and handles the received messages.
// It returns an error if any.
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		ch.QueueBind(
			q.Name,       // queue name
			topic,        // routing key
			"logs_topic", // exchange
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			log.Printf("Received a message: %s", payload.Name)
			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

// handlePayload handles the received payload based on its name.
// It performs different actions based on the payload name.
func handlePayload(payload Payload) {
	fmt.Printf("Handling payload: %s\n", payload.Name)
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Printf("Error logging event: %s", err)
		}
	case "auth":
		// authenticate
		// you can have as many cases as you want, as long as you write the logic
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	}
}

// logEvent logs the event by sending a POST request to the logger service.
// It returns an error if any.
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
