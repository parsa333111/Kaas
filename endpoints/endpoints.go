package endpoints

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/skye-tan/KaaS/monitoring"
)

func customLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/metrics" || c.Path() == "/healthz" {
			return next(c)
		}
		return middleware.Logger()(next)(c)
	}
}

func Start(listenAddress string) {
	e := echo.New()
	e.Use(customLogger)

	// Healthz Endpoint
	e.GET("/healthz", healthCheack)

	// Metrics Endpoint
	e.GET("/metrics", echoprometheus.NewHandlerWithConfig(
		echoprometheus.HandlerConfig{
			Gatherer: monitoring.Registry,
		}),
	)

	api := e.Group("/api")

	// Statistics Collector Middleware
	api.Use(monitoring.StatisticsCollectorMiddleware())

	// Deploy Endpoints
	api.POST("/deployment/create/custom", createCustomDeployment)
	api.POST("/deployment/create/postgres", createPotsgresDeployment)

	// Status Endpoints
	api.GET("/deployment/status/:deployment_name", getDeploymentStatus)
	api.GET("/deployment/status", getDeploymentsStatus)

	// Deployment's Health Endpoint
	api.GET("/deployment/health/:deployment_name", getDeploymentHealth)

	e.Logger.Fatal(e.Start(listenAddress))
}
