package services

import (
	"kasir/models"
	"kasir/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem, useLock bool) (*models.Transaction, error) {
	return s.repo.CreateTransaction(items, useLock)
}

func (s *TransactionService) GetTransactionSummary(startDate, endDate string) (map[string]interface{}, error) {
	return s.repo.GetTransactionSummary(startDate, endDate)
}
