package news

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

func TestConfigNews_Request(t *testing.T) {
	cn := ConfigNews{Language: "nl", Date: "2020-10-20", Source: "src", Limit: 100, Token: "token"}

	tests := []struct {
		name        string
		conf        *ConfigNews
		mockStatus  int
		mockBody    string
		wantMessage string
		wantData    *structs.ResponseNews
	}{
		{
			"ok - multiple articles",
			&cn,
			200,
			`{"totalResults":2,"articles":[{"title":"title 0","description":"description 0","url":"https://url0..."},{"title":"title 1","description":"description 1","url":"https://url1..."}]}`,
			"",
			&structs.ResponseNews{TotalResults: 2, Articles: []structs.Article{
				{Title: "title 0", Description: "description 0", URL: "https://url0..."},
				{Title: "title 1", Description: "description 1", URL: "https://url1..."},
			}},
		},
		{
			"ok - no articles",
			&cn,
			200,
			`{"totalResults":0,"articles":[]}`,
			"",
			&structs.ResponseNews{TotalResults: 0, Articles: []structs.Article{}},
		},
		{
			"go error (bad json)",
			&cn,
			200,
			`{"totalResults":1,"articles":[{"title":,"description":"description 0","url":"https://url0..."}]}`,
			"Reading JSON from news response body failed: invalid character ',' looking for beginning of value\n",
			nil,
		},
		{
			"bad request (incompatible parameters)",
			&cn,
			422,
			`{"status":"error","code":"parametersIncompatible","message":"You cannot mix the sources parameter with the country or category parameters."}`,
			`{"status":"error","code":"parametersIncompatible","message":"You cannot mix the sources parameter with the country or category parameters."}` + "\n",
			nil,
		},
		{
			"unauthorized (wrong token value)",
			&cn,
			401,
			`{"status":"error","code":"apiKeyInvalid","message":"Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key."}`,
			`{"status":"error","code":"apiKeyInvalid","message":"Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key."}` + "\n",
			nil,
		},
		{
			"unauthorized (token empty)",
			&cn,
			401,
			`{"status":"error","code":"apiKeyInvalid","message":"Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key."}`,
			`{"status":"error","code":"apiKeyInvalid","message":"Your API key is invalid or incorrect. Check your key, or go to https://newsapi.org to create a free API key."}` + "\n",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/top-headlines?language=%s&from=%s&sortBy=popularity&sources=%s&pageSize=%d&page=1&apiKey=%s",
					constants.NewsBaseURL, cn.Language, cn.Date, cn.Source, cn.Limit, cn.Token),
				httpmock.NewBytesResponder(tt.mockStatus, []byte(tt.mockBody)),
			)

			status, message, data := tt.conf.Request()

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

func TestConfigNews_RequestErrorHTTP(t *testing.T) {
	cn := ConfigNews{Language: "nl", Date: "2020-10-20", Source: "src", Limit: 100, Token: "token"}

	tests := []struct {
		name        string
		conf        *ConfigNews
		mockError   string
		wantStatus  int
		wantMessage string
		wantData    *structs.ResponseNews
	}{
		{
			"http error",
			&cn,
			"some error",
			0,
			fmt.Sprintf("News request failed: Get \"%s/...\": some error\n", constants.NewsBaseURL),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/top-headlines?language=%s&from=%s&sortBy=popularity&sources=%s&pageSize=%d&page=1&apiKey=%s",
					constants.NewsBaseURL, cn.Language, cn.Date, cn.Source, cn.Limit, cn.Token),
				httpmock.NewErrorResponder(errors.New(tt.mockError)),
			)

			status, message, data := tt.conf.Request()

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
// but this needs code change in the constructor.
func TestNewServiceNews(t *testing.T) {
	cn := ConfigNews{Language: "nl", Date: "2020-10-20", Source: "src", Limit: 100, Token: "token"}

	tests := []struct {
		name string
		conf *ConfigNews
	}{
		{"ok", &cn},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServiceNews(tt.conf.Date, tt.conf.Language, tt.conf.Source, tt.conf.Limit, tt.conf.Markup, tt.conf.Token)

			st := reflect.TypeOf(*got)
			_, exists := st.MethodByName("Request")
			if !exists {
				t.Error("Instance created by NewServiceNews() constructor does not have method Request()")
			}
		})
	}
}

func TestRun(t *testing.T) {
	articles := []structs.Article{
		{Title: "title 0", Description: "description 0", URL: "url 0"},
		{Title: "title 1", Description: "description 1", URL: "url 1"},
		{Title: "title 2", Description: "description 2", URL: "url 2"},
	}

	allArticlesText := "title 0\ndescription 0\nMore at url 0\n\ntitle 1\ndescription 1\nMore at url 1\n\ntitle 2\ndescription 2\nMore at url 2\n\n"
	allArticlesMarkup := "[link](url 0) title 0\n>description 0\n[link](url 1) title 1\n>description 1\n[link](url 2) title 2\n>description 2\n"

	tests := []struct {
		name        string
		conf        *ConfigNews
		mockStatus  int
		mockMessage string
		mockData    *structs.ResponseNews
		wantOut     string
	}{
		{
			"ok without markup (200 response)",
			&ConfigNews{Markup: false},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles), Articles: articles},
			allArticlesText,
		},
		{
			"ok with markup (200 response)",
			&ConfigNews{Markup: true},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles), Articles: articles},
			allArticlesMarkup,
		},
		{
			"display all articles if limit 100",
			&ConfigNews{Limit: 100},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles), Articles: articles},
			allArticlesText,
		},
		{
			"display all articles if limit 0",
			&ConfigNews{Limit: 0},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles), Articles: articles},
			allArticlesText,
		},
		{
			"display n articles if limit n is 0 < n < 10",
			&ConfigNews{Limit: 2},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles[0:2]), Articles: articles[0:2]},
			"title 0\ndescription 0\nMore at url 0\n\ntitle 1\ndescription 1\nMore at url 1\n\n",
		},
		{
			"display all articles if limit < 0",
			&ConfigNews{Limit: -5},
			200,
			"",
			&structs.ResponseNews{TotalResults: len(articles), Articles: articles},
			allArticlesText,
		},
		{
			"error output (non-200 response)",
			&ConfigNews{},
			400,
			"error 400",
			nil,
			"Error: error 400\n",
		},
		{
			"error output (0 response)",
			&ConfigNews{},
			0,
			"error 0",
			nil,
			"Error: error 0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := test.NewServiceNewsMock(tt.conf.Date, tt.conf.Language, tt.conf.Source, tt.conf.Limit, tt.conf.Markup, tt.conf.Token)
			cs.On("Request").Return(tt.mockStatus, tt.mockMessage, tt.mockData)

			out := &bytes.Buffer{}
			err := Run(out, cs, tt.conf)
			require.NoError(t, err, fmt.Sprintf("news.Run() failed with error '%s'", err))

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
