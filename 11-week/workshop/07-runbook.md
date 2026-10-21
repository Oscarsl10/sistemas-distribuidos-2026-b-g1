# CP8 - Runbook: `AuthRegisterFastBurn`

Alert: [06-alerts.yml](06-alerts.yml), more than 7.2 % of `POST /api/v1/auth/register` answered `5xx` in the last hour and in
the last 5 minutes, with at least 10 requests in the hour. Severity `page`, team `auth`. The SLO behind it is in
[05-slo.md](05-slo.md): 99.5 % in 30 days, a budget of 75 failed requests.

Commands run from the repository `csp-infra-postgres`, which composes the platform (project `csp-platform`; services
`api-gateway`, `auth-api`, `postgres`). `ENV` is `dev`, `qa` or `main`.

```bash
cd csp-infra-postgres
ENV=qa
alias dc='docker compose --env-file env/$ENV.env'
```

## 1. Symptoms

- Clients get `500`, `502` or `503` from `POST /api/v1/auth/register`; the sign-in portal does not let the person create an account.
- The gateway answers `503` with `SERVICE_UNAVAILABLE` when `csp-auth-api` does not answer; `csp-auth-api` answers `500`
  with the code `INTERNAL_ERROR` when a step fails inside.
- A `429` is not this alert: it is the gateway rate limit (10 requests per second per client, burst of 20).

## 2. Confirm

1. Is the alert real? The error ratio and the number of requests in the last hour (Prometheus on `127.0.0.1:9090`):

   ```promql
   sum(rate(auth_register_requests_total{outcome="error"}[1h])) / sum(rate(auth_register_requests_total[1h]))
   sum(increase(auth_register_requests_total[1h]))
   ```

   With about 21 requests per hour, fewer than 10 means the guard should have stopped the alert: look at the rule before acting.
2. Take the `traceId` of a failed answer (field of the error body, header `X-Correlation-Id`) and follow it through the logs:

   ```bash
   dc logs --no-color api-gateway auth-api | grep "<trace_id>"
   ```

   The gateway line says which status it returned and from which upstream; the `csp-auth-api` line with level `ERROR`
   names the step that failed (`register user failed`).
3. Look at the state of the dependencies and of the service:

   ```bash
   dc ps auth-api postgres
   curl -s -o /dev/null -w "%{http_code}\n" localhost:8000/api/v1/auth/health/ready
   ```

   A service that restarts every few seconds, or `postgres` not `healthy`, points to the database.
4. Did it start right after a release? Compare the time of the first error with the time the last version of `csp-auth-api` was deployed.

## 3. First mitigations, from the lowest risk to the highest

| # | Mitigation | Risk | When |
|---|---|---|---|
| 1 | Restart the service: `dc restart auth-api` | Low. The data of a registration is saved in PostgreSQL, not in the service; the requests in flight fail, and a registration is idempotent by its `Idempotency-Key`, so the client can send it again. | Errors without a clear cause in the log, or after a configuration change. |
| 2 | Recover PostgreSQL: if `postgres` is not `healthy`, read `dc logs postgres` and restart it with `dc restart postgres` | Medium. It is the single instance of the environment, shared with booking and with `csp-worker`: tell the owner of Booking before restarting it. | The `ERROR` lines of `csp-auth-api` are timeouts or connection errors to the database. |
| 3 | Go back to the previous release: check out the previous tag in the sibling repository `csp-auth-api`, then `./scripts/up.sh $ENV` to rebuild | High. The image is built from the source, so it needs a rebuild; the schema is not rolled back (Flyway migrations only go forward), so the previous version must work with the current schema. | The errors started right after a release and 1 and 2 did not help. |

Not a mitigation for this alert: turning off the Auth relay (`AUTH_OUTBOX_RELAY_ENABLED`). The registration only writes the
event to `auth.outbox_event` inside its transaction; the relay is asynchronous and publishes it afterwards, so it cannot cause these `5xx`.

Do not change the gateway `limit_req` for this alert either: it only limits clients with `429`, it does not fix a `5xx` of the service.

## 4. Escalation

| After | Who |
|---|---|
| Immediately | Oscar, owner of the Auth domain (`csp-auth-api`, `csp-auth-db`), first responder. |
| 30 minutes without a cause | The team channel, with the `traceId` of a failed request and what was already tried. |
| The cause is PostgreSQL, Redis or the platform composition | The owner of the other domain that shares the instance (Booking, Yeferson); `csp-infra-*` belong to the whole team. |
| 100 % of the error budget spent (75 failed requests in the window) | The leader is told and the policy of [05-slo.md](05-slo.md) applies: only fixes and rollbacks are merged. |

## 5. Verify the recovery

1. A registration through the gateway with a test account (synthetic email, never a real person) answers `201`:

   ```bash
   curl -i -X POST localhost:8000/api/v1/auth/register -H "Content-Type: application/json" \
     -H "Idempotency-Key: runbook-check-0001" -d '{...body of a test account...}'
   ```

2. In the last 15 minutes no request of `POST /api/v1/auth/register` answered `5xx` in the gateway log.
3. The error ratio of the last hour goes below 1 %, and `AuthRegisterFastBurn` is no longer firing (the 5-minute window clears first).
4. A line is added to the weekly status with the time, the cause and the budget spent, and a postmortem is written if more than 50 % of it is gone.
