// Package metrics defines the Prometheus collectors shared by the api and
// worker binaries and exposes them over HTTP.
package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry is the process registry; binaries register it once.
var Registry = prometheus.NewRegistry()

var (
	// HTTPRequests counts API requests by route, method and status class.
	HTTPRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "habarchy", Subsystem: "http", Name: "requests_total", Help: "HTTP requests",
	}, []string{"route", "method", "status"})
	// HTTPDuration measures request latency.
	HTTPDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "habarchy", Subsystem: "http", Name: "request_duration_seconds", Help: "HTTP request latency",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
	}, []string{"route", "method"})
	// MessagesTotal counts message outcomes by channel and status.
	MessagesTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "habarchy", Name: "messages_total", Help: "Message status transitions",
	}, []string{"channel", "status"})
	// ProviderCalls counts provider attempts by type and outcome.
	ProviderCalls = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "habarchy", Name: "provider_calls_total", Help: "Provider send attempts",
	}, []string{"type", "outcome"})
	// ProviderDuration measures provider latency.
	ProviderDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "habarchy", Name: "provider_call_duration_seconds", Help: "Provider call latency",
		Buckets: []float64{.05, .1, .25, .5, 1, 2, 5, 10, 30},
	}, []string{"type"})
	// WebhookDeliveries counts webhook attempts by outcome.
	WebhookDeliveries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "habarchy", Name: "webhook_deliveries_total", Help: "Webhook delivery attempts",
	}, []string{"outcome"})
)

func init() {
	Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		HTTPRequests, HTTPDuration, MessagesTotal, ProviderCalls, ProviderDuration, WebhookDeliveries,
	)
}

// Handler serves the registry in Prometheus text format.
func Handler() http.Handler {
	return promhttp.HandlerFor(Registry, promhttp.HandlerOpts{})
}

// FiberMiddleware records request counts and latency per route pattern
// (not per raw path, so ids do not explode cardinality).
func FiberMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		// fasthttp reuses request buffers; copy before/after Next so the
		// label strings are not overwritten by the next request.
		method := utils.CopyString(c.Method())
		err := c.Next()
		route := utils.CopyString(c.Route().Path)
		if route == "" || route == "/" && c.Path() != "/" {
			route = "unmatched"
		}
		status := c.Response().StatusCode()
		if err != nil {
			if fe, ok := err.(*fiber.Error); ok { //nolint:errorlint // fiber errors are concrete
				status = fe.Code
			} else if status < 400 {
				status = 500
			}
		}
		HTTPRequests.WithLabelValues(route, method, strconv.Itoa(status/100)+"xx").Inc()
		HTTPDuration.WithLabelValues(route, method).Observe(time.Since(start).Seconds())
		return err
	}
}

// Serve starts a tiny HTTP server exposing /metrics (worker, scheduler).
func Serve(addr string) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() { _ = srv.ListenAndServe() }()
	return srv
}
