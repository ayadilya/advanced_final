package repositories

import (
	"database/sql"
	"pharmacy-store/internal/domain/entities"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]entities.Product, error) {
	var products []entities.Product
	rows, err := r.db.Query("SELECT id, name, description, price, stock, category_id FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product entities.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
