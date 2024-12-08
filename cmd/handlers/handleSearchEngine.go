package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

type SearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type WebPageContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// HandleGoogleSearch performs Google search and returns top 10 results
func HandleGoogleSearch(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	// Format Google search URL
	searchURL := fmt.Sprintf("https://www.google.com/search?q=%s", strings.ReplaceAll(query, " ", "+"))

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request",
		})
	}

	// Set User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch search results",
		})
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse HTML",
		})
	}

	var results []SearchResult
	// Extract search results
	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		if i >= 10 { // Limit to top 10 results
			return
		}

		title := s.Find("h3").Text()
		description := s.Find("div.VwiC3b").Text()
		url, _ := s.Find("a").Attr("href")

		if title != "" && description != "" && url != "" {
			results = append(results, SearchResult{
				Title:       title,
				Description: description,
				URL:         url,
			})
		}
	})

	// Save results to JSON file
	_, err = json.MarshalIndent(results, "", "  ")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create JSON",
		})
	}

	return c.JSON(http.StatusOK, results)
}

// HandleWebScrape scrapes content from a selected search result
func HandleWebScrape(c echo.Context) error {
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "URL is required",
		})
	}

	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request",
		})
	}

	// Set User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch webpage",
		})
	}
	defer resp.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse HTML",
		})
	}

	// Extract content
	title := doc.Find("title").Text()
	content := doc.Find("body").Text()

	// Clean up content (remove extra whitespace)
	content = strings.Join(strings.Fields(content), " ")

	result := WebPageContent{
		Title:   title,
		Content: content,
		URL:     url,
	}

	return c.JSON(http.StatusOK, result)
}