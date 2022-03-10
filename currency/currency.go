package currency

import (
	"clingo/constants"
	"clingo/helpers"
	"clingo/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

// ServiceCurrency is an interface for ConfigCurrency struct
type ServiceCurrency interface {
	Request() (string, string, *structs.ResponseCurrency)
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

// Request is a method to send the HTTP call to the 3rd party currency API
func (cw *ConfigCurrency) Request() (string, string, *structs.ResponseCurrency) {
	currencyURL := fmt.Sprintf("%s/latest?apikey=%s&base_currency=%s", constants.CurrencyBaseURL, cw.Token, cw.From)
	resp, e1 := http.Get(currencyURL)
	if e1 != nil {
		message := strings.Replace(fmt.Sprintf("Currency request failed: %s", e1), currencyURL, constants.CurrencyBaseURL+"/...", 1)
		return "", message, nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		message := fmt.Sprintf("Failed to read currency response body: %s\n", e2)
		return resp.Status, message, nil
	}
	//log.Printf("Currency request from currency %s to currency %s responded with %s\n%s",
	//	cw.From, cw.To, resp.Status, string(body))

	if resp.Status != "200" {
		return resp.Status, string(body), nil
	}

	var currency structs.ResponseCurrency
	e3 := json.Unmarshal(body, &currency)
	if e3 != nil {
		return resp.Status, fmt.Sprintf("Reading currency response body failed: %s\n", e3), nil
	}

	return resp.Status, "", &currency
}

// GetCurrenciesInfo is a method that loads info about all supported currencies from a JSON file
// (this info is used further for validation and formatting purpose).
func (cw *ConfigCurrency) GetCurrenciesInfo() map[string]structs.DetailsCurrency {
	content := helpers.ReadJSON(constants.CurrencyDetailsJSONFilePath)
	var details map[string]structs.DetailsCurrency
	err := json.Unmarshal(content, &details)
	if err != nil {
		fmt.Printf(`Error loading JSON from "%s": %s`, constants.CurrencyDetailsJSONFilePath, err)
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
	//sc := *NewServiceCurrency(conf.From, conf.To, conf.Token)

	details := sc.GetCurrenciesInfo()
	validationError := sc.ValidateInputs(details, conf.To, conf.From)
	if validationError != "" {
		_, _ = fmt.Fprint(out, "", validationError+"\n")
		return nil
	}

	status, _, currency := sc.Request()

	if strings.HasPrefix(status, "200") {
		m := ""
		for _, c := range strings.Split(strings.ToUpper(conf.To), ",") {
			rate := sc.GetRate(currency, c)
			m += fmt.Sprintf(" = %.6f %s", rate, details[c].Symbol)
		}
		m = fmt.Sprintf("1 %s%s\n", details[strings.ToUpper(conf.From)].Symbol, m)
		_, _ = fmt.Fprint(out, "", m)
	}
	return nil
}
