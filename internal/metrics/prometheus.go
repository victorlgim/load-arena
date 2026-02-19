package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func Init() {
	prometheus.MustRegister(httpRequests, httpLatency)
}

type statusCapturer struct {
	http.ResponseWriter
	status int
}

func (s *statusCapturer) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func InstrumentHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusCapturer{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		path := r.URL.Path
		method := r.Method
		status := strconv.Itoa(sw.status)

		httpRequests.WithLabelValues(method, path, status).Inc()
		httpLatency.WithLabelValues(method, path).Observe(time.Since(start).Seconds())
	})
}
