package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

const tracerName = "gin-http-server"

func TracingMiddleware() gin.HandlerFunc {
	tracer := otel.Tracer(tracerName)
	propagator := otel.GetTextMapPropagator()

	return func(c *gin.Context) {
		ctx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		operation := fmt.Sprintf("HTTP %s %s", c.Request.Method, c.FullPath())

		ctx, span := tracer.Start(ctx, operation,
			trace.WithAttributes(
				semconv.HTTPMethod(c.Request.Method),
				semconv.HTTPRoute(c.FullPath()),
				semconv.HTTPClientIP(c.ClientIP()),
				attribute.String("http.host", c.Request.Host),
				attribute.String("http.scheme", c.Request.URL.Scheme),
				attribute.String("http.target", c.Request.URL.Path),
				attribute.String("http.proto", c.Request.Proto),
				attribute.String("http.referer", c.Request.Referer()),
			),
		)
		defer span.End()

		start := time.Now()
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		span.SetAttributes(
			semconv.HTTPStatusCode(c.Writer.Status()),
			attribute.String("http.response_content_type", c.Writer.Header().Get("Content-Type")),
			attribute.Int64("http.duration_ms", time.Since(start).Milliseconds()),
		)
		if c.Writer.Status() >= http.StatusBadRequest {
			span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", c.Writer.Status()))
			if len(c.Errors) > 0 {
				span.SetAttributes(attribute.String("error.details", c.Errors.String()))
			}
		} else {
			span.SetStatus(codes.Ok, "")
		}
		if userID, exists := c.Get("user_id"); exists {
			span.SetAttributes(attribute.String("business.user_id", fmt.Sprintf("%v", userID.(uuid.UUID).String())))
		}
	}
}
