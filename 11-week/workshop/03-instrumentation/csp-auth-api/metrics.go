package instrumentation

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
)

// Bucket limits in seconds. A registration hashes the password with bcrypt, so the first buckets
// are wider than they would be for a plain read.
var latencyBuckets = []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5}

type series struct{ route, outcome string }

type histogram struct {
	buckets []uint64
	count   uint64
	sum     float64
}

// Metrics keeps the two metrics of the registration in the Prometheus text format:
// auth_register_requests_total{route,outcome} and auth_register_request_duration_seconds{route}.
type Metrics struct {
	mu       sync.Mutex
	requests map[series]uint64
	latency  map[string]*histogram
}

func NewMetrics() *Metrics {
	return &Metrics{requests: map[series]uint64{}, latency: map[string]*histogram{}}
}

// Outcome groups an HTTP status in the label values success, client_error and error.
func Outcome(status int) string {
	switch {
	case status >= 500:
		return "error"
	case status >= 400:
		return "client_error"
	}
	return "success"
}

// Observe records one request of the route. The route is the pattern, never the raw path.
func (m *Metrics) Observe(route string, status int, seconds float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requests[series{route, Outcome(status)}]++
	h := m.latency[route]
	if h == nil {
		h = &histogram{buckets: make([]uint64, len(latencyBuckets))}
		m.latency[route] = h
	}
	for i, limit := range latencyBuckets {
		if seconds <= limit {
			h.buckets[i]++
		}
	}
	h.count++
	h.sum += seconds
}

// ServeHTTP answers GET /metrics.
func (m *Metrics) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var b strings.Builder
	b.WriteString("# TYPE auth_register_requests_total counter\n")
	keys := make([]series, 0, len(m.requests))
	for k := range m.requests {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].route+keys[i].outcome < keys[j].route+keys[j].outcome })
	for _, k := range keys {
		fmt.Fprintf(&b, "auth_register_requests_total{route=%q,outcome=%q} %d\n", k.route, k.outcome, m.requests[k])
	}
	b.WriteString("# TYPE auth_register_request_duration_seconds histogram\n")
	routes := make([]string, 0, len(m.latency))
	for r := range m.latency {
		routes = append(routes, r)
	}
	sort.Strings(routes)
	for _, r := range routes {
		h := m.latency[r]
		for i, limit := range latencyBuckets {
			fmt.Fprintf(&b, "auth_register_request_duration_seconds_bucket{route=%q,le=\"%g\"} %d\n", r, limit, h.buckets[i])
		}
		fmt.Fprintf(&b, "auth_register_request_duration_seconds_bucket{route=%q,le=\"+Inf\"} %d\n", r, h.count)
		fmt.Fprintf(&b, "auth_register_request_duration_seconds_sum{route=%q} %g\n", r, h.sum)
		fmt.Fprintf(&b, "auth_register_request_duration_seconds_count{route=%q} %d\n", r, h.count)
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	_, _ = w.Write([]byte(b.String()))
}
