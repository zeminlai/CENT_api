package handlers

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

const utmPassage = `Universiti Teknologi Malaysia (UTM) is a leading public research university in Malaysia, specializing in engineering, science, and technology. Established in 1904 as a technical school, UTM has grown into a premier institution offering a wide range of undergraduate and postgraduate programs. With campuses in Johor Bahru and Kuala Lumpur, it is renowned for its innovative research, industry collaboration, and global partnerships, making it a hub for academic excellence and technological advancements in the region.`

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
	results := []SearchResult{}

	// Show all UTM results for both link click and search query
	if query != "" || c.QueryParam("link") == "utm" {
		results = []SearchResult{
			{
				Title:       "About UTM",
				Description: utmPassage,
				URL:         "https://en.wikipedia.org/wiki/University_of_Technology_Malaysia",
			},
			{
				Title:       "UTM Johor Campus",
				Description: "UTM's main campus in Johor Bahru spans over 1,222 hectares of land, offering state-of-the-art facilities and a vibrant academic environment.",
				URL:         "https://www.utm.my/johor",
			},
			{
				Title:       "UTM Kuala Lumpur Campus",
				Description: "Located in the heart of Malaysia's capital, UTM's KL campus provides urban-centric education and research opportunities.",
				URL:         "https://www.utm.my/kl",
			},
			{
				Title:       "UTM International",
				Description: "UTM welcomes international students and collaborations, fostering a diverse and global academic community.",
				URL:         "https://www.utm.my/international",
			},
		}
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
