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
	rows, err := r.db.Query("SELECT id, name, price, stock, category_id FROM products ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *ProductRepository) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		`SELECT p.id, p.name, p.price, p.stock, p.category_id, COALESCE(c.name, '') as category_name 
		 FROM products p 
		 LEFT JOIN categories c ON p.category_id = c.id 
		 WHERE p.id = $1`, id).
		Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.CategoryName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Create(req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id, name, price, stock, category_id",
		req.Name, req.Price, req.Stock, req.CategoryID,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(id int, req models.ProductRequest) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		"UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5 RETURNING id, name, price, stock, category_id",
		req.Name, req.Price, req.Stock, req.CategoryID, id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	return err
}
