package usecase

import (
	"context"
	"sync"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/mocks"
	userMocks "github.com/ttagiyeva/entain/internal/user/mocks"
)

func TestProcess(t *testing.T) {
	user := &model.UserDao{
		ID:      gofakeit.UUID(),
		Balance: gofakeit.Float32(),
	}

	tr := &model.Transaction{
		UserID:        user.ID,
		TransactionID: gofakeit.UUID(),
		State:         "win",
		Amount:        gofakeit.Float32(),
	}

	testCases := []struct {
		name          string
		body          *model.Transaction
		buildStubs    func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository)
		checkResponse func(err error)
	}{
		{
			name: "OK",
			body: tr,
			buildStubs: func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				userRepo.EXPECT().UpdateUserBalance(gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			checkResponse: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "Existed transaction",
			body: tr,
			buildStubs: func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(true, nil)
			},
			checkResponse: func(err error) {
				require.Equal(t, model.ErrorTransactionAlreadyExists, err)
			},
		},
		{
			name: "User not found",
			body: tr,
			buildStubs: func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil).Times(1)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(nil, model.ErrorNotFound)

			},
			checkResponse: func(err error) {
				require.Equal(t, model.ErrorNotFound, err)
			},
		},
		{
			name: "Insufficient balance",
			body: tr,
			buildStubs: func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil).Times(1)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(nil, model.ErrorInsufficientBalance)
			},
			checkResponse: func(err error) {
				require.Equal(t, model.ErrorInsufficientBalance, err)
			},
		},
		{
			name: "Internal server",
			body: tr,
			buildStubs: func(trRepo *mocks.MockTransactionRepository, userRepo *userMocks.MockUserRepository) {
				trRepo.EXPECT().CheckExistance(gomock.Any(), tr.TransactionID).Return(false, nil)
				userRepo.EXPECT().GetUser(gomock.Any(), tr.UserID).Return(user, nil)
				userRepo.EXPECT().UpdateUserBalance(gomock.Any(), user).Return(nil).Times(1)
				trRepo.EXPECT().CreateTransaction(gomock.Any(), gomock.Any()).Return(model.ErrorInternalServer)
			},
			checkResponse: func(err error) {
				require.Equal(t, model.ErrorInternalServer, err)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer ctrl.Finish()

			trRepo := mocks.NewMockTransactionRepository(ctrl)
			userRepo := userMocks.NewMockUserRepository(ctrl)

			tc.buildStubs(trRepo, userRepo)

			usecase := New(nil, trRepo, userRepo)
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
		buildStubs func(trRepo *mocks.MockTransactionRepository, wg *sync.WaitGroup)
	}{
		{
			name: "OK",
			buildStubs: func(trRepo *mocks.MockTransactionRepository, wg *sync.WaitGroup) {
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

			trRepo := mocks.NewMockTransactionRepository(ctrl)
			userRepo := userMocks.NewMockUserRepository(ctrl)

			tc.buildStubs(trRepo, &wg)

			usecase := New(zap.NewNop().Sugar(), trRepo, userRepo)
			usecase.PostProcess(ctx)
			wg.Wait()
		})
	}

}
