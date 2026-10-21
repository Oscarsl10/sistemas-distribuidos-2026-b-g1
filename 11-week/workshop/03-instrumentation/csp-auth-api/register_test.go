package instrumentation

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	incomingTraceID = "9c1e7a42d5b84f03a6e1b27d38f0c915"
	incomingSpanID  = "1a8c4e92b6d70f35"
)

// register sends POST /api/v1/auth/register and returns what the next hop received.
func register(t *testing.T, traceparent string) (header http.Header, response *httptest.ResponseRecorder, logs string) {
	t.Helper()
	nextHop := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header = r.Header.Clone()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer nextHop.Close()

	var out bytes.Buffer
	server := NewServer(NewLogger(&out, "csp-auth-api", "3.0.0", "qa"), NewMetrics(), nextHop.Client(), nextHop.URL)
	req := httptest.NewRequest(http.MethodPost, registerRoute, nil)
	if traceparent != "" {
		req.Header.Set("traceparent", traceparent)
		req.Header.Set("X-Correlation-Id", "0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90")
	}
	response = httptest.NewRecorder()
	server.ServeHTTP(response, req)
	return header, response, out.String()
}

func TestTraceIdReceivedEqualsTraceIdSentDownstream(t *testing.T) {
	header, response, _ := register(t, "00-"+incomingTraceID+"-"+incomingSpanID+"-01")

	sent, ok := ParseTraceparent(header.Get("traceparent"))
	if !ok || sent.TraceID != incomingTraceID {
		t.Fatalf("next hop got trace_id %q, want %q", sent.TraceID, incomingTraceID)
	}
	if sent.SpanID == incomingSpanID {
		t.Fatal("the outbound call must have its own span id")
	}
	if got := header.Get("X-Correlation-Id"); got != "0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90" {
		t.Fatalf("correlation id lost: %q", got)
	}
	if answered, _ := ParseTraceparent(response.Header().Get("traceparent")); answered.TraceID != incomingTraceID {
		t.Fatal("the answer must carry the same trace_id")
	}
}

func TestStartsATraceWhenNoneArrives(t *testing.T) {
	header, response, _ := register(t, "")

	sent, ok := ParseTraceparent(header.Get("traceparent"))
	answered, _ := ParseTraceparent(response.Header().Get("traceparent"))
	if !ok || sent.TraceID != answered.TraceID {
		t.Fatalf("next hop %q and answer %q must share the new trace_id", sent.TraceID, answered.TraceID)
	}
	if response.Header().Get("X-Correlation-Id") == "" {
		t.Fatal("a correlation id must be created")
	}
}

func TestRejectsAMalformedTraceparent(t *testing.T) {
	for _, value := range []string{"garbage", "00-" + strings.Repeat("0", 32) + "-" + incomingSpanID + "-01", "00-ABC-" + incomingSpanID + "-01"} {
		if _, ok := ParseTraceparent(value); ok {
			t.Errorf("%q must not be accepted", value)
		}
	}
}

func TestLogLineFollowsTheContract(t *testing.T) {
	_, _, logs := register(t, "00-"+incomingTraceID+"-"+incomingSpanID+"-01")

	var line map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(logs)), &line); err != nil {
		t.Fatalf("not JSON: %v: %s", err, logs)
	}
	for _, field := range []string{"timestamp", "level", "service", "version", "environment", "trace_id", "span_id", "correlation_id", "event", "message", "duration_ms"} {
		if _, present := line[field]; !present {
			t.Errorf("field %q missing in %s", field, logs)
		}
	}
	if line["trace_id"] != incomingTraceID {
		t.Errorf("trace_id = %v", line["trace_id"])
	}
}

func TestMetricsCountTheRequestAndItsLatency(t *testing.T) {
	metrics := NewMetrics()
	metrics.Observe(registerRoute, http.StatusCreated, 0.12)
	metrics.Observe(registerRoute, http.StatusConflict, 0.03)
	recorder := httptest.NewRecorder()
	metrics.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/metrics", nil))

	body := recorder.Body.String()
	for _, want := range []string{
		`auth_register_requests_total{route="/api/v1/auth/register",outcome="success"} 1`,
		`auth_register_requests_total{route="/api/v1/auth/register",outcome="client_error"} 1`,
		`auth_register_request_duration_seconds_bucket{route="/api/v1/auth/register",le="0.05"} 1`,
		`auth_register_request_duration_seconds_count{route="/api/v1/auth/register"} 2`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("missing %q in:\n%s", want, body)
		}
	}
}
