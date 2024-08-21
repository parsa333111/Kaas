package monitoring

import (
	"time"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func StatisticsCollectorMiddleware() echo.MiddlewareFunc {
	return echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Registerer: Registry,
		BeforeNext: func(c echo.Context) {
			c.Set("start_time", time.Now())
		},
		AfterNext: func(c echo.Context, err error) {
			startTime := c.Get("start_time").(time.Time)
			requestDelay := time.Since(startTime)
			Statistics.RequestDelay.Add(requestDelay.Seconds())

			if err != nil {
				Statistics.Requests.WithLabelValues(Unsuccessful).Inc()
			} else {
				Statistics.Requests.WithLabelValues(Successful).Inc()
			}
		},
	})
}
