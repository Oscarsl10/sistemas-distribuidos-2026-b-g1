# CP4 - Golden signals and metric catalog

Journey: the registration of [01-journey.md](01-journey.md). One table per service of the journey, with the four golden
signals (latency, traffic, errors, saturation). Real elements used here: `csp-api-gateway`, `csp-auth-api`
(`POST /api/v1/auth/register`), the job `auth-outbox-relay` of `csp-worker` (table `auth.outbox_event`), the queue
`cine.ticketing_fulfillment.user-registered` of `csp-ticketing-api`.

State of each metric today: **built** means it is exported by the model of [CP3](03-instrumentation/HOW-TO-RUN.md);
**to build** means it is only specified here. The platform already runs the OpenTelemetry collector, Prometheus and
Grafana (`csp-infra-postgres/observability`), so a service only has to export the metrics.

Naming rules: `snake_case`, the unit as suffix (`_seconds`), `_total` for counters, the name starts with the service prefix.

## 1. `csp-api-gateway` (NGINX)

NGINX has no Prometheus metrics by itself; they would come from its JSON access log or from the exporter of NGINX.

| Signal | Metric | Type | Labels | PromQL for the dashboard | State |
|---|---|---|---|---|---|
| Latency | `gateway_request_duration_seconds` | histogram | `route`, `method` | `histogram_quantile(0.95, sum by (le) (rate(gateway_request_duration_seconds_bucket{route="/api/v1/auth/register"}[5m])))` | to build |
| Traffic | `gateway_requests_total` | counter | `route`, `method`, `status_class` | `sum by (route) (rate(gateway_requests_total[5m]))` | to build |
| Errors | `gateway_requests_total` | counter | `status_class="5xx"` | `sum(rate(gateway_requests_total{status_class="5xx"}[5m])) / sum(rate(gateway_requests_total[5m]))` | to build |
| Saturation | `gateway_rate_limited_total` | counter | `route` | `sum by (route) (rate(gateway_rate_limited_total[5m]))` (answers 429 when a client passes 10 req/s plus the burst of 20) | to build |

## 2. `csp-auth-api`

| Signal | Metric | Type | Labels | PromQL for the dashboard | State |
|---|---|---|---|---|---|
| Latency | `auth_register_request_duration_seconds` | histogram | `route` | `histogram_quantile(0.95, sum by (le) (rate(auth_register_request_duration_seconds_bucket[5m])))` | built |
| Traffic | `auth_register_requests_total` | counter | `route`, `outcome` | `sum(rate(auth_register_requests_total[5m]))` | built |
| Errors | `auth_register_requests_total` | counter | `outcome="error"` | `sum(rate(auth_register_requests_total{outcome="error"}[5m])) / sum(rate(auth_register_requests_total[5m]))` | built |
| Saturation | `auth_db_pool_connections_in_use` | gauge | none | `max(auth_db_pool_connections_in_use) / max(auth_db_pool_connections_max)` (the pool size is a setting of the service) | to build |

`outcome` has three values: `success` (2xx), `client_error` (4xx, for example `EMAIL_ALREADY_REGISTERED`) and `error` (5xx).
A `4xx` is a correct answer, so it does not count as an error for the SLO.

## 3. `csp-worker`, job `auth-outbox-relay`

| Signal | Metric | Type | Labels | PromQL for the dashboard | State |
|---|---|---|---|---|---|
| Latency | `outbox_relay_event_age_seconds` | gauge | `event_type` | `max by (event_type) (outbox_relay_event_age_seconds)` (age of the oldest event not published: how fresh `UserRegistered` is) | to build |
| Traffic | `outbox_relay_published_total` | counter | `event_type` | `sum by (event_type) (rate(outbox_relay_published_total[5m]))` | to build |
| Errors | `outbox_relay_failed_total` | counter | `event_type`, `reason` | `sum(rate(outbox_relay_failed_total[5m])) / (sum(rate(outbox_relay_failed_total[5m])) + sum(rate(outbox_relay_published_total[5m])))` | to build |
| Saturation | `outbox_relay_pending_events` | gauge | `event_type` | `sum(outbox_relay_pending_events)` (rows of `auth.outbox_event` waiting; the relay runs every 5 s, one instance) | to build |

## 4. `csp-ticketing-api`, consumer of `UserRegistered`

The service has only its README in `develop`, so everything here is specified for when the consumer exists.

| Signal | Metric | Type | Labels | PromQL for the dashboard | State |
|---|---|---|---|---|---|
| Latency | `ticketing_consumer_processing_duration_seconds` | histogram | `queue` | `histogram_quantile(0.95, sum by (le) (rate(ticketing_consumer_processing_duration_seconds_bucket[5m])))` | to build |
| Traffic | `ticketing_consumer_messages_total` | counter | `queue`, `outcome` | `sum by (queue) (rate(ticketing_consumer_messages_total[5m]))` | to build |
| Errors | `ticketing_consumer_messages_total` | counter | `outcome` in `requeued`, `dead_lettered` | `sum(rate(ticketing_consumer_messages_total{outcome=~"requeued\|dead_lettered"}[5m]))` | to build |
| Saturation | `ticketing_queue_messages_ready` | gauge | `queue` | `max by (queue) (ticketing_queue_messages_ready)` (limit of the queue: 10000 messages, ADR-036) | to build |

`outcome` of the consumer: `acked`, `duplicate` (rejected by `event_id`), `requeued`, `dead_lettered`.

## 5. Cardinality of the labels

| Label | Values | Series it adds |
|---|---|---|
| `route` | route pattern, about 8 routes of Auth in the gateway, 1 in the model | 8 |
| `method` | `GET`, `POST`, `PATCH` | 3 |
| `status_class` | `2xx`, `3xx`, `4xx`, `5xx` | 4 |
| `outcome` | 3 values in Auth, 4 in the consumer | 3 to 4 |
| `event_type` | `UserRegistered`, `UserEmailChanged`, `UserRoleChanged` | 3 |
| `queue` | one queue per consumer and event type, `cine.<service>.<event-type>` | 2 now |

Worst case for the register histogram: 1 route x (7 buckets + sum + count) = 9 series. The whole catalog stays under 300 series.

**Forbidden labels** (unbounded, they would create one series per user or per request): `user_id`, `email`, the raw URL
(`/users/7a3d.../role`), `event_id`, `trace_id`, the `Idempotency-Key`. The raw URL is replaced by the route pattern
(`/users/{userId}/role`); a single user or request is found in the logs by `trace_id`, not in the metrics.

## 6. Histogram buckets

Buckets of the register latency, in seconds: `0.05, 0.1, 0.25, 0.5, 1, 2.5`.

- A registration hashes the password with bcrypt on purpose, so a normal answer is slower than a read. I expect most of them
  below 0.5 s; this is an assumption to be checked with the first real measurement.
- The user waits for the answer on a form. Under 1 s feels immediate, and from 2.5 s I expect people to repeat the click
  (the `Idempotency-Key` protects against the duplicate, but it is still a bad experience). So `1` and `2.5` are the
  limits where the SLO of CP6 is decided.
- The small buckets (`0.05`, `0.1`) are kept so that a fast answer, for example a request refused by validation, is not
  counted in the same bucket as a real registration.
