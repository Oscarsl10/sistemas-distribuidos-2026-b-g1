package instrumentation

import (
	"log/slog"
	"net/http"
	"time"
)

const registerRoute = "/api/v1/auth/register"

// NewServer wires the model of csp-auth-api: POST /api/v1/auth/register, which calls the next hop
// of the journey, and GET /metrics, all behind the trace middleware.
func NewServer(logger *slog.Logger, metrics *Metrics, client *http.Client, nextHopURL string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST "+registerRoute, func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		status := http.StatusCreated
		resp, err := CallNextHop(r.Context(), client, nextHopURL)
		if err != nil {
			logger.ErrorContext(r.Context(), "next hop did not answer",
				"event", "nexthop.call.failed", "duration_ms", durationMillis(start), "error", err.Error())
			status = http.StatusBadGateway
		} else {
			resp.Body.Close()
		}
		metrics.Observe(registerRoute, status, time.Since(start).Seconds())
		w.WriteHeader(status)
	})
	mux.Handle("GET /metrics", metrics)
	return TraceMiddleware(logger, mux)
}
