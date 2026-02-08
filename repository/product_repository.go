package repository

import (
	"database/sql"
	"errors"
	"kasir-api/model"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAllProduct(name string) ([]model.Product, error) {
	query := "SELECT id, name, price, stock FROM products"
	args := []interface{}{}

	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (repo *ProductRepository) GetProductById(id int) (*model.Product, error) {
	query := "SELECT id, name, price, stock FROM products WHERE id = $1"
	row := repo.db.QueryRow(query, id)

	var p model.Product
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) CreateProduct(product *model.Product) error {
	query := "INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id"
	row := repo.db.QueryRow(query, product.Name, product.Price, product.Stock)

	err := row.Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *ProductRepository) UpdateProduct(product *model.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Product wasn't found")
	}

	return nil
}

func (repo *ProductRepository) DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	row, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return errors.New("Product wasn't found")
	}

	return err
}
