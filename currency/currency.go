package currency

import (
	"clingo/constants"
	"clingo/helpers"
	"clingo/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

// ServiceCurrency is an interface for ConfigCurrency struct
type ServiceCurrency interface {
	Request() (int, string, *structs.ResponseCurrency)
	GetCurrenciesInfo() map[string]structs.DetailsCurrency
	ValidateInputs(details map[string]structs.DetailsCurrency, to string, from string) string
	GetRate(rc *structs.ResponseCurrency, field string) float64
}

// ConfigCurrency is a struct to keep input parameters required for the HTTP request to currency API
type ConfigCurrency struct {
	From  string
	To    string
	Token string
}

// NewServiceCurrency is a constructor for ServiceCurrency
func NewServiceCurrency(from string, to string, token string) *ServiceCurrency {
	var conf ServiceCurrency = &ConfigCurrency{From: from, To: to, Token: token}
	return &conf
}

// Request is a method to send the HTTP call to the 3rd party currency API.
// Returns HTTP response status code (if available), error message or empty string, currency data structure or nil.
func (cw *ConfigCurrency) Request() (int, string, *structs.ResponseCurrency) {
	currencyURL := fmt.Sprintf("%s/latest?apikey=%s&base_currency=%s", constants.CurrencyBaseURL, cw.Token, cw.From)
	resp, e1 := http.Get(currencyURL)
	if e1 != nil {
		message := fmt.Sprintf("Currency request failed: %s\n", e1)
		return 0, strings.Replace(message, currencyURL, constants.CurrencyBaseURL+"/...", 1), nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		return resp.StatusCode, fmt.Sprintf("Failed to read currency response body: %s\n", e2), nil
	}
	log.Printf("Currency request from currency %s to currency %s responded with %s\n%s",
		cw.From, cw.To, resp.Status, string(body))

	if resp.StatusCode != 200 {
		return resp.StatusCode, string(body) + "\n", nil // TODO: return custom error message based on parsed body
	}

	var currency structs.ResponseCurrency
	e3 := json.Unmarshal(body, &currency)
	if e3 != nil {
		return resp.StatusCode, fmt.Sprintf("Reading currency response body failed: %s\n", e3), nil
	}

	return resp.StatusCode, "", &currency
}

// GetCurrenciesInfo is a method that loads info about all supported currencies from a JSON file
// (this info is used further for validation and formatting purpose).
func (cw *ConfigCurrency) GetCurrenciesInfo() map[string]structs.DetailsCurrency {
	content := helpers.ReadJSON(constants.CurrencyDetailsJSONFilePath)
	var details map[string]structs.DetailsCurrency
	err := json.Unmarshal(content, &details)
	if err != nil {
		fmt.Printf("Error loading JSON from \"%s\": %s\n", constants.CurrencyDetailsJSONFilePath, err)
		return nil
	}
	return details
}

// ValidateInputs is a method that takes values of provided CLI options
// and checks if the specified currencies codes are supported (i.e. exist in the external JSON file).
func (cw *ConfigCurrency) ValidateInputs(details map[string]structs.DetailsCurrency, to string, from string) string {
	cc := strings.Split(strings.ToUpper(to), ",")
	err := ""
	for _, c := range append(cc, strings.ToUpper(from)) {
		_, exists := details[c]
		if !exists {
			err += fmt.Sprintf("Value \"%s\" is not recognized as supported currency\n", c)
		}
	}
	return err
}

// GetRate is a method that takes a value of the given field using its name as a string,
// i.e. not as a property because property value is unknown (dynamic).
// If the currency does not present, return 0 rate.
func (cw *ConfigCurrency) GetRate(rc *structs.ResponseCurrency, field string) float64 {
	r := reflect.ValueOf(rc.Data)
	f := reflect.Indirect(r).FieldByName(field)
	if !f.IsValid() {
		return float64(0)
	}
	return f.Float()
}

// Run is a function to send an HTTP request to 3rd party Currency API and print the summary in case of success
func Run(out io.Writer, sc ServiceCurrency, conf *ConfigCurrency) error {
	output := ""

	details := sc.GetCurrenciesInfo()
	validationError := sc.ValidateInputs(details, conf.To, conf.From)
	if validationError != "" {
		_, _ = fmt.Fprint(out, "", validationError+"\n")
		return nil
	}

	status, message, currency := sc.Request()

	if status == 200 {
		for _, c := range strings.Split(strings.ToUpper(conf.To), ",") {
			rate := sc.GetRate(currency, c)
			output += fmt.Sprintf(" = %.6f %s", rate, details[c].Symbol)
		}
		output = fmt.Sprintf("1 %s%s\n", details[strings.ToUpper(conf.From)].Symbol, output)
	} else {
		output = fmt.Sprintf("Error: %s", message)
	}
	_, _ = fmt.Fprint(out, "", output)
	return nil
}
