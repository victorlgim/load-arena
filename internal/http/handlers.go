package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/victorlgim/load-arena/internal/workload"
)

type resp struct {
	OK       bool              `json:"ok"`
	Now      string            `json:"now"`
	Endpoint string            `json:"endpoint"`
	Params   map[string]any    `json:"params,omitempty"`
	Result   map[string]any    `json:"result,omitempty"`
	Error    string            `json:"error,omitempty"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/healthz",
	})
}

func Readyz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/readyz",
	})
}

func CPU(w http.ResponseWriter, r *http.Request) {
	n := intParam(r, "n", 30000, 1, 500000)
	algo := strParam(r, "algo", "sha")

	out := workload.CPU(n, algo)

	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/cpu",
		Params: map[string]any{
			"n":    n,
			"algo": algo,
		},
		Result: map[string]any{
			"hash": out,
		},
	})
}

func IO(w http.ResponseWriter, r *http.Request) {
	delay := intParam(r, "delay", 200, 0, 5000) // ms

	out := workload.IO(time.Duration(delay) * time.Millisecond)

	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/io",
		Params: map[string]any{
			"delay_ms": delay,
		},
		Result: map[string]any{
			"slept_ms": out.Milliseconds(),
		},
	})
}

func MEM(w http.ResponseWriter, r *http.Request) {
	mb := intParam(r, "mb", 50, 1, 1024)
	hold := intParam(r, "hold", 200, 0, 10000) 

	bytes := workload.MemAlloc(mb, time.Duration(hold)*time.Millisecond)

	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/mem",
		Params: map[string]any{
			"mb":       mb,
			"hold_ms":  hold,
			"requested": mb,
		},
		Result: map[string]any{
			"allocated_bytes": bytes,
		},
	})
}

func CHAOS(w http.ResponseWriter, r *http.Request) {
	rate := floatParam(r, "rate", 0.20, 0, 1) 
	mode := strParam(r, "mode", "http500")    

	out := workload.Chaos(rate, mode)

	if out.Failed {
		writeJSON(w, out.Status, resp{
			OK:       false,
			Now:      time.Now().Format(time.RFC3339Nano),
			Endpoint: "/chaos",
			Params: map[string]any{
				"rate": rate,
				"mode": mode,
			},
			Error: out.Message,
		})
		return
	}

	writeJSON(w, http.StatusOK, resp{
		OK:       true,
		Now:      time.Now().Format(time.RFC3339Nano),
		Endpoint: "/chaos",
		Params: map[string]any{
			"rate": rate,
			"mode": mode,
		},
		Result: map[string]any{
			"message": out.Message,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func intParam(r *http.Request, key string, def, min, max int) int {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

func floatParam(r *http.Request, key string, def, min, max float64) float64 {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def
	}
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return def
	}
	if f < min {
		return min
	}
	if f > max {
		return max
	}
	return f
}

func strParam(r *http.Request, key, def string) string {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def
	}
	return raw
}
