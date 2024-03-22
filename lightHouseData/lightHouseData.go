package lighthousedata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"io/ioutil"
	"log"
	"sync"
	"time"

	fetchsitedata "example.com/first/fetchSiteData"
	// lighthousedata "example.com/first/lightHouseData"
)

type PageSpeedOutput struct {
	LighthouseResult struct {
		Audits struct {
			FirstContentfulPaint struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
			} `json:"first-contentful-paint"`
			LargestContentfulPaint struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
			} `json:"largest-contentful-paint"`
			Interactive struct {
				DisplayValue string  `json:"displayValue"`
				Score        float64 `json:"score"`
				NumericValue float64 `json:"numericValue"`
			} `json:"interactive"`
		} `json:"audits"`
		Categories struct {
			Performance struct {
				Score float64 `json:"score"`
			} `json:"performance"`
		} `json:"categories"`
	} `json:"lighthouseResult"`
}

func GetlightHouseData(url string) (*PageSpeedOutput, error) {
	requestURL := fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s&key=AIzaSyApBC9gblaCzWrtEBgHnZkd_B37OF49BfM&category=PERFORMANCE&strategy=MOBILE", url)

	fmt.Println(requestURL)
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("referer", "https://www.builder.io/")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	// // Decode JSON response
	// var data map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
	// 	return nil, err
	// }

	// return data, nil

	// Decode JSON response
	var output PageSpeedOutput
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	return &output, nil

}

// LightHouseData returns the PageSpeed data for a given URL. It takes a string parameter url of type string and returns a map[string]interface{} and an error.
func LightHouseData(url string) (*PageSpeedOutput, error) {

	result, err := GetlightHouseData(url)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// fmt.Println("PageSpeed data:", result)

	return result, nil

}

/////

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"sync"
// 	"time"

// 	lighthousedata "example.com/first/lightHouseData"
// )

type Page struct {
	Href string `json:"href"`
	Tags string `json:"tags"`
	Size string `json:"size,omitempty"`
}

type PageSpeedData struct {
	PageTitle        string  `json:"page_title"`
	PerformanceScore float64 `json:"performance_score"`
	PageTags         string  `json:"page_tags"`
	SiteSize         string  `json:"site_size"`
	// SiteSize         string  `json:"site_size"`
	FCPDisplay string  `json:"fcp_display"`
	FCPScore   float64 `json:"fcp_score"`
	LCPDisplay string  `json:"lcp_display"`
	LCPScore   float64 `json:"lcp_score"`
	TTIDisplay string  `json:"tti_display"`
	TTIScore   float64 `json:"tti_score"`
	TTITime    float64 `json:"tti_time"`
	PageURL    string  `json:"page_url"`
}

func CollectPageSpeedData(page Page) (PageSpeedData, error) {
	var pageSpeedData PageSpeedData

	// Fetch PageSpeed data for the page
	pageInsights, err := LightHouseData(page.Href)
	// takescreenshot.TakeScreenShot(page.Href)
	// takeSnapShot := takescreenshot.TakeSnapShot(page.Href)

	// if takeSnapShot != nil {
	// 	fmt.Println(takeSnapShot)
	// }
	// fmt.Println(takeSnapShot)

	if err != nil {
		return pageSpeedData, fmt.Errorf("error fetching PageSpeed data for %s: %v", page.Href, err)
	}

	// Extracting the data
	// Fetch the site title
	siteTitle, err := fetchsitedata.FetchSiteTitle(page.Href)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(siteTitle)

	pageSpeedData.PageTitle = siteTitle
	pageSpeedData.PageTags = page.Tags
	pageSpeedData.SiteSize = page.Size
	pageSpeedData.PageURL = page.Href
	pageSpeedData.FCPDisplay = pageInsights.LighthouseResult.Audits.FirstContentfulPaint.DisplayValue
	pageSpeedData.FCPScore = pageInsights.LighthouseResult.Audits.FirstContentfulPaint.Score
	pageSpeedData.LCPDisplay = pageInsights.LighthouseResult.Audits.LargestContentfulPaint.DisplayValue
	pageSpeedData.LCPScore = pageInsights.LighthouseResult.Audits.LargestContentfulPaint.Score
	pageSpeedData.TTIDisplay = pageInsights.LighthouseResult.Audits.Interactive.DisplayValue
	pageSpeedData.TTIScore = pageInsights.LighthouseResult.Audits.Interactive.Score
	pageSpeedData.TTITime = pageInsights.LighthouseResult.Audits.Interactive.NumericValue
	pageSpeedData.PerformanceScore = pageInsights.LighthouseResult.Categories.Performance.Score

	return pageSpeedData, nil
}

func RunWriteGenerate() {
	startTime := time.Now()

	// Read the contents of pages.json
	data, err := ioutil.ReadFile("public/pages.json")
	if err != nil {
		log.Fatalf("Error reading pages.json: %v", err)
	}

	// Parse the JSON data into a slice of Page structs
	var pages []Page
	if err := json.Unmarshal(data, &pages); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Create a channel to send and receive PageSpeed data
	pageSpeedDataChan := make(chan PageSpeedData, len(pages))

	// Create a WaitGroup to wait for all Goroutines to finish
	var wg sync.WaitGroup

	// Define rate limiting parameters
	rateLimit := 20                                           // API rate limit per second
	requestInterval := time.Second / time.Duration(rateLimit) // Interval between requests
	semaphore := make(chan struct{}, rateLimit)

	// Process pages with rate limiting
	for _, page := range pages {
		semaphore <- struct{}{} // Acquire a token from the semaphore
		wg.Add(1)
		go func(page Page) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the token back to the semaphore
			pageSpeedData, err := CollectPageSpeedData(page)
			if err != nil {
				log.Printf("Error processing page %s: %v", page.Href, err)
				return
			}
			pageSpeedDataChan <- pageSpeedData
		}(page)
		time.Sleep(requestInterval) // Wait before processing the next page
	}

	// Close the channel when all Goroutines finish
	go func() {
		wg.Wait()
		close(pageSpeedDataChan)
	}()

	// Collect processed page speed data from the channel
	var allPageSpeedData []PageSpeedData
	for data := range pageSpeedDataChan {
		allPageSpeedData = append(allPageSpeedData, data)
	}

	// Marshal the page speed data slice into JSON format
	jsonData, err := json.MarshalIndent(allPageSpeedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write the JSON data to generate.json
	if err := ioutil.WriteFile("public/generate3.json", jsonData, 0644); err != nil {
		log.Fatalf("Error writing generate.json: %v", err)
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total execution time: %s\n", elapsedTime)
}

// type pageData struct {
// 	Href string `json:"href"`
// 	Size string `json:"size,omitempty"`
// 	Tags string `json:"tags"`
// }

func RunWriteGenerate_X(pagesJSONData []byte) {
	startTime := time.Now()

	// fmt.Println(pagesJsonData)

	var pagesData []Page

	// Unmarshal the JSON data into the pages slice
	err := json.Unmarshal(pagesJSONData, &pagesData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	// fmt.Println(jsonData)
	// fmt.Println(pagesData)

	// Create a channel to send and receive PageSpeed data
	pageSpeedDataChan := make(chan PageSpeedData, len(pagesData))

	// Create a WaitGroup to wait for all Goroutines to finish
	var wg sync.WaitGroup

	// Define rate limiting parameters
	rateLimit := 15                                           // API rate limit per second
	requestInterval := time.Second / time.Duration(rateLimit) // Interval between requests
	semaphore := make(chan struct{}, rateLimit)

	// Process pages with rate limiting
	for _, page := range pagesData {
		semaphore <- struct{}{} // Acquire a token from the semaphore
		wg.Add(1)
		go func(page Page) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release the token back to the semaphore
			pageSpeedData, err := CollectPageSpeedData(page)
			if err != nil {
				log.Printf("Error processing page %s: %v", page.Href, err)
				return
			}
			pageSpeedDataChan <- pageSpeedData
		}(page)
		time.Sleep(requestInterval) // Wait before processing the next page
	}

	// Close the channel when all Goroutines finish
	go func() {
		wg.Wait()
		close(pageSpeedDataChan)
	}()

	// Collect processed page speed data from the channel
	var allPageSpeedData []PageSpeedData
	for data := range pageSpeedDataChan {
		allPageSpeedData = append(allPageSpeedData, data)
	}

	// Marshal the page speed data slice into JSON format
	jsonData, err := json.MarshalIndent(allPageSpeedData, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Write the JSON data to generate.json
	if err := ioutil.WriteFile("public/generate3.json", jsonData, 0644); err != nil {
		log.Fatalf("Error writing generate.json: %v", err)
	}

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total execution time for lightHouseData.go : %s\n", elapsedTime)
}
