package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"go-api-base/internal/observability"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(body []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}

	n, err := rw.ResponseWriter.Write(body)
	rw.size += n
	return n, err
}

func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			writer := &responseWriter{
				ResponseWriter: w,
			}

			// Incrementa in-flight (ainda não sabemos a rota exata final, usamos genérica)
			observability.HTTPRequestsInFlight.WithLabelValues(r.Method, "unknown").Inc()
			defer observability.HTTPRequestsInFlight.WithLabelValues(r.Method, "unknown").Dec()

			next.ServeHTTP(writer, r)

			duration := time.Since(start)

			route := "unknown"
			routeCtx := chi.RouteContext(r.Context())
			if routeCtx != nil {
				route = routeCtx.RoutePattern()
			}
			if route == "" {
				route = "not_found"
			}

			observability.ObserveRequest(
				r.Method,
				route,
				writer.statusCode,
				duration,
				writer.size,
			)
		})
	}
}
