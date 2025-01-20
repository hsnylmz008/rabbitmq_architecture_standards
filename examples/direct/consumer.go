package direct

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

	// Exchange tanımla
	err = rmq.Channel.ExchangeDeclare(
		"direct_logs",
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

	// Durable queue oluştur
	q, err := rmq.Channel.QueueDeclare(
		"direct_logs_queue", // sabit queue ismi
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		amqp.Table{
			"x-dead-letter-exchange": "dlx",                      // DLX için
			"x-message-ttl":          int32(24 * 60 * 60 * 1000), // 24 saat TTL
		},
	)
	if err != nil {
		return err
	}

	// Tüm severity'leri dinle
	severities := []string{"error", "warning", "info"}
	for _, severity := range severities {
		err = rmq.Channel.QueueBind(
			q.Name,
			severity,
			"direct_logs",
			false,
			nil)
		if err != nil {
			return err
		}
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
			// Mesaj işleme
			err := processMessage(d)
			if err != nil {
				fmt.Printf("Hata: %v\n", err)
				// Retry mekanizması
				if d.Redelivered {
					// İkinci denemede de başarısız olursa DLX'e gönder
					d.Reject(false)
				} else {
					// İlk denemede başarısız olursa tekrar dene
					d.Reject(true)
				}
				continue
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Mesajlar bekleniyor. Çıkmak için CTRL+C'ye basın")
	<-forever

	return nil
}

func processMessage(d amqp.Delivery) error {
	// Mesaj işleme simülasyonu
	fmt.Printf("Mesaj işleniyor: %s\n", d.Body)

	// Hata durumu simülasyonu
	if string(d.Body) == "Hatalı mesaj" {
		return fmt.Errorf("mesaj işlenemedi")
	}

	// İşlem simülasyonu
	time.Sleep(time.Millisecond * 100)
	return nil
}
