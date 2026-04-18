package service

import (
	"qr-traceability/internal/models"
	"qr-traceability/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetProductByID(id uint) (*models.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductService) GetProductByCode(code string) (*models.Product, error) {
	return s.repo.GetByCode(code)
}

func (s *ProductService) ListProducts(limit, offset int) ([]models.Product, error) {
	return s.repo.List(limit, offset)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}
