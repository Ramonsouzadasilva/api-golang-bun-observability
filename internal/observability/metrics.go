package observability

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total de requisições HTTP.",
		},
		[]string{"method", "route", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duração das requisições HTTP.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	HTTPRequestsInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Número de requisições HTTP atualmente em processamento.",
		},
		[]string{"method", "route"},
	)

	HTTPResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Tamanho das respostas HTTP em bytes.",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B a ~1GB
		},
		[]string{"method", "route"},
	)

	BcryptDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "bcrypt_duration_seconds",
			Help:    "Duração do hash bcrypt.",
			Buckets: prometheus.DefBuckets,
		},
	)

	AuthTokensIssued = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_tokens_issued_total",
			Help: "Total de tokens de autenticação emitidos.",
		},
		[]string{"type"}, // access ou refresh
	)
)

func Register() {
	// Limpa o registro padrão para garantir que não haja duplicatas
	prometheus.DefaultRegisterer.Unregister(collectors.NewGoCollector())
	prometheus.DefaultRegisterer.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Registra coletores do Go runtime
	prometheus.MustRegister(collectors.NewGoCollector())
	prometheus.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	// Registra métricas da aplicação
	prometheus.MustRegister(
		HTTPRequestsTotal,
		HTTPRequestDuration,
		HTTPRequestsInFlight,
		HTTPResponseSize,
		BcryptDuration,
		AuthTokensIssued,
	)
}

func ObserveRequest(
	method string,
	route string,
	status int,
	duration time.Duration,
	size int,
) {
	HTTPRequestsTotal.WithLabelValues(
		method,
		route,
		strconv.Itoa(status),
	).Inc()

	HTTPRequestDuration.WithLabelValues(
		method,
		route,
	).Observe(duration.Seconds())

	HTTPResponseSize.WithLabelValues(
		method,
		route,
	).Observe(float64(size))
}
