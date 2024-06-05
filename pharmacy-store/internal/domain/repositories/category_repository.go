package repositories

import (
	"database/sql"
	"pharmacy-store/internal/domain/entities"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]entities.Category, error) {
	var categories []entities.Category
	rows, err := r.db.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *CategoryRepository) GetByID(id int) (*entities.Category, error) {
	var category entities.Category
	err := r.db.QueryRow("SELECT id, name FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Create(category *entities.Category) error {
	return r.db.QueryRow("INSERT INTO categories (name) VALUES ($1) RETURNING id", category.Name).Scan(&category.ID)
}

func (r *CategoryRepository) Update(category *entities.Category) error {
	_, err := r.db.Exec("UPDATE categories SET name = $1 WHERE id = $2", category.Name, category.ID)
	return err
}

func (r *CategoryRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	return err
}
