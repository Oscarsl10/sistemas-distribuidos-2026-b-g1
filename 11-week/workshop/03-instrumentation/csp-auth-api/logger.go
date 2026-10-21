package instrumentation

import (
	"context"
	"io"
	"log/slog"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000Z"

// NewLogger writes one JSON object per line with the fields of 02-log-contract.md: timestamp,
// level, service, version, environment, message, and trace_id, span_id and correlation_id taken
// from the context of each call. The caller adds event, duration_ms and the domain fields.
func NewLogger(w io.Writer, service, version, environment string) *slog.Logger {
	base := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.TimeKey:
				return slog.String("timestamp", a.Value.Time().UTC().Format(timestampLayout))
			case slog.MessageKey:
				return slog.String("message", a.Value.String())
			}
			return a
		},
	})
	handler := traceHandler{next: base.WithAttrs([]slog.Attr{
		slog.String("service", service),
		slog.String("version", version),
		slog.String("environment", environment),
	})}
	return slog.New(handler)
}

// traceHandler adds the trace fields of the context to every record.
type traceHandler struct{ next slog.Handler }

func (h traceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h traceHandler) Handle(ctx context.Context, r slog.Record) error {
	if c, ok := TraceFrom(ctx); ok {
		r.AddAttrs(
			slog.String("trace_id", c.TraceID),
			slog.String("span_id", c.SpanID),
			slog.String("correlation_id", c.CorrelationID),
		)
	}
	return h.next.Handle(ctx, r)
}

func (h traceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return traceHandler{next: h.next.WithAttrs(attrs)}
}

func (h traceHandler) WithGroup(name string) slog.Handler {
	return traceHandler{next: h.next.WithGroup(name)}
}

// durationMillis is the duration_ms field value.
func durationMillis(start time.Time) int64 { return time.Since(start).Milliseconds() }
