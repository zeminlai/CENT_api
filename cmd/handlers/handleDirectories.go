package handlers

import (
    "CENT_Notes/cmd/models"
    "CENT_Notes/cmd/repositories"
    "net/http"
    "strconv"
    "github.com/labstack/echo/v4"
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
    userId, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "Invalid user ID",
        })
    }

    directories, err := repositories.GetUserDirectories(userId)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    return c.JSON(http.StatusOK, directories)
} 