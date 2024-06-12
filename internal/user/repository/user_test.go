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
	"github.com/ttagiyeva/entain/internal/user/repository"
)

type userRepoTestSuite struct {
	suite.Suite
	testcontainers.Container
	db   *database.Postgres
	repo *repository.User
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, &userRepoTestSuite{})
}

func (u *userRepoTestSuite) SetupSuite() {
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
	u.Require().NoError(err)

	endpoint, err := container.Endpoint(ctx, "")
	u.Require().NoError(err)

	portStr := strings.Split(endpoint, ":")[1]

	port, err := strconv.Atoi(portStr)
	u.Require().NoError(err)

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

		u.Require().NoError(err)

		break
	}

	u.db = db

	u.repo = repository.New(u.db.Connection)
}

func (u *userRepoTestSuite) SetupTest() {
	if err := u.db.MigrateUp(); err != nil || errors.Is(err, migrate.ErrNoChange) {
		u.Require().NoError(err)
	}
}

func (u *userRepoTestSuite) TearDownTest() {
	u.NoError(u.db.MigrateDown())
}

func (u *userRepoTestSuite) TestCreateTransaction() {
	ctx := context.Background()

	us, err := u.repo.GetUser(ctx, "00000000-0000-0000-0000-000000000001")
	u.NoError(err)
	u.Equal("00000000-0000-0000-0000-000000000001", us.ID)

	us, err = u.repo.GetUser(ctx, faker.UUIDHyphenated())
	u.EqualError(err, fmt.Sprintf("repo.user: %v", model.ErrorNotFound))
	u.Nil(us)
}

func (u *userRepoTestSuite) TestUpdateUserBalance() {
	ctx := context.Background()

	user := &model.UserDao{
		ID:      "00000000-0000-0000-0000-000000000001",
		Balance: 100,
	}

	tx := u.db.Connection.MustBegin().Tx

	err := u.repo.UpdateUserBalance(tx, ctx, user)
	u.NoError(err)
	u.Equal(float32(100), user.Balance)

	user.Balance = -100
	err = u.repo.UpdateUserBalance(tx, ctx, user)
	u.EqualError(err, fmt.Sprintf("repo.user: %v", model.ErrorInsufficientBalance))

	tx = u.db.Connection.MustBegin().Tx
	user.ID = faker.UUIDHyphenated()
	err = u.repo.UpdateUserBalance(tx, ctx, user)
	u.Nil(err)

	tx.Rollback()
}
