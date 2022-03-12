package news

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

// ServiceNews is an interface for ConfigNews struct
type ServiceNews interface {
	Request() (int, string, *structs.ResponseNews)
}

// ConfigNews is a struct to keep input parameters required for the HTTP request to weather API
type ConfigNews struct {
	Date     string
	Language string
	Source   string
	Limit    int
	Markup   bool
	Token    string
}

// NewServiceNews is a constructor for ServiceNews
func NewServiceNews(date string, language string, source string, limit int, markup bool, token string) *ServiceNews {
	var conf ServiceNews = &ConfigNews{Date: date, Language: language, Source: source, Limit: limit, Markup: markup, Token: token}
	return &conf
}

// Request is a method to send the HTTP call to the 3rd party news API.
// Returns HTTP response status code (if available), error message or empty string, news data structure or nil.
func (cn *ConfigNews) Request() (int, string, *structs.ResponseNews) {
	// It seems like providing pageSize and page does not have any effect (10 is always a limit).
	// Keep parameters in the query in case news API will start working
	// as described at https://newsapi.org/docs/endpoints/top-headlines even for free accounts.
	newsURL := fmt.Sprintf("%s/top-headlines?language=%s&from=%s&sortBy=popularity&sources=%s&pageSize=%d&page=1&apiKey=%s",
		constants.NewsBaseURL, cn.Language, cn.Date, cn.Source, cn.Limit, cn.Token)
	resp, e1 := http.Get(newsURL)
	if e1 != nil {
		message := fmt.Sprintf("News request failed: %s\n", e1)
		return 0, strings.Replace(message, newsURL, constants.NewsBaseURL+"/...", 1), nil
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, e2 := ioutil.ReadAll(resp.Body)
	if e2 != nil {
		return resp.StatusCode, fmt.Sprintf("Failed to read news response body: %s\n", e2), nil
	}
	//log.Printf("News request for the date %s, language %s, and source %s responded with %s\n%s\n",
	//	cn.Date, cn.Language, cn.Source, resp.Status, string(body))

	if resp.StatusCode != 200 {
		return resp.StatusCode, string(body) + "\n", nil // TODO: return custom error message based on parsed body
	}

	var news structs.ResponseNews
	e3 := json.Unmarshal(body, &news)
	if e3 != nil {
		return resp.StatusCode, fmt.Sprintf("Reading JSON from news response body failed: %s\n", e3), nil
	}

	return resp.StatusCode, "", &news
}

// Run is a function to send an HTTP request to 3rd party News API and print the summary in case of success
func Run(out io.Writer, sn ServiceNews, conf *ConfigNews) error {
	output := ""
	status, message, news := sn.Request()

	i := 0 // to control limit explicitly since in news API page and pageSize parameters do not take effect.
	if status == 200 {
		for _, a := range news.Articles {
			i++
			if conf.Markup {
				output += fmt.Sprintf("[link](%s) %s\n>%s\n", a.URL, a.Title, a.Description)
			} else {
				output += fmt.Sprintf("%s\n%s\nMore at %s\n\n", a.Title, a.Description, a.URL)
			}
			if i == conf.Limit {
				break
			}
		}
	} else {
		output = fmt.Sprintf("Error: %s\n", message)
	}
	_, _ = fmt.Fprint(out, "", output)
	return nil
}
