package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id).
		Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(req models.CategoryRequest) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description",
		req.Name, req.Description,
	).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(id int, req models.CategoryRequest) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		"UPDATE categories SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description",
		req.Name, req.Description, id,
	).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
