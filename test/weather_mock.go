package test

import (
	"clingo/structs"

	"github.com/stretchr/testify/mock"
)

// ServiceWeatherMock is a struct to be used as a mock object in tests
type ServiceWeatherMock struct {
	mock.Mock
}

// NewServiceWeatherMock is a mock constructor for ServiceWeatherMock struct
func NewServiceWeatherMock(city string, token string) *ServiceWeatherMock {
	return new(ServiceWeatherMock)
}

// Request is a mock method for ServiceWeatherMock struct
func (m *ServiceWeatherMock) Request() (int, string, *structs.ResponseWeather) {
	args := m.Called()
	return args.Get(0).(int), args.Get(1).(string), args.Get(2).(*structs.ResponseWeather)
}

// GetEmoji is a mock method for ServiceWeatherMock struct
func (m *ServiceWeatherMock) GetEmoji(code int) string {
	args := m.Called(code)
	return args.Get(0).(string)
}
