package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

// GroqRequest represents the structure of the request body sent to Groq API
type GroqRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestBody represents the incoming request body
type RequestBody struct {
	Prompt string `json:"prompt"`
}

// GroqHandler handles the POST request for generating completions using Echo
func GroqHandler(c echo.Context) error {
	// Parse request body
	var reqBody RequestBody
	if err := c.Bind(&reqBody); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to parse request body",
		})
	}

	// Prepare request to Groq API
	groqReq := GroqRequest{
		Model: "llama3-70b-8192",
		Messages: []Message{
			{
				Role: "user",
				Content: ".please and must respond in html snippet format.only use div, bold, underline, " +
					"numbering and lists. you are a note taking assistant that only returns responses in html format." +
					"if the response is code, please do so but dont use <code> and <pre>." +
					"just answer the prompt. this is the prompt: " + reqBody.Prompt,
			},
		},
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(groqReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request",
		})
	}

	// Create request to Groq API
	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request",
		})
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GROQ_API_KEY")))

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate completion",
		})
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read response",
		})
	}

	// Parse Groq API response
	var groqResp map[string]interface{}
	if err := json.Unmarshal(body, &groqResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse response",
		})
	}

	// Extract content from response
	choices, ok := groqResp["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid response format",
		})
	}

	firstChoice, ok := choices[0].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid response format",
		})
	}

	message, ok := firstChoice["message"].(map[string]interface{})
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid response format",
		})
	}

	content, ok := message["content"].(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Invalid response format",
		})
	}

	// Return the response
	return c.JSON(http.StatusOK, content)
}
