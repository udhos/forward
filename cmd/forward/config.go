package main

import "github.com/udhos/boilerplate/envconfig"

type appConfig struct {
	debug                 bool
	logDriver             string // anything other than "zap" enables gin default logger
	jaegerURL             string
	applicationAddr       string
	healthAddr            string
	healthPath            string
	metricsAddr           string
	metricsPath           string
	metricsMaskPath       bool
	metricsNamespace      string
	metricsBucketsLatency []float64
}

func newConfig(roleSessionName string) appConfig {

	env := envconfig.NewSimple(roleSessionName)

	return appConfig{
		debug:                 env.Bool("DEBUG", true),
		logDriver:             env.String("LOG_DRIVER", ""), // anything other than "zap" enables gin default logger
		jaegerURL:             env.String("JAEGER_URL", "http://jaeger-collector:14268/api/traces"),
		applicationAddr:       env.String("LISTEN_ADDR", ":8080"),
		healthAddr:            env.String("HEALTH_ADDR", ":8888"),
		healthPath:            env.String("HEALTH_PATH", "/health"),
		metricsAddr:           env.String("METRICS_ADDR", ":3000"),
		metricsPath:           env.String("METRICS_PATH", "/metrics"),
		metricsMaskPath:       env.Bool("METRICS_MASK_PATH", true),
		metricsNamespace:      env.String("METRICS_NAMESPACE", ""),
		metricsBucketsLatency: env.Float64Slice("METRICS_BUCKETS_LATENCY", []float64{0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5, 10}),
	}
}
