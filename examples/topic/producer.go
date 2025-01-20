package topic

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"rabbitmq-examples/config"
	"strings"
	"time"
)

func StartProducer() error {
	rmq, err := config.NewRabbitMQ()
	if err != nil {
		return err
	}
	defer rmq.Close()

	err = rmq.Channel.ExchangeDeclare(
		"topic_logs",
		"topic",
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Farklı topic pattern'leri için mesajlar
	messages := map[string]string{
		"order.created.eu":     "Avrupa'da yeni sipariş oluşturuldu",
		"order.shipped.us":     "ABD'de sipariş kargoya verildi",
		"user.registered.asia": "Asya'da yeni kullanıcı kaydı",
		"order.cancelled.eu":   "Avrupa'da sipariş iptal edildi",
	}

	// Her mesaj için metadata ekle ve gönder
	for routingKey, message := range messages {
		err = rmq.SafePublish(
			"topic_logs",
			routingKey,
			amqp.Publishing{
				ContentType:  "text/plain",
				Body:         []byte(message),
				Persistent:   true,
				Timestamp:    time.Now(),
				MessageId:    uuid.New().String(),
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"message_type": "topic_event",
					"region":       routingKey[strings.LastIndex(routingKey, ".")+1:],
					"event_type":   strings.Split(routingKey, ".")[1],
				},
			})

		if err != nil {
			return fmt.Errorf("mesaj gönderilemedi [%s]: %v", routingKey, err)
		}

		fmt.Printf("Mesaj gönderildi: [%s] %s\n", routingKey, message)
		time.Sleep(time.Millisecond * 100) // Rate limiting
	}

	return nil
}
