package models

// Category represents a category in the kasir system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryRequest is used for create/update operations
type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
