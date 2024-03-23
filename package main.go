package main

// import (
// 	"bufio"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"strings"
// 	"sync"
// 	"time"

// 	lighthousedata "example.com/first/lightHouseData"
// 	takescreenshot "example.com/first/takeScreenShot"
// )

// func main() {
// 	startTime := time.Now()

// 	pagesJSONData, err := ioutil.ReadFile("public/pages.json")
// 	if err != nil {
// 		log.Fatalf("Error reading JSON file: %v", err)
// 	}

// 	// Print menu options
// 	fmt.Println("Choose an option:")
// 	fmt.Println("1) Capture Screenshots")
// 	fmt.Println("2) Generate Lighthouse Data")
// 	fmt.Println("3) Run Both")

// 	// Read user input
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Enter your choice (1/2/3): ")
// 	input, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatalf("Error reading user input: %v", err)
// 	}

// 	// Trim whitespace and convert to lower case
// 	input = strings.TrimSpace(input)
// 	input = strings.ToLower(input)

// 	// Run selected option
// 	switch input {
// 	case "1":
// 		runCaptureScreenshots(pagesJSONData)
// 	case "2":
// 		runGenerateLighthouseData(pagesJSONData)
// 	case "3":
// 		runBoth(pagesJSONData)
// 	default:
// 		fmt.Println("Invalid choice. Please choose 1, 2, or 3.")
// 	}

// 	endTime := time.Now()
// 	elapsedTime := endTime.Sub(startTime)
// 	fmt.Printf("Total execution time for main.go: %s\n", elapsedTime)
// }

// func runCaptureScreenshots(pagesJSONData []byte) {
// 	fmt.Println("Running Capture Screenshots...")
// 	takescreenshot.CaptureBatch(pagesJSONData)
// }

// func runGenerateLighthouseData(pagesJSONData []byte) {
// 	fmt.Println("Running Generate Lighthouse Data...")
// 	lighthousedata.RunWriteGenerate_X(pagesJSONData)
// }

// func runBoth(pagesJSONData []byte) {
// 	fmt.Println("Running Both Capture Screenshots and Generate Lighthouse Data...")
// 	var wg sync.WaitGroup
// 	wg.Add(2)
// 	go func() {
// 		defer wg.Done()
// 		takescreenshot.CaptureBatch(pagesJSONData)
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		lighthousedata.RunWriteGenerate_X(pagesJSONData)
// 	}()
// 	wg.Wait()
// }
