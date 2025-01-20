# RabbitMQ Go Örnekleri

Bu proje, RabbitMQ'nun farklı exchange tiplerini ve kullanım senaryolarını Go programlama dili ile örneklemektedir.

## Proje Yapısı
```
rabbitmq-examples/
├── config/
│   └── rabbitmq.go
├── examples/
│   ├── direct/
│   │   ├── producer.go
│   │   └── consumer.go
│   ├── fanout/
│   │   ├── producer.go
│   │   └── consumer.go
│   ├── topic/
│   │   ├── producer.go
│   │   └── consumer.go
│   └── dead-letter/
│       ├── producer.go
│       └── consumer.go
├── main.go
├── go.mod
└── README.md
```

## Kurulum

1. Projeyi klonlayın:
```bash
git clone https://github.com/kullanici/rabbitmq-examples
cd rabbitmq_architecture_standarts-examples
```

2. Bağımlılıkları yükleyin:
```bash
go mod download
```

3. RabbitMQ'nun çalıştığından emin olun:
```bash
docker run -d --name rabbitmq_architecture_standarts -p 5672:5672 -p 15672:15672 rabbitmq_architecture_standarts:management
```

## Kullanım

Farklı örnekleri çalıştırmak için aşağıdaki komutları kullanabilirsiniz:

1. Direct Exchange Örneği:
```bash
go run main.go -example=direct -mode=consumer
go run main.go -example=direct -mode=producer
```

2. Fanout Exchange Örneği:
```bash
go run main.go -example=fanout -mode=consumer
go run main.go -example=fanout -mode=producer
```

3. Topic Exchange Örneği:
```bash
go run main.go -example=topic -mode=consumer
go run main.go -example=topic -mode=producer
```

4. Dead Letter Exchange Örneği:
```bash
go run main.go -example=deadletter -mode=consumer
go run main.go -example=deadletter -mode=producer
``` # rabbitmq_architecture_standards
# rabbitmq_architecture_standards
