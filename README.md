# Kasir API

REST API for Point of Sale (POS) system using Go with Layered Architecture.

## 🛠️ Tech Stack

- **Language**: Go (Golang) 1.24+
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

Server will run at `http://localhost:8080`

## ⚙️ Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_CONN` | PostgreSQL connection string | `postgresql://user:pass@host:5432/db?sslmode=require` |

## 📚 API Documentation (Swagger)

Swagger UI: `http://localhost:8080/swagger/index.html`

## 📋 API Endpoints

### Health Check
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API info |
| GET | `/health` | Health check |

### Products
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/products` | Get all products |
| POST | `/products` | Create new product |
| GET | `/products/:id` | Get product by ID |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |

### Categories
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/categories` | Get all categories |
| POST | `/categories` | Create new category |
| GET | `/categories/:id` | Get category by ID |
| PUT | `/categories/:id` | Update category |
| DELETE | `/categories/:id` | Delete category |

## 📝 Example Requests

### Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Kopi Susu","price":15000,"stock":100,"category_id":1}'
```

### Create Category
```bash
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Beverages","description":"Various drinks"}'
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

## 🗄️ Database Schema

```sql
-- Categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT
);

-- Products table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    stock INTEGER NOT NULL,
    category_id INTEGER REFERENCES categories(id)
);
```

## 🔗 Deployment

Deploy to Railway or any platform that supports Go.

### Railway
1. Connect repository to Railway
2. Set environment variables (`PORT`, `DB_CONN`)
3. Automatic deployment

## 📄 License

MIT License
