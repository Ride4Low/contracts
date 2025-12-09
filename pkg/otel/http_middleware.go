package otel

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHTTPHandler wraps an http.Handler with OpenTelemetry instrumentation.
func OTelOperationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		operation := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

		// Wrap the next handler with otelhttp using the generated operation name
		otelHandler := otelhttp.NewHandler(next, operation)

		otelHandler.ServeHTTP(w, r)
	})
}

// GinMiddleware returns a Gin middleware that adds OpenTelemetry tracing.
// It uses the provided service name for the tracer.
func GinMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}
