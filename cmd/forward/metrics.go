package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	latencySpring *prometheus.HistogramVec
}

var (
	dimensionsSpring = []string{"method", "status", "uri"}
	metric           *metrics
)

func initMetrics(namespace string, latencyBuckets []float64) {
	if metric != nil {
		return
	}
	metric = newMetrics(namespace, latencyBuckets)
}

func newMetrics(namespace string, latencyBuckets []float64) *metrics {
	return &metrics{
		latencySpring: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_server_requests_seconds",
				Help:      "Spring-like server request duration in seconds.",
				Buckets:   latencyBuckets,
			},
			dimensionsSpring,
		),
	}
}

// middleware provides a gin middleware for exposing prometheus metrics.
func middlewareMetrics(metricsMaskPath bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		elap := time.Since(start)

		status := strconv.Itoa(c.Writer.Status())
		elapsedSeconds := float64(elap) / float64(time.Second)

		var path string
		if metricsMaskPath {
			path = c.FullPath() // masked path: GET /gateway/*gateway_name
			// FullPath returns a matched route full path. For not found routes returns an empty string.
			if path == "" {
				path = "NO_ROUTE_HANDLER"
			}
		} else {
			path = c.Request.URL.Path // actual path: GET /gateway/actual_gateway_name
		}

		metric.latencySpring.WithLabelValues(c.Request.Method, status, path).Observe(elapsedSeconds)
	}
}
