package transaction

import (
	"context"

	"github.com/ttagiyeva/entain/internal/domain"
)

// Repository is a repository interface for transaction.
//go:generate mockgen -source ./repository.go -mock_names Repository=MockTransactionRepository -package mocks -destination ../mocks/transactionRepository.mock.gen.go
type Repository interface {
	CreateTransaction(context.Context, *domain.Transaction) error
	CancelTransaction(ctx context.Context, id string) error
	CheckExistance(ctx context.Context, id string) (bool, error)
	GetLatestOddAndUncancelledTransactions(ctx context.Context, limit int) ([]*domain.Transaction, error)
}
