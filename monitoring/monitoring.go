package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	Requests      *prometheus.CounterVec
	Queries       *prometheus.CounterVec
	RequestDelay  prometheus.Gauge
	DatabaseDelay prometheus.Gauge
}

const (
	Successful   = "successful"
	Unsuccessful = "unsuccessful"
)

var Statistics *Metrics
var Registry *prometheus.Registry

func newMetrics() *Metrics {
	metrics := &Metrics{
		Requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "requests",
				Help: "Number of HTTP requests processed, partitioned by successful and unsuccessful requests.",
			},
			[]string{"status"},
		),
		Queries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "queries",
				Help: "Number of database queries processed, partitioned by successful and unsuccessful queris.",
			},
			[]string{"status"},
		),
		RequestDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "request_delay",
				Help: "Total delay spent on responding to the requests.",
			},
		),
		DatabaseDelay: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_delay",
				Help: "Total delay spent on waiting for the database.",
			},
		),
	}

	metrics.RequestDelay.Set(0)
	metrics.DatabaseDelay.Set(0)

	return metrics
}

func Initalize() {
	Registry = prometheus.NewRegistry()
	Statistics = newMetrics()

	if err := Registry.Register(Statistics.Requests); err != nil {
		panic(err)
	}
	if err := Registry.Register(Statistics.Queries); err != nil {
		panic(err)
	}
	if err := Registry.Register(Statistics.RequestDelay); err != nil {
		panic(err)
	}
	if err := Registry.Register(Statistics.DatabaseDelay); err != nil {
		panic(err)
	}
}
