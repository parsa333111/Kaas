package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/skye-tan/KaaS/tools/notepad/database"
)

// GET "/healthz"
func healthCheack(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// GET "/note"
func getNotes(c echo.Context) error {
	notes, ok := database.GetNotes()
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, requestFailure)
	}

	return c.JSON(http.StatusOK, notes)
}

// POST "/note"
func addNote(c echo.Context) error {
	content_type := c.Request().Header.Get(echo.HeaderContentType)
	if content_type != "application/json" {
		return echo.NewHTTPError(http.StatusBadRequest, invalidContentType)
	}

	content := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&content)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, invalidBodyFormat)
	}

	description, ok := content["description"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, missingData)
	}

	ok = database.AddNote(description)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, requestFailure)
	}

	return c.NoContent(http.StatusCreated)
}
