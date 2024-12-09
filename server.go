package main

import (
	"net/http"

	"CENT_Notes/cmd/handlers"
	"CENT_Notes/cmd/storage"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
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
	e.GET("/users/:userId/notes", handlers.GetUserNotes)
	e.PUT("/notes/:id", handlers.UpdateNote)
	e.DELETE("/notes/:id/users/:userId", handlers.DeleteNote)
	
	storage.InitDB()
	e.Logger.Fatal(e.Start(":1323"))
}
