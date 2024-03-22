package takescreenshot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

type Page struct {
	Href string `json:"href"`
	Tags string `json:"tags"`
	Size string `json:"size,omitempty"`
}

func TakeScreenShot(ctx context.Context, siteURL string) error {
	startTime := time.Now()

	// Navigate to the specified URL
	if err := chromedp.Run(ctx, chromedp.Navigate(siteURL)); err != nil {
		return fmt.Errorf("failed to navigate: %v", err)
	}

	if err := chromedp.Run(ctx, chromedp.EmulateViewport(1490, 930)); err != nil {
		return fmt.Errorf("failed to emulate viewport: %v", err)
	}

	// chromedp.EmulateViewport(1920, 1080)
	// Wait for the page to load
	if err := chromedp.Run(ctx, chromedp.Sleep(7*time.Second)); err != nil {
		return fmt.Errorf("failed to wait: %v", err)
	}

	// Take a screenshot of the entire page
	var buf []byte
	// if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
	// 	return fmt.Errorf("failed to take screenshot: %v", err)
	// }

	// var buf []byte
	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 50)); err != nil {
		return fmt.Errorf("failed to take fullshot: %v", err)
	}

	// hh := chromedp.FullScreenshot(&buf, 50)
	// println(hh.Do())

	// Sanitize the title to remove invalid characters
	sanitizedTitle := sanitizeFileName(siteURL)

	// Save the screenshot to a file
	file, err := os.Create("public/showcase/" + sanitizedTitle + ".webp")
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// if _, err := file.Write(buf); err != nil {
	// 	return fmt.Errorf("failed to write to file: %v", err)
	// }

	// Screenshot captured successfully
	fmt.Println("Screenshot saved to public/showcase/" + sanitizedTitle + ".png")

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Execution time for takeScreenShot for %s: %s\n", siteURL, elapsedTime)

	return nil
}

// sanitizeFileName replaces invalid characters in a filename with underscores
func sanitizeFileName(name string) string {
	// Replace 'https://' with an empty string
	name = regexp.MustCompile(`^https://`).ReplaceAllString(name, "")

	// Replace slashes ('/') with underscores ('_')
	name = regexp.MustCompile(`/`).ReplaceAllString(name, "_")

	// Replace dots ('.') with underscores ('_')
	name = regexp.MustCompile(`\.`).ReplaceAllString(name, "_")

	// Convert to lowercase
	name = regexp.MustCompile(`[A-Z]`).ReplaceAllStringFunc(name, func(s string) string {
		return string([]byte{s[0] + 32})
	})

	return name
}

func CaptureBatch(pagesJSONData []byte) {
	startTime := time.Now()

	var pagesData []Page

	// Unmarshal the JSON data into the pages slice
	if err := json.Unmarshal(pagesJSONData, &pagesData); err != nil {
		log.Fatalf("Error unmarshalling JSON data: %v", err)
	}

	// Limit the number of concurrent goroutines
	concurrencyLimit := 7
	sem := make(chan struct{}, concurrencyLimit)

	// Create a context for the entire batch
	ctx, cancel := chromedp.NewContext(context.Background())
	// chromedp.EmulateViewport(1920, 1080)
	// chromedp.Headless()
	defer cancel()

	// Iterate over pages and capture screenshots concurrently
	var wg sync.WaitGroup
	for _, page := range pagesData {
		sem <- struct{}{} // acquire semaphore
		wg.Add(1)
		go func(page Page) {
			defer func() {
				<-sem // release semaphore
				wg.Done()
			}()
			fmt.Println("--------------")
			fmt.Printf("started takescreenShot for: %s\n", page.Href)
			// Create a new context for each page
			pageCtx, pageCancel := chromedp.NewContext(ctx)
			defer pageCancel()

			// Capture screenshot
			if err := TakeScreenShot(pageCtx, page.Href); err != nil {
				log.Printf("Error taking screenshot for %s: %v\n", page.Href, err)
			}
			fmt.Printf("Done taking takescreenShot for: %s\n", page.Href)
			fmt.Println("--------------")
		}(page)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Total execution time for takescreenShot.go [CaptureBatch()] : %s\n", elapsedTime)
}
