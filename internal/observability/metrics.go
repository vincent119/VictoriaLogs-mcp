package observability

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all application metrics
type Metrics struct {
	// Tool call metrics
	ToolCallsTotal    *prometheus.CounterVec
	ToolCallDuration  *prometheus.HistogramVec
	ToolCallErrors    *prometheus.CounterVec

	// VictoriaLogs client metrics
	VLQueryTotal      *prometheus.CounterVec
	VLQueryDuration   *prometheus.HistogramVec
	VLQueryErrors     *prometheus.CounterVec

	// Policy metrics
	RateLimitHits     prometheus.Counter
	CircuitBreakerTrips prometheus.Counter
	AllowlistBlocks   prometheus.Counter
}

var (
	defaultMetrics *Metrics
)

// InitMetrics initializes metrics
func InitMetrics(namespace string) *Metrics {
	if namespace == "" {
		namespace = "vlmcp"
	}

	m := &Metrics{
		ToolCallsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tool_calls_total",
				Help:      "Total number of MCP tool calls",
			},
			[]string{"tool"},
		),
		ToolCallDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "tool_call_duration_seconds",
				Help:      "Duration of MCP tool calls in seconds",
				Buckets:   []float64{0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
			},
			[]string{"tool"},
		),
		ToolCallErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "tool_call_errors_total",
				Help:      "Total number of MCP tool call errors",
			},
			[]string{"tool", "error_type"},
		),
		VLQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "victorialogs_queries_total",
				Help:      "Total number of VictoriaLogs queries",
			},
			[]string{"endpoint"},
		),
		VLQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "victorialogs_query_duration_seconds",
				Help:      "Duration of VictoriaLogs queries in seconds",
				Buckets:   []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
			},
			[]string{"endpoint"},
		),
		VLQueryErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "victorialogs_query_errors_total",
				Help:      "Total number of VictoriaLogs query errors",
			},
			[]string{"endpoint", "status_code"},
		),
		RateLimitHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "rate_limit_hits_total",
				Help:      "Total number of rate limit hits",
			},
		),
		CircuitBreakerTrips: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "circuit_breaker_trips_total",
				Help:      "Total number of circuit breaker trips",
			},
		),
		AllowlistBlocks: promauto.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "allowlist_blocks_total",
				Help:      "Total number of allowlist blocks",
			},
		),
	}

	defaultMetrics = m
	return m
}

// GetMetrics returns the default metrics instance
func GetMetrics() *Metrics {
	return defaultMetrics
}

// RecordToolCall records a tool call
func (m *Metrics) RecordToolCall(tool string, duration time.Duration, err error) {
	m.ToolCallsTotal.WithLabelValues(tool).Inc()
	m.ToolCallDuration.WithLabelValues(tool).Observe(duration.Seconds())

	if err != nil {
		m.ToolCallErrors.WithLabelValues(tool, "error").Inc()
	}
}

// RecordVLQuery records a VictoriaLogs query
func (m *Metrics) RecordVLQuery(endpoint string, duration time.Duration, statusCode int) {
	m.VLQueryTotal.WithLabelValues(endpoint).Inc()
	m.VLQueryDuration.WithLabelValues(endpoint).Observe(duration.Seconds())

	if statusCode >= 400 {
		m.VLQueryErrors.WithLabelValues(endpoint, http.StatusText(statusCode)).Inc()
	}
}

// RecordRateLimitHit records a rate limit hit
func (m *Metrics) RecordRateLimitHit() {
	m.RateLimitHits.Inc()
}

// RecordCircuitBreakerTrip records a circuit breaker trip
func (m *Metrics) RecordCircuitBreakerTrip() {
	m.CircuitBreakerTrips.Inc()
}

// RecordAllowlistBlock records an allowlist block
func (m *Metrics) RecordAllowlistBlock() {
	m.AllowlistBlocks.Inc()
}

// Handler returns Prometheus HTTP handler
func Handler() http.Handler {
	return promhttp.Handler()
}
