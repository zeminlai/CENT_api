package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	"github.com/labstack/echo/v4"
)

type SearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Content     string `json:"content"` // Add content field
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

func limitWords(content string, limit int) string {
	words := strings.Fields(content)
	if len(words) <= limit {
		return content
	}
	return strings.Join(words[:limit], " ") + "..."
}

// HandleWebScrape scrapes content from a selected search result
func HandleWebScrape(c echo.Context) error {
	urlStr := c.QueryParam("url")
	if urlStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "URL is required",
		})
	}

	// Parse URL string to url.URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid URL",
		})
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Fetch the webpage
	resp, err := client.Get(urlStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch webpage",
		})
	}
	defer resp.Body.Close()

	// Parse using readability with correct URL type
	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse content",
		})
	}

	// Limit content to 200 words
	limitedContent := limitWords(article.TextContent, 500)

	return c.JSON(http.StatusOK, SearchResult{
		Title:       article.Title,
		Description: article.Excerpt,
		Content:     limitedContent,
		URL:         urlStr,
	})
}
