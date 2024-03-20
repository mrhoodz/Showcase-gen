package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// Create a new context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second) // Increased timeout to 30 seconds
	defer cancel()

	// Retry mechanism with maximum 3 retries
	var buf []byte
	var err error
	for i := 0; i < 3; i++ {
		if err = chromedp.Run(ctx, fullScreenshot("https://qwik-snowy.vercel.app/", 100, &buf)); err == nil {
			break
		}
		log.Printf("Attempt %d failed: %v. Retrying...", i+1, err)
		time.Sleep(3 * time.Second) // Wait 3 seconds before retrying
	}
	if err != nil {
		log.Fatal("Failed after multiple attempts:", err)
	}

	// Save the screenshot to a file
	if err := ioutil.WriteFile("realwebsdigital.png", buf, 0644); err != nil {
		log.Fatal(err)
	}
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(2 * time.Second),     // Give the page some time to load
		chromedp.EmulateViewport(1440, 980), // Set viewport size to 1440x980
		chromedp.FullScreenshot(res, quality),
	}
}
