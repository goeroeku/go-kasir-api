package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	rows, err := r.db.Query("SELECT id, nama, harga, stok, category_id FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		`SELECT p.id, p.nama, p.harga, p.stok, p.category_id, COALESCE(c.name, '') as category_name 
		 FROM products p 
		 LEFT JOIN categories c ON p.category_id = c.id 
		 WHERE p.id = $1`, id).
		Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID, &p.CategoryName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"INSERT INTO products (nama, harga, stok, category_id) VALUES ($1, $2, $3, $4) RETURNING id, nama, harga, stok, category_id",
		req.Nama, req.Harga, req.Stok, req.CategoryID,
	).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(id int, req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"UPDATE products SET nama = $1, harga = $2, stok = $3, category_id = $4 WHERE id = $5 RETURNING id, nama, harga, stok, category_id",
		req.Nama, req.Harga, req.Stok, req.CategoryID, id,
	).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
