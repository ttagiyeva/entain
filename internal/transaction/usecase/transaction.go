package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
	"github.com/ttagiyeva/entain/internal/user"
)

const (
	// Interval defines the interval with second for cancel transaction process.
	interval = 1
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
		return fmt.Errorf("failed to check transaction existance: %w", err)
	}

	if exist {
		return fmt.Errorf("failed because the transaction already exists: %w", model.ErrorTransactionAlreadyExists)
	}

	user, err := t.userRepo.GetUser(ctx, tr.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	switch tr.State {
	case "win":
		user.Balance += tr.Amount
	case "lost":
		user.Balance -= tr.Amount
	}

	if user.Balance < 0 {
		return fmt.Errorf("failed because balance of the user is not enough: %w", model.ErrorInsufficientBalance)
	}

	tx, err := t.db.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin a db tx: %w", err)
	}

	err = t.userRepo.UpdateUserBalance(tx, ctx, user)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return fmt.Errorf("failed to rollback the update user balance tx: %w %w", errTx, err)
		}

		return fmt.Errorf("failed to update user balance: %w", err)
	}

	trDao := model.TransactionToTransactionDao(tr)

	err = t.transactionRepo.CreateTransaction(tx, ctx, trDao)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return fmt.Errorf("failed to rollback the create transaction tx: %w %w", errTx, err)
		}

		return fmt.Errorf("failed to create the transaction: %w", err)
	}

	err = t.db.Commit(tx)
	if err != nil {
		return fmt.Errorf("failed to commit the db tx: %w", err)
	}

	return nil
}

// PostProcess cancels odd and uncancelled transactions in every interval.
func (t *Transaction) PostProcess(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * interval):
				transactions, err := t.transactionRepo.GetLatestOddAndUncancelledTransactions(ctx, 10)
				if err != nil {
					t.log.Error("failed to get latest odd and uncancelled transactions", "error", err)

					continue
				}

				for _, tr := range transactions {
					err := t.transactionRepo.CancelTransaction(ctx, tr.ID)
					if err != nil {
						t.log.Error("failed to cancel transaction", "error", err)

						continue
					}
				}
			}
		}
	}()
}
