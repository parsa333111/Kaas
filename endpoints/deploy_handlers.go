package endpoints

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skye-tan/KaaS/k8s_client"
	custom_messages "github.com/skye-tan/KaaS/utils"
)

// POST "/deployment/create/custom"
func createCustomDeployment(c echo.Context) error {
	var deploymentRequest k8s_client.CustomDeploymentRequest

	if err := c.Bind(&deploymentRequest); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, custom_messages.InvalidBody)
	}

	ok := k8s_client.CreateCustomDeployment(deploymentRequest)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, custom_messages.DeployFailure)
	}

	return c.NoContent(http.StatusCreated)
}

// POST "/deployment/create/postgres"
func createPotsgresDeployment(c echo.Context) error {
	var deploymentRequest k8s_client.PostgresDeploymentRequest

	if err := c.Bind(&deploymentRequest); err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, custom_messages.InvalidBody)
	}

	ok := k8s_client.CreatePostgresDeployment(deploymentRequest)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, custom_messages.DeployFailure)
	}

	return c.NoContent(http.StatusCreated)
}
