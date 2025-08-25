//||------------------------------------------------------------------------------------------------||
//|| Queue Package: API
//|| api.go
//||------------------------------------------------------------------------------------------------||

package queue

import (
	"time"

	"github.com/streadway/amqp"
)

//||------------------------------------------------------------------------------------------------||
//|| Publish Message (RabbitMQ)
//||------------------------------------------------------------------------------------------------||

func (q *RabbitMQWrapper) Publish(queue string, body []byte) error {
	_, err := q.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	return q.Channel.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

//||------------------------------------------------------------------------------------------------||
//|| Consume Message (RabbitMQ)
//||------------------------------------------------------------------------------------------------||

func (q *RabbitMQWrapper) Consume(queue string) (<-chan amqp.Delivery, error) {
	_, err := q.Channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return q.Channel.Consume(
		queue, // queue
		"",    // consumer
		true,  // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
}
