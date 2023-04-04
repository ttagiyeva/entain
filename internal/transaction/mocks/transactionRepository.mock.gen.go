// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ttagiyeva/entain/internal/model"
)

// MockTransactionRepository is a mock of Repository interface.
type MockTransactionRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionRepositoryMockRecorder
}

// MockTransactionRepositoryMockRecorder is the mock recorder for MockTransactionRepository.
type MockTransactionRepositoryMockRecorder struct {
	mock *MockTransactionRepository
}

// NewMockTransactionRepository creates a new mock instance.
func NewMockTransactionRepository(ctrl *gomock.Controller) *MockTransactionRepository {
	mock := &MockTransactionRepository{ctrl: ctrl}
	mock.recorder = &MockTransactionRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTransactionRepository) EXPECT() *MockTransactionRepositoryMockRecorder {
	return m.recorder
}

// CancelTransaction mocks base method.
func (m *MockTransactionRepository) CancelTransaction(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelTransaction", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// CancelTransaction indicates an expected call of CancelTransaction.
func (mr *MockTransactionRepositoryMockRecorder) CancelTransaction(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).CancelTransaction), ctx, id)
}

// CheckExistance mocks base method.
func (m *MockTransactionRepository) CheckExistance(ctx context.Context, id string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckExistance", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckExistance indicates an expected call of CheckExistance.
func (mr *MockTransactionRepositoryMockRecorder) CheckExistance(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckExistance", reflect.TypeOf((*MockTransactionRepository)(nil).CheckExistance), ctx, id)
}

// CreateTransaction mocks base method.
func (m *MockTransactionRepository) CreateTransaction(arg0 context.Context, arg1 *model.TransactionDao) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTransaction", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTransaction indicates an expected call of CreateTransaction.
func (mr *MockTransactionRepositoryMockRecorder) CreateTransaction(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTransaction", reflect.TypeOf((*MockTransactionRepository)(nil).CreateTransaction), arg0, arg1)
}

// GetLatestOddAndUncancelledTransactions mocks base method.
func (m *MockTransactionRepository) GetLatestOddAndUncancelledTransactions(ctx context.Context, limit int) ([]*model.TransactionDao, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestOddAndUncancelledTransactions", ctx, limit)
	ret0, _ := ret[0].([]*model.TransactionDao)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestOddAndUncancelledTransactions indicates an expected call of GetLatestOddAndUncancelledTransactions.
func (mr *MockTransactionRepositoryMockRecorder) GetLatestOddAndUncancelledTransactions(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestOddAndUncancelledTransactions", reflect.TypeOf((*MockTransactionRepository)(nil).GetLatestOddAndUncancelledTransactions), ctx, limit)
}
