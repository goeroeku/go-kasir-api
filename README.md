# Kasir API

REST API untuk sistem kasir (Point of Sale) menggunakan Go dengan Layered Architecture.

## 🛠️ Tech Stack

- **Bahasa**: Go (Golang) 1.24+
- **Database**: PostgreSQL (Supabase)
- **Driver**: pgx/v5
- **Config**: Viper
- **Documentation**: Swagger (swaggo/swag)
- **Architecture**: Layered (Handler → Service → Repository)

## 📁 Project Structure

```
kasir-api/
├── config/
│   └── config.go          # Configuration management (Viper)
├── database/
│   └── database.go        # Database connection (PostgreSQL)
├── handlers/
│   ├── category_handler.go
│   └── product_handler.go # HTTP handlers
├── services/
│   ├── category_service.go
│   └── product_service.go # Business logic
├── repositories/
│   ├── category_repository.go
│   └── product_repository.go # Database operations
├── models/
│   ├── category.go
│   └── product.go         # Data models
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml       # Swagger documentation
├── .env.example           # Environment template
├── .gitignore
├── go.mod
├── go.sum
├── main.go                # Application entry point
└── README.md
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24+
- PostgreSQL database (or Supabase)

### Setup

```bash
# Clone repository
git clone <repository-url>
cd kasir-api

# Copy environment template
cp .env.example .env

# Edit .env with your database credentials
# DB_CONN=postgresql://user:password@host:port/database?sslmode=require

# Install dependencies
go mod tidy

# Run server
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## ⚙️ Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_CONN` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db?sslmode=require` |

## 📚 API Documentation (Swagger)

Swagger UI: `http://localhost:8080/swagger/index.html`

## 📋 API Endpoints

### Health Check
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/` | API info |
| GET | `/health` | Health check |

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
  -d '{"nama":"Kopi Susu","harga":15000,"stok":100,"category_id":1}'
```

### Create Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Minuman","description":"Berbagai jenis minuman"}'
```

## 🏗️ Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │ HTTP
┌──────▼──────┐
│  Handlers   │  ← HTTP routing & request/response
└──────┬──────┘
       │
┌──────▼──────┐
│  Services   │  ← Business logic
└──────┬──────┘
       │
┌──────▼──────┐
│ Repositories│  ← Database operations
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │  ← PostgreSQL
└─────────────┘
```

## 🔗 Deployment

Deploy ke Railway atau platform lain yang mendukung Go.

### Railway
1. Connect repository ke Railway
2. Set environment variables (`PORT`, `DB_CONN`)
3. Deploy otomatis

## 📄 License

MIT License
