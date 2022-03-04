package helpers

import (
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCSV is a function to read contents of a CSV file and return a list of rows,
// where each row is a list of column values.
// Note: header row is processed the same way as the other rows.
func ReadCSV(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf(`Unable to read input CSV file "%s": %s`, filePath, err)
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		fmt.Printf(`Unable to parse file as CSV for "%s": %s`, filePath, err)
	}

	return records
}
