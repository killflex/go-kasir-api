package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAllProduct() ([]model.Product, error) {
	return s.repo.GetAllProduct()
}

func (s *ProductService) GetProductById(id int) (*model.Product, error) {
	return s.repo.GetProductById(id)
}

func (s *ProductService) CreateProduct(data *model.Product) error {
	return s.repo.CreateProduct(data)
}

func (s *ProductService) UpdateProduct(product *model.Product) error {
	return s.repo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.repo.DeleteProduct(id)
}
