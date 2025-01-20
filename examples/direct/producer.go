package direct

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
		"direct_logs", // exchange adı
		"direct",      // exchange tipi
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	// Batch işlem için mesajları hazırla
	severities := []string{"error", "warning", "info"}
	for _, severity := range severities {
		message := fmt.Sprintf("Bu bir %s mesajıdır", severity)

		// Persistent mesaj
		err = rmq.SafePublish(
			"direct_logs",
			severity,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(message),
				Persistent:   true, // Mesaj kalıcılığı
				Timestamp:    time.Now(),
				MessageId:    uuid.New().String(),
				DeliveryMode: amqp.Persistent,
			})

		if err != nil {
			return fmt.Errorf("mesaj gönderilemedi: %v", err)
		}

		fmt.Printf("Mesaj gönderildi: %s\n", message)
	}

	return nil
}
