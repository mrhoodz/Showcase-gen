package readjson

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type pageData struct {
	Href string `json:"href"`
	Size string `json:"size,omitempty"`
	Tags string `json:"tags"`
}

func ReadJSON(fileName string) ([]pageData, error) {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Initialize a variable to store the JSON data
	var jsonData []byte

	// Read the file line by line
	for scanner.Scan() {
		// Append each line to the jsonData variable
		jsonData = append(jsonData, scanner.Bytes()...)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %v", err)
	}

	// Unmarshal the JSON data into a slice of pageData structs
	var pages []pageData
	if err := json.Unmarshal(jsonData, &pages); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return pages, nil
}
