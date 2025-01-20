package fanout

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

	// Durable queue oluştur
	q, err := rmq.Channel.QueueDeclare(
		"fanout_queue", // sabit isim
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		amqp.Table{
			"x-dead-letter-exchange": "dlx",
			"x-message-ttl":          int32(24 * 60 * 60 * 1000), // 24 saat TTL
			"x-max-length":           int32(10000),               // Maximum kuyruk uzunluğu
			"x-overflow":             "reject-publish",           // Kuyruk dolunca yeni mesajları reddet
		},
	)
	if err != nil {
		return err
	}

	err = rmq.Channel.QueueBind(
		q.Name,
		"", // fanout için routing key gerekmez
		"logs_fanout",
		false,
		nil)
	if err != nil {
		return err
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
					d.Reject(false) // DLX'e gönder
					fmt.Printf("Mesaj DLX'e gönderildi: %s\n", d.MessageId)
				} else {
					d.Reject(true) // Tekrar kuyruğa al
					fmt.Printf("Mesaj tekrar kuyruğa alındı: %s\n", d.MessageId)
				}
				continue
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Fanout mesajları bekleniyor. Çıkmak için CTRL+C'ye basın")
	<-forever

	return nil
}

func processMessage(d amqp.Delivery) error {
	// Mesaj bilgilerini logla
	fmt.Printf("Mesaj alındı: ID=%s, Timestamp=%v, Headers=%v\n",
		d.MessageId, d.Timestamp, d.Headers)

	// Mesaj işleme simülasyonu
	time.Sleep(time.Millisecond * 100)

	// Hata simülasyonu (her 5 mesajdan birinde hata)
	if time.Now().UnixNano()%5 == 0 {
		return fmt.Errorf("işlem hatası simülasyonu")
	}

	fmt.Printf("Mesaj başarıyla işlendi: %s\n", d.Body)
	return nil
}
