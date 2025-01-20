package deadletter

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

	// Ana exchange'i oluştur
	err = rmq.Channel.ExchangeDeclare(
		"main_exchange",
		"direct",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Dead Letter Exchange'i oluştur
	err = rmq.Channel.ExchangeDeclare(
		"dlx_exchange",
		"direct",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Test mesajları hazırla
	messages := []struct {
		body    string
		headers amqp.Table
	}{
		{
			body: "Normal mesaj 1",
			headers: amqp.Table{
				"message_type": "normal",
				"priority":     1,
			},
		},
		{
			body: "İşlenecek uzun mesaj",
			headers: amqp.Table{
				"message_type": "long_process",
				"priority":     2,
			},
		},
		{
			body: "Hatalı mesaj",
			headers: amqp.Table{
				"message_type": "error_prone",
				"priority":     3,
			},
		},
	}

	// Mesajları gönder
	for _, msg := range messages {
		err = rmq.SafePublish(
			"main_exchange",
			"task",
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(msg.body),
				Persistent:   true,
				Timestamp:    time.Now(),
				MessageId:    uuid.New().String(),
				DeliveryMode: amqp.Persistent,
				Headers:      msg.headers,
				AppId:        "dead_letter_producer",
			})

		if err != nil {
			return fmt.Errorf("mesaj gönderilemedi: %v", err)
		}
		fmt.Printf("Mesaj gönderildi: %s (Headers: %v)\n", msg.body, msg.headers)
		time.Sleep(time.Millisecond * 100) // Rate limiting
	}

	return nil
}
