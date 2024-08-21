package endpoints

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skye-tan/KaaS/k8s_client"
	custom_messages "github.com/skye-tan/KaaS/utils"
)

// GET "/deployment/status/:deployment_name"
func getDeploymentStatus(c echo.Context) error {
	deploymentName := c.Param("deployment_name")

	deploymentStatus, ok := k8s_client.GetDeploymentStatus(deploymentName)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, custom_messages.FetchFailure)
	}

	return c.JSON(http.StatusOK, deploymentStatus)
}

// GET "/deployment/status"
func getDeploymentsStatus(c echo.Context) error {
	deploymentStatuses, ok := k8s_client.GetDeploymentsStatus()
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, custom_messages.FetchFailure)
	}

	return c.JSON(http.StatusOK, deploymentStatuses)
}
