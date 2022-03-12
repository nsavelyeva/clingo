package jokes

import (
	"clingo/constants"
	"clingo/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// ServiceJokes is an interface for ConfigJokes struct
type ServiceJokes interface {
	Request() (int, string, *structs.ResponseJokes)
}

// ConfigJokes is a struct to keep input parameters required for the HTTP request to weather API
type ConfigJokes struct {
	Emoji bool
	Token string
}

// NewServiceJokes is a constructor for ServiceJokes
func NewServiceJokes(token string) *ServiceJokes {
	var conf ServiceJokes = &ConfigJokes{Token: token}
	return &conf
}

// Request is a method to send the HTTP call to the 3rd party jokes API.
// Returns HTTP response status code (if available), error message or empty string, jokes data structure or nil.
func (cj *ConfigJokes) Request() (int, string, *structs.ResponseJokes) {
	jokesURL := fmt.Sprintf("%s/Any?format=json&type=single&blacklistFlags=nsfw,racist",
		constants.JokesBaseURL)
	req, e1 := http.NewRequest("GET", jokesURL, nil)
	if e1 != nil {
		return 0, fmt.Sprintf("Failed to build jokes request: %s\n", e1), nil
	}

	req.Header.Add("x-rapidapi-host", "jokeapi-v2.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", cj.Token)

	resp, e2 := http.DefaultClient.Do(req)
	if e2 != nil {
		return 0, fmt.Sprintf("Jokes request failed: %s\n", e2), nil
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e3 := ioutil.ReadAll(resp.Body)
	if e3 != nil {
		return resp.StatusCode, fmt.Sprintf("Failed to read jokes response body: %s\n", e3), nil
	}
	//log.Printf("Jokes request responded with %s\n%s", resp.Status, string(body))

	if resp.StatusCode != 200 {
		return resp.StatusCode, fmt.Sprintf("Jokes request responded with %s\n%s\n", resp.Status, body), nil // TODO: return custom error message based on parsed body
	}

	var joke structs.ResponseJokes
	e4 := json.Unmarshal(body, &joke)
	if e4 != nil {
		return resp.StatusCode, fmt.Sprintf("Reading JSON from jokes response body failed: %s\n", e4), nil
	}

	return resp.StatusCode, "", &joke
}

// Run is a function to send an HTTP request to 3rd party Jokes API and print the summary in case of success
func Run(out io.Writer, sj ServiceJokes, conf *ConfigJokes) error {
	output := ""
	status, message, jokes := sj.Request()

	if status == 200 {
		output = jokes.Joke + "\n"
		if conf.Emoji {
			output = fmt.Sprintf(":rolling_on_the_floor_laughing: %s", output)
		}
	} else {
		output = fmt.Sprintf("Error: %s\n", message)
	}
	_, _ = fmt.Fprint(out, "", output)
	return nil
}
