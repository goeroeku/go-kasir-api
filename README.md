# Kasir API

REST API sederhana untuk sistem kasir (Point of Sale) menggunakan Go standard library.

## 🛠️ Tech Stack

- **Bahasa**: Go (Golang)
- **Library**: Standard Library (`net/http`, `encoding/json`)
- **Documentation**: Swagger (swaggo/swag)
- **Storage**: In-Memory (slice)

## 🚀 Quick Start

```bash
# Clone repository
git clone <repository-url>
cd kasir-api

# Run server
go run main.go

# Build binary
go build -o kasir-api main.go
./kasir-api
```

Server akan berjalan di `http://localhost:8080`

## 📚 API Documentation (Swagger)

Swagger UI tersedia di: `http://localhost:8080/swagger/index.html`

![Swagger UI](docs/swagger-ui.png)

## 📋 API Endpoints

### Health Check
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Cek status API |

### Produk
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/produk` | Get semua produk |
| POST | `/produk` | Create produk baru |
| GET | `/produk/:id` | Get produk by ID |
| PUT | `/produk/:id` | Update produk |
| DELETE | `/produk/:id` | Delete produk |

### Categories
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/categories` | Get semua kategori |
| POST | `/categories` | Create kategori baru |
| GET | `/categories/:id` | Get kategori by ID |
| PUT | `/categories/:id` | Update kategori |
| DELETE | `/categories/:id` | Delete kategori |

## 📝 Contoh Request

### Create Produk
```bash
curl -X POST http://localhost:8080/produk \
  -H "Content-Type: application/json" \
  -d '{"nama":"Kopi Susu","harga":15000,"stok":100}'
```

### Create Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Minuman","description":"Berbagai jenis minuman"}'
```

### Get All Produk
```bash
curl http://localhost:8080/produk
```

### Update Produk
```bash
curl -X PUT http://localhost:8080/produk/1 \
  -H "Content-Type: application/json" \
  -d '{"nama":"Kopi Susu Gula Aren","harga":18000,"stok":90}'
```

### Delete Produk
```bash
curl -X DELETE http://localhost:8080/produk/1
```

## 📦 Data Models

### Produk
```json
{
  "id": 1,
  "nama": "Kopi Susu",
  "harga": 15000,
  "stok": 100
}
```

### Category
```json
{
  "id": 1,
  "name": "Minuman",
  "description": "Berbagai jenis minuman"
}
```

## 📁 Project Structure

```
kasir-api/
├── docs/
│   ├── docs.go        # Swagger generated docs
│   ├── swagger.json
│   └── swagger.yaml
├── .gitignore
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## 🔗 Deployment

Deploy ke Zeabur atau Railway untuk hosting gratis.
