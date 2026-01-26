package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Kasir API
// @version 1.0
// @description REST API untuk sistem kasir (Point of Sale)
// @host localhost:8080
// @BasePath /

// Produk adalah struct untuk data produk
type Produk struct {
	ID    int    `json:"id" example:"1"`
	Nama  string `json:"nama" example:"Kopi Susu"`
	Harga int    `json:"harga" example:"15000"`
	Stok  int    `json:"stok" example:"100"`
}

// Category adalah struct untuk data kategori
type Category struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name" example:"Minuman"`
	Description string `json:"description" example:"Berbagai jenis minuman"`
}

// MessageResponse adalah struct untuk response message
type MessageResponse struct {
	Message string `json:"message" example:"Operation successful"`
}

// HealthResponse adalah struct untuk health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Kasir API is running"`
}

// In-memory storage
var produkList []Produk
var nextID = 1
var categoryList []Category
var nextCategoryID = 1

// WelcomeResponse adalah struct untuk response welcome
type WelcomeResponse struct {
	Name        string   `json:"name" example:"Kasir API"`
	Version     string   `json:"version" example:"1.0.0"`
	Description string   `json:"description" example:"REST API untuk sistem kasir"`
	Docs        string   `json:"docs" example:"/swagger/index.html"`
	Endpoints   []string `json:"endpoints"`
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
		Version:     "1.0.0",
		Description: "REST API untuk sistem kasir (Point of Sale)",
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

// getAllProdukHandler godoc
// @Summary Get all produk
// @Description Mengambil semua daftar produk
// @Tags Produk
// @Produce json
// @Success 200 {array} Produk
// @Router /produk [get]
func getAllProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produkList)
}

// createProdukHandler godoc
// @Summary Create produk
// @Description Membuat produk baru
// @Tags Produk
// @Accept json
// @Produce json
// @Param produk body Produk true "Data produk"
// @Success 201 {object} Produk
// @Failure 400 {string} string "Invalid request body"
// @Router /produk [post]
func createProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var produk Produk
	if err := json.NewDecoder(r.Body).Decode(&produk); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	produk.ID = nextID
	nextID++
	produkList = append(produkList, produk)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produk)
}

// getProdukByIDHandler godoc
// @Summary Get produk by ID
// @Description Mengambil produk berdasarkan ID
// @Tags Produk
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} Produk
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /produk/{id} [get]
func getProdukByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, produk := range produkList {
		if produk.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// updateProdukHandler godoc
// @Summary Update produk
// @Description Mengupdate produk berdasarkan ID
// @Tags Produk
// @Accept json
// @Produce json
// @Param id path int true "Produk ID"
// @Param produk body Produk true "Data produk"
// @Success 200 {object} Produk
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /produk/{id} [put]
func updateProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedProduk Produk
	if err := json.NewDecoder(r.Body).Decode(&updatedProduk); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, produk := range produkList {
		if produk.ID == id {
			updatedProduk.ID = id
			produkList[i] = updatedProduk
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedProduk)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// deleteProdukHandler godoc
// @Summary Delete produk
// @Description Menghapus produk berdasarkan ID
// @Tags Produk
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /produk/{id} [delete]
func deleteProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, produk := range produkList {
		if produk.ID == id {
			produkList = append(produkList[:i], produkList[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageResponse{
				Message: fmt.Sprintf("Produk dengan ID %d berhasil dihapus", id),
			})
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

func produkHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) == 2 || (len(pathParts) == 3 && pathParts[2] == "") {
		switch r.Method {
		case http.MethodGet:
			getAllProdukHandler(w, r)
		case http.MethodPost:
			createProdukHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			getProdukByIDHandler(w, r)
		case http.MethodPut:
			updateProdukHandler(w, r)
		case http.MethodDelete:
			deleteProdukHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// ==================== CATEGORY HANDLERS ====================

// getAllCategoriesHandler godoc
// @Summary Get all categories
// @Description Mengambil semua daftar kategori
// @Tags Categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]
func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryList)
}

// createCategoryHandler godoc
// @Summary Create category
// @Description Membuat kategori baru
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body Category true "Data kategori"
// @Success 201 {object} Category
// @Failure 400 {string} string "Invalid request body"
// @Router /categories [post]
func createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	category.ID = nextCategoryID
	nextCategoryID++
	categoryList = append(categoryList, category)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// getCategoryByIDHandler godoc
// @Summary Get category by ID
// @Description Mengambil kategori berdasarkan ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 404 {string} string "Category tidak ditemukan"
// @Router /categories/{id} [get]
func getCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for _, category := range categoryList {
		if category.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// updateCategoryHandler godoc
// @Summary Update category
// @Description Mengupdate kategori berdasarkan ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Data kategori"
// @Success 200 {object} Category
// @Failure 404 {string} string "Category tidak ditemukan"
// @Router /categories/{id} [put]
func updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedCategory Category
	if err := json.NewDecoder(r.Body).Decode(&updatedCategory); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i, category := range categoryList {
		if category.ID == id {
			updatedCategory.ID = id
			categoryList[i] = updatedCategory
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedCategory)
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// deleteCategoryHandler godoc
// @Summary Delete category
// @Description Menghapus kategori berdasarkan ID
// @Tags Categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {string} string "Category tidak ditemukan"
// @Router /categories/{id} [delete]
func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, category := range categoryList {
		if category.ID == id {
			categoryList = append(categoryList[:i], categoryList[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(MessageResponse{
				Message: fmt.Sprintf("Category dengan ID %d berhasil dihapus", id),
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) == 2 || (len(pathParts) == 3 && pathParts[2] == "") {
		switch r.Method {
		case http.MethodGet:
			getAllCategoriesHandler(w, r)
		case http.MethodPost:
			createCategoryHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			getCategoryByIDHandler(w, r)
		case http.MethodPut:
			updateCategoryHandler(w, r)
		case http.MethodDelete:
			deleteCategoryHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func main() {
	// Inisialisasi data awal produk
	produkList = []Produk{
		{ID: 1, Nama: "Kopi Susu", Harga: 15000, Stok: 100},
		{ID: 2, Nama: "Teh Manis", Harga: 8000, Stok: 150},
		{ID: 3, Nama: "Roti Bakar", Harga: 12000, Stok: 50},
	}
	nextID = 4

	// Inisialisasi data awal kategori
	categoryList = []Category{
		{ID: 1, Name: "Minuman", Description: "Berbagai jenis minuman"},
		{ID: 2, Name: "Makanan", Description: "Berbagai jenis makanan"},
	}
	nextCategoryID = 3

	// Setup routing
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkHandler)
	http.HandleFunc("/categories", categoryHandler)
	http.HandleFunc("/categories/", categoryHandler)

	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Start server
	port := "8080"
	fmt.Printf("🚀 Kasir API berjalan di http://localhost:%s\n", port)
	fmt.Printf("📚 Swagger UI: http://localhost:%s/swagger/index.html\n", port)
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET  /health          - Health check")
	fmt.Println("   -- Produk --")
	fmt.Println("   GET  /produk          - Get all produk")
	fmt.Println("   POST /produk          - Create produk")
	fmt.Println("   GET  /produk/:id      - Get produk by ID")
	fmt.Println("   PUT  /produk/:id      - Update produk")
	fmt.Println("   DELETE /produk/:id    - Delete produk")
	fmt.Println("   -- Categories --")
	fmt.Println("   GET  /categories      - Get all categories")
	fmt.Println("   POST /categories      - Create category")
	fmt.Println("   GET  /categories/:id  - Get category by ID")
	fmt.Println("   PUT  /categories/:id  - Update category")
	fmt.Println("   DELETE /categories/:id - Delete category")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
