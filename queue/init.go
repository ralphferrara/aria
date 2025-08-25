//||------------------------------------------------------------------------------------------------||
//|| Queue Package: Initialization
//|| init.go
//||------------------------------------------------------------------------------------------------||

package queue

//||------------------------------------------------------------------------------------------------||
//|| Import
//||------------------------------------------------------------------------------------------------||

import (
	"fmt"

	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Queue: Globals (RabbitMQ)
//||------------------------------------------------------------------------------------------------||

var (
	Rabbit = map[string]*RabbitMQWrapper{}
)

//||------------------------------------------------------------------------------------------------||
//|| Queue: Init - Connects all queues from config
//||------------------------------------------------------------------------------------------------||

func Init() error {
	cfg := config.GetConfig()
	for name, queueCfg := range cfg.Queue {
		switch queueCfg.Backend {
		case "rabbitmq":
			conn, ch, err := connectRabbit(queueCfg)
			if err != nil {
				return fmt.Errorf("failed to connect to rabbitmq '%s': %w", name, err)
			}
			Rabbit[name] = &RabbitMQWrapper{
				Name:    name,
				Conn:    conn,
				Channel: ch,
			}
		default:
			return fmt.Errorf("unsupported queue backend: %s", queueCfg.Backend)
		}
	}
	return nil
}
