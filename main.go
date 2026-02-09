package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "kasir-api/docs"

	"kasir-api/config"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Kasir API
// @version 2.0
// @description REST API untuk sistem kasir (Point of Sale) dengan Layered Architecture
// @host localhost:8080
// @BasePath /

// WelcomeResponse represents the welcome message
type WelcomeResponse struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Docs        string   `json:"docs"`
	Endpoints   []string `json:"endpoints"`
}

// HealthResponse represents health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// welcomeHandler godoc
// @Summary Welcome
// @Description Menampilkan informasi welcome dan daftar endpoint
// @Tags Info
// @Produce json
// @Success 200 {object} WelcomeResponse
// @Router / [get]
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(WelcomeResponse{
		Name:        "Kasir API",
		Version:     "2.0.0",
		Description: "REST API untuk sistem kasir (Point of Sale) - Layered Architecture",
		Docs:        "/swagger/index.html",
		Endpoints: []string{
			"GET  /health     - Health check",
			"GET  /produk     - Get all produk",
			"POST /produk     - Create produk",
			"GET  /produk/:id - Get produk by ID",
			"PUT  /produk/:id - Update produk",
			"DELETE /produk/:id - Delete produk",
			"GET  /categories     - Get all categories",
			"POST /categories     - Create category",
			"GET  /categories/:id - Get category by ID",
			"PUT  /categories/:id - Update category",
			"DELETE /categories/:id - Delete category",
		},
	})
}

// healthHandler godoc
// @Summary Health check
// @Description Mengecek apakah API berjalan dengan baik
// @Tags Health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Message: "Kasir API is running",
	})
}

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database
	database.InitDB()
	defer database.CloseDB()

	// Check if database is available
	if database.DB == nil {
		log.Println("Running in demo mode without database")
		log.Println("Set DB_CONN in .env to connect to Postgres")

		// Setup routes without database (demo mode)
		setupDemoRoutes()
	} else {
		// Dependency Injection with database
		setupDatabaseRoutes()
	}

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Start server
	port := config.AppConfig.Port
	fmt.Printf("🚀 Kasir API v2.0 berjalan di http://localhost:%s\n", port)
	fmt.Printf("📚 Swagger UI: http://localhost:%s/swagger/index.html\n", port)
	fmt.Println("📋 Architecture: Layered (Handler → Service → Repository)")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupDatabaseRoutes() {
	// Initialize repositories
	productRepo := repositories.NewProductRepository(database.DB)
	categoryRepo := repositories.NewCategoryRepository(database.DB)

	// Initialize services
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routing
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/produk", productHandler.Handler)
	http.HandleFunc("/produk/", productHandler.Handler)
	http.HandleFunc("/categories", categoryHandler.Handler)
	http.HandleFunc("/categories/", categoryHandler.Handler)
}

func setupDemoRoutes() {
	// Demo mode without database - keep basic routes
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/health", healthHandler)

	// Demo handlers that return sample data
	http.HandleFunc("/produk", demoProductHandler)
	http.HandleFunc("/produk/", demoProductHandler)
	http.HandleFunc("/categories", demoCategoryHandler)
	http.HandleFunc("/categories/", demoCategoryHandler)
}

func demoProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "nama": "Kopi Susu", "harga": 15000, "stok": 100, "category_id": 1},
		{"id": 2, "nama": "Teh Manis", "harga": 8000, "stok": 150, "category_id": 1},
		{"id": 3, "nama": "Roti Bakar", "harga": 12000, "stok": 50, "category_id": 2},
	})
}

func demoCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "name": "Minuman", "description": "Berbagai jenis minuman"},
		{"id": 2, "name": "Makanan", "description": "Berbagai jenis makanan"},
	})
}
