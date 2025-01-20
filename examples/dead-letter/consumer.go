package deadletter

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

	// Dead Letter Queue oluştur
	dlq, err := rmq.Channel.QueueDeclare(
		"dead_letter_queue",
		true, // durable
		false,
		false,
		false,
		amqp.Table{
			"x-max-length":  int32(1000),
			"x-message-ttl": int32(7 * 24 * 60 * 60 * 1000), // 7 gün
			"x-queue-mode":  "lazy",                         // Disk öncelikli
			"x-overflow":    "reject-publish",
		},
	)
	if err != nil {
		return err
	}

	// Dead Letter Queue'yu DLX'e bağla
	err = rmq.Channel.QueueBind(
		dlq.Name,
		"dead_letter",
		"dlx_exchange",
		false,
		nil)
	if err != nil {
		return err
	}

	// Ana kuyruğu oluştur (DLX yapılandırması ile)
	args := amqp.Table{
		"x-dead-letter-exchange":    "dlx_exchange",
		"x-dead-letter-routing-key": "dead_letter",
		"x-message-ttl":             int32(5 * 60 * 1000), // 5 dakika
		"x-max-length":              int32(10000),
		"x-overflow":                "reject-publish",
	}
	mainQueue, err := rmq.Channel.QueueDeclare(
		"main_queue",
		true, // durable
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return err
	}

	// Ana kuyruğu main exchange'e bağla
	err = rmq.Channel.QueueBind(
		mainQueue.Name,
		"task",
		"main_exchange",
		false,
		nil)
	if err != nil {
		return err
	}

	// Ana kuyruktan mesajları tüket
	msgs, err := rmq.Channel.Consume(
		mainQueue.Name,
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

	// Dead Letter Queue'dan mesajları tüket
	dlMsgs, err := rmq.Channel.Consume(
		dlq.Name,
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

	// Ana kuyruk işleyicisi
	go func() {
		for d := range msgs {
			err := processMessage(d)
			if err != nil {
				fmt.Printf("Hata: %v\n", err)
				// Retry politikası
				if d.Headers["x-retry-count"] == nil {
					// İlk hata - tekrar dene
					d.Headers["x-retry-count"] = 1
					d.Reject(true)
				} else if d.Headers["x-retry-count"].(int) < 3 {
					// 3 denemeye kadar tekrar dene
					d.Headers["x-retry-count"] = d.Headers["x-retry-count"].(int) + 1
					d.Reject(true)
				} else {
					// 3 denemeden sonra DLX'e gönder
					fmt.Printf("Maksimum deneme sayısına ulaşıldı, DLX'e gönderiliyor: %s\n", d.MessageId)
					d.Reject(false)
				}
				continue
			}
			d.Ack(false)
		}
	}()

	// Dead Letter Queue işleyicisi
	go func() {
		for d := range dlMsgs {
			fmt.Printf("DLQ'da mesaj alındı: ID=%s, Headers=%v, Body=%s\n",
				d.MessageId, d.Headers, d.Body)

			// DLQ mesajlarını işle ve logla
			logDeadLetter(d)

			d.Ack(false)
		}
	}()

	fmt.Println("Mesajlar bekleniyor. Çıkmak için CTRL+C'ye basın")
	<-forever

	return nil
}

func processMessage(d amqp.Delivery) error {
	fmt.Printf("Mesaj işleniyor: ID=%s, Type=%v\n",
		d.MessageId, d.Headers["message_type"])

	// Mesaj tipine göre işlem
	switch d.Headers["message_type"] {
	case "error_prone":
		return fmt.Errorf("hata simülasyonu")
	case "long_process":
		time.Sleep(time.Second * 2) // Uzun işlem simülasyonu
	}

	fmt.Printf("Mesaj başarıyla işlendi: %s\n", d.Body)
	return nil
}

func logDeadLetter(d amqp.Delivery) {
	// Gerçek uygulamada bu bilgiler veritabanına veya log sistemine yazılabilir
	fmt.Printf("Dead Letter Log:\n"+
		"MessageID: %s\n"+
		"Timestamp: %v\n"+
		"RoutingKey: %s\n"+
		"Headers: %v\n"+
		"Body: %s\n"+
		"ReplyTo: %s\n"+
		"AppId: %s\n"+
		"-------------------\n",
		d.MessageId, d.Timestamp, d.RoutingKey,
		d.Headers, d.Body, d.ReplyTo, d.AppId)
}
