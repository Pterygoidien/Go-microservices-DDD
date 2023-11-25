package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares a topic exchange with the given name and properties.
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topics", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
}

// declareRandomQueue declares a random queue with the given properties.
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // auto-deleted
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}
