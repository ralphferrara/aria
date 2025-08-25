//||------------------------------------------------------------------------------------------------||
//|| Queue Package: Helpers
//|| consumer.go
//||------------------------------------------------------------------------------------------------||

package queue

func (q *RabbitMQWrapper) ConsumeQueue(queue string, handler func([]byte)) error {
	deliveries, err := q.Consume(queue)
	if err != nil {
		return err
	}
	go func() {
		for d := range deliveries {
			handler(d.Body)
		}
	}()
	return nil
}
