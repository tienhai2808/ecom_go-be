package repository

import "backend/internal/model"

type ProductRepository interface {
	GetAllProducts() ([]*model.Product, error)
}