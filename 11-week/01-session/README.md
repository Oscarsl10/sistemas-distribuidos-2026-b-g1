# Week 11 - Session 1: Observability (logs, metrics and distributed tracing)

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.
Scope: the Auth service (`csp-auth-api`) behind the API gateway (`csp-api-gateway`), with the platform
composition in `csp-infra-postgres`. Everything marked **exists** was read in the repositories (`develop`);
everything marked **to build** is a proposal that is not implemented yet.

## 1. The three pillars in the platform today

| Pillar | State | Evidence |
|---|---|---|
| Logs | **exists**, partial | `csp-auth-api` writes JSON logs with `log/slog` (`cmd/auth-api/main.go`). The gateway writes one JSON access log line per request with `correlationId` (`csp-api-gateway/nginx/nginx.conf`). |
| Correlation / trace id | **exists** | The gateway creates or reuses `X-Correlation-Id`, forwards it to the upstream (`nginx/conf.d/20-server.conf`) and returns it in the response (`nginx/snippets/headers.conf`). `csp-auth-api` reuses the id when it is printable ASCII of at most 128 characters, otherwise it generates a UUID v4 (`internal/adapter/in/httpapi/correlation.go`, covered by `correlation_test.go`). |
| Metrics | **platform only** | The platform composition runs an OpenTelemetry collector (OTLP on 4317/4318) that exposes metrics to Prometheus (`:8889`), and Grafana with Prometheus as datasource (`csp-infra-postgres/observability/`). No service exports metrics yet. |
| Traces | **platform only** | The collector accepts OTLP traces but only prints them with the `debug` exporter. No service creates spans yet. |

## 2. What a request leaves behind today

`POST /api/v1/auth/login` through the gateway on port 8000:

1. The gateway assigns the correlation id, writes its access log line with `correlationId`, and forwards the
   request with `X-Correlation-Id`.
2. `csp-auth-api` reads the same id (`traceID := correlationID(r)`) and returns it as `traceId` in every error body
   (the shared `ErrorResponse` of `_shared.yaml`).
3. When something fails, the service logs it with the same value: `slog.Error("login user failed", "traceId", traceID, "error", err)`
   (`internal/adapter/in/httpapi/errors.go`). The same pattern exists for register, refresh, role change, user list, logout and purge.

So one failed request can already be followed by hand: take the `traceId` from the error body (or the
`X-Correlation-Id` response header), find the gateway line with that `correlationId`, then the service line with that `traceId`.

## 3. Gaps found

| # | Gap | Effect |
|---|---|---|
| G1 | The service logs only errors and startup messages; there is no access log with status and duration per request. | A slow or successful request leaves no trace in the service. |
| G2 | The id is called `correlationId` in the gateway and `traceId` in the service. | Searching both logs with one query needs both names. |
| G3 | The correlation id is not in the context of the use cases, only in the HTTP adapter. | Logs written below the adapter (persistence, outbox, Redis) cannot carry it. |
| G4 | No service exports request count, error count or latency. | The golden signals cannot be drawn in Grafana, even though the stack is running. |
| G5 | No spans: nothing sends OTLP to the collector. | There is no timed view of a request across gateway, auth, PostgreSQL and Redis. |

## 4. Instrumentation to build (checklist of the activity)

1. **Structured logs with a propagated `trace_id`**: one access-log line per request in `csp-auth-api`
   (`method`, `route`, `status`, `duration_ms`, `trace_id`), using the same field name as the gateway; carry the id in the
   `context.Context` so the use cases and adapters log it too (closes G1, G2, G3). Never log passwords, tokens or the
   `Idempotency-Key` value.
2. **Basic metrics (rate, errors, latency)**: an HTTP middleware that records `http_server_requests_total{route,method,status}`
   and `http_server_request_duration_seconds{route}` (histogram), exported by OTLP to the collector that already exists (closes G4).
   The route label is the route pattern (`/users/{userId}/role`), never the raw path, to keep the cardinality low (norm 5.3.10).
3. **One distributed trace**: OpenTelemetry HTTP server instrumentation in `csp-auth-api` and `traceparent` forwarding in the gateway,
   so a login shows the spans `gateway` -> `auth POST /login` -> `PostgreSQL` -> `bcrypt` (closes G5). The collector pipeline
   would export traces to a store instead of `debug`.

## 5. Trace of the login (expected shape once built)

```
POST /api/v1/auth/login                      total
 |- gateway (nginx): rate limit + route      ~1 ms
 `- csp-auth-api POST /login
     |- find user by email (PostgreSQL)
     |- bcrypt compare                       the expected bottleneck: bcrypt is deliberately slow
     `- store refresh token (PostgreSQL)
```

The bcrypt comparison is the span expected to dominate, which is why the latency SLO of session 2 is defined on the
whole request and not on the database calls. The durations are not measured yet; this diagram is the hypothesis the
trace must confirm.

## 6. Self-check (answers)

- **Monitoring vs observability:** monitoring watches questions known in advance (is the health route up?); observability lets me ask a new
  question about a request (why was this login slow for this user?) from the data the service already emits.
- **Why a correlation id:** it is the only key that joins the gateway line, the service line and the error the client sees.
- **Alert on symptoms:** error rate and p99 latency of login, not CPU of the container (see session 2).
- **Deliberately not done here:** the SLOs and alerts belong to session 2.
