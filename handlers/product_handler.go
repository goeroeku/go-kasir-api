package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"kasir-api/models"
	"kasir-api/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// GetAll godoc
// @Summary Get all produk
// @Description Mengambil semua daftar produk
// @Tags Produk
// @Produce json
// @Success 200 {array} models.Product
// @Router /produk [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Create godoc
// @Summary Create produk
// @Description Membuat produk baru
// @Tags Produk
// @Accept json
// @Produce json
// @Param produk body models.ProductRequest true "Data produk"
// @Success 201 {object} models.Product
// @Router /produk [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.Create(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// GetByID godoc
// @Summary Get produk by ID
// @Description Mengambil produk berdasarkan ID
// @Tags Produk
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} models.Product
// @Failure 404 {string} string "Produk tidak ditemukan"
// @Router /produk/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(id)
	if err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Update godoc
// @Summary Update produk
// @Description Mengupdate produk berdasarkan ID
// @Tags Produk
// @Accept json
// @Produce json
// @Param id path int true "Produk ID"
// @Param produk body models.ProductRequest true "Data produk"
// @Success 200 {object} models.Product
// @Router /produk/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req models.ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.service.Update(id, req)
	if err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Delete godoc
// @Summary Delete produk
// @Description Menghapus produk berdasarkan ID
// @Tags Produk
// @Produce json
// @Param id path int true "Produk ID"
// @Success 200 {object} map[string]string
// @Router /produk/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path)
	if id == 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		http.Error(w, "Produk tidak ditemukan", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Produk dengan ID %d berhasil dihapus", id),
	})
}

// Handler routes requests to appropriate method handlers
func (h *ProductHandler) Handler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")

	if len(pathParts) == 2 || (len(pathParts) == 3 && pathParts[2] == "") {
		switch r.Method {
		case http.MethodGet:
			h.GetAll(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		switch r.Method {
		case http.MethodGet:
			h.GetByID(w, r)
		case http.MethodPut:
			h.Update(w, r)
		case http.MethodDelete:
			h.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func extractID(path string) int {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return id
}
