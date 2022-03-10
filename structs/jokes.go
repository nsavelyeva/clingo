package structs

// ResponseJokes is a struct to store successful HTTP response from jokes API
type ResponseJokes struct {
	Error    bool   `json:"error"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Joke     string `json:"joke"`
	Flags    Flags  `json:"flags"`
	ID       int    `json:"id"`
	Safe     bool   `json:"safe"`
	Lang     string `json:"lang"`
}

// Flags is a sub-struct of ResponseJokes struct
type Flags struct {
	Nsfw      bool `json:"nsfw"`
	Religious bool `json:"religious"`
	Racist    bool `json:"racist"`
	Sexist    bool `json:"sexist"`
	Political bool `json:"political"`
	Explicit  bool `json:"explicit"`
}
