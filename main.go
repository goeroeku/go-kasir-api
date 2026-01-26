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

// In-memory storage
var produkList []Produk
var nextID = 1

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

func main() {
	// Inisialisasi data awal (opsional)
	produkList = []Produk{
		{ID: 1, Nama: "Kopi Susu", Harga: 15000, Stok: 100},
		{ID: 2, Nama: "Teh Manis", Harga: 8000, Stok: 150},
		{ID: 3, Nama: "Roti Bakar", Harga: 12000, Stok: 50},
	}
	nextID = 4

	// Setup routing
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkHandler)

	// Start server
	port := "8080"
	fmt.Printf("🚀 Kasir API berjalan di http://localhost:%s\n", port)
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET  /health     - Health check")
	fmt.Println("   GET  /produk     - Get all produk")
	fmt.Println("   POST /produk     - Create produk")
	fmt.Println("   GET  /produk/:id - Get produk by ID")
	fmt.Println("   PUT  /produk/:id - Update produk")
	fmt.Println("   DELETE /produk/:id - Delete produk")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
