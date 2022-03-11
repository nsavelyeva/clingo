package currency

import (
	"bytes"
	"clingo/constants"
	"clingo/structs"
	"clingo/test"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func Test_GetRate(t *testing.T) {
	responseData := &structs.ResponseCurrency{
		Data: &structs.Data{
			USD: 1.00,
			RUB: 0.30,
			EUR: 1.50,
		},
	}

	cc := ConfigCurrency{From: "from", To: "to", Token: "token"}

	tests := []struct {
		name     string
		currency string
		want     float64
	}{
		{"Currency found (capitalized)", "EUR", 1.5},
		{"Currency not found (non-capitalized)", "eur", 0},
		{"Currency not found (wrong currency value)", "foo", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.GetRate(responseData, tt.currency); got != tt.want {
				t.Errorf("GetRate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ValidateInputs(t *testing.T) {
	var details = map[string]structs.DetailsCurrency{
		"USD": {Symbol: "$", Name: "US Dollar"},
		"RUB": {Symbol: "₽", Name: "Russian Ruble"},
		"EUR": {Symbol: "€", Name: "Euro"},
	}

	cc := ConfigCurrency{From: "from", To: "to", Token: "token"}

	tests := []struct {
		name string
		from string
		to   string
		want string
	}{
		{"Validate uppercase from-currency", "USD", "EUR", ""},
		{"Validate lowercase from-currency", "usd", "EUR", ""},
		{"Validate uppercase to-currency", "USD", "EUR", ""},
		{"Validate lowercase to-currency", "USD", "eur", ""},
		{"Validate mixed to-currency multiple", "USD", "EUR,RUB,USD", ""},
		{"Validate wrong from-currency", "PLN", "EUR", "Value \"PLN\" is not recognized as supported currency\n"},
		{"Validate wrong to-currency", "USD", "BYN", "Value \"BYN\" is not recognized as supported currency\n"},
		{"Validate wrong to-currencies in list of to-currencies", "usd", "EUR,PLN,BYN,USD", "Value \"PLN\" is not recognized as supported currency\nValue \"BYN\" is not recognized as supported currency\n"},
		{"Validate empty from-currency", "", "EUR", "Value \"\" is not recognized as supported currency\n"},
		{"Validate empty to-currency", "USD", "", "Value \"\" is not recognized as supported currency\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cc.ValidateInputs(details, tt.to, tt.from); got != tt.want {
				t.Errorf("ValidateInputs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigCurrency_Request(t *testing.T) {
	tests := []struct {
		name        string
		from        string
		to          string
		token       string
		mockStatus  int
		mockBody    string
		wantMessage string
		wantData    *structs.ResponseCurrency
	}{
		{
			"ok: base currency capitalized",
			"EUR",
			"usd",
			"some_token",
			200,
			`{"data":{"USD":0.359306}}`,
			"",
			&structs.ResponseCurrency{Data: &structs.Data{USD: 0.359306}},
		},
		{
			"ok: base currency in small letters",
			"eur",
			"usd,rub,eur",
			"some_token",
			200,
			`{"data":{"USD":0.359306,"RUB":0.112025,"EUR":1.0}}`,
			"",
			&structs.ResponseCurrency{Data: &structs.Data{USD: 0.359306, RUB: 0.112025, EUR: 1.0}},
		},
		{
			"go error (bad json)",
			"eur",
			"usd",
			"some_token",
			200,
			`{"data":{"USD":0.359306,"RUB":0.112025,"EUR":1.0}`,
			"Reading currency response body failed: unexpected end of JSON input\n",
			nil,
		},
		{
			"bad request (wrong base currency value)",
			"foo",
			"usd",
			"some_token",
			422,
			`{"message":"The selected base currency is invalid.","errors":{"base_currency":["The selected base currency is invalid."]}}`,
			`{"message":"The selected base currency is invalid.","errors":{"base_currency":["The selected base currency is invalid."]}}` + "\n",
			nil,
		},
		{
			"ok (base currency empty - defaults to USD)",
			"",
			"usd",
			"some_token",
			200,
			`{"data":{"USD":1.0,"RUB":0.112025,"EUR":0.359306}}`,
			"",
			&structs.ResponseCurrency{Data: &structs.Data{EUR: 0.359306, RUB: 0.112025, USD: 1.0}},
		},
		{
			"unauthorized (wrong token value)",
			"eur",
			"usd",
			"foo",
			429,
			`{"message":"API rate limit exceeded"}`,
			`{"message":"API rate limit exceeded"}` + "\n", // TODO: improve
			nil,
		},
		{
			"unauthorized (token empty)",
			"eur",
			"usd",
			"",
			429,
			`{"message":"API rate limit exceeded"}`,
			`{"message":"API rate limit exceeded"}` + "\n",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/latest?apikey=%s&base_currency=%s", constants.CurrencyBaseURL, tt.token, tt.from),
				httpmock.NewBytesResponder(tt.mockStatus, []byte(tt.mockBody)),
			)

			cc := ConfigCurrency{From: tt.from, To: tt.to, Token: tt.token}
			status, message, data := cc.Request()

			if status != tt.mockStatus {
				t.Errorf("Request() status got = %v, want %v", status, tt.mockStatus)
			}
			if message != tt.wantMessage {
				t.Errorf("Request() message got = %v, want %v", message, tt.wantMessage)
			}
			if !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("Request() data got = %v, want %v", data, tt.wantData)
			}
		})
	}
}

func TestConfigCurrency_RequestErrorHTTP(t *testing.T) {
	cc := ConfigCurrency{From: "from", To: "to", Token: "token"}

	tests := []struct {
		name        string
		mockError   string
		wantStatus  int
		wantMessage string
		wantData    *structs.ResponseCurrency
	}{
		{
			"http error",
			"some error",
			0,
			fmt.Sprintf("Currency request failed: Get \"%s/...\": some error\n", constants.CurrencyBaseURL),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/latest?apikey=%s&base_currency=%s", constants.CurrencyBaseURL, "token", "from"),
				httpmock.NewErrorResponder(errors.New(tt.mockError)),
			)

			status, message, data := cc.Request()

			if status != tt.wantStatus {
				t.Errorf("Request() status got = %v, want %v", status, tt.wantStatus)
			}
			if message != tt.wantMessage {
				t.Errorf("Request() message got = %v, want %v", message, tt.wantMessage)
			}
			if !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("Request() data got = %v, want %v", data, tt.wantData)
			}
		})
	}
}

// Test constructor is successful for different combinations of parameters values,
// i.e. the created instance has Request() method (at least).
// Note, no need to loop over all methods because it is handled by compilation.
// TODO: empty parameters (from, to, token) should be rejected (hence, 2nd, 3rd and 4th tests will need to be updated then),
// but this needs code change in the constructor.
func TestNewServiceCurrency(t *testing.T) {
	tests := []struct {
		name  string
		from  string
		to    string
		token string
	}{
		{"ok", "eur", "usd,rub", "token"},
		{"empty from is ok", "", "eur", "token"},
		{"empty to is ok", "usd", "", "token"},
		{"empty token is ok", "rub", "eur", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServiceCurrency(tt.from, tt.to, tt.token)

			st := reflect.TypeOf(*got)
			_, exists := st.MethodByName("Request")
			if !exists {
				t.Error("Instance created by NewServiceWeather() constructor does not have method Request()")
			}
		})
	}
}

func TestRun(t *testing.T) {
	var details = map[string]structs.DetailsCurrency{
		"USD": {Symbol: "$", Name: "US Dollar"},
		"RUB": {Symbol: "₽", Name: "Russian Ruble"},
		"EUR": {Symbol: "€", Name: "Euro"},
	}

	mockData := &structs.ResponseCurrency{
		Data: &structs.Data{
			USD: 1.00,
			RUB: 0.30,
			EUR: 1.50,
		},
	}

	tests := []struct {
		name         string
		from         string
		to           string
		token        string
		mockValidate string
		mockStatus   int
		mockMessage  string
		mockData     structs.ResponseCurrency
		wantOut      string
	}{
		{
			"ok (200 response)",
			"rub",
			"eur",
			"token",
			"",
			200,
			"",
			*mockData,
			"1 ₽ = 1.500000 €\n",
		},
		{
			"error output (non-200 response)",
			"eur",
			"usd",
			"token",
			"validation error",
			400,
			"",
			*mockData,
			"validation error\n",
		},
		{
			"error output (0 response)",
			"eur",
			"usd",
			"token",
			"validation error",
			400,
			"",
			*mockData,
			"validation error\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := test.NewServiceCurrencyMock(tt.from, tt.to, tt.token)
			cs.On("Request").Return(tt.mockStatus, tt.mockMessage, &tt.mockData)
			cs.On("GetRate", &tt.mockData, "EUR").Return(tt.mockData.Data.EUR)
			cs.On("GetCurrenciesInfo").Return(details)
			cs.On("ValidateInputs", details, tt.to, tt.from).Return(tt.mockValidate)

			conf := ConfigCurrency{From: tt.from, To: tt.to, Token: tt.token}
			out := &bytes.Buffer{}
			err := Run(out, cs, &conf)
			require.NoError(t, err, fmt.Sprintf("currency.Run() failed with error '%s'", err))

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
