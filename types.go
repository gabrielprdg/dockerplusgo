package main

import (
	"time"

	"github.com/google/uuid"
)

// `json:type` is a struct tag
type CreateRequestProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Product struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Date_create time.Time `json:"date_create"`
}

type ProductCreate struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Date_create time.Time `json:"date_create"`
}

func NewProduct(name string, desc string, price float64) *ProductCreate {
	return &ProductCreate{
		Name:        name,
		Description: desc,
		Price:       price,
		Date_create: time.Now().UTC(),
	}
}
