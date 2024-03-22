package getlighthousedata

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func GetlightHouseData(url string) (map[string]interface{}, error) {
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

	// Decode JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

// // LightHouseData returns the PageSpeed data for a given URL. It takes a string parameter url of type string and returns a map[string]interface{} and an error.
// func LightHouseData(url string) (map[string]interface{}, error) {

// 	result, err := GetlightHouseData(url)
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return nil, err
// 	}
// 	fmt.Println("PageSpeed data:", result)

// 	return result, nil

// }
