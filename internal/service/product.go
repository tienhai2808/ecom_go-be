package service

import "backend/internal/model"

type ProductService interface {
	GetAllProducts() ([]*model.Product, error)
}