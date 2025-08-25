//||------------------------------------------------------------------------------------------------||
//|| Queue Package: Tests
//|| queue_test.go
//||------------------------------------------------------------------------------------------------||

package queue

import (
	"testing"
	"time"

	"github.com/ralphferrara/aria/config"
)

//||------------------------------------------------------------------------------------------------||
//|| Helper: Init Test Config
//||------------------------------------------------------------------------------------------------||

func initTestConfig(t *testing.T) *config.Config {
	cfg, err := config.Init("../config.sample.json")
	if err != nil {
		t.Fatalf("failed to load test config: %v", err)
	}
	return cfg
}

//||------------------------------------------------------------------------------------------------||
//|| Test: RabbitMQ Publish & Consume (Loopback)
//||------------------------------------------------------------------------------------------------||

func TestRabbitMQPublishConsume(t *testing.T) {
	// Load config and init
	_ = initTestConfig(t)
	if err := Init(); err != nil {
		t.Fatalf("queue init failed: %v", err)
	}
	for name, q := range Rabbit {
		t.Run("Rabbit_"+name, func(t *testing.T) {
			queueName := "unit_test_queue"
			msgBody := []byte("test message " + name)
			// Publish
			if err := q.Publish(queueName, msgBody); err != nil {
				t.Fatalf("publish failed: %v", err)
			}
			// Consume (expect at least one message within a short period)
			deliveries, err := q.Consume(queueName)
			if err != nil {
				t.Fatalf("consume failed: %v", err)
			}
			select {
			case msg := <-deliveries:
				if string(msg.Body) != string(msgBody) {
					t.Fatalf("message mismatch: want %q, got %q", msgBody, msg.Body)
				}
			case <-time.After(3 * time.Second):
				t.Fatalf("timed out waiting for message")
			}
		})
	}
}
