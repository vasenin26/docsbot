package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// RequestsTotal counts requests processed by the server labeled by path, method and HTTP status code.
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "docsbot_requests_total",
			Help: "Total number of processed requests",
		},
		[]string{"path", "method", "code"},
	)
)

func init() {
	// Register metrics with Prometheus default registry.
	prometheus.MustRegister(RequestsTotal)
}
