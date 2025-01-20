package main

import (
	"flag"
	"fmt"
	"os"
	"rabbitmq-examples/examples/dead-letter"
	"rabbitmq-examples/examples/direct"
	"rabbitmq-examples/examples/fanout"
	"rabbitmq-examples/examples/topic"
)

func main() {
	exampleType := flag.String("example", "direct", "Çalıştırılacak örnek (direct, fanout, topic, deadletter)")
	mode := flag.String("mode", "consumer", "Çalıştırma modu (consumer, producer)")
	flag.Parse()

	var err error

	switch *exampleType {
	case "direct":
		if *mode == "consumer" {
			err = direct.StartConsumer()
		} else {
			err = direct.StartProducer()
		}
	case "fanout":
		if *mode == "consumer" {
			err = fanout.StartConsumer()
		} else {
			err = fanout.StartProducer()
		}
	case "topic":
		if *mode == "consumer" {
			err = topic.StartConsumer()
		} else {
			err = topic.StartProducer()
		}
	case "deadletter":
		if *mode == "consumer" {
			err = deadletter.StartConsumer()
		} else {
			err = deadletter.StartProducer()
		}
	default:
		fmt.Printf("Geçersiz örnek tipi: %s\n", *exampleType)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Hata: %v\n", err)
		os.Exit(1)
	}
}
