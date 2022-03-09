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
	Request() (string, *structs.ResponseCurrency)
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
func (cw *ConfigCurrency) Request() (string, *structs.ResponseCurrency) {
	currencyUrl := fmt.Sprintf("%s/latest?apikey=%s&base_currency=%s", constants.CurrencyBaseURL, cw.Token, cw.From)
	resp, e1 := http.Get(currencyUrl)
	if e1 != nil {
		fmt.Printf("Currency request failed: %s\n", e1)
		return e1.Error(), nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		fmt.Printf("Failed to read currency response body: %s\n", e2)
		return e2.Error(), nil
	}
	//log.Printf("Currency request from currency %s to currency %s responded with %s\n%s",
	//	cw.From, cw.To, resp.Status, string(body))

	var currency structs.ResponseCurrency
	e3 := json.Unmarshal(body, &currency)
	if e3 != nil {
		fmt.Printf("Reading currency response body failed: %s\n", e3)
		return e3.Error(), nil
	}

	return resp.Status, &currency
}

// GetRate is a method that takes a value of the given field using its name as a string,
// i.e. not as a property because property value is unknown (dynamic)
func (cw *ConfigCurrency) GetRate(rc *structs.ResponseCurrency, field string) float64 {
	r := reflect.ValueOf(rc.Data)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Float()
}

// Run is a function to send an HTTP request to 3rd party Currency API and print the summary in case of success
func Run(out io.Writer, conf *ConfigCurrency) error {
	content := helpers.ReadJSON(constants.CurrencyDetailsJSONFilePath)
	var details map[string]structs.DetailsCurrency
	err := json.Unmarshal(content, &details)
	if err != nil {
		fmt.Printf(`Error loading JSON from "%s": %s`, constants.CurrencyDetailsJSONFilePath, err)
		return nil
	}

	cc := strings.Split(strings.ToUpper(conf.To), ",")
	for _, c := range append(cc, strings.ToUpper(conf.From)) {
		_, exists := details[c]
		if !exists {
			fmt.Printf("Value \"%s\" is not recognized as supported currency\n", c)
			return nil
		}
	}

	sc := *NewServiceCurrency(conf.From, conf.To, conf.Token)
	status, currency := sc.Request()

	if strings.HasPrefix(status, "200") {
		m := ""
		for _, c := range cc {
			rate := sc.GetRate(currency, c)
			m += fmt.Sprintf(" = %.6f %s", rate, details[c].Symbol)
		}
		m = fmt.Sprintf("1 %s%s\n", details[strings.ToUpper(conf.From)].Symbol, m)
		_, _ = fmt.Fprint(out, "", m)
	}
	return nil
}
