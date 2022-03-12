package jokes

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

func TestConfigJokes_Request(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		mockStatus  int
		mockBody    string
		wantMessage string
		wantData    *structs.ResponseJokes
	}{
		{
			"ok",
			"some_token",
			200,
			`{"joke": "Ha-ha-ha."}`,
			"",
			&structs.ResponseJokes{Joke: "Ha-ha-ha."},
		},
		{
			"go error (bad json)",
			"some_token",
			200,
			`{"joke": Ha-ha-ha."}`,
			"Reading JSON from jokes response body failed: invalid character 'H' looking for beginning of value\n",
			nil,
		},
		{
			"unauthorized (wrong token value)",
			"some_token",
			403,
			`{"message":"You are not subscribed to this API."}`,
			"Jokes request responded with 403\n" + `{"message":"You are not subscribed to this API."}` + "\n",
			nil,
		},
		{
			"unauthorized (token empty)",
			"",
			401,
			`{"message":"Invalid API key. Go to https:\/\/docs.rapidapi.com\/docs\/keys for more info."}`,
			"Jokes request responded with 401\n" + `{"message":"Invalid API key. Go to https:\/\/docs.rapidapi.com\/docs\/keys for more info."}` + "\n",
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/Any?format=json&type=single&blacklistFlags=nsfw,racist", constants.JokesBaseURL),
				httpmock.NewBytesResponder(tt.mockStatus, []byte(tt.mockBody)),
			)

			cw := ConfigJokes{Token: tt.token}
			status, message, data := cw.Request()

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

func TestConfigJokes_RequestErrorHTTP(t *testing.T) {
	cw := ConfigJokes{Token: "token"}

	tests := []struct {
		name        string
		mockError   string
		wantStatus  int
		wantMessage string
		wantData    *structs.ResponseJokes
	}{
		{
			"http error",
			"some error",
			0,
			fmt.Sprintf("Jokes request failed: Get \"%s/Any?format=json&type=single&blacklistFlags=nsfw,racist\": some error\n", constants.JokesBaseURL),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/Any?format=json&type=single&blacklistFlags=nsfw,racist", constants.JokesBaseURL),
				httpmock.NewErrorResponder(errors.New(tt.mockError)),
			)

			status, message, data := cw.Request()

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
// TODO: empty parameters (token) should be rejected (hence, 2nd test will need to be updated then),
// but this needs code change in the constructor.
func TestNewServiceJokes(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"ok", "token"},
		{"empty token is ok", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServiceJokes(tt.token)

			st := reflect.TypeOf(*got)
			_, exists := st.MethodByName("Request")
			if !exists {
				t.Error("Instance created by NewServiceJokes() constructor does not have method Request()")
			}
		})
	}
}

func TestRun(t *testing.T) {
	mockData := &structs.ResponseJokes{
		Joke: "Ha-ha!",
	}

	tests := []struct {
		name        string
		token       string
		mockEmoji   bool
		mockStatus  int
		mockMessage string
		mockData    *structs.ResponseJokes
		wantOut     string
	}{
		{
			"ok without emoji(200 response)",
			"token",
			false,
			200,
			"",
			mockData,
			mockData.Joke + "\n",
		},
		{
			"ok with emoji (200 response)",
			"token",
			true,
			200,
			"",
			mockData,
			fmt.Sprintf(":rolling_on_the_floor_laughing: %s\n", mockData.Joke),
		},
		{
			"error output (non-200 response)",
			"token",
			false,
			400,
			"error 400",
			nil,
			"Error: error 400\n",
		},
		{
			"error output (0 response)",
			"token",
			false,
			0,
			"error 0",
			nil,
			"Error: error 0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := test.NewServiceJokesMock(tt.token)
			js.On("Request").Return(tt.mockStatus, tt.mockMessage, tt.mockData)

			conf := ConfigJokes{Emoji: tt.mockEmoji, Token: tt.token}
			out := &bytes.Buffer{}
			err := Run(out, js, &conf)
			require.NoError(t, err, fmt.Sprintf("jokes.Run() failed with error '%s'", err))

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
