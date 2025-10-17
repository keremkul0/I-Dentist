# OpenTelemetry Collector + Jaeger ile Tracing Nasıl Kurulur?

## 🚀 Hızlı Başlangıç

### Adım 1: Docker Compose ile Servisleri Başlatın

Proje dizininde (`dental-clinic-system/`) aşağıdaki komutu çalıştırın:

```bash
docker-compose up -d jaeger otel-collector
```

Bu komut şunları başlatır:
- **Jaeger**: Trace'leri görselleştirmek için UI
- **OpenTelemetry Collector**: Trace'leri toplamak ve Jaeger'a iletmek için

### Adım 2: Servislerin Çalıştığını Kontrol Edin

```bash
docker ps
```

Şu container'ları görmelisiniz:
- `jaeger`
- `otel-collector`

### Adım 3: Servis Durumunu Kontrol Edin

**Jaeger UI:**
- URL: http://localhost:16686
- Tarayıcınızda açın, UI görünüyorsa Jaeger çalışıyor demektir

**OpenTelemetry Collector Health:**
```bash
curl http://localhost:8888/metrics
```

### Adım 4: Go Uygulamanızı Başlatın

Backend dizininde:

```bash
cd backend
go run main.go
```

veya derleme yapıp çalıştırın:

```bash
go build -o app
./app
```

### Adım 5: Test Edin

Uygulamanıza birkaç istek gönderin:

```bash
# Login denemesi
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"test123"}'

# Başka endpoint'ler
curl http://localhost:8080/api/clinics
```

### Adım 6: Trace'leri Görüntüleyin

1. Tarayıcınızda http://localhost:16686 adresini açın
2. Sol üstteki "Service" dropdown'dan **"dental-clinic-system"** seçin
3. "Find Traces" butonuna tıklayın
4. Trace'lerinizi göreceksiniz! 🎉

## 📊 Trace'leri Anlamak

Jaeger UI'de göreceğiniz bilgiler:

### Ana Ekran (Trace List)
- **Trace ID**: Her isteğin benzersiz kimliği
- **Duration**: İsteğin toplam süresi
- **Spans**: Trace içindeki işlem sayısı
- **Service**: dental-clinic-system

### Trace Detayları
Bir trace'e tıkladığınızda:
- **Timeline**: İşlemlerin zaman çizelgesi
- **Spans**: Her işlemin detayı
  - HTTP method (GET, POST, vb.)
  - URL ve route
  - Status code
  - Hata mesajları (varsa)
  - Süre

## 🔍 Örnek Kullanım Senaryoları

### 1. Yavaş İstekleri Bulmak
- Jaeger UI'de "Min Duration" filtresini kullanın
- Örnek: 500ms'den uzun süren istekleri görmek için "500ms" yazın
- Hangi endpoint'lerin yavaş olduğunu görebilirsiniz

### 2. Hata Trace'lerini Görmek
- "Tags" bölümünden "error=true" filtresi ekleyin
- Sadece hatalı istekleri göreceksiniz

### 3. Belirli Bir Endpoint'i İzlemek
- "Tags" bölümünden "http.route=/api/clinics" gibi filtreler ekleyin

## 🏗️ Mimari

```
Go App (localhost:8080)
    ↓ (OTLP HTTP)
OpenTelemetry Collector (localhost:4318)
    ↓ (Processing, Batching)
Jaeger (localhost:14250)
    ↓
Jaeger UI (localhost:16686) ← Buradan görüntülüyorsunuz
```

## 📝 Yapılandırma Detayları

### otel-collector-config.yaml
Oluşturduğum yapılandırma dosyası:

```yaml
receivers:
  otlp:
    protocols:
      grpc: 0.0.0.0:4317
      http: 0.0.0.0:4318

processors:
  batch: # Trace'leri gruplar, performans için
    timeout: 10s
    send_batch_size: 1024
  
  memory_limiter: # Bellek kullanımını sınırlar
    limit_mib: 512

exporters:
  jaeger: # Jaeger'a gönder
    endpoint: jaeger:14250
  
  logging: # Console'a da yazdır (debug için)
    loglevel: debug
```

### docker-compose.yml
Eklediğim servisler:

**Jaeger:**
- Port 16686: Web UI
- Port 4318: OTLP HTTP (direkt Jaeger'a göndermek için)
- Port 14250: Collector'dan trace almak için

**OTel Collector:**
- Port 4318: OTLP HTTP (uygulamanızdan trace alır)
- Port 4317: OTLP gRPC
- Port 8888: Collector metrikleri

## 🛠️ Sorun Giderme

### Problem: Trace'ler görünmüyor

**1. Container'ları kontrol edin:**
```bash
docker ps
docker logs jaeger
docker logs otel-collector
```

**2. Collector loglarını inceleyin:**
```bash
docker logs -f otel-collector
```

Başarılı trace gönderimi görmelisiniz:
```
2025-01-17T... info    TracesExporter  {"kind": "exporter", ...}
```

**3. Uygulamanızın loglarını kontrol edin:**
Go uygulamanız başlarken şunu görmelisiniz:
```
{"level":"info","message":"Initializing OpenTelemetry tracer..."}
{"level":"info","message":"OpenTelemetry tracer initialized successfully"}
```

### Problem: Connection refused

**Endpoint'i kontrol edin:**
- Uygulamanız Docker içinde değilse: `localhost:4318`
- Uygulamanız Docker içindeyse: `otel-collector:4318`

Şu anki yapılandırmanızda `localhost:4318` doğru (uygulama Docker dışında çalışıyor).

### Problem: Port çakışması

Eğer 4318 portu kullanılıyorsa:

```bash
# Windows
netstat -ano | findstr :4318

# Linux/Mac
lsof -i :4318
```

Başka bir şey kullanıyorsa, docker-compose.yml'de portu değiştirin:
```yaml
otel-collector:
  ports:
    - "14318:4318"  # Host:Container
```

Ve application.yml'de:
```yaml
otlpEndpoint: "localhost:14318"
```

## 🎯 İleri Seviye: Özel Span'ler

Kodunuzda özel span'ler ekleyerek daha detaylı trace'ler oluşturabilirsiniz:

### Örnek 1: Service Method'da Span

```go
// application/clinicService/clinicService.go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func (s *ClinicService) GetClinic(ctx context.Context, id string) (*Clinic, error) {
    tracer := otel.Tracer("dental-clinic-system")
    ctx, span := tracer.Start(ctx, "ClinicService.GetClinic")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("clinic.id", id),
    )
    
    clinic, err := s.repository.FindByID(ctx, id)
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    span.SetAttributes(
        attribute.String("clinic.name", clinic.Name),
    )
    
    return clinic, nil
}
```

### Örnek 2: Database Query Span

```go
// infrastructure/repository/clinicRepository/clinicRepository.go
func (r *Repository) FindByID(ctx context.Context, id string) (*Clinic, error) {
    tracer := otel.Tracer("dental-clinic-system")
    ctx, span := tracer.Start(ctx, "Database.FindClinicByID")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("db.system", "postgresql"),
        attribute.String("db.operation", "SELECT"),
        attribute.String("db.table", "clinics"),
    )
    
    var clinic Clinic
    err := r.db.WithContext(ctx).First(&clinic, "id = ?", id).Error
    
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    return &clinic, nil
}
```

Bu şekilde Jaeger'da şöyle bir hiyerarşi göreceksiniz:

```
POST /api/clinics/:id
  ├─ ClinicService.GetClinic (15ms)
  │   └─ Database.FindClinicByID (12ms)
  └─ Response (3ms)
```

## 📈 Production İçin Öneriler

### 1. Sampling Ekleyin
Tüm istekleri trace etmek yerine %10'unu trace edin:

`infrastructure/observability/otel.go` dosyasında:
```go
sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.1)), // %10 sampling
```

### 2. Resource Attributes Ekleyin
```go
resource.WithAttributes(
    semconv.ServiceNameKey.String(config.ServiceName),
    semconv.ServiceVersionKey.String(config.ServiceVersion),
    semconv.DeploymentEnvironmentKey.String(config.Environment),
    semconv.HostNameKey.String(hostname),
    attribute.String("service.instance.id", instanceID),
)
```

### 3. Collector'ı Scale Edin
docker-compose.yml'de:
```yaml
otel-collector:
  deploy:
    replicas: 3
```

## 🔗 Faydalı Linkler

- Jaeger UI: http://localhost:16686
- Collector Metrics: http://localhost:8888/metrics
- Collector Health: http://localhost:8888/health

## 🎉 Tebrikler!

Artık distributed tracing sisteminiz hazır! Jaeger UI'de trace'lerinizi görebilir, performans sorunlarını tespit edebilir ve uygulamanızın davranışını daha iyi anlayabilirsiniz.

**Sorularınız için:** TRACING_README.md dosyasına bakın veya logları inceleyin.

