package http

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"github.com/ttagiyeva/entain/internal/constants"
	"github.com/ttagiyeva/entain/internal/model"
	"github.com/ttagiyeva/entain/internal/transaction/mocks"
	"go.uber.org/zap"
)

// TestTransactionHandler_Process tests the transaction handler process method.
func TestTransactionHandler_Process(t *testing.T) {
	testCases := []struct {
		name          string
		body          []byte
		buildStubs    func(trUsecase *mocks.MockUsecase)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid request body",
			body: []byte(`{"transactionId":"1","state":"win","amount":"1"}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid state",
			body: []byte(`{"transactionId":"1","state":"won","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid amount",
			body: []byte(`{"transactionId":"1","state":"win","amount":-1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Transaction already exists",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(model.ErrorTransactionAlreadyExists)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusConflict, recorder.Code)
			},
		},
		{
			name: "User not found",
			body: []byte(`{"transactionId":"1","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(model.ErrorNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "Internal server error",
			body: []byte(`{"transactionId":"0","state":"win","amount":1}`),
			buildStubs: func(trUsecase *mocks.MockUsecase) {
				trUsecase.EXPECT().Process(gomock.Any(), gomock.Any()).Return(fmt.Errorf("unexpected error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			trUsecase := mocks.NewMockUsecase(ctrl)
			tc.buildStubs(trUsecase)

			handler := NewHandler(zap.NewNop().Sugar(), trUsecase)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/users/1/transactions", bytes.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(constants.SourceType, "game")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler.Process(c)

			tc.checkResponse(rec)
		})
	}
}
