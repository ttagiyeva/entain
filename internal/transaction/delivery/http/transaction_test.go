package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/mocks"
)

// TestTransactionHandler_Process tests the transaction handler process method.
func TestTransactionHandler_Process(t *testing.T) {
	testCases := []struct {
		name          string
		body          []byte
		SourceType    string
		buildStubs    func(trUsecase *mocks.MockUsecase)
		expectedCode  int
		expectedError model.Error
	}{
		{
			name: "OK",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:       "Invalid request body",
			body:       []byte(`{"transactionId":"1","state":"win","amount":"1"}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {},
			expectedError: model.Error{
				Code:    http.StatusBadRequest,
				Message: model.ErrorBadRequest,
			},
		},
		{
			name:       "Invalid source type",
			body:       []byte(`{"transactionId":"1","state":"win","amount":1}`),
			SourceType: "test",
			buildStubs: func(trUsecase *mocks.MockUsecase) {},
			expectedError: model.Error{
				Code:    http.StatusBadRequest,
				Message: "Value of the SourceType field must be one of 'game server payment'",
			},
		},
		{
			name:       "Invalid state",
			body:       []byte(`{"transactionId":"1","state":"won","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {},
			expectedError: model.Error{
				Code:    http.StatusBadRequest,
				Message: "Value of the State field must be one of 'win lost'",
			},
		},
		{
			name:       "Invalid amount",
			body:       []byte(`{"transactionId":"1","state":"win","amount":-1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {},
			expectedError: model.Error{
				Code:    http.StatusBadRequest,
				Message: "Value of the Amount field must be greater than 0",
			},
		},
		{
			name:       "Invalid transactionId",
			body:       []byte(`{"state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {},
			expectedError: model.Error{
				Code:    http.StatusBadRequest,
				Message: "TransactionID field is required",
			},
		},
		{
			name: "Transaction already exists",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(model.ErrorTransactionAlreadyExists)
			},
			expectedError: getError(model.ErrorTransactionAlreadyExists),
		},
		{
			name: "Insufficient Balance error",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(model.ErrorInsufficientBalance)
			},
			expectedError: getError(model.ErrorInsufficientBalance),
		},
		{
			name: "User not found",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(model.ErrorUserNotFound)
			},
			expectedError: getError(model.ErrorUserNotFound),
		},
		{
			name: "Internal server error",
			body: []byte(`{"transactionId":"0","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(fmt.Errorf("unexpected error"))
			},
			expectedError: getError(fmt.Errorf("unexpected error")),
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			trUsecase := mocks.NewMockUsecase(ctrl)
			tc.buildStubs(trUsecase)

			handler := NewHandler(slog.Default(), trUsecase)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users/1/transactions", bytes.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.SourceType == "" {
				req.Header.Set(SourceType, "game")
			} else {
				req.Header.Set(SourceType, tc.SourceType)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues("1")

			err := handler.Process(c)
			require.NoError(t, err)

			if tc.expectedCode != 0 {
				require.Equal(t, tc.expectedCode, rec.Code)
			} else {
				require.Equal(t, tc.expectedError.Code, rec.Code)

				expectedErr, err := json.Marshal(tc.expectedError)
				require.NoError(t, err)

				require.JSONEq(t, rec.Body.String(), string(expectedErr))
			}
		})
	}
}
