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
	rows, err := r.db.Query("SELECT id, nama, harga, stok FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow("SELECT id, nama, harga, stok FROM products WHERE id = $1", id).
		Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"INSERT INTO products (nama, harga, stok) VALUES ($1, $2, $3) RETURNING id, nama, harga, stok",
		req.Nama, req.Harga, req.Stok,
	).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(id int, req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"UPDATE products SET nama = $1, harga = $2, stok = $3 WHERE id = $4 RETURNING id, nama, harga, stok",
		req.Nama, req.Harga, req.Stok, id,
	).Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
