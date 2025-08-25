//||------------------------------------------------------------------------------------------------||
//|| Queue Package: Helpers
//|| helpers.go
//||------------------------------------------------------------------------------------------------||

package queue

import (
	"fmt"

	"github.com/ralphferrara/aria/config"
	"github.com/streadway/amqp"
)

//||------------------------------------------------------------------------------------------------||
//|| Build RabbitMQ DSN
//||------------------------------------------------------------------------------------------------||

func buildRabbitDSN(cfg config.QueueInstanceConfig) string {
	if cfg.User != "" && cfg.Password != "" {
		return fmt.Sprintf(
			"amqp://%s:%s@%s:%d%s",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost, // <-- fix here
		)
	}
	return fmt.Sprintf(
		"amqp://%s:%d%s",
		cfg.Host, cfg.Port, cfg.Vhost, // <-- fix here
	)
}

//||------------------------------------------------------------------------------------------------||
//|| Connect RabbitMQ
//||------------------------------------------------------------------------------------------------||

func connectRabbit(cfg config.QueueInstanceConfig) (*amqp.Connection, *amqp.Channel, error) {
	dsn := buildRabbitDSN(cfg)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}
	return conn, ch, nil
}
