package structs

import "time"

// ResponseNews is a struct to store successful HTTP response from news API
type ResponseNews struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

// Article is a struct, a list element of Articles
type Article struct {
	Source      Source    `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

// Source is a sub-struct of Article struct
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
