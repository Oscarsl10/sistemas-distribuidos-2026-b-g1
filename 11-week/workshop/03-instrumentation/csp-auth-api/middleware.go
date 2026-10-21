package instrumentation

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const headerCorrelationID = "X-Correlation-Id"

// TraceMiddleware reads the traceparent that the gateway forwards, or starts a trace when there is
// none, gives this service its own span, answers with the same headers and writes one access-log
// line per request.
func TraceMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		trace, ok := ParseTraceparent(r.Header.Get(headerTraceparent))
		if !ok {
			trace = NewTraceContext()
		}
		trace = trace.Child()
		trace.CorrelationID = r.Header.Get(headerCorrelationID)
		if trace.CorrelationID == "" {
			trace.CorrelationID = newUUID()
		}
		w.Header().Set(headerTraceparent, trace.Traceparent())
		w.Header().Set(headerCorrelationID, trace.CorrelationID)

		ctx := WithTrace(r.Context(), trace)
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(recorder, r.WithContext(ctx))

		logger.InfoContext(ctx, r.Method+" "+r.URL.Path,
			"event", "http.request.completed",
			"duration_ms", durationMillis(start),
			"http_status", recorder.status)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

// newUUID is a random (version 4) UUID, like the correlation id csp-auth-api generates.
func newUUID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
