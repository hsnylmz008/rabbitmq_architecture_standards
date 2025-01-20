# RabbitMQ Mimari Standartları: Kapsamlı Bir Rehber

![RabbitMQ Banner](https://www.rabbitmq.com/img/logo-rabbitmq.svg)

## Giriş
Modern yazılım dünyasında dağıtık sistemler ve mikroservis mimarileri giderek yaygınlaşıyor. Bu sistemlerin en önemli bileşenlerinden biri olan mesaj kuyruk sistemleri, uygulamalarımızın güvenilir ve ölçeklenebilir olmasını sağlıyor. RabbitMQ, bu alanda en çok tercih edilen çözümlerden biri.

## RabbitMQ Nedir?
RabbitMQ, açık kaynaklı bir mesaj broker'ıdır. AMQP (Advanced Message Queuing Protocol) protokolünü temel alır ve güvenilir, hızlı ve esnek bir mesajlaşma altyapısı sunar.

## Temel Kavramlar

### 1. Producer (Üretici)
- Mesajları üreten ve RabbitMQ'ya gönderen uygulamalar
- Mesajları doğrudan kuyruklara değil, exchange'lere gönderir
- Mesaj gönderirken routing key kullanır

```mermaid
graph LR
    style Producer fill:#fff,stroke:#FF6600
    style RabbitMQ fill:#FF6600,stroke:#FF6600,color:white
    Producer[Üretici] -->|Publish| RabbitMQ[RabbitMQ]
```

### 2. Consumer (Tüketici)
- Kuyruktan mesajları işleyen uygulamalar
- Mesajları senkron veya asenkron işleyebilir
- Acknowledgment mekanizması ile mesaj işleme garantisi sağlar

### 3. Exchange Tipleri

```mermaid
graph TD
    style Exchange fill:#FF6600,stroke:#FF6600,color:white
    style Direct fill:#fff,stroke:#FF6600
    style Fanout fill:#fff,stroke:#FF6600
    style Topic fill:#fff,stroke:#FF6600
    style Headers fill:#fff,stroke:#FF6600

    Exchange[Exchange Tipleri] --> Direct[Direct Exchange]
    Exchange --> Fanout[Fanout Exchange]
    Exchange --> Topic[Topic Exchange]
    Exchange --> Headers[Headers Exchange]
```

#### a) Direct Exchange
```mermaid
graph LR
    style Producer fill:#fff,stroke:#FF6600
    style Exchange fill:#FF6600,stroke:#FF6600,color:white
    style Q1 fill:#fff,stroke:#FF6600
    style Q2 fill:#fff,stroke:#FF6600

    Producer[Üretici] -->|routing_key=error| Exchange[Direct Exchange]
    Exchange -->|binding=error| Q1[Error Kuyruğu]
    Exchange -->|binding=info| Q2[Info Kuyruğu]
```

#### b) Fanout Exchange
```mermaid
graph LR
    style Producer fill:#fff,stroke:#FF6600
    style Exchange fill:#FF6600,stroke:#FF6600,color:white
    style Q1 fill:#fff,stroke:#FF6600
    style Q2 fill:#fff,stroke:#FF6600
    style Q3 fill:#fff,stroke:#FF6600

    Producer[Üretici] --> Exchange[Fanout Exchange]
    Exchange --> Q1[Kuyruk 1]
    Exchange --> Q2[Kuyruk 2]
    Exchange --> Q3[Kuyruk 3]
```

#### c) Topic Exchange
```mermaid
graph LR
    style Producer fill:#fff,stroke:#FF6600
    style Exchange fill:#FF6600,stroke:#FF6600,color:white
    style Q1 fill:#fff,stroke:#FF6600
    style Q2 fill:#fff,stroke:#FF6600

    Producer[Üretici] -->|order.created.eu| Exchange[Topic Exchange]
    Exchange -->|order.#| Q1[Tüm Siparişler]
    Exchange -->|*.created.*| Q2[Yeni Oluşturulanlar]
```

### 4. Dead Letter Exchange
```mermaid
graph LR
    style MainQ fill:#fff,stroke:#FF6600
    style DLX fill:#FF6600,stroke:#FF6600,color:white
    style DLQ fill:#fff,stroke:#FF6600

    MainQ[Ana Kuyruk] -->|Hata/TTL| DLX[Dead Letter Exchange]
    DLX --> DLQ[Dead Letter Queue]
```

## Best Practices

### 1. Mesaj Dayanıklılığı
- Önemli mesajlar için durable queue kullanın
- persistent=true ile mesajları disk'e yazın
- Publisher confirms kullanın

### 2. Performans Optimizasyonu
- Prefetch count ayarı
- Bulk mesaj işleme
- Connection pooling
- Channel yönetimi

### 3. Hata Yönetimi
- Dead Letter Exchange (DLX) kullanımı
- Retry mekanizması
- Circuit breaker pattern
- Monitoring ve alerting

## Örnek Kod

### Direct Exchange Örneği
```csharp
// Producer
channel.ExchangeDeclare("direct_logs", ExchangeType.Direct);
channel.BasicPublish(
    exchange: "direct_logs",
    routingKey: "error",
    body: Encoding.UTF8.GetBytes(message)
);

// Consumer
channel.ExchangeDeclare("direct_logs", ExchangeType.Direct);
channel.QueueBind(queueName, "direct_logs", "error");
```

### Topic Exchange Örneği
```csharp
// Producer
channel.ExchangeDeclare("topic_logs", ExchangeType.Topic);
channel.BasicPublish(
    exchange: "topic_logs",
    routingKey: "order.created.eu",
    body: Encoding.UTF8.GetBytes(message)
);

// Consumer
channel.ExchangeDeclare("topic_logs", ExchangeType.Topic);
channel.QueueBind(queueName, "topic_logs", "order.*.eu");
```

## Sonuç
RabbitMQ, doğru kullanıldığında sistemlerinize güvenilirlik, ölçeklenebilirlik ve esneklik kazandırır. Bu rehberde anlattığımız mimari standartlar ve best practice'ler, RabbitMQ tabanlı sistemlerinizi daha etkili bir şekilde tasarlamanıza yardımcı olacaktır.

#rabbitmq #messagequeue #microservices #software-architecture #dotnet