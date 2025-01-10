package handlers

import (
	"CENT_Notes/cmd/models"
	"CENT_Notes/cmd/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	aai "github.com/AssemblyAI/assemblyai-go-sdk"
	"github.com/labstack/echo/v4"
)

type TranscriptionResponse struct {
	Text string `json:"text"`
}

func HandleAudioToNote(c echo.Context) error {
	log.Println("Starting audio to note conversion process")
	
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

	// After transcription
	log.Printf("Transcription result: %s", transcription)

	// Process with Groq to format the note
	groqResponse, err := processWithGroq(transcription)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Groq processing failed: %v", err),
		})
	}

	// After Groq processing
	log.Printf("Formatted note from Groq: %s", groqResponse)

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

	// After note creation
	log.Printf("Note created successfully with ID: %d", newNote.Id)

	return c.JSON(http.StatusOK, newNote)
}

func transcribeAudio(audioData []byte) (string, error) {
	log.Println("Starting audio transcription with AssemblyAI")
	tempFile, err := os.CreateTemp("", "audio-*")
	if err != nil {
			return "", fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up temp file when done
	
	// Write audio data to temp file
	if _, err := tempFile.Write(audioData); err != nil {
			return "", fmt.Errorf("failed to write audio data to temp file: %v", err)
	}
	tempFile.Close()
	f,err := os.Open(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to open temp file: %v", err)
	}
	apiKey := os.Getenv("ASSEMBLY_API_KEY")
	client := aai.NewClient(apiKey)

	// Create transcription request
	transcript, err := client.Transcripts.TranscribeFromReader(context.Background(), f, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create transcript: %v", err)
	}

	return *transcript.Text, nil
}

func processWithGroq(text string) (string, error) {
	log.Println("Starting Groq processing")
	
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

	// After response
	log.Printf("Groq API response status: %d", resp.StatusCode)

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