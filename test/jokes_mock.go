package test

import (
	"clingo/structs"

	"github.com/stretchr/testify/mock"
)

// ServiceJokesMock is a struct to be used as a mock object in tests
type ServiceJokesMock struct {
	mock.Mock
}

// NewServiceJokesMock is a mock constructor for ServiceJokesMock struct
func NewServiceJokesMock(token string) *ServiceJokesMock {
	return new(ServiceJokesMock)
}

// Request is a mock method for ServiceJokesMock struct
func (m *ServiceJokesMock) Request() (int, string, *structs.ResponseJokes) {
	args := m.Called()
	return args.Get(0).(int), args.Get(1).(string), args.Get(2).(*structs.ResponseJokes)
}
