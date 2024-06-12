package usecase

import (
	"context"
	"time"

	"github.com/ttagiyeva/entain/internal/constants"
	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction"
	"github.com/ttagiyeva/entain/internal/user"
	"go.uber.org/zap"
)

// Transaction is a structure which manages transaction usecase.
type Transaction struct {
	log             *zap.SugaredLogger
	transactionRepo transaction.Repository
	userRepo        user.Repository
	db              transaction.Database
}

// New creates a new transaction usecase.
func New(log *zap.SugaredLogger, r transaction.Repository, u user.Repository, d *database.Postgres) *Transaction {
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
		return err
	}

	if exist {
		return model.ErrorTransactionAlreadyExists
	}

	user, err := t.userRepo.GetUser(ctx, tr.UserID)
	if err != nil {
		return err
	}

	switch tr.State {
	case "win":
		user.Balance += tr.Amount
	case "lost":
		user.Balance -= tr.Amount
	}

	if user.Balance < 0 {
		return model.ErrorInsufficientBalance
	}

	tx, err := t.db.BeginTx(ctx)
	if err != nil {
		return model.ErrorInternalServer
	}

	err = t.userRepo.UpdateUserBalance(tx, ctx, user)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return model.ErrorInternalServer
		}

		return err
	}

	trDao := model.TransactionToTransactionDao(tr)

	err = t.transactionRepo.CreateTransaction(tx, ctx, trDao)
	if err != nil {
		errTx := t.db.Rollback(tx)
		if errTx != nil {
			return model.ErrorInternalServer
		}

		return err
	}

	err = t.db.Commit(tx)
	if err != nil {
		return model.ErrorInternalServer
	}

	return nil
}

// PostProcess cancels Every N minutes 10 latest odd records.
func (t *Transaction) PostProcess(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				t.log.Info("transaction usecase post process is stopped")

				return
			case <-time.After(time.Second * constants.Interval):
				transactions, err := t.transactionRepo.GetLatestOddAndUncancelledTransactions(ctx, 10)
				if err != nil {
					t.log.Errorf("error while getting latest odd and uncancelled transactions: %v", err)

					continue
				}

				for _, tr := range transactions {
					err := t.transactionRepo.CancelTransaction(ctx, tr.ID)
					if err != nil {
						t.log.Errorf("error while cancelling transaction: %v", err)

						continue
					}

				}
			}
		}
	}()
}
