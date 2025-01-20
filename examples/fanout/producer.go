package fanout

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"rabbitmq-examples/config"
	"time"
)

func StartProducer() error {
	rmq, err := config.NewRabbitMQ()
	if err != nil {
		return err
	}
	defer rmq.Close()

	// Exchange tanımla
	err = rmq.Channel.ExchangeDeclare(
		"logs_fanout",
		"fanout",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Batch işlem için mesajları hazırla
	messages := make([]string, 5)
	for i := range messages {
		messages[i] = fmt.Sprintf("Broadcast mesaj #%d - %s", i+1, time.Now().Format("15:04:05"))
	}

	// Mesajları toplu gönder
	for _, message := range messages {
		err = rmq.SafePublish(
			"logs_fanout",
			"", // fanout için routing key gerekmez
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(message),
				Persistent:   true,
				Timestamp:    time.Now(),
				MessageId:    uuid.New().String(),
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"message_type": "broadcast",
					"producer_id":  "fanout_producer_1",
				},
			})

		if err != nil {
			return fmt.Errorf("mesaj gönderilemedi: %v", err)
		}
		fmt.Printf("Broadcast mesaj gönderildi: %s\n", message)
		time.Sleep(time.Millisecond * 100) // Rate limiting
	}

	return nil
}
