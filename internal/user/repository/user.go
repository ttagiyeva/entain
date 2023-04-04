package repository

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/ttagiyeva/entain/internal/model"
)

// User is the repository for users.
type User struct {
	log  *zap.SugaredLogger
	conn *sqlx.DB
}

// New returns a new User object.
func New(log *zap.SugaredLogger, conn *sqlx.DB) *User {
	return &User{
		log:  log,
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
		WHERE id = $1;
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
		a.log.Errorf("error while getting user: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrorNotFound
		}

		return nil, model.ErrorInternalServer
	}

	return user, nil
}

// UpdateUserBalance updates the balance of a user.
func (a *User) UpdateUserBalance(ctx context.Context, user *model.UserDao) error {
	query := `
		UPDATE users
		SET balance = $1
		WHERE id = $2
		RETURNING id, balance;
	`
	tx, err := a.conn.Begin()
	if err != nil {
		a.log.Errorf("error while updating balance: %v", err)

		return model.ErrorInternalServer
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		user.Balance,
		user.ID,
	)

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			return model.ErrorInternalServer
		}

		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Constraint == "check_positive" {
				return model.ErrorInsufficientBalance
			}
		}

		return model.ErrorInternalServer
	}

	err = tx.Commit()
	if err != nil {
		return model.ErrorInternalServer
	}

	return nil
}
