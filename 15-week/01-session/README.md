# Week 15 - Session 1: Testing at scale: contract, e2e, load and chaos

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: apply the testing pyramid to a distributed system, use real dependencies in the integration tests, keep contract tests between the
services, cover the money paths end to end, run a load test against the SLO and prove the resilience with one chaos experiment.

**Scope of this work.** Nothing is applied to the `csp-*` repositories from this session. The project already has part of these tests, and Section 1 reports
their results as they are written in the pull requests of the repositories: **I did not run them again today**, except the 16 tests of the Redis circuit breaker and the 7 of the bulkhead of Section 2, and the tests of Section 3, which I ran today on a copy of the release stack. Each row says on which commit it was run.
The pieces the project does not have are done as a rehearsal in [rehearsal/](rehearsal/), with a stand-in of a checkout and of a payment provider:
an integration test with Testcontainers, a contract test, an end-to-end test, a load test against an SLO, one chaos experiment and a proposal of the workflow that would run them.
Section 3 has the tests on the copy of the release stack, and Section 5 the real output of the rehearsal.

## 1. The pyramid for a distributed system, and what the project has

| Layer | What exists in the project | Result (source) |
|---|---|---|
| Unit | `go test ./...` in `csp-auth-api`; `mvn verify` in `csp-workflow`; `ng test` in `csp-auth-portal` | `csp-auth-api` release head: 1196 tests passed, contract tests included. `csp-auth-portal`: 289 of 289. `csp-workflow` on `7e160ee`: 118 tests, 0 failures ([csp-docs#282](https://github.com/code-corhuila/csp-docs/issues/282) closing comment, 2026-11-04; [csp-workflow#13](https://github.com/code-corhuila/csp-workflow/pull/13)) |
| Integration with real dependencies | `go test -tags integration ./...` in `csp-auth-api` against PostgreSQL and Redis; `JdbcSagaStoreTest` of `csp-workflow` against PostgreSQL (`TEST_DATABASE_URL`); `db-ci.yml` of each `-db` repository (migrate, test, rollback) | `csp-auth-api`: 1449 passed, 1 skipped (file modes on Windows), 0 failed. `csp-auth-db`: the 11 steps of `db-ci.yml` passed, 45 migrations with rollbacks. A full-stack run with a real Mailpit and PostgreSQL 16 passed 19 of 19 checks (same closing comment). `csp-workflow#12`: 104 tests, 8 of them `JdbcSagaStoreTest` against the real database |
| Contract | The responses of `csp-auth-api` are validated against `auth-service.yaml`; the portal checks its own requests against it (`npm run test:contract`) | The portal check matches 14 operations of `auth-service.yaml` 1.17.0 (closing comment of #282) |
| Gateway | `tests/smoke.sh` and `tests/integration.sh` of `csp-api-gateway` | Release head: 97 smoke checks ok, 0 failures; 307 integration checks passed, 0 failed (closing comment of #282) |
| End to end (money path) | `tests/e2e/up.sh` and `tests/e2e/book-seats.sh` of `csp-workflow`: the saga that holds and confirms the seats, against the real `csp-booking-api` and `csp-catalog-api` through the gateway | `TOTAL passed=89 failed=0`, twice, on `c635586` ([csp-workflow#15](https://github.com/code-corhuila/csp-workflow/pull/15)); `TOTAL passed=105 failed=0`, twice, on `76256cb` ([csp-workflow#17](https://github.com/code-corhuila/csp-workflow/pull/17)) |
| Chaos (partial) | The same end-to-end script kills the container of the workflow in the middle of a saga and makes the compensation fail on purpose; the Redis circuit breaker and the bulkhead of `csp-auth-api` have unit tests | See Section 2: 16 breaker tests and 7 bulkhead tests passed today on `a998029` |
| Load | Nothing in the repositories | Section 3 (a copy of the release stack) and Section 5 (rehearsal) |

Where the real dependencies come from: the project does not use the Testcontainers library. Real PostgreSQL and Redis come from containers that the scripts or the
workflow start, or from the persistent development environment. The rehearsal does use Testcontainers, because it is the tool of the session. For validating the project my rule is to use the persistent
development and qa instances and not throwaway containers, so using Testcontainers there would be a decision to take with the team first. I searched the repositories for `testcontainers`, `k6`, `gatling`, `jmeter` and `chaos` in the
workflow, build and test files and found none. The contract tests are OpenAPI checks, not consumer-driven contracts with Pact.

## 2. Chaos in the project: what the end-to-end test already injects

Both experiments are in `tests/e2e` of `csp-workflow`, and they ran in a **disposable copy of the platform** (Compose project `csp-e2e`) built from `origin/develop`, not in the persistent environments.

| Fault injected | What was verified | Source |
|---|---|---|
| The workflow container is killed in the middle of a saga | The restarted service finishes the saga from where it stopped, with the same key `<sagaId>:<step>` for each step, and the steps keep the original `X-Correlation-Id` | [csp-workflow#15](https://github.com/code-corhuila/csp-workflow/pull/15), [#17](https://github.com/code-corhuila/csp-workflow/pull/17) |
| A trigger of the test database makes every update of a reservation fail, so Booking answers `500` to the confirm and to the release | The saga ends `FAILED` with `failedStep` `confirm-reservation`, each call was tried three times with the same key, and the hold stays `HELD` until it expires | [csp-workflow#17](https://github.com/code-corhuila/csp-workflow/pull/17) |
| A recording proxy delays a confirm until the hold expires in Booking | The confirmation fails and the saga compensates it | [csp-workflow#15](https://github.com/code-corhuila/csp-workflow/pull/15) |

**A circuit breaker with its fault tests.** `csp-auth-api` has one Circuit Breaker, around the calls to Redis of the attempt limits (`attemptsBreaker`, ADR-048): every call waits at most
a timeout, and after a failure Redis is not called again until a pause has passed. The unit tests that inject the failure (a store that does not answer, a store that fails, a cancelled request that must not
pause it, and the store tried again when the pause ends) were run today on `release/3.0.0`, commit `a998029`, clean working tree:
`go test ./internal/adapter/out/redis -run 'Pause|Paused|Timeout|Breaker|GivenUp|GivesUp|NotCalledDuring' -count=1 -v`, **16 `--- PASS`, 0 `--- FAIL`**, `ok` in 2.8 s. It is a unit-level test with a fake store, not an experiment on the running platform.

**A bulkhead under load.** `csp-auth-api` also has a Bulkhead (`internal/adapter/in/httpapi/bulkhead.go`, norm 5.3.10): it caps how many requests of a route run at once, makes the next ones wait a short, bounded time and then
refuses them with `429`, so a slow dependency cannot drain the connection pool of the other routes. Its tests ran today on the same commit `a998029`:
`go test ./internal/adapter/in/httpapi -run Bulkhead -count=1 -v`, **7 `--- PASS`, 0 `--- FAIL`**, `ok` in 4.4 s, including `TestRegisterAndLoginRefuseWith429WhenTheBulkheadIsFull` and
`TestVerifyEmailChangeIsNotLimitedByTheBulkhead`. They are unit tests, not a load run: nobody has sent the traffic that fills the bulkhead on the real service.

What is **not** injected: network latency, dropped messages and the death of Booking or Catalog in the middle of the traffic. I searched the repositories for `circuit breaker`, `resilience4j`, `gobreaker` and `hystrix`: the only breaker is the Redis one in
`csp-auth-api`. The breaker that protects a call to another service, which is what the second half of the session asks for, is only in the rehearsal. Section 3 confirms on the running stack that the gateway has none: when `csp-auth-api` is killed, it answers every request with a `503` after about 2 seconds.

## 3. Real tests on a copy of the release stack

**How it was done.** The stack `csp-release` (release/3.0.0 of every repository, persistent, with the data of the team) was stopped by me, and I did not touch its volumes. I made a copy as a second Compose project (`csp-w15`),
outside the repositories, with its own random passwords and keys and **empty data**, so the accounts of my teammates were not copied. I reused the images already built for `csp-release`, tagged with the new name and without rebuilding, and started only what the tests
need: PostgreSQL, Redis, Mailpit, `csp-auth-api` and `csp-api-gateway`. The port of `csp-auth-api` was published only in the copy, to measure it without the gateway. One synthetic user was created through the normal flow
(register `201`, link in Mailpit, confirm `201`, login `200`). When the tests ended I removed the copy (containers, volumes, image tags and its secrets) and checked that the 24 containers, the 4 volumes and the 15 images of `csp-release` were still there.
The scripts are in [real-stack/](real-stack/) (`seed.py` and `tests.py`), and they read the synthetic credentials from the environment (`W15_EMAIL` and `W15_PASSWORD`). They expect the copy as described: Compose project `csp-w15`, the gateway on port 8000, Mailpit on 8025, and the port 18081 published only for `csp-auth-api`.

**The limits that decide what can be measured**, read in the clone of `csp-api-gateway` (`nginx/`) and in the defaults of `csp-auth-api` (`internal/config`):

| Where | Limit |
|---|---|
| Gateway, every route | 10 requests per second per client address, burst 20, answer `429` |
| Gateway, `login`, `register` and `register/confirm` | their own bucket: 2 per second, burst 3 (ADR-060) |
| Gateway, `forgot-password` | 1 per second, burst 2 |
| `csp-auth-api`, passwords | bcrypt cost 12, at most 4 hashes at once, a wait of 100 ms and then `503` |
| `csp-auth-api`, failed logins | 10 failures of one address lock it for 300 s (`429`); Redis is called with a 300 ms timeout and paused for 10 s after a failure |

The SLO used here is the real one, from ADR-061 (D3): p95 of `POST /login` and `POST /register` at most 1 second. The traffic the workshop estimates is about 500 registrations a day.

### A. The credential bucket of the gateway (the criterion of story #281)

```text
== A. the credential bucket of the gateway (2 r/s, burst 3) on POST /login, valid credentials
12 requests in a row: 200 200 200 200 200 200 200 429 429 429 429 429 | Retry-After of a 429: 1
sustained 5 req/s for 10 s: sent 50 | answers {200: 22, 429: 28} | expected about 4 + 2 x 10 = 24 accepted
```

The first requests pass (the burst plus the refill while the bcrypt of each login runs, about 0.3 s), and then every request gets `429` with `Retry-After: 1`. At 5 requests per second, 22 of 50 were accepted, which is the 2 per second of the bucket.

### B. The per-client limit of the gateway on a cheap route

```text
== B. the per-client limit of the gateway (10 r/s, burst 20) on GET /jwks
  rate  sent   200   429 other   p95 ms
     8    48    48     0     0       27
    15    90    77    13     0       23
    30   180    77   103     0       26
```

The gateway accepts about 10 requests per second per address plus the burst of 20 (about 77 of 90 and of 180 in 6 seconds) and refuses the rest with `429`. Below that, the latency of an answer is about 25 ms. So **from one address nothing can load the services through the gateway**: the real capacity is measured in C without it.

### C. Capacity of `csp-auth-api` without the gateway (login with valid credentials)

```text
== C. capacity of csp-auth-api, direct (no gateway), POST /login with valid credentials; SLO of ADR-061: p95 <= 1000 ms
 rate  sent   200  429  503 other    p50    p95    p99 p95 of 200  SLO
    5    50    50    0    0     0    222    252    265        252  met=True
   10   100   100    0    0     0    228    252    266        252  met=True
   15   150   150    0    0     0    232    255    267        255  met=True
   20   200   175    0   25     0    301    339    364        344  met=False
   30   300   177    0  123     0    300    342    353        346  met=False
   40   400   176    0  224     0    130    353    366        362  met=False
```

- The service serves about **17.6 logins per second**, which is what the hash bound gives: 4 hashes at once, each of about 0.22 s. Up to 15 per second every login is answered `200` with p95 of about 255 ms. At 20 per second, 175 of 200 are served and 25 are refused; at 40 per second, 176 of 400.
- What it refuses, it refuses with `503` and the code `SERVICE_UNAVAILABLE` and `Retry-After: 1`, and the logins it accepts keep the same latency (the p95 of the `200` stays between 252 and 362 ms). **The SLO of latency (p95 at most 1 s) is met at every rate; what falls is the share of logins served**: 44 % at 40 per second.
- The criterion of story #281 says that the refusal must come in less than 200 ms. In the log of the service (1,240 logins, all this test and a burst of 40 at once), the `503` took **p50 101 ms, p95 140 ms and at most 725 ms**, and the `200`, p50 284 ms and p95 331 ms. So the criterion holds at the 95th percentile and not for the slowest ones in a burst. My client measured a p95 of 662 ms for the same `503`; the difference is my Python client with 40 threads, not the service.

### D. Chaos: Redis stopped while the logins go on

```text
== D. chaos: Redis stopped while the logins go on (direct to csp-auth-api)
redis stopped: Exited (0) Less than a second ago
15 wrong-password logins with Redis down: 401 401 401 401 401 401 401 401 401 401 401 401 401 401 401 | total 3.9s | first request 885 ms, the rest p95 218 ms
3 valid logins with Redis down: 200 200 200
redis back: Up 7 seconds (healthy) | waiting 12 s for the pause of the breaker (10 s) to end
control, Redis up again, 13 wrong-password logins: 401 401 401 401 401 401 401 401 401 401 429 429 429 (the lock of the 11th attempt is the criterion of story #280)
WARN lines of auth-api in the last 2 minutes: 214 | messages: ['account login failure counter paused after a Redis failure, failed logins are not slowed meanwhile', 'login address locked', 'login attempt limit paused after a Redis failure, the login goes on without it', 'password hashing refused: too many hashes at once']
```

- **Redis down.** The 15 wrong-password logins were all `401`: **none was locked**, because the lock needs Redis and the limiter fails open. The first request took 885 ms (it waited for the Redis timeout), and the rest 218 ms at p95, because the breaker paused the calls to Redis. The 3 valid logins were `200`. This is the threat of the row S of the STRIDE table (the lock fails open when Redis does not answer), seen on the running service: someone who can make Redis fail can try passwords without a limit.
- **Control, Redis back.** After the pause of 10 s, 10 wrong passwords answered `401` and the 11th, `429`: **the criterion 1 of story #280 holds**.
- **The audit event of the lock (criterion 2 of #280)** is in the log of the service as `login address locked`, with `audit: login_locked`, `clientHash` (a SHA-256 of the address), `attempts: 10`, `lockedAt`, `lockSeconds: 300` and `traceId`. It has no email, no password and no address in clear.
- The service wrote the warning of the paused limiter once for each pause (2 and 2 in the 5 minutes, because I ran the experiment twice after fixing a mistake in my script) and 396 warnings `password hashing refused: too many hashes at once`, one for each `503` of C and of the burst of 40.

### E. Chaos: `csp-auth-api` killed in the middle of the traffic through the gateway

```text
== E. chaos: csp-auth-api killed in the middle of the traffic through the gateway (GET /jwks at 8 req/s)
  at t=12.0 s: docker start (restart policy: on-failure) | state before: exited
      window  sent   200   502   504 other  p95 ms
 before kill    48    48     0     0     0      25
      killed    48     0     0     0    48    2133
 after start    48    44     0     0     4    1552
   recovered    64    64     0     0     0      21
first 200 after docker start: 0.9 s
answers during the outage: 48, p95 2133 ms, max 2136 ms
```

- **There is no fast failure.** During the 6 seconds of the outage every request got `503` (the gateway turns `502`, `503` and `504` of the upstream into one `503`: `error_page 502 503 504 =503`), and each took about **2.0 s** to fail: `resolver_timeout` and `proxy_connect_timeout` are both 2 s. At 8 requests per second that is about 16 requests waiting at the same time. The rehearsal of Section 5 answers the same case in 1 ms because its breaker is open.
- **The recovery is automatic and quick.** The first `200` came 0.9 s after `docker start`, with no action on the gateway, which resolves the name of the service with the resolver of Docker (`resolver 127.0.0.11 valid=10s`) and found it again. A few requests right after the start still failed (4 in 6 s).
- **The restart policy did not help here.** The container has `on-failure`, and 6 seconds after `docker kill` it was still `exited`: it only came back when I started it.

### What these results say

| STRIDE row or story | Result on the running stack |
|---|---|
| S, criterion 1 of #280 | The lock of the 11th failed attempt works with Redis; with Redis down the limiter fails open |
| S, criterion 2 of #280 | The audit event of the lock has the expected fields and none of the sensitive ones |
| D, criterion 1 of #281 | The gateway refuses with `429` and `Retry-After` above 2 logins per second per address |
| D, criterion 2 of #281 | The hash bound refuses with `503 SERVICE_UNAVAILABLE` at p95 140 ms; the slowest took 725 ms |
| Resilience | No circuit breaker between the gateway and a dead service: each request fails after 2 s; the recovery is automatic |

## 4. The rehearsal

A small Go module with no dependencies, in [rehearsal/](rehearsal/). It does not touch any `csp-*` repository.

| File | What it is |
|---|---|
| `payments.go` | The payment provider: a real HTTP server on a real port. It takes 50 ms per payment and serves 20 at the same time (a capacity of 400 payments per second). It can be killed and started again on the same address |
| `checkout.go` | The money path: reserve a seat, pay, confirm; if the payment fails, release the seat (the compensation of a saga). A circuit breaker opens after 5 failures in a row, refuses for 500 ms and then lets one probe through |
| `contract/payments.contract.json` | What the checkout expects from the provider: `200`, `application/json`, the fields `id` and `status`, and the status `APPROVED` |
| `store.go` | The seats in PostgreSQL: table `seat_hold` with a primary key and a named `CHECK` on the status (the naming conventions of the project), and `Hold`, `Confirm`, `Release` and `Counts` |
| `integration_test.go` | The integration test: it starts a real PostgreSQL 16 in a container with Testcontainers and removes it at the end; it is skipped, not failed, when Docker is not running |
| `scale_test.go` | The contract, end-to-end, load and chaos tests below |
| `ci.yml` | A proposal of the workflow: contract and end to end, and integration, on every pull request; load and chaos once a week |

The tests:

- **Integration with a real database.** The seat store runs against PostgreSQL 16 in a container. Four cases: hold, confirm and release move the seats; the primary key refuses a second hold of the same order; the named `CHECK` refuses a status outside the list; and a seat that is no longer held cannot be moved. The three last cases are the ones a mock of the store accepts and a real database refuses.
- **Contract.** The checkout (consumer) writes what it expects in the JSON file, and the test checks that the provider answers exactly that.
- **End to end.** One order through the real HTTP entry and the real provider: the status is `201`, one seat is confirmed and none stays held.
- **Load.** An open-loop generator (it keeps the rate whatever the answers) sends orders for 2 seconds at 100, 200, 300, 400 and 600 requests per second. The SLO of the rehearsal is p95 at most 250 ms and every order confirmed, with a target of 200 requests per second. The test must fail if the target is not met, and reports the first rate that breaks the SLO.
- **Chaos.** At 100 orders per second for 5 seconds, the provider is killed at second 1 and started again at second 3. The test checks that the orders in flight are compensated, that the breaker opens and answers fast (`503`, p95 at most 50 ms), that the system recovers by itself, and that no seat stays held and the confirmed seats are exactly the orders answered `201`.

## 5. Real output of the rehearsal

Go 1.27.0 on Windows 11, 2026-11-11. Two commands, from `15-week/01-session/rehearsal`: `go test -count=1 -v -run 'Contract|E2E|Load|Chaos' ./...` for the four tests that need no Docker, and
`go test -count=1 -v -run Integration ./...` for the one that starts a container. The block below is the first one; I ran it again later the same day and the four tests passed again, with numbers that differ by a few milliseconds.

```text
=== RUN   TestContractPaymentsAnswersWhatCheckoutExpects
    scale_test.go:132: contract met: 200 application/json, fields [id status], status APPROVED
--- PASS: TestContractPaymentsAnswersWhatCheckoutExpects (0.02s)
=== RUN   TestE2EPlaceOrderReservePayConfirm
    scale_test.go:146: order confirmed: status 201, free 999999, held 0, confirmed 1
--- PASS: TestE2EPlaceOrderReservePayConfirm (0.02s)
=== RUN   TestLoadFindsTheBreakingPointAgainstTheSLO
    scale_test.go:154: rate         sent     ok       p50       p95       p99  SLO
    scale_test.go:168: 100           198    198      51ms      57ms      64ms  met=true
    scale_test.go:168: 200           364    364      52ms      63ms     185ms  met=true
    scale_test.go:168: 300           588    588      51ms      55ms      59ms  met=true
    scale_test.go:168: 400           784    784      55ms      67ms      72ms  met=true
    scale_test.go:168: 600          1104   1104     419ms     822ms     859ms  met=false
    scale_test.go:170: SLO: p95 <= 250ms and every order confirmed. Target 200 req/s met; first rate that breaks it: 600 req/s
--- PASS: TestLoadFindsTheBreakingPointAgainstTheSLO (11.15s)
=== RUN   TestChaosKillThePaymentsServiceMidFlow
    scale_test.go:211: phase                   confirmed  compensated  rejected rejected p95
    scale_test.go:217: before the kill                95            4         0           0s
    scale_test.go:217: provider down                   0            4       196          1ms
    scale_test.go:217: provider back                  91            0         9          1ms
    scale_test.go:217: recovered                     101            0         0           0s
    scale_test.go:240: no seat is left held: free 999713, held 0, confirmed 287 (= orders answered 201: 287)
--- PASS: TestChaosKillThePaymentsServiceMidFlow (5.06s)
PASS
ok  	cinesync/w15-session1-scale	20.586s
```

How to read it:

- **Load.** The provider can serve 400 payments per second, and the test finds the limit where the design says: the SLO holds up to 400 and breaks at 600 (p95 822 ms, and the requests wait for a free slot). The generator sent less than the nominal rate (for example 364 orders in 2 seconds at a nominal 200, an effective 182 per second, probably because of the coarse timer tick of Windows, which I did not check), so the real rate of the last step was about 550 per second.
- **Chaos.** The 4 orders compensated "before the kill" are the ones in flight when the provider died (100 per second times 50 ms is about 5). The breaker opened after 5 failures and, during the outage, 196 orders were refused at once (p95 1 ms) instead of waiting for a timeout; the 4 compensated ones are the probes that the breaker lets through every 500 ms. Within the first second after the restart, 91 orders were confirmed and 9 were still refused, and in the last two seconds everything was confirmed. No seat stayed held.

**Integration output** (the banner of Testcontainers is left out). The container was started once, used for the four cases and removed; afterwards no container or network of the test was left.

```text
=== RUN   TestIntegrationSeatStoreAgainstARealPostgres
=== RUN   TestIntegrationSeatStoreAgainstARealPostgres/hold,_confirm_and_release_move_the_seats
    integration_test.go:63: counts after hold a, b, c, confirm a and release b: map[CONFIRMED:1 HELD:1 RELEASED:1]
=== RUN   TestIntegrationSeatStoreAgainstARealPostgres/the_primary_key_refuses_a_second_hold_of_the_same_order
    integration_test.go:71: second hold of order a: ErrAlreadyHeld (SQLSTATE 23505, pk_seat_hold)
=== RUN   TestIntegrationSeatStoreAgainstARealPostgres/the_named_CHECK_refuses_a_status_outside_the_list
    integration_test.go:79: insert with status BOGUS: pq: new row for relation "seat_hold" violates check constraint "ck_seat_hold_status" (23514)
=== RUN   TestIntegrationSeatStoreAgainstARealPostgres/a_seat_that_is_no_longer_held_cannot_be_moved
    integration_test.go:91: confirm and release of the released order b, and confirm of an unknown order: ErrNotFound
--- PASS: TestIntegrationSeatStoreAgainstARealPostgres (15.91s)
    --- PASS: TestIntegrationSeatStoreAgainstARealPostgres/hold,_confirm_and_release_move_the_seats (0.02s)
    --- PASS: TestIntegrationSeatStoreAgainstARealPostgres/the_primary_key_refuses_a_second_hold_of_the_same_order (0.04s)
    --- PASS: TestIntegrationSeatStoreAgainstARealPostgres/the_named_CHECK_refuses_a_status_outside_the_list (0.00s)
    --- PASS: TestIntegrationSeatStoreAgainstARealPostgres/a_seat_that_is_no_longer_held_cannot_be_moved (0.01s)
PASS
ok  	cinesync/w15-session1-scale	16.246s
```

The two constraints are the point of the test: the database itself answers `23505` for a duplicate order and `23514` for a status that is not in the list, which no in-memory fake would do.

**The tests catch faults.** I changed a copy of the rehearsal twice, and the chaos test failed both times, as it must:

- Without the release of the seat after a failed payment, `held` stayed at 9: `seats: free 999705, held 9, confirmed 286, orders answered 201: 286`.
- With a breaker that never opens (threshold 1,000,000), `during the outage: compensated 189 (at most 30), rejected 0 (at least 100), rejected p95 0s (at most 50ms)`: every order waited for its own failure instead of being refused fast.

## 6. What this does not prove

- The rehearsal is one machine, with the checkout, the provider and the generator in the same process on the loopback interface. Real services add network, a database and other processes, so the numbers say how the design behaves, not how fast CineSync Platform is.
- The SLO of the rehearsal (p95 250 ms, a target of 200 per second) is my assumption to have a number; it is not an SLO of the project. The tests of Section 3 use the real one, from ADR-061.
- The tests of Section 3 are one machine: Docker Desktop with 8 CPUs, where the services and my load generator (a Python script) share the processor. The generator adds its own delay under many threads (it showed in the `503` of a burst), so I took the times that matter from the log of the service.
- Section 3 measures `csp-auth-api` without the gateway, which no client can do: through the gateway one address gets 2 logins per second. The 17.6 per second is the capacity of the service for all the clients together, not for one.
- Each step of C lasted 10 seconds, and each experiment of D and E ran once. The data was empty and there was one synthetic user, so the database did not have the size of a real one. I did not test the confirmation of a registration or the money path on the copy.
- The images were those already built for `csp-release` from the clones on `release/3.0.0`; I did not rebuild them or check the commit they were built from.
- A step lasts 2 seconds. A real load test runs for minutes to find leaks and slow growth, and it uses a tool such as k6 (not installed on this machine) that reports the rate it really sent.
- The integration test needs Docker. I ran it once, before knowing that another stack was running in Docker on this machine, and no container was left. I did not run it again, so the output above is from that one run.
- `ci.yml` was checked only as YAML (it parses, with the triggers `pull_request` and `schedule` and three jobs). I did not run `actionlint` on it, because that needs a container image and Docker was busy with another stack. It never ran on GitHub, where Actions is locked by billing.
- The contract of the rehearsal is a JSON file read by a test, not Pact. It shows the idea of a consumer-driven contract, not the broker and the verification of a provider in CI.
- The results of Section 1 come from the pull requests, on the commits they name, and they were not run again today. GitHub Actions is locked by billing, so none of them ran in CI.
- The money path of the project is the saga that holds and confirms seats; there is no payment step in it, and the payment of the rehearsal is the example of the session.

## 7. Self-check (answers)

1. **In the distributed pyramid you have...** many unit tests, fewer integration and contract tests, and few end-to-end tests, with load and chaos on top.
2. **Testcontainers are used to...** run a real database or broker in a container during the integration tests, and remove it afterwards. The rehearsal does it with PostgreSQL 16, and the database itself refuses a duplicate order and a status outside the list.
3. **Load testing exists to...** find the capacity and latency limits against the SLOs before the users do. The rehearsal found the limit at the rate where its design said, and on the real stack the limit of `csp-auth-api` was 17.6 logins per second, what its hash bound gives.
4. **Chaos engineering means...** injecting a failure to check that the system degrades gracefully. The rehearsal kills the provider and the system refuses fast and recovers; on the real stack, with `csp-auth-api` killed, the gateway refuses every request after 2 seconds and recovers by itself in under 1 second.
5. **Contract tests protect...** the compatibility between producers and consumers across services.
6. **"All unit tests green but dead at 200 req/s" is fixed by...** a load test and a chaos test before the release. A load test finds the limit, and a chaos test shows what happens when a dependency dies.

The page has no answer key; these answers come from its own explanations.
