package jokes

import (
	"clingo/constants"
	"clingo/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// ServiceJokes is an interface for ConfigJokes struct
type ServiceJokes interface {
	Request() (string, *structs.ResponseJokes)
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

// Request is a method to send the HTTP call to the 3rd party jokes API
func (cj *ConfigJokes) Request() (string, *structs.ResponseJokes) {
	jokesURL := fmt.Sprintf("%s/joke/Any?format=json&type=single&blacklistFlags=nsfw,racist",
		constants.JokesBaseURL)
	req, e1 := http.NewRequest("GET", jokesURL, nil)
	if e1 != nil {
		fmt.Printf("Failed to build jokes request: %s", e1)
		return e1.Error(), nil
	}

	req.Header.Add("x-rapidapi-host", "jokeapi-v2.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", cj.Token)

	resp, e2 := http.DefaultClient.Do(req)
	if e2 != nil {
		fmt.Printf("Jokes request failed: %s\n", e2)
		return e2.Error(), nil
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e3 := ioutil.ReadAll(resp.Body)
	if e3 != nil {
		fmt.Printf("Failed to read jokes response body: %s\n", e3)
		return e3.Error(), nil
	}
	//log.Printf("Jokes request responded with %s\n%s", resp.Status, string(body))

	var joke structs.ResponseJokes
	e4 := json.Unmarshal(body, &joke)
	if e4 != nil {
		fmt.Printf("Reading weather response body failed: %s\n", e4)
		return e4.Error(), nil
	}

	return resp.Status, &joke
}

// Run is a function to send an HTTP request to 3rd party Jokes API and print the summary in case of success
func Run(out io.Writer, conf *ConfigJokes) error {
	sj := *NewServiceJokes(conf.Token)
	status, jokes := sj.Request()

	if strings.HasPrefix(status, "200") {
		message := jokes.Joke
		if conf.Emoji == true {
			message = ":rolling_on_the_floor_laughing: " + message
		}

		_, _ = fmt.Fprint(out, "", message+"\n")
	}
	return nil
}
