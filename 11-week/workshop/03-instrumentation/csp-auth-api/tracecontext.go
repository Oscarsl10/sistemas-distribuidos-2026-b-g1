// Package instrumentation is a minimal model of the observability of csp-auth-api for the
// journey "client registers an account" (POST /api/v1/auth/register).
package instrumentation

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

const headerTraceparent = "traceparent"

// TraceContext is the W3C trace context of one step of a request.
type TraceContext struct {
	TraceID       string
	SpanID        string
	CorrelationID string
}

type contextKey struct{}

// ParseTraceparent reads "00-<trace_id>-<span_id>-<flags>". It refuses a malformed value
// or an all-zero id, so the caller starts a new trace instead of trusting it.
func ParseTraceparent(value string) (TraceContext, bool) {
	parts := strings.Split(strings.TrimSpace(value), "-")
	if len(parts) != 4 || parts[0] != "00" || !isHex(parts[1], 32) || !isHex(parts[2], 16) {
		return TraceContext{}, false
	}
	return TraceContext{TraceID: parts[1], SpanID: parts[2]}, true
}

// NewTraceContext starts a trace with random ids.
func NewTraceContext() TraceContext {
	return TraceContext{TraceID: randomHex(16), SpanID: randomHex(8)}
}

// Child keeps the trace and the correlation id and gives a new span id.
func (c TraceContext) Child() TraceContext {
	c.SpanID = randomHex(8)
	return c
}

// Traceparent is the value of the header traceparent for this span.
func (c TraceContext) Traceparent() string {
	return fmt.Sprintf("00-%s-%s-01", c.TraceID, c.SpanID)
}

// WithTrace stores the trace context in the context.
func WithTrace(ctx context.Context, c TraceContext) context.Context {
	return context.WithValue(ctx, contextKey{}, c)
}

// TraceFrom returns the trace context stored by WithTrace.
func TraceFrom(ctx context.Context) (TraceContext, bool) {
	c, ok := ctx.Value(contextKey{}).(TraceContext)
	return c, ok
}

func isHex(value string, length int) bool {
	if len(value) != length || strings.Trim(value, "0") == "" {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil && value == strings.ToLower(value)
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
