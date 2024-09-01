package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateProduct(*ProductCreate) (*Product, error)
	DeleteProduct(int) error
	UpdateProduct(*Product) error
	GetProductById(int) (*Product, error)
	GetProducts() ([]*Product, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connString := "user=postgres dbname=product_db password=root sslmode=disable"
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.CreateProductTable()
}

func (s *PostgresStorage) CreateProductTable() error {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS products(
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(50),
		description VARCHAR(100),
		price FLOAT,
		date_create TIMESTAMP
	)`

	_, err := s.db.Query(query)
	return err

}

func (s *PostgresStorage) CreateProduct(product *ProductCreate) (*Product, error) {
	query := `
		INSERT INTO products (name, description, price, date_create) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, price, date_create
	`

	prd := &Product{}
	rows, err := s.db.Query(query, product.Name, product.Description, product.Price, product.Date_create)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&prd.Id, &prd.Name, &prd.Description, &prd.Price, &prd.Date_create)
		if err != nil {
			return nil, err
		}
	}

	return prd, nil

}

func (s *PostgresStorage) UpdateProduct(*Product) error {
	return nil
}

func (s *PostgresStorage) GetProductById(id int) (*Product, error) {
	return nil, nil
}

func (s *PostgresStorage) GetProducts() ([]*Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := []*Product{}
	for rows.Next() {
		product, err := scanIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *PostgresStorage) DeleteProduct(id int) error {
	return nil
}

func scanIntoProduct(rows *sql.Rows) (*Product, error) {
	product := new(Product)
	err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Date_create)

	return product, err
}
