package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	otel "go-product-service/drivers/tracing"
	ioOtel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// TraceMiddleware is function to create middleware trace
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := ioOtel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		c.Request = c.Request.WithContext(ctx)

		ctx, span := otel.Start(c, fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()))
		defer span.End()

		traceID := span.SpanContext().TraceID().String()
		span.SetAttributes(attribute.String("traceID", traceID))
		ctx = context.WithValue(ctx, "traceID", traceID)

		w := NewResponseBodyWriter(c.Writer)
		c.Writer = w
		c.Request = c.Request.WithContext(ctx)

		// continue to next handler
		c.Next()

		// get response and set to span information
		span.SetAttributes(
			attribute.Int("http.response.status", c.Writer.Status()),
			attribute.String("http.response.body", w.body.String()))

		w.PutBack()
	}
}
