package main

import "github.com/udhos/boilerplate/envconfig"

type appConfig struct {
	debug            bool
	logDriver        string // anything other than "zap" enables gin default logger
	jaegerURL        string
	applicationAddr  string
	healthAddr       string
	healthPath       string
	metricsAddr      string
	metricsPath      string
	metricsMaskPath  bool
	metricsNamespace string
}

func newConfig(roleSessionName string) appConfig {

	env := envconfig.NewSimple(roleSessionName)

	return appConfig{
		debug:            env.Bool("DEBUG", true),
		logDriver:        env.String("LOG_DRIVER", ""), // anything other than "zap" enables gin default logger
		jaegerURL:        env.String("JAEGER_URL", "http://jaeger-collector:14268/api/traces"),
		applicationAddr:  env.String("LISTEN_ADDR", ":8080"),
		healthAddr:       env.String("HEALTH_ADDR", ":8888"),
		healthPath:       env.String("HEALTH_PATH", "/health"),
		metricsAddr:      env.String("METRICS_ADDR", ":3000"),
		metricsPath:      env.String("METRICS_PATH", "/metrics"),
		metricsMaskPath:  env.Bool("METRICS_MASK_PATH", true),
		metricsNamespace: env.String("METRICS_NAMESPACE", ""),
	}
}
