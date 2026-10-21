package instrumentation

import (
	"context"
	"net/http"
)

// CallNextHop calls the next hop of the journey and passes the trace on: the same trace_id and
// correlation id as the request being served, and a new span id for the outbound call. In the real
// platform the next hop is the message that csp-worker publishes to cine.events, which carries the
// same two values as AMQP headers; here it is an HTTP call so the propagation can be tested.
func CallNextHop(ctx context.Context, client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	if trace, ok := TraceFrom(ctx); ok {
		out := trace.Child()
		req.Header.Set(headerTraceparent, out.Traceparent())
		req.Header.Set(headerCorrelationID, out.CorrelationID)
	}
	return client.Do(req)
}
