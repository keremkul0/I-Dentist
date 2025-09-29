# I-Dentist / Dental Clinic System

Bu repo, bir diş kliniği yönetim sistemi için Go tabanlı bir backend, Postgres ve Redis altyapısı ile isteğe bağlı HashiCorp Vault entegrasyonu içerir.

- Backend dili: Go (Gorilla/Mux, GORM, Viper, Zerolog)
- Altyapı: Postgres, Redis, Docker Compose
- Güvenlik: HashiCorp Vault (JWT secret ve e-posta parolası gibi sırlar için)

Ana kaynaklar `dental-clinic-system/` dizini altındadır.

## İçindekiler
- Özellikler
- Mimari ve dizin yapısı
- Önkoşullar
- Hızlı başlangıç (lokalde çalışma)
- Vault kurulumu ve sır yönetimi (opsiyonel ama tavsiye edilir)
- Test çalıştırma
- GoLand (JetBrains GoLand) ile çalıştırma ve test
- Docker ile derleme ve çalıştırma
- Sorun giderme
- Güvenlik notları

---

## Özellikler
- Kullanıcı/rol yönetimi, hasta, randevu ve prosedür API’leri
- JWT tabanlı kimlik doğrulama ve yetkilendirme
- Redis ile token/oturum yönetimi
- GORM ile Postgres veritabanı erişimi ve otomatik migrasyon
- Arkaplanda cron ile JWT temizliği
- Viper ile profil bazlı konfigürasyon (ENV=qa|prod)
- Vault ile sırların yönetilmesi (JWT secret, e-posta parolası)

## Mimari ve dizin yapısı (özet)
```
I-Dentist/
  ├─ README.md               ← Bu dosya
  ├─ .gitignore
  └─ dental-clinic-system/
       ├─ backend/
       │   ├─ main.go        ← Uygulama girişi
       │   ├─ go.mod         ← Go modülü (module: dental-clinic-system)
       │   ├─ Dockerfile
       │   ├─ resources/
       │   │   ├─ application_example.yml ← Örnek konfigürasyon
       │   │   └─ application.yml         ← Gerçek konfigürasyon (siz oluşturun)
       │   ├─ api/, application/, infrastructure/, middleware/, models/, ...
       │   └─ validations/                ← Testler (ör. validationsUser_test.go)
       ├─ docker-compose.yml  ← Postgres ve Redis
       ├─ vault-compose.yml   ← (İsteğe bağlı) Vault’u Docker ile çalıştırmak için
       ├─ vault-config/       ← Örnek Vault konfigürasyonu
       └─ Email_HTMLs/        ← E-posta şablonları
```

## Önkoşullar
- Go (önerilen: 1.22+). Not: Dockerfile ve go.mod içinde 1.25 referansı var; yerel ortamda 1.22+ kullanın. Eğer derleme hatası alırsanız Go sürümünü 1.22 veya 1.23 ile deneyin ve Dockerfile’ı buna göre uyarlayın.
- Docker ve Docker Compose
- Git

## Hızlı başlangıç (lokalde çalışma)
1) Altyapıyı başlatın (Postgres ve Redis):
```
cd dental-clinic-system
docker compose up -d
```

2) Backend konfigürasyonunu oluşturun:
- `dental-clinic-system/backend/resources/application_example.yml` dosyasını kopyalayıp `application.yml` olarak adlandırın ve değerleri yerelinize göre güncelleyin.
- Profil seçimi için ortam değişkeni olarak `ENV=qa` kullanılır (boşsa varsayılan `qa`).

3) (Opsiyonel, fakat tavsiye edilir) Vault’u başlatın ve sırları yükleyin; bkz. aşağıda “Vault kurulumu”. Vault kullanmıyorsanız `application.yml` içindeki `vault` kısmındaki değerler dev/test için boş/yerel kalsın ve `ReadConfigFromVault` kısmının başarısız olmaması için Vault erişimi sağlanmalıdır. Kod Vault’u zorunlu kıldığı için lokal geliştirmede de Vault’u çalıştırmanız önerilir.

4) Backend’i çalıştırın:
```
cd dental-clinic-system\backend
set ENV=qa
go run .
```
- Sunucu varsayılan olarak `:8080` portunda açılır (resources/application.yml → server.port).

## Vault kurulumu ve sır yönetimi
Vault’u Docker Compose ile çalıştırabilirsiniz.

1) Vault’u ayağa kaldırın:
```
cd dental-clinic-system
docker compose -f vault-compose.yml up -d
```

2) Vault’u ilk kez başlatma ve unseal/root token alma:
```
docker exec -it vault sh
# kapsül içinde:
vault operator init -key-shares=1 -key-threshold=1
# ÇIKAN "Unseal Key" ve "Initial Root Token" değerlerini güvenle saklayın.
# Ardından unseal yapın:
vault operator unseal <UNSEAL_KEY>
# root token ile giriş yapın:
vault login <ROOT_TOKEN>
```

3) Gerekli sırları yazın (backend bunları okur):
- JWT anahtarı: path `secret/jwt_token`, key `token`
- E-posta parolası: path `secret/email_password`, key `password`
```
vault kv put secret/jwt_token token="your-jwt-secret"
vault kv put secret/email_password password="your-mail-password"
```

4) Backend konfigürasyonunda (`backend/resources/application.yml`) Vault bilgilerinizi belirtin:
- `vault.addr` (örn. `http://127.0.0.1:8200`)
- `vault.initialRootToken` ve `vault.unsealKeys`

Not: Uygulama açılışında Vault’a bağlanır, gerekirse unseal eder, sonra yukarıdaki iki sır değerini alır.

## Test çalıştırma
Testler `backend` modülü altındadır. Modül kökü `backend/` dizinidir, bu yüzden testleri bu dizinden çalıştırmalısınız.

- Tüm testler:
```
cd dental-clinic-system\backend
go test ./... -v
```

- Sadece bir paket veya dosya (örnek):
```
cd dental-clinic-system\backend
go test ./validations -run TestUserValidations -v
```

Sık görülen hata ve çözüm:
- Hata: `pattern ./dental-clinic-system/backend/...: directory prefix dental-clinic-system\backend does not contain main module or its selected dependencies`
- Sebep: Komutu repo kökünden çalıştırmak ve modül kökü olmayan bir yolu desende kullanmak.
- Çözüm: Önce `cd dental-clinic-system\backend` deyin ve sonra `go test ./...` çalıştırın.
- Alternatif: Repo köküne bir `go.work` dosyası ekleyin:
```
cd <repo-kökü>
go work init .\dental-clinic-system\backend
# Sonra kökten de şunu çalıştırabilirsiniz:
go test .\dental-clinic-system\backend\... -v
```

## GoLand ile çalıştırma (Run/Debug) ve test
- Uygulamayı çalıştırma:
  1) Run/Debug Configurations → New → Go Build
  2) Name: Backend (qa)
  3) Run kind: File; File: `dental-clinic-system/backend/main.go`
  4) Working directory: `dental-clinic-system/backend`
  5) Environment: `ENV=qa`
  6) Apply → OK; Sağ üstteki yeşil Run tuşu ile çalıştırın.

- Tüm testleri çalıştırma:
  1) Run/Debug Configurations → New → Go Test
  2) Test kind: All packages in directory
  3) Directory: `dental-clinic-system/backend`
  4) "Recursive" işaretli olsun (alt paketler dahil)
  5) Environment: gerekirse `ENV=qa`
  6) Apply → OK ve Run.

## Docker ile derleme ve çalıştırma
Backend Dockerfile’ı `backend/` dizinindedir. İmajı derlemek ve konteyneri çalıştırmak için:
```
cd dental-clinic-system\backend
docker build -t i-dentist-api .
# Çalıştırırken ENV=qa ve 8080 portunu açın
docker run --rm -p 8080:8080 -e ENV=qa --name i-dentist-api i-dentist-api
```
Not: Dockerfile `golang:1.25-alpine` tabanı kullanıyor. Eğer bu imaj tag’i mevcut değilse, Dockerfile’ı yerelde kullandığınız bir Go sürümüne (örn. `golang:1.22-alpine`) güncelleyerek yeniden deneyin.

## Sorun giderme
- Bağlantı/Vault hataları: Vault çalıştığından, unseal edildiğinden ve `application.yml` içindeki Vault ayarlarının doğru olduğundan emin olun. `VAULT_ADDR`/`addr` uyumuna dikkat edin.
- Veritabanı hataları: Postgres konteyneriniz çalışıyor olmalı (docker ps). `resources/application.yml` içindeki `database.dns`’i yerelinize göre güncelleyin.
- Redis hataları: Redis konteyneri çalışıyor olmalı; `addr` ve `password` doğru olsun.
- Port çakışmaları: `server.port` veya Docker port eşlemelerini değiştirin.
- Go mod/sürüm uyumsuzlukları: Go 1.22+ ile derlemeyi deneyin; hatada `go env` çıktısını kontrol ederek `GOMOD` ve mod kökünü doğrulayın.

## Güvenlik notları
- Sırlar (DB şifreleri, JWT secret, e-posta parolası) VCS’e (Git) kesinlikle commit edilmemelidir. `application_example.yml`’i temel alarak yerelde `application.yml` oluşturun ve `.gitignore` ile hariç tutun.
- Mevcut `backend/resources/application.yml` içinde gerçek IP veya sır varsa bunları derhal kaldırın/rotasyon yapın. Vault kullanıyorsanız sırları yalnızca Vault’a yazın.
- Üretim ortamında TLS kullanın ve Vault’u HTTPS ile çalıştırın.

---

## Lisans
Bu projenin lisansını belirtin (ör. MIT, Apache-2.0). Yoksa bu kısmı kaldırabilirsiniz.
