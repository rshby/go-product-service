package tracing

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"log"
	"runtime"
)

func newOtlpGrpcExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(ctx, otlptracegrpc.WithEndpoint("localhost:4317"), otlptracegrpc.WithInsecure())
}

// NewTraceProvider is function to
func NewTraceProvider(ctx context.Context, exp ...sdktrace.SpanExporter) func() {
	if len(exp) == 0 {
		exp = make([]sdktrace.SpanExporter, 1)
		exporter, err := newOtlpGrpcExporter(ctx)
		if err != nil {
			log.Fatal(err)
		}
		exp[0] = exporter
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-product-service"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp[0]),
		sdktrace.WithResource(r))

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // Propagasi trace-id, span-id, dll
		propagation.Baggage{},      // Propagasi metadata tambahan (key-value pairs)
	))

	return func() {
		_ = tp.Shutdown(ctx)
	}
}

// Start is function to create new span
func Start(ctx context.Context, opts ...string) (context.Context, trace.Span) {
	pc, _, _, _ := runtime.Caller(1)
	c, ok := ctx.(*gin.Context)
	if ok {
		ctx = c.Request.Context()
	}

	functionName := runtime.FuncForPC(pc).Name()
	if len(opts) > 0 {
		functionName = opts[0]
	}
	return otel.Tracer("go-product-service").Start(ctx, functionName)
}
