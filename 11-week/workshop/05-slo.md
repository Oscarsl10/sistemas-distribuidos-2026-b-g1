# CP6 - SLIs, SLOs and error budget

Journey: the registration of [01-journey.md](01-journey.md), `POST /api/v1/auth/register` on `csp-auth-api` through
`csp-api-gateway`, and the event `UserRegistered` that `csp-worker` publishes from `auth.outbox_event`. The PromQL uses the
metrics of [04-metrics.md](04-metrics.md): `auth_register_requests_total`, `auth_register_request_duration_seconds` and
`outbox_relay_event_age_seconds`.

None of these metrics is exported by the real services yet, so the SLOs are targets that cannot be measured today. They
are written so that the day the metrics exist, the queries can be copied as they are.

## 1. Traffic estimate

There is no real traffic: the platform is a course project that only runs in `dev` and in demos. I chose a pilot of one cinema:
**500 registrations per day**, that is 15,000 in a 30-day window, or 0.006 requests per second on average.

Why this number: I assume that a single cinema sells a few hundred seats a day, and only people who have no account have to register, so 500 new accounts
a day is already a busy day. On a premiere day I expect up to three times more (1,500), but the SLOs are computed on the average, because the window is 30 days.

A consequence of such low traffic: a five-minute window sees about two requests, so one failure moves the ratio a lot.
That is why the windows of the alerts in CP7 start at one hour.

## 2. The three SLOs

| | Availability | Latency | Freshness |
|---|---|---|---|
| **What the user feels** | The registration answers without a server error. | The form answers fast enough. | The new account reaches Ticketing soon, so the welcome message and the contact projection are not late. |
| **SLI formula** | requests not answered with `5xx` / all requests | requests answered in 1 second or less / all requests | minutes in which the oldest unpublished `UserRegistered` is 60 seconds old or less / all minutes |
| **Target** | 99.5 % | 95 % (p95 under 1 s) | 99 % of the minutes |
| **Window** | 30 days, rolling | 30 days, rolling | 30 days, rolling |

An answer `4xx` (for example `409 EMAIL_ALREADY_REGISTERED`) is a correct answer, so it counts as good for availability.

### PromQL of each SLI

```promql
# Availability
sum(increase(auth_register_requests_total{outcome!="error"}[30d]))
  / sum(increase(auth_register_requests_total[30d]))

# Latency: the bucket le="1" is the limit chosen in CP4
sum(increase(auth_register_request_duration_seconds_bucket{le="1"}[30d]))
  / sum(increase(auth_register_request_duration_seconds_count[30d]))

# Freshness: a gauge, so the SLI is the share of minutes in which it was within the limit
avg_over_time((max(outbox_relay_event_age_seconds{event_type="UserRegistered"}) <= bool 60)[30d:1m])
```

The freshness SLI is time based because in CP4 I defined `outbox_relay_event_age_seconds` as a gauge, not as a histogram. The
relay runs every 5 seconds, so 60 seconds leaves room for twelve runs before the SLO is missed.

## 3. Error budget

The budget is 100 % minus the target. With 15,000 registrations in 30 days (43,200 minutes):

| SLO | Budget | In requests | In minutes |
|---|---|---|---|
| Availability 99.5 % | 0.5 % | 75 failed requests per month | 216 minutes (3 h 36 min) of full outage |
| Latency 95 % | 5 % | 750 requests slower than 1 s | 2,160 minutes (36 h) with every request slow |
| Freshness 99 % | 1 % | not counted in requests | 432 minutes (7 h 12 min) with the oldest event older than 60 s |

With this traffic 75 failed requests is a small number: ten errors in one afternoon already spend 13 % of the availability
budget. I accept that, because a registration that fails is a lost user.

## 4. Error budget policy

The policy applies to the availability SLO, the one that decides if people can create an account. The owner of the Auth domain
decides when it starts and ends, and writes it in the weekly status.

| Budget spent in the window | What the team does |
|---|---|
| **50 %** (37 failed requests, or 108 minutes) | Normal work continues. The failures are reviewed with their `trace_id` and written in the weekly status. Every PR that changes `register` or the outbox says in its description how it is rolled back. |
| **75 %** (56 failed requests, or 162 minutes) | Changes to the registration path (contract, schema, relay) are frozen, except fixes. Any other change in `csp-auth-api` needs a second review. |
| **100 %** (75 failed requests, or 216 minutes) | Feature work in Auth stops. Only the fixes for the cause of the failures and rollbacks are merged, until the 30-day window recovers. A postmortem is written and the leader is told. |

What the team can really do today is limited: the CI is locked and there is no on-call rota, so "stop" means that the
owner does not merge new stories of Auth, not that a system blocks the deploy.
