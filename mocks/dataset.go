// Code generated by MockGen. DO NOT EDIT.
// Source: dataset.go

// Package mock_domoapi is a generated GoMock package.
package mock_domoapi

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockRequestHandlerService is a mock of RequestHandlerService interface
type MockRequestHandlerService struct {
	ctrl     *gomock.Controller
	recorder *MockRequestHandlerServiceMockRecorder
}

// MockRequestHandlerServiceMockRecorder is the mock recorder for MockRequestHandlerService
type MockRequestHandlerServiceMockRecorder struct {
	mock *MockRequestHandlerService
}

// NewMockRequestHandlerService creates a new mock instance
func NewMockRequestHandlerService(ctrl *gomock.Controller) *MockRequestHandlerService {
	mock := &MockRequestHandlerService{ctrl: ctrl}
	mock.recorder = &MockRequestHandlerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRequestHandlerService) EXPECT() *MockRequestHandlerServiceMockRecorder {
	return m.recorder
}

// Handler mocks base method
func (m *MockRequestHandlerService) Handler(req *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handler", req)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handler indicates an expected call of Handler
func (mr *MockRequestHandlerServiceMockRecorder) Handler(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handler", reflect.TypeOf((*MockRequestHandlerService)(nil).Handler), req)
}