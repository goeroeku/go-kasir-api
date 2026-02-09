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
// @description REST API for Point of Sale (POS) system with Layered Architecture
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
// @Description Display welcome information and endpoint list
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
		Description: "REST API for Point of Sale (POS) system - Layered Architecture",
		Docs:        "/swagger/index.html",
		Endpoints: []string{
			"GET  /health       - Health check",
			"GET  /products     - Get all products",
			"POST /products     - Create product",
			"GET  /products/:id - Get product by ID",
			"PUT  /products/:id - Update product",
			"DELETE /products/:id - Delete product",
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
// @Description Check if the API is running properly
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
	fmt.Printf("🚀 Kasir API v2.0 running at http://localhost:%s\n", port)
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
	transactionRepo := repositories.NewTransactionRepository(database.DB)
	reportRepo := repositories.NewReportRepository(database.DB)

	// Initialize services
	productService := services.NewProductService(productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	reportService := services.NewReportService(reportRepo)

	// Initialize handlers
	productHandler := handlers.NewProductHandler(productService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	reportHandler := handlers.NewReportHandler(reportService)

	// Setup routing
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/health", healthHandler)

	// Product Routes
	http.HandleFunc("/products", productHandler.Handler)
	http.HandleFunc("/products/", productHandler.Handler)

	// Category Routes
	http.HandleFunc("/categories", categoryHandler.Handler)
	http.HandleFunc("/categories/", categoryHandler.Handler)

	// Transaction Routes
	http.HandleFunc("/transactions", transactionHandler.Handler)

	// Report Routes
	http.HandleFunc("/reports/today", reportHandler.GetReportToday)
	http.HandleFunc("/reports", reportHandler.GetReportCustom)
}

func setupDemoRoutes() {
	// Demo mode without database - keep basic routes
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/health", healthHandler)

	// Demo handlers that return sample data
	http.HandleFunc("/products", demoProductHandler)
	http.HandleFunc("/products/", demoProductHandler)
	http.HandleFunc("/categories", demoCategoryHandler)
	http.HandleFunc("/categories/", demoCategoryHandler)
}

func demoProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "name": "Kopi Susu", "price": 15000, "stock": 100, "category_id": 1},
		{"id": 2, "name": "Teh Manis", "price": 8000, "stock": 150, "category_id": 1},
		{"id": 3, "name": "Roti Bakar", "price": 12000, "stock": 50, "category_id": 2},
	})
}

func demoCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{
		{"id": 1, "name": "Beverages", "description": "Various drinks"},
		{"id": 2, "name": "Food", "description": "Various food items"},
	})
}
