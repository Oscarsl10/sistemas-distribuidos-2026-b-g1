# How to run the instrumentation of `csp-auth-api` (CP3)

A minimal Go module that models the register journey of `csp-auth-api`: a middleware that reads or creates the trace
context, a JSON logger that follows [02-log-contract.md](../02-log-contract.md), propagation of the trace to the next hop
and the two metrics of the registration on `/metrics`. It uses only the Go standard library.
It is a model for the workshop: it does not change the real project.

Tested on Windows 11 with Go 1.27.0 (the module asks for Go 1.24 or newer).

## Install and test

```bash
cd 11-week/workshop/03-instrumentation/csp-auth-api
go vet ./...
go test -v ./...
```

Real output of my machine:

```
=== RUN   TestTraceIdReceivedEqualsTraceIdSentDownstream
--- PASS: TestTraceIdReceivedEqualsTraceIdSentDownstream (0.01s)
=== RUN   TestStartsATraceWhenNoneArrives
--- PASS: TestStartsATraceWhenNoneArrives (0.00s)
=== RUN   TestRejectsAMalformedTraceparent
--- PASS: TestRejectsAMalformedTraceparent (0.00s)
=== RUN   TestLogLineFollowsTheContract
--- PASS: TestLogLineFollowsTheContract (0.00s)
=== RUN   TestMetricsCountTheRequestAndItsLatency
--- PASS: TestMetricsCountTheRequestAndItsLatency (0.00s)
PASS
ok  	csp-auth-api-instrumentation	2.242s
```

`TestTraceIdReceivedEqualsTraceIdSentDownstream` is the proof asked by the workshop: the `trace_id` received by the
service is the same one the next hop receives, with a new `span_id`.

## Run and call it

```bash
go run ./cmd/demo
```

It starts the model on `127.0.0.1:8081` and a stub of the next hop on `127.0.0.1:8082`. In another terminal:

```bash
curl -i -X POST localhost:8081/api/v1/auth/register \
  -H "traceparent: 00-9c1e7a42d5b84f03a6e1b27d38f0c915-1a8c4e92b6d70f35-01" \
  -H "X-Correlation-Id: 0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90"
curl localhost:8081/metrics
```

Answer of the first call (same `trace_id`, new `span_id`, same correlation id):

```
HTTP/1.1 201 Created
Traceparent: 00-9c1e7a42d5b84f03a6e1b27d38f0c915-2b83d67cce6fe6c3-01
X-Correlation-Id: 0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90
```

Logs written by the demo (stdout). The first line is the stub of the next hop: it received the same `trace_id` with a
different `span_id`.

```
{"timestamp":"2026-10-21T00:34:27.267Z","level":"INFO","message":"next hop received the message","service":"csp-auth-api","version":"3.0.0","environment":"dev","event":"nexthop.received","traceparent":"00-9c1e7a42d5b84f03a6e1b27d38f0c915-e52167a8b4e2c92a-01"}
{"timestamp":"2026-10-21T00:34:27.267Z","level":"INFO","message":"POST /api/v1/auth/register","service":"csp-auth-api","version":"3.0.0","environment":"dev","event":"http.request.completed","duration_ms":1,"http_status":201,"trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"2b83d67cce6fe6c3","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90"}
```

Answer of `/metrics` after that one request:

```
# TYPE auth_register_requests_total counter
auth_register_requests_total{route="/api/v1/auth/register",outcome="success"} 1
# TYPE auth_register_request_duration_seconds histogram
auth_register_request_duration_seconds_bucket{route="/api/v1/auth/register",le="0.05"} 1
auth_register_request_duration_seconds_bucket{route="/api/v1/auth/register",le="+Inf"} 1
auth_register_request_duration_seconds_sum{route="/api/v1/auth/register"} 0.0018321
auth_register_request_duration_seconds_count{route="/api/v1/auth/register"} 1
```

The output above keeps only the buckets `0.05` and `+Inf`; the real answer also lists `0.1`, `0.25`, `0.5`, `1` and `2.5`.

## Files

| File | What it does |
|---|---|
| `tracecontext.go` | Reads, creates and formats `traceparent`; keeps it in the `context`. |
| `logger.go` | `slog` JSON logger with the fields of the log contract; adds `trace_id`, `span_id` and `correlation_id` from the context. |
| `middleware.go` | Reads or creates the trace, answers with the headers and writes one access-log line per request. |
| `client.go` | Calls the next hop with the same `trace_id` and a new `span_id`. |
| `metrics.go` | Counter by route and outcome and latency histogram, in the Prometheus text format. |
| `register.go`, `cmd/demo` | The model of `POST /api/v1/auth/register` and a runnable demo. |
