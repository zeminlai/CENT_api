package handlers

import (
    "CENT_Notes/cmd/models"
    "CENT_Notes/cmd/repositories"
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
    "fmt"
    "CENT_Notes/cmd/storage"
)

func CreateDirectory(c echo.Context) error {
    dir := models.Directory{}
    if err := c.Bind(&dir); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid directory data",
        })
    }

    newDir, err := repositories.CreateDirectory(dir)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return c.JSON(http.StatusCreated, newDir)
}

func GetUserDirectories(c echo.Context) error {
    fmt.Println("=== Starting GetUserDirectories ===")
    
    // Test DB connection
    db := storage.GetDB()
    err := db.Ping()
    if err != nil {
        fmt.Printf("Database connection error: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": "Database connection error",
        })
    }
    fmt.Println("Database connection successful")

    userId, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        fmt.Printf("Invalid user ID: %v\n", err)
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid user ID",
        })
    }

    fmt.Printf("Fetching directories for user ID: %d\n", userId)

    // Test query
    var count int
    err = db.QueryRow("SELECT COUNT(*) FROM directories WHERE user_id = $1", userId).Scan(&count)
    if err != nil {
        fmt.Printf("Count query error: %v\n", err)
    } else {
        fmt.Printf("Found %d directories for user %d\n", count, userId)
    }

    directories, err := repositories.GetUserDirectories(userId)
    if err != nil {
        fmt.Printf("Error fetching directories: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    fmt.Printf("Successfully retrieved %d directories\n", len(directories))
    return c.JSON(http.StatusOK, directories)
}

func GetDirectory(c echo.Context) error {
    directoryId, err := strconv.Atoi(c.Param("directoryId"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid directory ID",
        })
    }

    fmt.Printf("Fetching directory ID: %d\n", directoryId)

    directory, err := repositories.GetDirectory(directoryId)
    if err != nil {
        fmt.Printf("Error fetching directory: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, directory)
} 