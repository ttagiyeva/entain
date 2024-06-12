package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ttagiyeva/entain/internal/constants"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
	"github.com/ttagiyeva/entain/internal/user"
)

// Transaction is a structure which manages transaction usecase.
type Transaction struct {
	log             *slog.Logger
	transactionRepo transaction.Repository
	userRepo        user.Repository
	db              transaction.Database
}

// New creates a new transaction usecase.
func New(log *slog.Logger, r transaction.Repository, u user.Repository, d transaction.Database) *Transaction {
	return &Transaction{
		log:             log,
		transactionRepo: r,
		userRepo:        u,
		db:              d,
	}
}

// Process processes a transaction.
func (t *Transaction) Process(ctx context.Context, tr *model.Transaction) error {
	exist, err := t.transactionRepo.CheckExistance(ctx, tr.TransactionID)
	if err != nil {
		return fmt.Errorf("usecase.transaction.Process.CheckExistance: %w", err)
	}

	if exist {
		return fmt.Errorf("usecase.transaction.Process.Exist: %w", model.ErrorTransactionAlreadyExists)
	}

	user, err := t.userRepo.GetUser(ctx, tr.UserID)
	if err != nil {
		return fmt.Errorf("usecase.transaction.Process.GetUser: %w", err)
	}

	switch tr.State {
	case "win":
		user.Balance += tr.Amount
	case "lost":
		user.Balance -= tr.Amount
	}

	if user.Balance < 0 {
		return fmt.Errorf("usecase.transaction.Process.Balance: %w", model.ErrorInsufficientBalance)
	}

	tx, err := t.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("usecase.transaction.Process.BeginTx: %w", err)
	}

	err = t.userRepo.UpdateUserBalance(tx, ctx, user)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return fmt.Errorf("usecase.transaction.Process.Rollback: %w %w", errTx, err)
		}

		return fmt.Errorf("usecase.transaction.Process.UpdateUserBalance: %w", err)
	}

	trDao := model.TransactionToTransactionDao(tr)

	err = t.transactionRepo.CreateTransaction(tx, ctx, trDao)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return fmt.Errorf("usecase.transaction.Process.Rollback: %w %w", errTx, err)
		}

		return fmt.Errorf("usecase.transaction.Process.CreateTransaction: %w", err)
	}

	err = t.db.Commit(tx)
	if err != nil {
		return fmt.Errorf("usecase.transaction.Process.Commit: %w", err)
	}

	return nil
}

// PostProcess cancels Every N minutes 10 latest odd records.
func (t *Transaction) PostProcess(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * constants.Interval):
				transactions, err := t.transactionRepo.GetLatestOddAndUncancelledTransactions(ctx, 10)
				if err != nil {
					t.log.Error("usecase.transaction.PostProcess.GetLatestTransactions", "error", err)

					continue
				}

				for _, tr := range transactions {
					err := t.transactionRepo.CancelTransaction(ctx, tr.ID)
					if err != nil {
						t.log.Error("usecase.transaction.PostProcess.CancelTransaction", "error", err)

						continue
					}

				}
			}
		}
	}()
}
