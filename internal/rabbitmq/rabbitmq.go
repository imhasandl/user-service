package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Constants related to RabbitMQ configuration
const (
	// ExchangeName is the default exchange name used for RabbitMQ communications
	ExchangeName = "notification.topic"
	QueueName    = "notification_service_queue"
	RoutingKey   = "message-service.notification"
)

// RabbitMQ represents a connection to a RabbitMQ server and provides methods
// for publishing and consuming messages.
type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ instance and establishes connection with the RabbitMQ server.
// It initializes the exchange and returns the ready-to-use RabbitMQ instance.
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Printf("can't connect to rabbit mq: %v", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("can't connect to the channel: %v", err)
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

// Close gracefully closes the RabbitMQ connection and related channels.
func (r *RabbitMQ) Close() error {
	if r.Channel != nil {
		if err := r.Channel.Close(); err != nil {
			log.Printf("error closing channel: %v", err)
			return err
		}
	}
	if r.Conn != nil {
		if err := r.Conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
			return err
		}
	}
	return nil
}
