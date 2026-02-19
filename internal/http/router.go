package http

import (
	"net/http"

	"github.com/victorlgim/load-arena/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/readyz", Readyz)

	mux.HandleFunc("/cpu", CPU)
	mux.HandleFunc("/io", IO)
	mux.HandleFunc("/mem", MEM)
	mux.HandleFunc("/chaos", CHAOS)

	mux.Handle("/metrics", promhttp.Handler())

	var h http.Handler = mux
	h = RequestID(h)
	h = Logging(h)
	h = metrics.InstrumentHTTP(h)

	return h
}
