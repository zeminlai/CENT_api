package handlers

import (
	"CENT_Notes/cmd/models"
	"CENT_Notes/cmd/repositories"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func HandleAudioToNote(c echo.Context) error {
	// Get the uploaded file
	file, err := c.FormFile("audio")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "No audio file provided",
		})
	}

	// Read file content
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Read file bytes
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to read file",
		})
	}

	// Send to AssemblyAI API
	transcription, err := transcribeAudio(fileBytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Transcription failed: %v", err),
		})
	}

	// Process with Groq to format the note
	groqResponse, err := processWithGroq(transcription)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Groq processing failed: %v", err),
		})
	}

	// Create a new note
	userId := c.QueryParam("userId")
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}
	
	note := models.Note{
		Title:   "Audio Note",
		Content: groqResponse,
		UserId:  userIdInt,
	}

	newNote, err := repositories.CreateNote(note)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create note",
		})
	}

	return c.JSON(http.StatusOK, newNote)
}

func transcribeAudio(audioData []byte) (string, error) {
	apiKey := os.Getenv("ASSEMBLY_API_KEY")
	url := "https://api.assemblyai.com/v2/transcript"

	// Create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(audioData))
	if err != nil {
		return "", err
	}

	// Set headers
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "audio/mpeg")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse response
	var result TranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Text, nil
}

func processWithGroq(text string) (string, error) {
	groqReq := GroqRequest{
		Model: "llama3-70b-8192",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Please format this transcription as a well-structured note: " + text,
			},
		},
	}

	jsonData, err := json.Marshal(groqReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GROQ_API_KEY")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Extract the formatted text from Groq's response
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				return content, nil
			}
		}
	}

	return "", fmt.Errorf("failed to parse Groq response")
} 