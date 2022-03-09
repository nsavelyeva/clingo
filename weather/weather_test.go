package weather

import (
	"bytes"
	"clingo/constants"
	"clingo/structs"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"testing"
)

// Slicing records example: records[:][0:1] is a list containing a single element (header row)
var records = [][]string {
	{"code", "day", "night", "icon", "emoji"},
	{"1000", "Sunny", "Clear", "113", ":sunny:"},
	{"1003", "Partly cloudy", "Partly cloudy", "116", ":sun_behind_cloud:"},
}

func TestConfigWeather_GetEmoji(t *testing.T) {
	tests := []struct {
		name   string
		records [][]string
		code   int
		want   string
	}{
		{"Emoji found", records, 1003, ":sun_behind_cloud:"},
		{"Emoji not found", records, 1333, ""},
		{"Emoji in only header row", records[:][0:1], 1333, ""},
		{"Emoji in empty array", [][]string{}, 1003, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEmoji(tt.records, tt.code); got != tt.want {
				t.Errorf("GetEmoji() = %v, want %v", got, tt.want)
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
		wantStatus  string
		wantMessage string
		wantData    *structs.ResponseWeather
	}{
		{
			"ok",
			"Amsterdam",
			"some_token",
			200,
			`{"Location":{"name":"Amsterdam"}}`,
			"200",
			"",
			&structs.ResponseWeather{Location: &structs.Location{Name: "Amsterdam"}},
		},
		{
			"bad request (wrong city value)",
			"city",
			"some_token",
			400,
			`{"error":{"code":1006,"message":"No matching location found."}}`,
			"400",
			`{"error":{"code":1006,"message":"No matching location found."}}`,
			nil,
		},
		{
			"bad request (city empty)",
			"",
			"some_token",
			400,
			`{"error":{"code":1003,"message":"Parameter q is missing."}}`,
			"400",
			`{"error":{"code":1003,"message":"Parameter q is missing."}}`,
			nil,
		},
		{
			"unauthorized (wrong token value)",
			"Amsterdam",
			"some_token",
			401,
			`{"error":{"code":2006,"message":"API key is invalid."}}`,
			"401",
			`{"error":{"code":2006,"message":"API key is invalid."}}`,
			nil,
		},
		{
			"unauthorized (token empty)",
			"Amsterdam",
			"",
			401,
			`{"error":{"code":1002,"message":"API key is invalid or not provided."}}`,
			"401",
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
				fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", constants.WeatherBaseURL, "token", "city"),
				httpmock.NewStringResponder(tt.mockStatus, tt.mockBody),
			)

			cw := ConfigWeather{City: "city", Token: "token"}
			status, message, data := cw.Request()

			if status != tt.wantStatus {
				t.Errorf("Request() status got = %v, want %v", status, tt.wantStatus)
			}
			if !reflect.DeepEqual(data, tt.wantData) {
				t.Errorf("Request() data got = %v, want %v", data, tt.wantData)
			}
			if !reflect.DeepEqual(message, tt.wantMessage) {
				t.Errorf("Request() message got = %v, want %v", message, tt.mockBody)
			}
		})
	}
}

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

// TODO: Switch to weather mock
func TestRun(t *testing.T) {
	mockBody := `{
"location":{"name":"city"},
"current":{"temp_c":1.0,
  "feelslike_c":0.1,
  "win_dir":"N",
  "wind_kph":5.0,
  "pressure_mb":11.2,
  "humidity":90,
  "uv":3.0,
  "condition":{"text":"mock weather","code":1153}
  }
}`
	wantOutput := "city: :rain_cloud: mock weather, t 1.0C (feels like 0.1C), wind  5.00 km/h (1.4 m/s), pressure 11.2 mb, humidity 90, UV 3.0\n"

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		"GET",
		fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", constants.WeatherBaseURL, "token", "city"),
		httpmock.NewStringResponder(200, mockBody),
	)
	tests := []struct {
		name    string
		city    string
		token   string
		wantOut string
		wantErr bool
	}{
		{"ok", "city", "token", wantOutput, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			conf := ConfigWeather{City: tt.city, Token: tt.token}
			require.NoError(t,
				os.Chdir(".."),
				fmt.Sprintf("error changing to the temporary test directory %s", ".."),
			)

			err := Run(out, &conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("Run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
