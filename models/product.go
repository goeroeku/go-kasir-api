package models

// Product represents a product in the kasir system
type Product struct {
	ID           int    `json:"id"`
	Nama         string `json:"nama"`
	Harga        int    `json:"harga"`
	Stok         int    `json:"stok"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"`
}

// ProductRequest is used for create/update operations
type ProductRequest struct {
	Nama       string `json:"nama"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
	CategoryID int    `json:"category_id"`
}
