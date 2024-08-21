package endpoints

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	requestFailure     = "request failed"
	invalidContentType = "invalid content type"
	invalidBodyFormat  = "invalid body format"
	missingData        = "missing data"
)

func customLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return middleware.Logger()(next)(c)
	}
}

func Start(listenAddress string) {
	e := echo.New()
	e.Use(customLogger)

	// Healthz Endpoint
	e.GET("/healthz", healthCheack)

	// Note Endpoints
	e.GET("/note", getNotes)
	e.POST("/note", addNote)

	e.Logger.Fatal(e.Start(listenAddress))
}
