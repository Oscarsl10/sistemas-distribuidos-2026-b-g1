# Week 15 - Session 1: Testing at scale: contract, e2e, load and chaos

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: apply the testing pyramid to a distributed system, use real dependencies in the integration tests, keep contract tests between the
services, cover the money paths end to end, run a load test against the SLO and prove the resilience with one chaos experiment.

**Scope of this work.** Nothing is applied to the `csp-*` repositories from this session. The project already has part of these tests, and Section 1 reports
their results as they are written in the pull requests of the repositories: **I did not run them again today**, except the 16 tests of the Redis circuit breaker and the 7 of the bulkhead of Section 2. Each row says on which commit it was run.
The two pieces the project does not have, a load test and a chaos experiment, are done as a rehearsal in [rehearsal/](rehearsal/), with a stand-in of a checkout
and of a payment provider. Section 4 has the real output of the rehearsal.

## 1. The pyramid for a distributed system, and what the project has

| Layer | What exists in the project | Result (source) |
|---|---|---|
| Unit | `go test ./...` in `csp-auth-api`; `mvn verify` in `csp-workflow`; `ng test` in `csp-auth-portal` | `csp-auth-api` release head: 1196 tests passed, contract tests included. `csp-auth-portal`: 289 of 289. `csp-workflow` on `7e160ee`: 118 tests, 0 failures ([csp-docs#282](https://github.com/code-corhuila/csp-docs/issues/282) closing comment, 2026-11-04; [csp-workflow#13](https://github.com/code-corhuila/csp-workflow/pull/13)) |
| Integration with real dependencies | `go test -tags integration ./...` in `csp-auth-api` against PostgreSQL and Redis; `JdbcSagaStoreTest` of `csp-workflow` against PostgreSQL (`TEST_DATABASE_URL`); `db-ci.yml` of each `-db` repository (migrate, test, rollback) | `csp-auth-api`: 1449 passed, 1 skipped (file modes on Windows), 0 failed. `csp-auth-db`: the 11 steps of `db-ci.yml` passed, 45 migrations with rollbacks. A full-stack run with a real Mailpit and PostgreSQL 16 passed 19 of 19 checks (same closing comment). `csp-workflow#12`: 104 tests, 8 of them `JdbcSagaStoreTest` against the real database |
| Contract | The responses of `csp-auth-api` are validated against `auth-service.yaml`; the portal checks its own requests against it (`npm run test:contract`) | The portal check matches 14 operations of `auth-service.yaml` 1.17.0 (closing comment of #282) |
| Gateway | `tests/smoke.sh` and `tests/integration.sh` of `csp-api-gateway` | Release head: 97 smoke checks ok, 0 failures; 307 integration checks passed, 0 failed (closing comment of #282) |
| End to end (money path) | `tests/e2e/up.sh` and `tests/e2e/book-seats.sh` of `csp-workflow`: the saga that holds and confirms the seats, against the real `csp-booking-api` and `csp-catalog-api` through the gateway | `TOTAL passed=89 failed=0`, twice, on `c635586` ([csp-workflow#15](https://github.com/code-corhuila/csp-workflow/pull/15)); `TOTAL passed=105 failed=0`, twice, on `76256cb` ([csp-workflow#17](https://github.com/code-corhuila/csp-workflow/pull/17)) |
| Chaos (partial) | The same end-to-end script kills the container of the workflow in the middle of a saga and makes the compensation fail on purpose; the Redis circuit breaker and the bulkhead of `csp-auth-api` have unit tests | See Section 2: 16 breaker tests and 7 bulkhead tests passed today on `a998029` |
| Load | Nothing | Section 4 (rehearsal) |

Where the real dependencies come from: the project does not use the Testcontainers library. Real PostgreSQL and Redis come from containers that the scripts or the
workflow start, or from the persistent development environment. I searched the repositories for `testcontainers`, `k6`, `gatling`, `jmeter` and `chaos` in the
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
`csp-auth-api`. The breaker that protects a call to another service, which is what the second half of the session asks for, is only in the rehearsal.

## 3. The rehearsal

A small Go module with no dependencies, in [rehearsal/](rehearsal/). It does not touch any `csp-*` repository.

| File | What it is |
|---|---|
| `payments.go` | The payment provider: a real HTTP server on a real port. It takes 50 ms per payment and serves 20 at the same time (a capacity of 400 payments per second). It can be killed and started again on the same address |
| `checkout.go` | The money path: reserve a seat, pay, confirm; if the payment fails, release the seat (the compensation of a saga). A circuit breaker opens after 5 failures in a row, refuses for 500 ms and then lets one probe through |
| `contract/payments.contract.json` | What the checkout expects from the provider: `200`, `application/json`, the fields `id` and `status`, and the status `APPROVED` |
| `scale_test.go` | The four tests below |

The tests:

- **Contract.** The checkout (consumer) writes what it expects in the JSON file, and the test checks that the provider answers exactly that.
- **End to end.** One order through the real HTTP entry and the real provider: the status is `201`, one seat is confirmed and none stays held.
- **Load.** An open-loop generator (it keeps the rate whatever the answers) sends orders for 2 seconds at 100, 200, 300, 400 and 600 requests per second. The SLO of the rehearsal is p95 at most 250 ms and every order confirmed, with a target of 200 requests per second. The test must fail if the target is not met, and reports the first rate that breaks the SLO.
- **Chaos.** At 100 orders per second for 5 seconds, the provider is killed at second 1 and started again at second 3. The test checks that the orders in flight are compensated, that the breaker opens and answers fast (`503`, p95 at most 50 ms), that the system recovers by itself, and that no seat stays held and the confirmed seats are exactly the orders answered `201`.

## 4. Real output

Go 1.27.0 on Windows 11, 2026-11-11. Command: `cd 15-week/01-session/rehearsal && go test -count=1 -v ./...`.

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

**The tests catch faults.** I changed a copy of the rehearsal twice, and the chaos test failed both times, as it must:

- Without the release of the seat after a failed payment, `held` stayed at 9: `seats: free 999705, held 9, confirmed 286, orders answered 201: 286`.
- With a breaker that never opens (threshold 1,000,000), `during the outage: compensated 189 (at most 30), rejected 0 (at least 100), rejected p95 0s (at most 50ms)`: every order waited for its own failure instead of being refused fast.

## 5. What this does not prove

- The rehearsal is one machine, with the checkout, the provider and the generator in the same process on the loopback interface. Real services add network, a database and other processes, so the numbers say how the design behaves, not how fast CineSync Platform is.
- The SLO of the rehearsal (p95 250 ms, a target of 200 per second) is my assumption to have a number; it is not an SLO of the project. The project has no load figure: nobody has run a load test on its services.
- A step lasts 2 seconds. A real load test runs for minutes to find leaks and slow growth, and it uses a tool such as k6 (not installed on this machine) that reports the rate it really sent.
- The contract of the rehearsal is a JSON file read by a test, not Pact. It shows the idea of a consumer-driven contract, not the broker and the verification of a provider in CI.
- The results of Section 1 come from the pull requests, on the commits they name, and they were not run again today. GitHub Actions is locked by billing, so none of them ran in CI.
- The money path of the project is the saga that holds and confirms seats; there is no payment step in it, and the payment of the rehearsal is the example of the session.

## 6. Self-check (answers)

1. **In the distributed pyramid you have...** many unit tests, fewer integration and contract tests, and few end-to-end tests, with load and chaos on top.
2. **Testcontainers are used to...** run a real database or broker in a container during the integration tests. The project gets the same result with containers started by its scripts.
3. **Load testing exists to...** find the capacity and latency limits against the SLOs before the users do. The rehearsal found the limit at the rate where its design said.
4. **Chaos engineering means...** injecting a failure to check that the system degrades gracefully. The rehearsal kills the provider and the system refuses fast and recovers.
5. **Contract tests protect...** the compatibility between producers and consumers across services.
6. **"All unit tests green but dead at 200 req/s" is fixed by...** a load test and a chaos test before the release. A load test finds the limit, and a chaos test shows what happens when a dependency dies.

The page has no answer key; these answers come from its own explanations.
