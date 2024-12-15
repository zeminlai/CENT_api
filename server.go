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
	e.POST("/llm/groq", handlers.GroqHandler)
	storage.InitDB()
	e.Logger.Fatal(e.Start(":1323"))
}
