package test

import (
	"clingo/structs"

	"github.com/stretchr/testify/mock"
)

// ServiceCurrencyMock is a struct to be used as a mock object in tests
type ServiceCurrencyMock struct {
	mock.Mock
}

// NewServiceCurrencyMock is a mock constructor for ServiceCurrencyMock struct
func NewServiceCurrencyMock(from string, to string, token string) *ServiceCurrencyMock {
	return new(ServiceCurrencyMock)
}

// Request is a mock method for ServiceCurrencyMock struct
func (m *ServiceCurrencyMock) Request() (string, string, *structs.ResponseCurrency) {
	args := m.Called()
	return args.Get(0).(string), args.Get(1).(string), args.Get(2).(*structs.ResponseCurrency)
}

// GetCurrenciesInfo is a mock method for ServiceCurrencyMock struct
func (m *ServiceCurrencyMock) GetCurrenciesInfo() map[string]structs.DetailsCurrency {
	args := m.Called()
	return args.Get(0).(map[string]structs.DetailsCurrency)
}

// ValidateInputs is a mock method for ServiceCurrencyMock struct
func (m *ServiceCurrencyMock) ValidateInputs(details map[string]structs.DetailsCurrency, to string, from string) string {
	args := m.Called(details, to, from)
	return args.Get(0).(string)
}

// GetRate is a mock method for ServiceCurrencyMock struct
func (m *ServiceCurrencyMock) GetRate(rc *structs.ResponseCurrency, field string) float64 {
	args := m.Called(rc, field)
	return args.Get(0).(float64)
}
