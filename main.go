package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Produk adalah struct untuk data produk
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// Category adalah struct untuk data kategori
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage
var produkList []Produk
var nextID = 1
var categoryList []Category
var nextCategoryID = 1

// healthHandler - endpoint untuk cek kesehatan API
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"message": "Kasir API is running",
	})
}

// getAllProdukHandler - mengambil semua produk
func getAllProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produkList)
}

// createProdukHandler - membuat produk baru
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

// getProdukByIDHandler - mengambil produk berdasarkan ID
func getProdukByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ekstrak ID dari URL path
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

	// Cari produk berdasarkan ID
	for _, produk := range produkList {
		if produk.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// updateProdukHandler - mengupdate produk berdasarkan ID
func updateProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ekstrak ID dari URL path
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

	// Cari dan update produk
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

// deleteProdukHandler - menghapus produk berdasarkan ID
func deleteProdukHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ekstrak ID dari URL path
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

	// Cari dan hapus produk
	for i, produk := range produkList {
		if produk.ID == id {
			produkList = append(produkList[:i], produkList[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Produk dengan ID %d berhasil dihapus", id),
			})
			return
		}
	}

	http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
}

// produkHandler - handler untuk /produk endpoint dengan berbagai method
func produkHandler(w http.ResponseWriter, r *http.Request) {
	// Cek apakah ada ID di path
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) == 2 || (len(pathParts) == 3 && pathParts[2] == "") {
		// /produk - tanpa ID
		switch r.Method {
		case http.MethodGet:
			getAllProdukHandler(w, r)
		case http.MethodPost:
			createProdukHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		// /produk/{id} - dengan ID
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

// getAllCategoriesHandler - mengambil semua kategori
func getAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categoryList)
}

// createCategoryHandler - membuat kategori baru
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

// getCategoryByIDHandler - mengambil kategori berdasarkan ID
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

// updateCategoryHandler - mengupdate kategori berdasarkan ID
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

// deleteCategoryHandler - menghapus kategori berdasarkan ID
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
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Category dengan ID %d berhasil dihapus", id),
			})
			return
		}
	}

	http.Error(w, "Category tidak ditemukan", http.StatusNotFound)
}

// categoryHandler - handler untuk /categories endpoint dengan berbagai method
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
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkHandler)
	http.HandleFunc("/categories", categoryHandler)
	http.HandleFunc("/categories/", categoryHandler)

	// Start server
	port := "8080"
	fmt.Printf("🚀 Kasir API berjalan di http://localhost:%s\n", port)
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
