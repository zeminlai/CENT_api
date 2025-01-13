package main

import (
	"CENT_Notes/cmd/handlers"
	"CENT_Notes/cmd/storage"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize database first
	storage.InitDB()

	e := echo.New()

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/users", handlers.CreateUser)
	e.GET("/users/:id", handlers.GetUser)
	e.PUT("/users/:id", handlers.UpdateUser)
	e.DELETE("/users/:id", handlers.DeleteUser)

	// Web Scrape 
	e.GET("/search", handlers.HandleGoogleSearch)
	e.GET("/scrape", handlers.HandleWebScrape)
	
	// Notes routes
	e.POST("/notes", handlers.CreateNote)
	e.GET("/notes/:id", handlers.GetNote)
	e.GET("/directory/:directoryId/notes", handlers.GetNotesByDirectory)
	e.GET("/users/:userId/notes", handlers.GetUserNotes)
	e.PUT("/notes/:id", handlers.UpdateNote)
	e.DELETE("/notes/:id/users/:userId", handlers.DeleteNote)
	
	e.POST("/llm/groq", handlers.GroqHandler)

	// Directories routes
	e.POST("/directories", handlers.CreateDirectory)
	e.GET("/users/:userId/directories", handlers.GetUserDirectories)
	e.GET("/directories/:directoryId", handlers.GetDirectory)

	// Audio to note
	e.POST("/audio-note", handlers.HandleAudioToNote)

	e.Logger.Fatal(e.Start(":1323"))
}
