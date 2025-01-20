package topic

import (
	"fmt"
	"github.com/streadway/amqp"
	"rabbitmq-examples/config"
	"time"
)

func StartConsumer() error {
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

	// Durable queue oluştur
	q, err := rmq.Channel.QueueDeclare(
		"topic_queue",
		true, // durable
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "dlx",
			"x-message-ttl":          int32(24 * 60 * 60 * 1000),
			"x-max-length":           int32(10000),
			"x-overflow":             "reject-publish",
		},
	)
	if err != nil {
		return err
	}

	// Topic pattern'lerini dinle
	patterns := []string{"order.*.*", "*.*.eu"}
	for _, pattern := range patterns {
		err = rmq.Channel.QueueBind(
			q.Name,
			pattern,
			"topic_logs",
			false,
			nil)
		if err != nil {
			return err
		}
		fmt.Printf("'%s' pattern'i dinleniyor\n", pattern)
	}

	msgs, err := rmq.Channel.Consume(
		q.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			err := processMessage(d)
			if err != nil {
				fmt.Printf("Hata: %v\n", err)
				if d.Redelivered {
					d.Reject(false) // DLX'e gönder
				} else {
					d.Reject(true) // Tekrar dene
				}
				continue
			}
			d.Ack(false)
		}
	}()

	fmt.Println("Topic mesajları bekleniyor. Çıkmak için CTRL+C'ye basın")
	<-forever

	return nil
}

func processMessage(d amqp.Delivery) error {
	// Mesaj detaylarını logla
	fmt.Printf("Mesaj alındı: ID=%s, RoutingKey=%s, Headers=%v\n",
		d.MessageId, d.RoutingKey, d.Headers)

	// İşlem simülasyonu
	time.Sleep(time.Millisecond * 100)

	// Region'a göre özel işlem
	region := d.Headers["region"].(string)
	switch region {
	case "eu":
		fmt.Printf("Avrupa bölgesi mesajı işleniyor: %s\n", d.Body)
	case "us":
		fmt.Printf("ABD bölgesi mesajı işleniyor: %s\n", d.Body)
	case "asia":
		fmt.Printf("Asya bölgesi mesajı işleniyor: %s\n", d.Body)
	}

	return nil
}
