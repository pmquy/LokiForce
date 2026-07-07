package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func Metrics(appName string) gin.HandlerFunc {
	meter := otel.GetMeterProvider().Meter(appName)

	// Create instruments
	requestDuration, err := meter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP server requests."),
		metric.WithUnit("s"),
	)
	if err != nil {
		otel.Handle(err)
	}

	activeRequests, err := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP server requests."),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		otel.Handle(err)
	}

	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		activeRequests.Add(ctx, 1)
		defer activeRequests.Add(ctx, -1)

		c.Next()

		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		duration := time.Since(start).Seconds()

		attrs := attribute.NewSet(
			attribute.String("http.route", path),
			attribute.String("http.request.method", method),
			attribute.Int("http.response.status_code", status),
			attribute.String("service.name", appName),
		)

		requestDuration.Record(ctx, duration, metric.WithAttributeSet(attrs))
	}
}
