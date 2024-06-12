package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/ttagiyeva/entain/internal/model"
)

// Transaction is a structure which manages transaction repository.
type Transaction struct {
	conn *sqlx.DB
}

// New returns a new Transaction object.
func New(conn *sqlx.DB) *Transaction {
	return &Transaction{
		conn: conn,
	}
}

// CreateTransaction creates a new transaction.
func (t *Transaction) CreateTransaction(tx *sql.Tx, ctx context.Context, transaction *model.TransactionDao) error {
	query := `
		INSERT INTO transactions ( 
			user_id,
			transaction_id,
			source_type,
			state,
			amount
		) VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		transaction.UserID,
		transaction.TransactionID,
		transaction.SourceType,
		transaction.State,
		transaction.Amount,
	).Scan(&transaction.ID)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Constraint == "unique_transaction_id" {
				return fmt.Errorf("repo.transaction: %w", model.ErrorTransactionAlreadyExists)
			}
		}

		return fmt.Errorf("repo.transaction: %w", err)
	}

	return nil
}

// CancelTransaction cancels a transaction by id.
func (t *Transaction) CancelTransaction(ctx context.Context, id string) error {
	query := `
		UPDATE transactions
		SET cancelled = true, cancelled_at = NOW()
		WHERE id = $1
		RETURNING id;
	`
	tx, err := t.conn.Begin()
	if err != nil {
		return fmt.Errorf("repo.transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return fmt.Errorf("repo.transaction: %w %w", errTx, err)
		}

		return fmt.Errorf("repo.transaction: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("repo.transaction: %w", err)
	}

	return nil
}

// CheckExistance checks existance of transaction in database
func (t *Transaction) CheckExistance(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM transactions
			WHERE transaction_id = $1
		);
	`
	var exists bool

	err := t.conn.QueryRowContext(
		ctx,
		query,
		&id,
	).Scan(
		&exists,
	)

	if err != nil {
		return false, fmt.Errorf("repo.transaction: %w", err)
	}

	return exists, nil
}

// GetLatestOddAndUncancelledTransactions returns the latest odd transactions with a limit.
// Odd records definition was unclear, it means odd amount, transactionId or id etc., so I haven't implemented it.
func (t *Transaction) GetLatestOddAndUncancelledTransactions(ctx context.Context, limit int) ([]*model.TransactionDao, error) {
	query := `
		SELECT id,
			user_id,
			transaction_id,
			source_type,
			state,
			amount,
			created_at,
			cancelled 
		 FROM transactions
			WHERE cancelled = false
			ORDER BY created_at DESC
			LIMIT $1
	`
	rows, err := t.conn.QueryContext(
		ctx,
		query,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("repo.transaction: %w", err)
	}

	defer rows.Close()

	transactions := []*model.TransactionDao{}

	for rows.Next() {
		transaction := &model.TransactionDao{}
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.TransactionID,
			&transaction.SourceType,
			&transaction.State,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.Cancelled,
		)
		if err != nil {
			return nil, fmt.Errorf("repo.transaction: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
