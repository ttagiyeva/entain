package transaction

import (
	"context"

	"github.com/ttagiyeva/entain/internal/model"
)

// Repository is a repository interface for transaction.
//
//go:generate mockgen -source ./repository.go -mock_names Repository=MockTransactionRepository -package mocks -destination mocks/transactionRepository.mock.gen.go
type Repository interface {
	CreateTransaction(context.Context, *model.TransactionDao) error
	CancelTransaction(ctx context.Context, id string) error
	CheckExistance(ctx context.Context, id string) (bool, error)
	GetLatestOddAndUncancelledTransactions(ctx context.Context, limit int) ([]*model.TransactionDao, error)
}
