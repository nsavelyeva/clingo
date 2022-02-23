package slack

import (
	"encoding/json"
	"fmt"
	"io"
)

// Config is a struct to store input parameters just for a demo purpose
type Config struct {
	String string
	Float  float64
	Bool   bool
}

// Run is a function to dump Config struct into JSON format and display it just for a demo purpose
func Run(out io.Writer, conf Config) error {
	b, err := json.MarshalIndent(conf, "", " ")
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(out, "Slack config:\n%s", b)

	return nil
}
