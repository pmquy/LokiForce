package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Logger(log *slog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()

		c.Next()

		latency := time.Since(start)

		requestID, _ := c.Get("request_id")

		span := trace.SpanFromContext(c.Request.Context())

		ctx := span.SpanContext()

		traceID := ctx.TraceID().String()

		spanID := ctx.SpanID().String()

		log.Info(
			"http_request",

			slog.String("request_id", requestID.(string)),

			slog.String("trace_id", traceID),

			slog.String("span_id", spanID),

			slog.String("method", c.Request.Method),

			slog.String("path", c.FullPath()),

			slog.Int("status", c.Writer.Status()),

			slog.String("ip", c.ClientIP()),

			slog.Int("size", c.Writer.Size()),

			slog.Duration("latency", latency),
		)
	}
}
