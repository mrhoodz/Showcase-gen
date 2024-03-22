package fetchsitedata

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func FetchSiteTitle(siteURL string) (string, error) {
	// Make an HTTP GET request to example.com
	resp, err := http.Get(siteURL)
	if err != nil {
		return "", fmt.Errorf("error fetching site data: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML response
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error parsing HTML: %v", err)
	}

	// Find the title tag in the HTML document
	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
		}
	}
	findTitle(doc)

	// Return the title
	return title, nil
}
