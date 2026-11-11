# Week 15 - Session 2: Release, MVP 3 operability checklist, demo under stress and retrospective

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: promote and tag the production-grade release, verify it against the operability checklist, demo it under load and failure,
and hold a retrospective that feeds the final presentation.

**Scope of this work.** It is a rehearsal and an audit. **Nothing is applied to the project**: no `csp-*` repository was changed, no tag was created,
no pull request was touched. The release itself is not mine to cut: the release pull requests wait for the professor and the leader (the team rule),
so what I can do honestly is (1) rehearse the mechanics in a throwaway repository, (2) read the state of the checklist in the repositories and on GitHub,
and (3) reuse the results of [session 1](../01-session/README.md), which ran load and chaos on a copy of the release stack.
Every statement about the project was read from GitHub or from the repositories (branch `develop`) on 2026-11-11.

## 1. Where the release is today

The page asks to promote `develop -> qa -> main` and tag `v3.0.0`. In this project the cut is `release/3.0.0`, one pull request per repository to `main` (norm 6.3.2).

| Repository | Release PR | State | Review | Size (added / deleted) |
|---|---|---|---|---|
| `csp-auth-api` | #83 | open | approved | 36,887 / 72 |
| `csp-auth-db` | #33 | open | approved | 595 / 3 |
| `csp-auth-portal` | #49 | open | approved | 4,767 / 254 |
| `csp-booking-api` | #61 | open | approved | 3,539 / 126 |
| `csp-workflow` | #26 | open | approved | 5,252 / 4 |
| `csp-api-gateway` | #6 | open | approved | 1,639 / 5 |
| `csp-ticketing-api` | #84 | open | review required, blocked | 9,129 / 6 |

- Six of the seven are approved; the ticketing one still needs its review. The merge state of the six is `UNSTABLE`, which is what the billing lock of GitHub Actions
  (since 2026-10-14) produces: no check runs, so none can be green.
- **No repository has a `v3` tag.** I checked the seven with the tags API. The last release is `v2.0.0`. So the "tag `v3.0.0`" step is not done, and the order of the
  project is: merge `csp-auth-db`, then the APIs, then the portal, and tag after the merge (a tag is never moved).
- The size column is raw lines of the whole release, not the 400-line cap of a single PR: a release is a promotion of many stories.

## 2. The rehearsal: cut, gate, tag, deploy and roll back

[rehearsal/release.sh](rehearsal/release.sh) does the mechanics in a temporary directory with a local bare remote, and deletes it at the end.

1. Cut `release/3.0.0` from `main` and apply one commit per story with `git cherry-pick -x`.
2. A red gate must block the tag; after the fix the gate is green.
3. Merge to `main` with a merge commit, set the version, create the annotated tag `v3.0.0`, push it.
4. Deploy the tag, then deploy the previous tag as the rollback, with no rebuild and no new commit.

```text
1. cut release/3.0.0 from main and apply one commit per story
  ok   three story commits on the release branch
  ok   every commit carries the cherry picked from trail
2. the gate decides
  ok   a red gate blocks the tag
  ok   no tag exists after the red gate
  ok   the gate is green again
3. merge to main, tag and publish
  ok   v3.0.0 is an annotated tag
  ok   the remote has the tag
  ok   a published tag is not moved (a second tag with the same name)
  ok   main contains the three stories
4. deploy the tag and roll back to the previous one (the same artifact, no rebuild)
  ok   the deployed version is 3.0.0
  ok   the rollback deployed 0.1.0 (the content of v2.0.0)
  ok   the rollback did not create a commit

checks passed: 12, failed: 0
```

What it proves: the order of the steps and the three rules that make a release reproducible (the gate blocks the tag, the tag is annotated and not moved,
the rollback is a previous artifact). What it does not prove: a real deployment. The "deploy" here copies a file; the project deploys with Docker Compose,
and the rebuild-versus-promote gap found in week 14 is still open.

## 3. The operability checklist against the project

Each item of the page, with the evidence I found and what is missing. "Not found" means I looked in the places named and saw nothing; it does not mean it cannot exist elsewhere.

| Item of the page | Evidence | Gap |
|---|---|---|
| Acceptance criteria met, CI green | The release PRs are approved by review. The CI cannot run (billing lock), so the pipelines were validated locally by the team before each merge. | **No green CI exists for the release.** The page requires it. |
| Coverage and security scans in CI | In week 14 I read ten workflows: no coverage floor, no secret scan, no dependency scan. Not re-read today. | Both scans still missing, as of week 14 |
| Zero-downtime deploy | `csp-docs` (`05-architecture/deployment.md`) describes a rolling deploy and `kubectl rollout undo`. The running stack is Docker Compose. | The documented strategy is not the one that runs. Not rehearsed on the stack. |
| Rehearsed rollback | Section 2 rehearses it on a throwaway repository. | Not rehearsed on the real stack or with the real images |
| Dashboards with the golden signals | Prometheus and Grafana run in the release stack; `csp-infra-postgres` holds `prometheus.yml` and `grafana-datasources.yaml` (the only datasource is Prometheus). | I did not open a dashboard or find a dashboard definition in the repository. Not verified. |
| A real alert with a runbook | `csp-infra-postgres/observability/alerts/` has `auth-latency.yml` (p95 of login and register above the 1 s SLO of ADR-061, for 15 minutes, with a guard of at least 10 requests) and `worker-outbox-relay.yml` (four alerts). Each has a `.test.yml` and a `runbook_url`. The auth runbook has a section for the latency alert. | Strong item. Not seen firing on the stack. |
| Distributed tracing | `otel-collector.yaml` receives OTLP and has a `traces` pipeline. | **Its only trace exporter is `debug`**: spans are printed by the collector and not stored. There is no Jaeger or Tempo, so a trace cannot be looked at. |
| Authentication and authorization per service | JWT RS256 validated per service and the service audience of ADR-056 (read in session 1 and week 13). | Not re-tested for every service today |
| No secrets in git | The release stack loads secrets from git-ignored `env/` and `keys/`. | No scan proves it (see the second row) |
| Load test meets the SLO | Session 1, on a copy of the release stack: p95 of the logins that were served stayed at most 362 ms at every rate, below the 1 s SLO. The service served about 17.6 logins per second and refused the rest with `503`. | One service, one machine, 10 s steps |
| Chaos shows graceful degradation | Session 1: with `csp-auth-api` killed the gateway answered `503` after 2 s and recovered by itself in under 1 s; with Redis down the lock failed open. The workflow service has e2e tests of saga compensation and crash recovery. | Degradation exists but is not fast (no circuit breaker at the gateway) and one case fails open |
| Consistency under failure (outbox, idempotency) | Outbox relay with its metrics and alerts; `Idempotency-Key` required on register and confirm (a 400 without it, seen in session 1). | Not re-tested under failure today |
| Versioned `v3.0.0`, CHANGELOG, ADRs, runbooks | CHANGELOG, ADRs and a runbook per service exist in `csp-docs`. | **The tag does not exist** (Section 1) |

The score, counted by me: of the 13 rows, 3 have solid evidence that I saw (alert with runbook, load, chaos with the saga), and the others are partial, unproven or missing.
The page warns that the release is judged on operability; the honest reading is that the project has the pieces and lacks three proofs:
a green pipeline, a trace backend, and a rollback on the real stack.

## 4. The demo under stress, as a script

The page asks to drive the journey, open dashboards and traces, load lightly, kill a service, and show alert, degradation and recovery.
This is the order I would use, with what each step shows and what I already have as evidence.

| Step | Action | What the audience sees | Evidence I have |
|---|---|---|---|
| 1 | Register, confirm from Mailpit, login through the gateway | The journey works end to end | Session 1, `real-stack/seed.py` |
| 2 | Open Grafana and Prometheus | The auth latency histogram and the outbox metrics | Datasource only; dashboards not verified |
| 3 | Load at 5, 10 and 15 logins per second | p95 stays under 1 s | Session 1 table C: p95 at most about 255 ms up to 15 per second |
| 4 | Load above capacity (30 per second) | The service refuses with `503` and `Retry-After`, the served logins keep their latency | Session 1 table C |
| 5 | Kill `csp-auth-api` | `503` for 2 s each, then automatic recovery | Session 1 section E |
| 6 | Show the alert | `AuthCredentialsLatencyP95AboveSlo` needs 15 minutes above the SLO, so it **cannot fire live in a demo** | I would show its unit test (`auth-latency.test.yml`) and the runbook instead |
| 7 | Roll back | Redeploy the previous tag | Section 2, only on a throwaway repository |

Two honest limits of the demo: the alert cannot be shown firing in a short slot (its `for` is 15 minutes, chosen to avoid paging on noise), and the traces cannot be
shown at all until a backend stores them.

## 5. Retrospective toward the final presentation

**What made the system reliable.**
- Contracts and an ADR per decision: a limit such as the credential bucket of the gateway (ADR-060) or the latency SLO (ADR-061) is written, so a test can compare it with the running system. In session 1 every number I measured matched the configuration I had read.
- Bounded resources: the hash bulkhead refuses with `503` instead of queueing, and the Redis breaker keeps the logins answering when Redis fails.
- Compensation and recovery in the workflow service, with tests that kill it in the middle of a saga.
- Alerts that have a test and a runbook, and an SLO of the real traffic rather than a guess.

**What I would do differently.**
- Install the pipeline gates that cost one job each (coverage floor, dependency scan, secret scan) before the features, not after. The billing lock showed the other cost: every control that depends on a hosted service stops when the service stops, so the local validation steps must be scripted from the start.
- Make the tag part of the pipeline: a release that is merged but not tagged is not reproducible.
- Store traces: a collector that only prints them is the same as no tracing.
- Fail fast at the gateway (a breaker or a shorter connect timeout) instead of 2 seconds per request, and decide the lock fail-open case in an ADR.

**What the metrics show.** I have these figures, from my own measurements and from GitHub, and no others:
- SLO of login and register: p95 at most 1 s. Measured: at most 362 ms for served requests; 17.6 logins per second of capacity (session 1).
- Recovery after killing a service: first `200` 0.9 s after the restart.
- Delivery: 138 commits in the week 14 evidence (`14-week/hu-status/commits.md`). I do not have the lead time or the deployment frequency of the release: **the release has not been deployed to a tagged `main`**, so those DORA figures do not exist yet, and I will not estimate them.

## 6. What this does not prove

- That the release is production-grade. The checklist above says three proofs are missing, and the tag does not exist.
- That the checks of Section 3 are current for every repository: scans and workflows are as I read them in week 14.
- That a rollback works on the real stack. Only the mechanics were rehearsed.
- The demo was not run live for this session: it is a script built on the results of session 1.

## 7. Self-check (answers)

1. MVP 3 must add operability: observability, zero-downtime deploy, security and proven resilience. Not only features.
2. The final release is tagged `v3.0.0` on `main`, through the gated pipeline. In this project it is still pending: the release PRs are open.
3. A production-grade demo shows live metrics and traces, load within the SLO and a graceful failure. Not only the happy path.
4. Resilience is credible only if it was load- and chaos-tested rather than assumed (session 1 did both on a copy of the release stack).
5. The checklist includes a rehearsed rollback, an alert with a runbook, security and resilience.
6. Corte 3 also includes the certification and final project components; the `hu-status` of each member feeds the personal grade.
