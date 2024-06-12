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

// User is the repository for users.
type User struct {
	conn *sqlx.DB
}

// New returns a new User object.
func New(conn *sqlx.DB) *User {
	return &User{
		conn: conn,
	}
}

// GetUser returns a user by id.
func (a *User) GetUser(ctx context.Context, id string) (*model.UserDao, error) {
	query := `
		SELECT
			id,
			balance
		FROM users
		WHERE id = $1 FOR UPDATE;
	`
	user := &model.UserDao{}

	err := a.conn.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Balance,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("repo.user.GetUser: %w", model.ErrorNotFound)
		}

		return nil, fmt.Errorf("repo.user.GetUser: %w", err)
	}

	return user, nil
}

// UpdateUserBalance updates the balance of a user.
func (a *User) UpdateUserBalance(tx *sql.Tx, ctx context.Context, user *model.UserDao) error {
	query := `
		UPDATE users
		SET balance = $1
		WHERE id = $2
		RETURNING id, balance;
	`
	_, err := tx.ExecContext(
		ctx,
		query,
		user.Balance,
		user.ID,
	)

	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Constraint == "users_balance_check" {
				return fmt.Errorf("repo.user.UpdateUserBalance: %w", model.ErrorInsufficientBalance)
			}
		}

		return fmt.Errorf("repo.user.UpdateUserBalance: %w", err)
	}

	return nil
}
