package helpers

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ReadJSON is a function to read contents of a JSON file and return its content as a byte value,
// which can be unmarshalled with json.Unmarshal() using the proper struct to load JSON data into.
func ReadJSON(filePath string) []byte {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf(`Unable to read input JSON file "%s": %s`, filePath, err)
	}

	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf(`Unable to parse file as JSON for "%s": %s`, filePath, err)
	}

	return byteValue
}
