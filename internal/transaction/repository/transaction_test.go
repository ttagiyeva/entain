package repository_test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"github.com/ttagiyeva/entain/internal/config"
	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/repository"
)

type transactionRepoTestSuite struct {
	suite.Suite
	testcontainers.Container
	db   *database.Postgres
	repo *repository.Transaction
}

func TestTransactionRepoTestSuite(t *testing.T) {
	suite.Run(t, &transactionRepoTestSuite{})
}

func (t *transactionRepoTestSuite) SetupSuite() {
	ctx := context.Background()

	cfg := config.Config{
		DB: config.DB{
			Host:     "localhost",
			Port:     5432,
			User:     "root",
			Password: "root",
			Name:     "postgres",
		},
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     cfg.DB.User,
			"POSTGRES_PASSWORD": cfg.DB.Password,
			"POSTGRES_DB":       cfg.DB.Name,
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	t.Require().NoError(err)

	endpoint, err := container.Endpoint(ctx, "")
	t.Require().NoError(err)

	portStr := strings.Split(endpoint, ":")[1]

	port, err := strconv.Atoi(portStr)
	t.Require().NoError(err)

	cfg.DB.Port = uint16(port)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in database connection", r)
		}
	}()

	db := database.NewPostgres()

	// 10 iterations to wait for the database to be ready.
	for i := 0; i < 10; i++ {
		err = db.Connect(ctx, &cfg)
		if err != nil {
			time.Sleep(time.Millisecond * 500)

			continue
		}

		t.Require().NoError(err)

		break
	}

	t.db = db

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
	ctx := context.Background()

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

	t.NoError(t.repo.CreateTransaction(tx, ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := t.repo.CreateTransaction(tx, ctx, transaction)
	t.EqualError(err, fmt.Sprintf("repo.transaction.CreateTransaction: %v", model.ErrorTransactionAlreadyExists))
}

func (t *transactionRepoTestSuite) TestCancelTransaction() {
	ctx := context.Background()
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

	t.NoError(t.repo.CreateTransaction(tx, ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	t.NoError(t.repo.CancelTransaction(ctx, transaction.ID))

	err = t.repo.CancelTransaction(ctx, faker.UUIDHyphenated())
	t.Nil(err)
}

func (t *transactionRepoTestSuite) TestCheckExistance() {
	ctx := context.Background()
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

	t.NoError(t.repo.CreateTransaction(tx, ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	ok, err := t.repo.CheckExistance(ctx, transaction.TransactionID)
	t.Equal(true, ok)
	t.NoError(err)

	ok, err = t.repo.CheckExistance(ctx, faker.UUIDHyphenated())
	t.Equal(false, ok)
	t.NoError(err)
}

func (t *transactionRepoTestSuite) TestGetLatestOddAndUncancelledTransactions() {
	ctx := context.Background()
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

	t.NoError(t.repo.CreateTransaction(tx, ctx, transaction))
	t.NotEqual(0, transaction.ID)

	err := tx.Commit()
	t.NoError(err)

	transactions, err := t.repo.GetLatestOddAndUncancelledTransactions(ctx, 10)
	t.NoError(err)
	t.NotEmpty(transactions)
	t.Equal(1, len(transactions))
}
