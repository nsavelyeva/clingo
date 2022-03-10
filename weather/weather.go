package weather

import (
	"clingo/constants"
	"clingo/helpers"
	"clingo/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// ServiceWeather is an interface for ConfigWeather struct
type ServiceWeather interface {
	Request() (string, string, *structs.ResponseWeather)
	GetEmoji(int) string
}

// ConfigWeather is a struct to keep input parameters required for the HTTP request to weather API
type ConfigWeather struct {
	City  string
	Token string
}

// NewServiceWeather is a constructor for ServiceWeather
func NewServiceWeather(city string, token string) *ServiceWeather {
	var conf ServiceWeather = &ConfigWeather{City: city, Token: token}
	return &conf
}

// Request is a method to send the HTTP call to the 3rd party weather API
func (cw *ConfigWeather) Request() (string, string, *structs.ResponseWeather) {
	weatherURL := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no", constants.WeatherBaseURL, cw.Token, cw.City)
	resp, e1 := http.Get(weatherURL)
	if e1 != nil {
		message := strings.Replace(fmt.Sprintf("Weather request failed: %s", e1), weatherURL, constants.WeatherBaseURL+"/...", 1)
		return "", message, nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		return resp.Status, fmt.Sprintf("Failed to read weather response body: %s", e2), nil
	}
	//log.Printf("Weather request for the city %s responded with %s\n%s",
	//	cw.City, resp.Status, string(body))

	if resp.Status != "200" {
		return resp.Status, string(body), nil
	}

	var weather structs.ResponseWeather
	e3 := json.Unmarshal(body, &weather)
	if e3 != nil {
		return resp.Status, fmt.Sprintf("Reading JSON from weather response body failed: %s", e3), nil
	}

	return resp.Status, "", &weather
}

// GetEmoji is a method to read CSV data of weather conditions and find emoji there by its code
func (cw *ConfigWeather) GetEmoji(code int) string {
	records := helpers.ReadCSV(constants.WeatherConditionsCSVFilePath)
	emoji := FindEmoji(records, code)
	return emoji
}

// FindEmoji is a function to pick up the emoji corresponding to the provided weather condition code
func FindEmoji(records [][]string, code int) string {
	emoji := ""
	for _, value := range records {
		if value[0] == strconv.Itoa(code) {
			emoji = value[4]
			break
		}
	}
	return emoji
}

// Run is a function to send an HTTP request to 3rd party Weather API and print the summary in case of success
func Run(out io.Writer, sw ServiceWeather) error {
	status, _, weather := sw.Request()

	if strings.HasPrefix(status, "200") {
		emoji := sw.GetEmoji(weather.Current.Condition.Code)
		ms := weather.Current.WindKph * 1000 / 3600

		output := fmt.Sprintf("%s: %s %s, t %.1fC (feels like %.1fC), wind %s %.2f km/h (%.1f m/s), pressure %.1f mb, humidity %d, UV %.1f\n",
			weather.Location.Name, emoji, weather.Current.Condition.Text,
			weather.Current.TempC, weather.Current.FeelslikeC,
			weather.Current.WindDir, weather.Current.WindKph, ms,
			weather.Current.PressureMb, weather.Current.Humidity, weather.Current.Uv)

		_, _ = fmt.Fprint(out, "", output)
	}
	return nil
}
