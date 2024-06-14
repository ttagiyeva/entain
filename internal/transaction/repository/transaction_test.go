package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/repository"
	"github.com/ttagiyeva/entain/internal/util"
)

type transactionRepoTestSuite struct {
	suite.Suite
	testcontainers.Container
	db   *database.Postgres
	repo *repository.Transaction
	ctx  context.Context
}

func TestTransactionRepoTestSuite(t *testing.T) {
	suite.Run(t, &transactionRepoTestSuite{})
}

func (t *transactionRepoTestSuite) SetupSuite() {
	t.ctx = context.Background()
	t.db = util.CreateTestContainer(t.ctx, &t.Suite)
	t.repo = repository.New(t.db.Connection)
}

func (t *transactionRepoTestSuite) SetupTest() {
	if err := t.db.MigrateUp(); err != nil || errors.Is(err, migrate.ErrNoChange) {
		t.Require().NoError(err)
	}
}

func (t *transactionRepoTestSuite) TearDownTest() {
	t.NoError(t.db.MigrateDown())
}

func (t *transactionRepoTestSuite) TestCreateTransaction() {
	transaction := &model.TransactionDao{
		UserID:        "00000000-0000-0000-0000-000000000001",
		TransactionID: faker.UUIDHyphenated(),
		SourceType:    "game",
		State:         "win",
		Amount:        0.0,
		CreatedAt:     time.Now(),
		Cancelled:     false,
	}

	tx := t.db.Connection.MustBegin().Tx
	defer tx.Commit()

	t.NoError(t.repo.CreateTransaction(tx, t.ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := t.repo.CreateTransaction(tx, t.ctx, transaction)
	t.Equal(true, errors.Is(err, model.ErrorTransactionAlreadyExists))
}

func (t *transactionRepoTestSuite) TestCancelTransaction() {
	transaction := &model.TransactionDao{
		UserID:        "00000000-0000-0000-0000-000000000001",
		TransactionID: faker.UUIDHyphenated(),
		SourceType:    "game",
		State:         "win",
		Amount:        0.0,
		CreatedAt:     time.Now(),
		Cancelled:     false,
	}

	tx := t.db.Connection.MustBegin().Tx

	t.NoError(t.repo.CreateTransaction(tx, t.ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	t.NoError(t.repo.CancelTransaction(t.ctx, transaction.ID))

	err = t.repo.CancelTransaction(t.ctx, faker.UUIDHyphenated())
	t.Nil(err)
}

func (t *transactionRepoTestSuite) TestCheckExistance() {
	transaction := &model.TransactionDao{
		UserID:        "00000000-0000-0000-0000-000000000001",
		TransactionID: faker.UUIDHyphenated(),
		SourceType:    "game",
		State:         "win",
		Amount:        0.0,
		CreatedAt:     time.Now(),
		Cancelled:     false,
	}

	tx := t.db.Connection.MustBegin().Tx

	t.NoError(t.repo.CreateTransaction(tx, t.ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	ok, err := t.repo.CheckExistance(t.ctx, transaction.TransactionID)
	t.Equal(true, ok)
	t.NoError(err)

	ok, err = t.repo.CheckExistance(t.ctx, faker.UUIDHyphenated())
	t.Equal(false, ok)
	t.NoError(err)
}

func (t *transactionRepoTestSuite) TestGetLatestOddAndUncancelledTransactions() {
	transaction := &model.TransactionDao{
		UserID:        "00000000-0000-0000-0000-000000000001",
		TransactionID: faker.UUIDHyphenated(),
		SourceType:    "game",
		State:         "win",
		Amount:        0.0,
		CreatedAt:     time.Now(),
		Cancelled:     false,
	}

	tx := t.db.Connection.MustBegin().Tx

	t.NoError(t.repo.CreateTransaction(tx, t.ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	transactions, err := t.repo.GetLatestOddAndUncancelledTransactions(t.ctx, 10)
	t.NoError(err)
	t.NotEmpty(transactions)
	t.Equal(1, len(transactions))
}
