// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecase.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ttagiyeva/entain/internal/model"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// PostProcess mocks base method.
func (m *MockUsecase) PostProcess(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PostProcess", ctx)
}

// PostProcess indicates an expected call of PostProcess.
func (mr *MockUsecaseMockRecorder) PostProcess(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostProcess", reflect.TypeOf((*MockUsecase)(nil).PostProcess), ctx)
}

// Process mocks base method.
func (m *MockUsecase) Process(arg0 context.Context, arg1 *model.Transaction) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockUsecaseMockRecorder) Process(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockUsecase)(nil).Process), arg0, arg1)
}
