package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type SearchResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Google API response types
type GoogleSearchResponse struct {
	Kind              string       `json:"kind"`
	URL               SearchInfo   `json:"url"`
	Queries           QueryData    `json:"queries"`
	Context           SearchInfo   `json:"context"`
	SearchInformation SearchInfo   `json:"searchInformation"`
	Items             []SearchItem `json:"items"`
}

type SearchInfo map[string]interface{}

type QueryData map[string][]SearchInfo

type SearchItem struct {
	Kind             string     `json:"kind"`
	Title            string     `json:"title"`
	HTMLTitle        string     `json:"htmlTitle"`
	Link             string     `json:"link"`
	DisplayLink      string     `json:"displayLink"`
	Snippet          string     `json:"snippet"`
	HTMLSnippet      string     `json:"htmlSnippet"`
	FormattedURL     string     `json:"formattedUrl"`
	HTMLFormattedURL string     `json:"htmlFormattedUrl"`
	PageMap          PageMapData `json:"pagemap,omitempty"`
}

type PageMapData map[string]interface{}

type WebPageContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	URL     string `json:"url"`
}

// HandleGoogleSearch performs Google search using the Google Custom Search API
func HandleGoogleSearch(c echo.Context) error {
	query := c.QueryParam("q")
	fmt.Println("HandleGoogleSearch called with query:", query)

	if query == "" {
		fmt.Println("Error: Search query is required")
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Search query is required",
		})
	}

	// Use Google Custom Search API
	results, err := fetchGoogleSearchResults(query)
	if err != nil {
		fmt.Println("Error fetching search results:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch search results: %v", err),
		})
	}

	fmt.Println("Search successful. Found", len(results), "results")
	fmt.Println("Sending response with status 200")
	return c.JSON(http.StatusOK, results)
}

// fetchGoogleSearchResults retrieves search results from Google Custom Search API
func fetchGoogleSearchResults(searchTerm string) ([]SearchResult, error) {
	maxResults := 10
	client := &http.Client{}

	// Build the Google Custom Search API URL
	baseURL := "https://www.googleapis.com/customsearch/v1"
	params := url.Values{}
	params.Add("key", "AIzaSyDaU1kU_RSH1RWgkLXvYmU_umVUfIEU8bU") // Your API key
	params.Add("cx", "57bf992fba6c74a45")                        // Your custom search engine ID
	params.Add("q", searchTerm)

	// Build the final URL
	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Parse the JSON response
	var googleResponse GoogleSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&googleResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Convert to our SearchResult format
	var results []SearchResult
	for _, item := range googleResponse.Items {
		if len(results) >= maxResults {
			break
		}

		// Debug info
		fmt.Printf("Found result: Title: %s, URL: %s\n", item.Title, item.Link)
		
		// Convert API response to our model
		result := SearchResult{
			Title:       item.Title,
			Description: item.Snippet,
			URL:         item.Link,
		}
		results = append(results, result)
	}
	return results, nil
}

// HandleWebScrape scrapes content from a selected search result
func HandleWebScrape(c echo.Context) error {
	url := c.QueryParam("url")
	fmt.Println("HandleWebScrape called with URL:", url)

	if url == "" {
		fmt.Println("Error: URL is required")
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

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read webpage content",
		})
	}
	
	// Basic extraction of title and content using simple string operations
	bodyString := string(bodyBytes)
	title := extractTitle(bodyString)
	content := extractBodyContent(bodyString)
	
	// Clean up content (remove extra whitespace)
	content = strings.Join(strings.Fields(content), " ")
	result := WebPageContent{
		Title:   title,
		Content: content,
		URL:     url,
	}

	fmt.Println("Web scrape successful. Title length:", len(title), "Content length:", len(content))

	fmt.Println("Sending response with status 200")
	return c.JSON(http.StatusOK, result)
}

// Helper functions for HTML parsing without goquery
func extractTitle(html string) string {
	titleStart := strings.Index(html, "<title")
	if titleStart == -1 {
		return ""
	}
	
	// Find the closing > of the title tag
	titleStart = strings.Index(html[titleStart:], ">")
	if titleStart == -1 {
		return ""
	}
	titleStart += 1 // Move past the > character
	
	// Find the end of the title
	titleEnd := strings.Index(html[titleStart:], "</title>")
	if titleEnd == -1 {
		return ""
	}
	
	// Extract the title text
	return html[titleStart : titleStart+titleEnd]
}

func extractBodyContent(html string) string {
	bodyStart := strings.Index(html, "<body")
	if bodyStart == -1 {
		// If no body tag is found, return a portion of the document
		if len(html) > 1000 {
			return html[:1000] // Return first 1000 characters
		}
		return html
	}
	
	// Find the closing > of the body tag
	bodyTagEnd := strings.Index(html[bodyStart:], ">")
	if bodyTagEnd == -1 {
		return ""
	}
	bodyStart += bodyTagEnd + 1
	
	// Find the end of the body
	bodyEnd := strings.Index(html[bodyStart:], "</body>")
	if bodyEnd == -1 {
		// If no closing body tag, just return from body start to the end
		return html[bodyStart:]
	}
	
	// Get raw body content 
	bodyContent := html[bodyStart : bodyStart+bodyEnd]
	
	// Simple HTML tag removal (this is basic, not a complete solution)
	return stripHTMLTags(bodyContent)
}

func stripHTMLTags(html string) string {
	insideTag := false
	var result strings.Builder
	
	for _, r := range html {
		if r == '<' {
			insideTag = true
			continue
		}
		if r == '>' {
			insideTag = false
			result.WriteRune(' ') // Add space where tags were
			continue
		}
		if !insideTag {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}
