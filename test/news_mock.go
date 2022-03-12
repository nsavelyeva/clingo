package test

import (
	"clingo/structs"

	"github.com/stretchr/testify/mock"
)

// ServiceNewsMock is a struct to be used as a mock object in tests
type ServiceNewsMock struct {
	mock.Mock
}

// NewServiceNewsMock is a mock constructor for ServiceNewsMock struct
func NewServiceNewsMock(date string, language string, source string, limit int, markup bool, token string) *ServiceNewsMock {
	return new(ServiceNewsMock)
}

// Request is a mock method for ServiceNewsMock struct
func (m *ServiceNewsMock) Request() (int, string, *structs.ResponseNews) {
	args := m.Called()
	return args.Get(0).(int), args.Get(1).(string), args.Get(2).(*structs.ResponseNews)
}
