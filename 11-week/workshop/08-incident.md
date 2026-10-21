# CP9 - Incident drill and blameless postmortem

Drill on the registration journey of [01-journey.md](01-journey.md). The critical hop marked in CP1 is hop 3: `csp-auth-api` writing to
PostgreSQL. **Incident of the drill:** on 2026-10-21 from 14:00 to 14:40 UTC the shared PostgreSQL instance answers each statement in 3 seconds for 25 % of
the registrations. Everything else works. The numbers come from [05-slo.md](05-slo.md) (500 registrations per day, budget of 75 failed requests),
[06-alerts.yml](06-alerts.yml) and [07-runbook.md](07-runbook.md). It is a drill: the cause is assumed, not observed.

## 1. Why a slow database becomes an error

Facts from the configuration of the repositories: the registration runs, in order, bcrypt and then these statements in one transaction: look up the
idempotency key, check the email, insert the user, insert the key, insert the `UserRegistered` event. The service waits up to 5 seconds for each statement
(`APP_AUTH_DATABASE_QUERY_TIMEOUT`), so 3 seconds is slow but not an error for the service. The gateway waits only 10 seconds in total (`proxy_read_timeout 10s`) and answers `503`
with `SERVICE_UNAVAILABLE` after that. Three statements of 3 seconds take about 9 seconds and the fourth does not fit: a request hit by the slow database is cut by the gateway.
The expected rate of 25 % is applied to the 14 requests of the 40 minutes (500 per day is 0.35 per minute); in this drill 4 of the 14 were hit.

## 2. Trace of one slow request (14:14:02)

`trace_id` `b83f6d1e4a9c47e2905d71c3a4e8f026`. The duration of bcrypt is an assumption (0.25 s), as in CP4.

```
POST /api/v1/auth/register                         csp-api-gateway     10,004 ms  503
  proxy to csp-auth-api                            csp-api-gateway     10,001 ms  read timeout of 10 s
    POST /register                                 csp-auth-api        10,011 ms  canceled when the gateway closes the connection
      validate command and idempotency key         csp-auth-api             2 ms
      bcrypt hash of the password                  csp-auth-api           250 ms
      transaction register user                    csp-auth-api         9,759 ms  rolled back, nothing is saved
        find idempotency key                       PostgreSQL           3,004 ms  slow
        check that the email is not registered     PostgreSQL           3,002 ms  slow
        insert app_user and user_role              PostgreSQL           3,001 ms  slow
        insert idempotency_key                     PostgreSQL             745 ms  canceled
```

The `outbox_event` insert is never reached, so no `UserRegistered` is created for a failed registration.

## 3. Four log lines of that trace (CP2 format)

```json
{"timestamp":"2026-10-21T14:14:05.256Z","level":"WARN","service":"csp-auth-api","version":"3.0.0","environment":"qa","trace_id":"b83f6d1e4a9c47e2905d71c3a4e8f026","span_id":"4c9e1f7a20b63d85","correlation_id":"6d2a8f40-1b7c-4e93-a5d0-92c8e1f3b764","event":"db.statement.slow","message":"find idempotency key took 3 s","duration_ms":3004}
{"timestamp":"2026-10-21T14:14:08.258Z","level":"WARN","service":"csp-auth-api","version":"3.0.0","environment":"qa","trace_id":"b83f6d1e4a9c47e2905d71c3a4e8f026","span_id":"4c9e1f7a20b63d85","correlation_id":"6d2a8f40-1b7c-4e93-a5d0-92c8e1f3b764","event":"db.statement.slow","message":"email check took 3 s","duration_ms":3002}
{"timestamp":"2026-10-21T14:14:12.004Z","level":"ERROR","service":"csp-api-gateway","version":"3.0.0","environment":"qa","trace_id":"b83f6d1e4a9c47e2905d71c3a4e8f026","span_id":"a7d3052e8c1f9b46","correlation_id":"6d2a8f40-1b7c-4e93-a5d0-92c8e1f3b764","event":"http.request.completed","message":"POST /api/v1/auth/register upstream did not answer in 10 s","duration_ms":10004,"http_status":503}
{"timestamp":"2026-10-21T14:14:12.011Z","level":"ERROR","service":"csp-auth-api","version":"3.0.0","environment":"qa","trace_id":"b83f6d1e4a9c47e2905d71c3a4e8f026","span_id":"4c9e1f7a20b63d85","correlation_id":"6d2a8f40-1b7c-4e93-a5d0-92c8e1f3b764","event":"register.canceled","message":"request canceled by the gateway, transaction rolled back","duration_ms":10011}
```

The slow-statement line of the third statement is not shown. The two `db.statement.slow` lines are what the trace would show after action item A1: during the incident they did not exist, and the owner only had the last two lines (the `503` of the gateway and the error of `csp-auth-api`), which say that the request was cut but not which step was slow.

## 4. What the SLIs reach

| | In the 40 minutes of the incident | In the 30 days (15,000 requests) |
|---|---|---|
| Requests / failed (`503`) | 14 / 4 | 15,000 / 4 |
| Availability SLI | 10 / 14 = 71.4 % (target 99.5 %) | 14,996 / 15,000 = 99.973 % |
| Error ratio and burn rate | 28.6 %, burn rate 28.6 / 0.5 = 57 | |
| Latency SLI (answers in 1 s or less) | 10 / 14 = 71.4 % | 99.973 % (target 95 %) |
| Availability budget spent | 4 / 75 = **5.3 %**, equal to 11.5 of the 216 minutes | |

Freshness is not affected in the drill: the failed registrations never wrote an event. The policy of CP6 does not start (under 50 %).

## 5. Which alert fires, and when

| Minute | Event |
|---|---|
| 6 (14:06) | First failed registration. Error ratio of the last hour: 1 / 21 = 4.8 %, under 7.2 %: **no alert**. |
| 14 (14:14) | Second failed registration (the trace above). 2 / 21 = 9.5 % in the hour and the 5-minute window has the failure; 21 requests in the hour pass the guard of 10. |
| **16 (14:16)** | **`AuthRegisterFastBurn` fires** (page) after its `for: 2m`. It is the first and the only alert during the incident. |
| 51 (14:51) | `AuthRegisterSlowBurn` fires (ticket): the fourth failure at 14:36 puts the 6-hour ratio at 4 / 125 = 3.2 %, above 3 %, and the rule waits 15 minutes. It arrives 11 minutes after the recovery, so it serves as a reminder to write the postmortem, not to detect. |

With so little traffic the alert needs two failed requests: it fires 10 minutes after the incident starts only because the second failure came early.

## 6. Which runbook step stops the bleeding

The step that stops it is **mitigation 2, recover PostgreSQL**: in `dc logs postgres` the long statement of another domain is visible, and
stopping it (or restarting the instance, after telling the owner of Booking) brings the answers back to normal at minute 40. Mitigation 1 (`dc restart auth-api`) was tried first, at minute 24, and
changed nothing, because the slowness is in the database, not in the service.

## 7. Blameless postmortem

**Summary.** For 40 minutes, 4 of 14 registrations failed with `503` because the shared PostgreSQL answered slowly. 5.3 % of the monthly budget was spent. No data was lost: the failed transactions were rolled back and the clients could repeat the registration with the same `Idempotency-Key`.

**Timeline (UTC).** 14:00 PostgreSQL starts to be slow. 14:06 first `503`. 14:14 second `503`. 14:16 `AuthRegisterFastBurn` pages. 14:19 the owner follows the `traceId` of a failed answer and finds the `503` of the gateway and the error of `csp-auth-api`, but not which step was slow. 14:24 `auth-api` is restarted, no effect. 14:34 the owner reads the logs of `postgres` and finds a long statement from another schema. 14:38 the owner tells the owner of Booking and stops the statement. 14:40 the answers are normal. 14:51 the slow alert fires and is closed as already handled.

**Impact.** 4 failed requests of 14 (28.6 %), 4 of the 75 allowed (5.3 %), 11.5 of the 216 budget minutes; the people affected could not create an account at that moment.

**Root cause (assumed for the drill).** A long statement run by hand on the shared PostgreSQL instance took the disk for 40 minutes; Auth, Booking and the worker use the same instance (Annex J).

**Contributing factors.** (1) The service has no log of a slow statement, so the owner needed about 15 minutes (14:19 to 14:34) to see where the time went. (2) The runbook starts with the restart of the service even when the logs show the database is slow. (3) The gateway timeout (10 s) is shorter than the sum of the statements, so a slow database becomes an error and not just a slow answer. (4) One instance for all the domains and no per-role statement limit.

**Action items.**

1. **A1. As the owner of Auth, I want a `WARN` line `db.statement.slow` with the `trace_id` for every statement of the registration that takes more than 1 second, so that I can find the slow step with one query.**
   - Given a statement that takes 1.5 s, when the registration runs, then one `WARN` line is written with `event`, `duration_ms` and the same `trace_id` as the request.
   - Given statements under 1 s, then no such line is written.
   - A unit test with a fake slow port covers both cases, and no value of the request (email, password, key) appears in the line.
2. **A2. As the person on call, I want the runbook to say "if the slow statements are in the log, go straight to PostgreSQL", so that I do not lose time restarting the service.**
   - The runbook of CP8 has this rule before mitigation 1, and a step to list the statements that are running (`pg_stat_activity`) before restarting `postgres`.
   - A drill of the same incident, following only the runbook, stops the slowness within 10 minutes of the page.
3. **A3. As the owner of the platform composition, I want a `statement_timeout` for each login role of the shared PostgreSQL instance (`auth_app`, `booking_app`, `worker_app` and the administrator used by hand), so that a long statement of one user cannot slow the other domains.**
   - Given a statement over the limit, when the role runs it, then PostgreSQL cancels it and the error is logged.
   - The registration of Auth keeps its p95 under 1 s while a long statement of another role is cancelled, checked in a test of the composition.
