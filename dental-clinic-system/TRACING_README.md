# OpenTelemetry Tracing Kurulumu

Bu projede OpenTelemetry (OTEL) kullanarak distributed tracing desteği eklenmiştir.

## Yapılan Değişiklikler

### 1. Yeni Dosyalar
- `infrastructure/observability/otel.go` - OpenTelemetry tracer başlatma ve yapılandırma
- `middleware/tracingMiddleware/tracingMiddleware.go` - Fiber için custom tracing middleware
- `infrastructure/config/config_model.go` - TracingConfig yapısı eklendi

### 2. Yapılandırma
`resources/application.yml` dosyasına tracing yapılandırması eklendi:

```yaml
tracing:
  enabled: true
  serviceName: "dental-clinic-system"
  serviceVersion: "1.0.0"
  environment: "qa"
  otlpEndpoint: "localhost:4318" # OTLP HTTP endpoint
```

### 3. Özellikler
- ✅ HTTP istekleri otomatik olarak trace edilir
- ✅ Trace ID ve Span ID'ler otomatik oluşturulur
- ✅ HTTP method, URL, status code gibi metadata otomatik kaydedilir
- ✅ W3C Trace Context propagation desteği
- ✅ Graceful shutdown ile trace'lerin güvenli şekilde gönderilmesi

## Tracing Backend Kurulumu

Trace'leri görüntülemek için bir backend gereklidir. En popüler seçenekler:

### Jaeger (Önerilen)
Docker ile Jaeger çalıştırma:

```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -p 16686:16686 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest
```

Jaeger UI: http://localhost:16686

### OpenTelemetry Collector + Jaeger

`docker-compose.yml` dosyanıza ekleyebileceğiniz yapılandırma:

```yaml
version: '3.8'

services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    container_name: jaeger
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    ports:
      - "16686:16686"  # Jaeger UI
      - "4318:4318"    # OTLP HTTP receiver
      - "4317:4317"    # OTLP gRPC receiver
    networks:
      - tracing-network

  otel-collector:
    image: otel/opentelemetry-collector:latest
    container_name: otel-collector
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./otel-collector-config.yaml:/etc/otel-collector-config.yaml
    ports:
      - "4318:4318"    # OTLP HTTP receiver
      - "4317:4317"    # OTLP gRPC receiver
    depends_on:
      - jaeger
    networks:
      - tracing-network

networks:
  tracing-network:
    driver: bridge
```

### Grafana Tempo
Daha gelişmiş bir çözüm için Grafana Tempo kullanabilirsiniz.

## Kullanım

### Tracing'i Aktif/Pasif Etme

`application.yml` dosyasında `enabled` değerini değiştirin:

```yaml
tracing:
  enabled: false  # Tracing'i kapatmak için
```

### Kod İçinde Manuel Span Oluşturma

Herhangi bir service veya handler içinde:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func (s *Service) SomeMethod(ctx context.Context) error {
    tracer := otel.Tracer("dental-clinic-system")
    ctx, span := tracer.Start(ctx, "SomeMethod")
    defer span.End()
    
    // İşlem metadata'sı ekle
    span.SetAttributes(
        attribute.String("user.id", userID),
        attribute.Int("record.count", count),
    )
    
    // İşlemleriniz...
    
    // Hata durumunda
    if err != nil {
        span.RecordError(err)
        return err
    }
    
    return nil
}
```

### Database Query Tracing (İsteğe Bağlı)

GORM için tracing eklemek isterseniz:

```go
import (
    "go.opentelemetry.io/otel"
    "gorm.io/gorm"
)

// GORM callback ekle
db.Callback().Create().Before("gorm:create").Register("otel:before_create", beforeCreate)
db.Callback().Create().After("gorm:create").Register("otel:after_create", afterCreate)

func beforeCreate(db *gorm.DB) {
    tracer := otel.Tracer("dental-clinic-system")
    ctx, span := tracer.Start(db.Statement.Context, "gorm.create")
    db.Statement.Context = ctx
    db.Set("otel:span", span)
}

func afterCreate(db *gorm.DB) {
    if span, ok := db.Get("otel:span"); ok {
        span.(trace.Span).End()
    }
}
```

## Trace Görüntüleme

1. Uygulamayı başlatın
2. Jaeger UI'yi açın: http://localhost:16686
3. Service dropdown'dan "dental-clinic-system" seçin
4. "Find Traces" butonuna tıklayın
5. İstediğiniz trace'i seçerek detaylarını görün

## Trace Bilgileri

Her HTTP isteği için şu bilgiler kaydedilir:
- `http.method` - HTTP metodu (GET, POST, vb.)
- `http.url` - Tam URL
- `http.target` - Path
- `http.host` - Host adı
- `http.scheme` - HTTP/HTTPS
- `http.user_agent` - Client user agent
- `http.route` - Route pattern
- `http.status_code` - Response status code

## Performans

- Tracing minimal performans etkisine sahiptir (~1-2% overhead)
- Trace'ler asenkron olarak gönderilir
- Batch processing ile ağ trafiği optimize edilir
- Production'da sampling yapılabilir (şu an tüm istekler trace ediliyor)

## Sampling (İsteğe Bağlı)

Production'da tüm istekleri trace etmek yerine %10 sampling yapmak isterseniz, `infrastructure/observability/otel.go` dosyasında:

```go
// Bu satırı değiştirin:
sdktrace.WithSampler(sdktrace.AlwaysSample()),

// Bu şekilde:
sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // %10 sampling
```

## Sorun Giderme

### Trace'ler görünmüyor
- Jaeger'ın çalıştığından emin olun: `docker ps`
- `otlpEndpoint` yapılandırmasının doğru olduğundan emin olun
- Uygulamanın loglarını kontrol edin

### Connection refused hatası
- OTLP endpoint'in erişilebilir olduğundan emin olun
- Docker network yapılandırmasını kontrol edin
- Firewall kurallarını kontrol edin

## Kaynaklar

- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Jaeger Documentation](https://www.jaegertracing.io/docs/)
- [W3C Trace Context](https://www.w3.org/TR/trace-context/)

