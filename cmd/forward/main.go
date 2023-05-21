// Package main implements the forward tool.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/udhos/boilerplate/boilerplate"
	"github.com/udhos/gateboard/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const version = "0.1.0"

type application struct {
	me            string
	serverMain    *serverGin
	serverHealth  *serverGin
	serverMetrics *serverGin
	tracer        trace.Tracer
	config        appConfig
}

func main() {

	var showVersion bool
	flag.BoolVar(&showVersion, "version", showVersion, "show version")
	flag.Parse()

	me := filepath.Base(os.Args[0])

	{
		v := boilerplate.LongVersion(me + " version=" + version)
		if showVersion {
			fmt.Print(v)
			fmt.Println()
			return
		}
		log.Print(v)
	}

	app := &application{
		me:     me,
		config: newConfig(me),
	}

	//
	// initialize tracing
	//

	{
		tp, errTracer := tracing.TracerProvider(app.me, app.config.jaegerURL)
		if errTracer != nil {
			log.Fatalf("tracer provider: %v", errTracer)
		}

		// Register our TracerProvider as the global so any imported
		// instrumentation in the future will default to using it.
		otel.SetTracerProvider(tp)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Cleanly shutdown and flush telemetry when the application exits.
		defer func(ctx context.Context) {
			// Do not make the application hang when it is shutdown.
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				log.Fatalf("trace shutdown: %v", err)
			}
		}(ctx)

		tracing.TracePropagation()

		app.tracer = tp.Tracer(fmt.Sprintf("%s-main", me))
	}

	//
	// init application
	//

	initApplication(app, app.config.applicationAddr)

	//
	// start application server
	//

	go func() {
		log.Printf("application server: listening on %s", app.config.applicationAddr)
		err := app.serverMain.server.ListenAndServe()
		log.Printf("application server: exited: %v", err)
	}()

	//
	// start health server
	//

	app.serverHealth = newServerGin(app.config.healthAddr)

	log.Printf("registering route: %s %s", app.config.healthAddr, app.config.healthPath)
	app.serverHealth.router.GET(app.config.healthPath, func(c *gin.Context) {
		c.String(http.StatusOK, "health ok")
	})

	go func() {
		log.Printf("health server: listening on %s", app.config.healthAddr)
		err := app.serverHealth.server.ListenAndServe()
		log.Printf("health server: exited: %v", err)
	}()

	//
	// start metrics server
	//

	app.serverMetrics = newServerGin(app.config.metricsAddr)

	go func() {
		prom := promhttp.Handler()
		app.serverMetrics.router.GET(app.config.metricsPath, func(c *gin.Context) {
			prom.ServeHTTP(c.Writer, c.Request)
		})
		log.Printf("metrics server: listening on %s %s", app.config.metricsAddr, app.config.metricsPath)
		err := app.serverMetrics.server.ListenAndServe()
		log.Printf("metrics server: exited: %v", err)
	}()

	//
	// handle graceful shutdown
	//

	shutdown(app)
}

func initApplication(app *application, addr string) {

	initMetrics(app.config.metricsNamespace)

	app.serverMain = newServerGin(addr)
	app.serverMain.router.Use(middlewareMetrics(app.config.metricsMaskPath))
	app.serverMain.router.Use(otelgin.Middleware(app.me))
	app.serverMain.router.Use(gin.Logger())

	//
	// register application routes
	//

	const route = "/forward"

	log.Printf("registering route: %s %s", addr, route)

	app.serverMain.router.Any(route, func(c *gin.Context) { forward(c, app) })
}

func shutdown(app *application) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("received signal '%v', initiating shutdown", sig)

	const timeout = 5 * time.Second
	app.serverHealth.shutdown(timeout)
	app.serverMetrics.shutdown(timeout)
	app.serverMain.shutdown(timeout)

	log.Printf("exiting")
}
