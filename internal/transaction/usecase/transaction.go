package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/domain"
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
}

// New creates a new transaction usecase.
func New(log *zap.SugaredLogger, r transaction.Repository, u user.Repository) transaction.Usecase {
	return &Transaction{
		log:             log,
		transactionRepo: r,
		userRepo:        u,
	}
}

// Process processes a transaction.
func (t *Transaction) Process(ctx context.Context, tr *model.Transaction) error {
	exist, err := t.transactionRepo.CheckExistance(ctx, tr.TransactionID)
	if err != nil {
		return err
	}

	if !exist {
		user, err := t.userRepo.GetUser(ctx, tr.UserID)
		if err != nil {
			return err
		}

		switch tr.State {
		case model.Win:
			user.Balance += tr.Amount
		case model.Lost:
			user.Balance -= tr.Amount
		}

		if user.Balance < 0 {
			return model.ErrorInsufficientBalance
		}

		err = t.userRepo.UpdateUserBalance(ctx, user)
		if err != nil {
			return err
		}

		transac := &domain.Transaction{
			UserID:        tr.UserID,
			TransactionID: tr.TransactionID,
			SourceType:    fmt.Sprint(tr.SourceType),
			State:         fmt.Sprint(tr.State),
			Amount:        tr.Amount,
		}

		err = t.transactionRepo.CreateTransaction(ctx, transac)
		if err != nil {
			return err
		}

		return nil
	} else {
		return model.ErrorTransactionAlreadyExists
	}

}

// PostProcess cancels Every N minutes 10 latest odd records.
func (t *Transaction) PostProcess(ctx context.Context) {
	go func() {

		for {
			select {
			case <-ctx.Done():
				t.log.Info("transaction usecase post process is stopped")
				return

			case <-time.After(time.Second * config.Interval):

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
			return
		}

	}()
}
