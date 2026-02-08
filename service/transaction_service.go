package service

import (
	"kasir-api/model"
	"kasir-api/repository"
)

type TransactionService struct {
	repo *repository.TransactionRepository
}

func NewTransactionService(repo *repository.TransactionRepository) *TransactionService{
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CheckoutItem(items []model.CheckoutItem) (*model.Transaction, error) {
	return s.repo.CheckoutItem(items)
}