package weather

import (
	"bytes"
	"clingo/constants"
	"clingo/structs"
	"clingo/test"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func Test_FindEmoji(t *testing.T) {
	// Slicing records example: records[:][0:1] is a list containing a single element (header row)
	var records = [][]string{
		{"code", "day", "night", "icon", "emoji"},
		{"1000", "Sunny", "Clear", "113", ":sunny:"},
		{"1003", "Partly cloudy", "Partly cloudy", "116", ":sun_behind_cloud:"},
	}

	tests := []struct {
		name    string
		records [][]string
		code    int
		want    string
	}{
		{"Emoji found", records, 1003, ":sun_behind_cloud:"},
		{"Emoji not found", records, 1333, ""},
		{"Emoji in only header row", records[:][0:1], 1333, ""},
		{"Emoji in empty array", [][]string{}, 1003, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindEmoji(tt.records, tt.code); got != tt.want {
				t.Errorf("FindEmoji() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigWeather_Request(t *testing.T) {
	tests := []struct {
		name        string
		city        string
		token       string
		mockStatus  int
		mockBody    string
		wantMessage string
		wantData    *structs.ResponseWeather
	}{
		{
			"ok",
			"Amsterdam",
			"some_token",
			200,
			`{"Location":{"name":"Amsterdam"}}`,
			"",
			&structs.ResponseWeather{Location: &structs.Location{Name: "Amsterdam"}},
		},
		{
			"go error (bad json)",
			"Amsterdam",
			"some_token",
			200,
			`{"Location":{"name":Amsterdam"}}`,
			"Reading JSON from weather response body failed: invalid character 'A' looking for beginning of value",
			nil,
		},
		{
			"bad request (wrong city value)",
			"city",
			"some_token",
			400,
			`{"error":{"code":1006,"message":"No matching location found."}}`,
			`{"error":{"code":1006,"message":"No matching location found."}}`, // TODO: improve
			nil,
		},
		{
			"bad request (city empty)",
			"",
			"some_token",
			400,
			`{"error":{"code":1003,"message":"Parameter q is missing."}}`,
			`{"error":{"code":1003,"message":"Parameter q is missing."}}`,
			nil,
		},
		{
			"unauthorized (wrong token value)",
			"Amsterdam",
			"some_token",
			401,
			`{"error":{"code":2006,"message":"API key is invalid."}}`,
			`{"error":{"code":2006,"message":"API key is invalid."}}`,
			nil,
		},
		{
			"unauthorized (token empty)",
			"Amsterdam",
			"",
			401,
			`{"error":{"code":1002,"message":"API key is invalid or not provided."}}`,
			`{"error":{"code":1002,"message":"API key is invalid or not provided."}}`,
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", constants.WeatherBaseURL, tt.token, tt.city),
				httpmock.NewBytesResponder(tt.mockStatus, []byte(tt.mockBody)),
			)

			cw := ConfigWeather{City: tt.city, Token: tt.token}
			status, message, data := cw.Request()

			if status != tt.mockStatus {
				t.Errorf("Request() status got = %v, want %v", status, tt.mockStatus)
			}
			if strings.TrimRight(message, "\n") != tt.wantMessage {
				t.Errorf("Request() message got = %v, want %v", message, tt.wantMessage)
			}
			if !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("Request() data got = %v, want %v", data, tt.wantData)
			}
		})
	}
}

func TestConfigWeather_RequestErrorHTTP(t *testing.T) {
	cw := ConfigWeather{City: "city", Token: "token"}

	tests := []struct {
		name        string
		mockError   string
		wantStatus  int
		wantMessage string
		wantData    *structs.ResponseWeather
	}{
		{
			"http error",
			"some error",
			0,
			fmt.Sprintf("Weather request failed: Get \"%s/...\": some error\n", constants.WeatherBaseURL),
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()
			httpmock.RegisterResponder(
				"GET",
				fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", constants.WeatherBaseURL, "token", "city"),
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
// TODO: empty parameters (city, token) should be rejected (hence, 2nd and 3rd tests will need to be updated then),
// but this needs code change in the constructor.
func TestNewServiceWeather(t *testing.T) {
	tests := []struct {
		name  string
		city  string
		token string
	}{
		{"ok", "city", "token"},
		{"empty city is ok", "", "token"},
		{"empty token is ok", "city", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewServiceWeather(tt.city, tt.token)

			st := reflect.TypeOf(*got)
			_, exists := st.MethodByName("Request")
			if !exists {
				t.Error("Instance created by NewServiceWeather() constructor does not have method Request()")
			}
		})
	}
}

func TestRun(t *testing.T) {
	mockData := &structs.ResponseWeather{
		Location: &structs.Location{
			Name: "city",
		},
		Current: &structs.Current{
			TempC:      1.0,
			Condition:  structs.Condition{Text: "mock weather", Code: 1153},
			WindKph:    5.0,
			WindDir:    "N",
			PressureMb: 11.2,
			Humidity:   90,
			FeelslikeC: 0.1,
			Uv:         3.0,
		},
	}

	tests := []struct {
		name        string
		city        string
		token       string
		mockEmoji   string
		mockStatus  int
		mockMessage string
		mockData    structs.ResponseWeather
		wantOut     string
	}{
		{
			"ok (200 response)",
			"city",
			"token",
			":clingo_weather:",
			200,
			"",
			*mockData,
			"city: :clingo_weather: mock weather, t 1.0C (feels like 0.1C), wind N 5.00 km/h (1.4 m/s), pressure 11.2 mb, humidity 90, UV 3.0\n",
		},
		{
			"error output (non-200 response)",
			"city",
			"token",
			"",
			400,
			"error 400",
			*mockData,
			"Error: error 400",
		},
		{
			"error output (0 response)",
			"city",
			"token",
			"",
			0,
			"error 0",
			*mockData,
			"Error: error 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ws := test.NewServiceWeatherMock(tt.city, tt.token)
			ws.On("Request").Return(tt.mockStatus, tt.mockMessage, &tt.mockData)
			ws.On("GetEmoji", mockData.Current.Condition.Code).Return(tt.mockEmoji)

			out := &bytes.Buffer{}
			err := Run(out, ws)
			require.NoError(t, err, fmt.Sprintf("weather.Run() failed with error '%s'", err))

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
