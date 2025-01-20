package config

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	mu      sync.Mutex // Channel işlemleri için mutex
}

// Connection pool için basit bir implementasyon
var (
	instance *RabbitMQ
	once     sync.Once
)

func NewRabbitMQ() (*RabbitMQ, error) {
	var err error
	once.Do(func() {
		instance, err = connect()
	})

	if err != nil {
		return nil, err
	}

	return instance, nil
}

func connect() (*RabbitMQ, error) {
	// Bağlantı için retry mekanizması
	var conn *amqp.Connection
	var err error

	for i := 0; i < 3; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return nil, fmt.Errorf("RabbitMQ bağlantısı kurulamadı: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Kanal oluşturulamadı: %v", err)
	}

	// Global channel ayarları
	err = ch.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, fmt.Errorf("Qos ayarlanamadı: %v", err)
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}

// SafePublish thread-safe publish işlemi
func (r *RabbitMQ) SafePublish(exchange, routingKey string, msg amqp.Publishing) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Publisher confirms aktif et
	err := r.Channel.Confirm(false)
	if err != nil {
		return fmt.Errorf("publisher confirms aktif edilemedi: %v", err)
	}

	confirms := r.Channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	err = r.Channel.Publish(
		exchange,
		routingKey,
		true,  // mandatory
		false, // immediate
		msg,
	)
	if err != nil {
		return err
	}

	// Publish onayını bekle
	if confirmed := <-confirms; !confirmed.Ack {
		return fmt.Errorf("mesaj onaylanmadı")
	}

	return nil
}
