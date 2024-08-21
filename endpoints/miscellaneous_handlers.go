package endpoints

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/skye-tan/KaaS/database"
	"github.com/skye-tan/KaaS/monitoring"
	custom_messages "github.com/skye-tan/KaaS/utils"
)

// GET "/healthz"
func healthCheack(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// GET "/deployment/health/:deployment_name"
func getDeploymentHealth(c echo.Context) error {
	deploymentName := c.Param("deployment_name")

	startTime := time.Now()

	healthCheck, ok := database.GetDeploymentHealth(deploymentName)

	delayDelay := time.Since(startTime)
	monitoring.Statistics.DatabaseDelay.Add(delayDelay.Seconds())

	if !ok {
		monitoring.Statistics.Queries.WithLabelValues(monitoring.Unsuccessful).Inc()
		return echo.NewHTTPError(http.StatusInternalServerError, custom_messages.RequestFailure)
	}
	monitoring.Statistics.Queries.WithLabelValues(monitoring.Successful).Inc()

	return c.JSON(http.StatusOK, healthCheck)
}
