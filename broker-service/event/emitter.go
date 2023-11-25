// Package event provides functionality for emitting events using AMQP protocol.
package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter represents an event emitter.
type Emitter struct {
	connection *amqp.Connection
}

// setup initializes the event emitter by declaring the exchange.
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	return declareExchange(channel)
}

// Push sends an event with the specified severity to the AMQP exchange.
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	log.Println("Pushing event", event, "with severity", severity)

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		return err
	}
	return nil
}

// NewEventEmitter creates a new event emitter with the specified AMQP connection.
func NewEventEmitter(connection *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: connection,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
