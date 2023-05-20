package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func newSpanGin(c *gin.Context, caller string, app *application) (context.Context, trace.Span) {
	ctx := c.Request.Context()
	return newSpan(ctx, caller, app)
}

func newSpan(ctx context.Context, caller string, app *application) (context.Context, trace.Span) {
	if app.tracer == nil {
		return ctx, nil
	}
	newCtx, span := app.tracer.Start(ctx, caller)
	return newCtx, span
}
