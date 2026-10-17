# Week 11 - Session 2: SLIs, SLOs, error budget and alerting

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.
Service: Auth (`csp-auth-api`) behind the API gateway. This is a plan: the metrics it needs are the ones listed
in [session 1](../01-session/README.md), section 4, and are **not exported yet**. Until they are built, the SLIs
cannot be measured and the numbers below are targets, not results.

## 1. SLOs (3, with window)

| SLO | SLI (what is measured) | Target | Window |
|---|---|---|---|
| **Availability of sign-in** | Share of `POST /login` and `POST /register` requests that do not answer `5xx` (a `4xx` such as wrong password or duplicate email is a correct answer) | 99.5% | 30 days rolling |
| **Latency of sign-in** | Share of `POST /login` requests answered in less than 800 ms, measured at the gateway | 95% (p95 < 800 ms) | 30 days rolling |
| **Freshness of `UserRegistered`** | Time from the commit of a registration to the event leaving the outbox | 99% within 60 s | 30 days rolling |

Notes on the choices:
- Login latency is higher than a generic API target because the password is compared with bcrypt on purpose; 800 ms is
  a first guess to be corrected with the first real p95.
- The freshness SLO covers the asynchronous flow of the story (outbox, norm 5.3.11). It needs the outbox relay to expose the age of the
  oldest pending row; the relay is a decision of the platform ADRs, so this SLO stays a proposal until that metric exists.
- No vanity metric (CPU, memory, number of pods) is an SLO.

## 2. Error budget and policy

Budget of the availability SLO: 100% - 99.5% = **0.5%** of the requests in 30 days.
With 100 000 sign-in requests a month that is 500 failed requests; this is the figure to recompute with the real traffic.

| Budget spent in the window | Policy |
|---|---|
| < 50% | Normal work. Releases of the Auth service go on as planned. |
| 50% - 100% | Every change to login, register or refresh needs a second review and a rollback note in its PR. |
| >= 100% | Feature work on the Auth service stops. Only reliability fixes (the cause of the failures) and rollbacks are merged until the budget recovers. |

The owner of the Auth service decides when the freeze starts and ends and records it in the weekly status.

## 3. Alert (symptom based)

| Field | Value |
|---|---|
| Name | `AuthSignInErrorBudgetBurn` |
| Trigger | Share of `5xx` on `POST /login` above **5%** during **5 minutes** (budget burning about ten times faster than allowed) |
| Severity | Page (urgent): users cannot sign in, and every protected area of the platform depends on it |
| Owner | Oscar (`Oscarsl10`), owner of the Auth domain |
| Notification | Grafana alert to the team channel |
| Runbook | Section 4 of this document |

It is a symptom alert (what the user feels), not a cause alert: it does not look at CPU, memory or the database separately.
A second, lower alert on p95 latency above the SLO for 15 minutes would be added after the first measurements, to avoid
alert fatigue before the real values are known.

## 4. Runbook: sign-in fails with 5xx

**Symptoms:** clients get `500` or `503` from `/api/v1/auth/login`; the portal shows that sign-in is not available; the gateway
answers `503` with `SERVICE_UNAVAILABLE` when the upstream is down.

**Confirm:**
1. Take a failing response and read its `traceId` (or the `X-Correlation-Id` header).
2. Find that id in the gateway access log (`correlationId`) to see the status the gateway returned and from which upstream.
3. Find the same value in the `csp-auth-api` log (`traceId`): the error line (`login user failed`, `refresh session failed`) carries the cause.
4. Check the health routes of the service (`GET /api/v1/auth/health` and `/health/ready`) and whether the container is restarting.

**First mitigations:**
- Database not reachable or not ready: the service restarts until PostgreSQL answers; check the PostgreSQL container of the platform and
  the database credentials of the Auth domain.
- Refresh and logout failing but login fine: check the Redis instance of the Auth state; the service logs when the revoked-token list is unavailable.
- A failure that started with a release: roll back the Auth image to the previous tag (releases are tagged), then investigate.
- `429` is not an outage: it is the gateway rate limit (10 requests per second per client with a burst of 20).

**Escalate:** if the cause is not found in 30 minutes, tell the team channel and the owners of the platform composition
(`csp-infra-postgres`, `csp-infra-redis`).

**Verify recovery:** a successful `POST /login` with a test user through the gateway, the `5xx` share back under 1% for 15 minutes, and the
alert resolved.

## 5. Observability stories (drafts, not filed on the board)

They are written to be tested; they are not created as issues because the team rule is not to open new stories without approval.

**OBS-AUTH-A - Every Auth request is logged with its trace id**
```gherkin
Scenario 1: Request with a correlation id
  Given the gateway forwards a request with X-Correlation-Id "abc-123"
  When  csp-auth-api answers POST /login
  Then  the service writes one JSON log line with trace_id "abc-123", method, route, status and duration_ms
  And   the line contains no password, token or Idempotency-Key value

Scenario 2: Request without a correlation id
  Given a request reaches csp-auth-api without X-Correlation-Id
  When  the service answers
  Then  it generates a UUID v4 and uses it in the log line and in the traceId of any error body
```

**OBS-AUTH-B - Auth exports the golden signals**
```gherkin
Scenario 1: Rate and errors
  Given csp-auth-api exports metrics to the OpenTelemetry collector
  When  a client calls POST /login and the service answers 401 and then 200
  Then  http_server_requests_total increases by one for status 401 and one for status 200, labelled with the route pattern

Scenario 2: Latency
  Given the same setup
  When  a request is answered
  Then  http_server_request_duration_seconds records its duration in the histogram of the route
```

**OBS-AUTH-C - Sign-in availability alert**
```gherkin
Scenario 1: The alert fires
  Given the share of 5xx answers of POST /login is above 5% for 5 minutes
  Then  the alert AuthSignInErrorBudgetBurn fires to its owner with the link to the runbook

Scenario 2: The alert does not fire on client errors
  Given many POST /login answers are 401 (wrong password)
  When  no answer is a 5xx
  Then  the alert does not fire
```

## 6. Self-check (answers)

- **SLI / SLO / error budget:** the SLI is the measured share of good requests, the SLO is the target for it (99.5% in 30 days), the error budget is what is left to fail (0.5%).
- **What to alert on:** the 5xx share and the p99/p95 latency of login, which users feel; not CPU.
- **When the budget runs out:** feature work on Auth stops and only reliability fixes are merged (policy in section 2).
