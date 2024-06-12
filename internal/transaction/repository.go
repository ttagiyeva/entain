package transaction

import (
	"context"
	"database/sql"

	"github.com/ttagiyeva/entain/internal/model"
)

// Repository is a repository interface for transaction.
//
//go:generate mockgen -source ./repository.go -package mocks -destination mocks/transactionRepository.mock.gen.go
type Repository interface {
	CreateTransaction(tx *sql.Tx, ctx context.Context, tr *model.TransactionDao) error
	CancelTransaction(ctx context.Context, id string) error
	CheckExistance(ctx context.Context, id string) (bool, error)
	GetLatestOddAndUncancelledTransactions(ctx context.Context, limit int) ([]*model.TransactionDao, error)
}

//go:generate mockgen -source ./repository.go -package mocks -destination mocks/transactionRepository.mock.gen.go
type Database interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	Rollback(tx *sql.Tx) error
	Commit(tx *sql.Tx) error
}
