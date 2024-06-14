package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/mocks"
	userMocks "github.com/ttagiyeva/entain/internal/user/mocks"
)

func TestProcess(t *testing.T) {
	user := &model.UserDao{
		ID:      gofakeit.UUID(),
		Balance: 10,
	}

	tx := &sql.Tx{}
	dummyErr := errors.New("dummy error")

	tr := &model.Transaction{
		UserID:        user.ID,
		TransactionID: gofakeit.UUID(),
		State:         "lost",
		Amount:        1,
	}

	testCases := []struct {
		name          string
		body          *model.Transaction
		buildStubs    func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase)
		checkResponse func(err error)
	}{
		{
			name: "OK",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(tx, gomock.Any(), gomock.Any()).Return(nil).Times(1)
				db.EXPECT().Commit(tx).Return(nil).Times(1)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "CheckExistance error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, dummyErr)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "Existed transaction",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(true, nil)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, model.ErrorTransactionAlreadyExists))
			},
		},
		{
			name: "User not found",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil).Times(1)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(nil, model.ErrorUserNotFound)

			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, model.ErrorUserNotFound))
			},
		},
		{
			name: "Insufficient balance",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil).Times(1)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(&model.UserDao{
					ID:      user.ID,
					Balance: 0,
				}, nil)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, model.ErrorInsufficientBalance))
			},
		},
		{
			name: "BeginTx error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(nil, dummyErr)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "UpdateUserBalance error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(dummyErr)
				db.EXPECT().Rollback(tx).Return(nil).Times(1)

			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "Rollback of UpdateUserBalance error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(dummyErr)
				db.EXPECT().Rollback(tx).Return(errors.New("rollback error")).Times(1)

			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "CreateTransaction error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(tx, gomock.Any(), gomock.Any()).Return(dummyErr)
				db.EXPECT().Rollback(tx).Return(nil).Times(1)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "Rollback of CreateTransaction error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(tx, gomock.Any(), gomock.Any()).Return(dummyErr)
				db.EXPECT().Rollback(tx).Return(errors.New("rollback error"))
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
		{
			name: "Commit error",
			body: tr,
			buildStubs: func(trRepo *mocks.MockRepository, userRepo *userMocks.MockUserRepository, db *mocks.MockDatabase) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				db.EXPECT().BeginTx(gomock.Any()).Return(tx, nil).Times(1)
				userRepo.EXPECT().UpdateUserBalance(tx, gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(tx, gomock.Any(), gomock.Any()).Return(nil).Times(1)
				db.EXPECT().Commit(tx).Return(dummyErr)
			},
			checkResponse: func(err error) {
				require.Equal(t, true, errors.Is(err, dummyErr))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			trRepo := mocks.NewMockRepository(ctrl)
			userRepo := userMocks.NewMockUserRepository(ctrl)
			db := mocks.NewMockDatabase(ctrl)

			tc.buildStubs(trRepo, userRepo, db)

			usecase := New(nil, trRepo, userRepo, db)
			err := usecase.Process(context.Background(), tr)

			tc.checkResponse(err)

		})
	}

}

func TestPostProcess(t *testing.T) {
	transactions := []*model.TransactionDao{
		{
			ID:            gofakeit.UUID(),
			UserID:        gofakeit.UUID(),
			TransactionID: gofakeit.UUID(),
			State:         "win",
			Amount:        gofakeit.Float32(),
		},
	}

	testCases := []struct {
		name       string
		buildStubs func(trRepo *mocks.MockRepository, wg *sync.WaitGroup)
	}{
		{
			name: "OK",
			buildStubs: func(trRepo *mocks.MockRepository, wg *sync.WaitGroup) {
				trRepo.EXPECT().GetLatestOddAndUncancelledTransactions(gomock.Any(), gomock.Any()).Return(transactions, nil).Do(func(arg0, ar1 interface{}) {
					defer wg.Done()
				})
				trRepo.EXPECT().CancelTransaction(gomock.Any(), gomock.Any()).Return(nil).Do(func(arg0, ar1 interface{}) {
					defer wg.Done()
				})
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			ctrl := gomock.NewController(t)

			wg := sync.WaitGroup{}
			wg.Add(2)

			defer ctrl.Finish()

			trRepo := mocks.NewMockRepository(ctrl)
			userRepo := userMocks.NewMockUserRepository(ctrl)
			db := mocks.NewMockDatabase(ctrl)

			tc.buildStubs(trRepo, &wg)

			usecase := New(slog.Default(), trRepo, userRepo, db)
			usecase.PostProcess(ctx)
			wg.Wait()
		})
	}
}
