package transaction

import (
	"context"

	"github.com/ttagiyeva/entain/internal/model"
)

// Usecase is a usecase interface for transaction.
//go:generate mockgen -source ./usecase.go -mock_names Repository=MockTransactionUsecase -package mocks -destination ../mocks/transactionUsecase.mock.gen.go
type Usecase interface {
	Process(context.Context, *model.Transaction) error
	PostProcess(ctx context.Context)
}
