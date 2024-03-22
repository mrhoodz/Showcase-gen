package readfile

import (
	"bufio"
	"fmt"
	"os"
)

// "href":
// "size":
// "tags":
// type pagesData struct {
// 	href string
// 	size string
// 	tags string
// }

func ReadFile() (string, error) {

	// Open the file
	file, err := os.Open("public/example.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return "", fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	// data := []pagesData{}
	var fileData string
	// Read the file line by line
	for scanner.Scan() {
		// Print the line
		//  data = append(data, )
		fileData = fileData + scanner.Text() + "\n"
		fmt.Println(fileData)
		// fileData
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}

	return fileData, nil
}
