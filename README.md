# TCKN Validation API

Bu API, TC Kimlik Numarası doğrulama işlemlerini NVI (Nüfus ve Vatandaşlık İşleri) SOAP servisi üzerinden gerçekleştirir.

## Özellikler

- TC Kimlik Numarası algoritma kontrolü
- NVI SOAP servisi üzerinden kimlik doğrulama
- REST API endpoint'i
- Türkçe karakter desteği

## Kurulum

```bash
git clone https://github.com/KULLANICI_ADINIZ/tckn-api.git
cd tckn-api
go mod download
```

## Kullanım

Servisi başlatmak için:

```bash
go run .
```

API endpoint'i:

```bash
curl -X POST http://localhost:8080/validate \
  -H "Content-Type: application/json" \
  -d '{
    "tckn": "11111111111",
    "ad": "AD",
    "soyad": "SOYAD",
    "dogumYili": 1990
  }'
```

## API Detayları

### POST /validate

Request body:
```json
{
  "tckn": "11111111111",
  "ad": "AD",
  "soyad": "SOYAD",
  "dogumYili": 1990
}
```

Response:
```json
{
  "valid": true
}
```

## Lisans

MIT 