package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func New(url string, queueName string) (*Client, error) {
	var conn *amqp091.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp091.Dial(url)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to connect RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Failed to open channel: %w", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // nama antrean
		true,      // durable (tetap ada meskipun server restart)
		false,     // delete when unused (hapus kalau tidak dipakai)
		false,     // exclusive (bisa diakses oleh banyak consumer)
		false,     // no-wait (tunggu respons server)
		nil,       // arguments (opsi tambahan)
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to declare a queue: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (c *Client) Publish(body string) error {
	err := c.channel.Publish(
		"",           // Exchange
		c.queue.Name, // Routing key
		false,        // Mandatory
		false,        // Immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)

	if err != nil {
		return fmt.Errorf("Failed to publish: %w", err)
	}

	return nil
}

func (c *Client) Consume() {
	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",    // Consumer tag
		true,  // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)

	if err != nil {
		log.Printf("Failed to consume messages: %v", err)
		return
	}

	for msg := range msgs {
		log.Printf("Received message: %s", string(msg.Body))
	}
}

func (c *Client) Close() error {
	c.channel.Close()
	c.conn.Close()

	return nil
}
