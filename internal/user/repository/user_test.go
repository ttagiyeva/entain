package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"

	"github.com/ttagiyeva/entain/internal/database"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/user/repository"
	"github.com/ttagiyeva/entain/internal/util"
)

type userRepoTestSuite struct {
	suite.Suite
	testcontainers.Container
	db   *database.Postgres
	repo *repository.User
	ctx  context.Context
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, &userRepoTestSuite{})
}

func (u *userRepoTestSuite) SetupSuite() {
	u.ctx = context.Background()
	u.db = util.CreateTestContainer(u.ctx, &u.Suite)
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
	us, err := u.repo.GetUser(u.ctx, "00000000-0000-0000-0000-000000000001")
	u.NoError(err)
	u.Equal("00000000-0000-0000-0000-000000000001", us.ID)

	us, err = u.repo.GetUser(u.ctx, faker.UUIDHyphenated())
	u.Equal(true, errors.Is(err, model.ErrorUserNotFound))
	u.Nil(us)
}

func (u *userRepoTestSuite) TestUpdateUserBalance() {
	user := &model.UserDao{
		ID:      "00000000-0000-0000-0000-000000000001",
		Balance: 100,
	}

	tx := u.db.Connection.MustBegin().Tx

	err := u.repo.UpdateUserBalance(tx, u.ctx, user)
	u.NoError(err)
	u.Equal(float32(100), user.Balance)

	user.Balance = -100
	err = u.repo.UpdateUserBalance(tx, u.ctx, user)
	u.Equal(true, errors.Is(err, model.ErrorInsufficientBalance))

	tx = u.db.Connection.MustBegin().Tx
	user.ID = faker.UUIDHyphenated()
	err = u.repo.UpdateUserBalance(tx, u.ctx, user)
	u.Nil(err)

	err = tx.Rollback()
	u.NoError(err)
}
