# Week 13 - Session 2: Threat model, security backlog and security testing in CI (MVP 3)

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: model the threats of the system with STRIDE over its data flow, turn the findings into a prioritized security backlog with
testable criteria, and plan the security gates of the pipeline for MVP 3.

**Scope of this work.** This session is a plan, and it is **not applied to the project** by this document: no `csp-*` repository is changed
because of it. The work of the session was done in the week 13 workshop, so this page summarizes it and adds what the workshop did not
cover: the plan of the security gates for every repository of MVP 3. Every statement about the project was read in the repositories (branch `develop`).

## 1. Attack surface and trust boundaries

The data flow of the platform has seven trust boundaries (internet to the gateway, gateway to each API, services to their databases, Redis, the broker, the SMTP
relay and the service tokens of the worker). The map, with the nine assets and the entry points of my boundary, is in
[01-attack-surface.md](../workshop/01-attack-surface.md). Threats concentrate where data crosses a boundary, so I chose the one that every client crosses to
reach the domain I own: **TB2a, `csp-api-gateway` to `csp-auth-api`**.

## 2. STRIDE on the boundary

The full table (context, likelihood x impact, existing control and missing mitigation, ranking, residual risk and references) is in the team repository:
[stride-Oscarsl10.md](https://github.com/code-corhuila/csp-docs/blob/security/stride-Oscarsl10/05-architecture/security/stride-Oscarsl10.md), pull request
[csp-docs#278](https://github.com/code-corhuila/csp-docs/pull/278). The ranking, by risk and then by impact:

| Rank | STRIDE | Threat | Risk | Priority |
|---|---|---|---|---|
| 1 | S | Stolen credentials from many addresses, or with Redis down | 6 | P1 |
| 2 | D | Flood of `/login` and `/register` that exhausts bcrypt | 6 | P1 |
| 3 | I | Enumeration of emails through `409 EMAIL_ALREADY_REGISTERED` | 6 | P1 |
| 4 | R | Role and password changes without a lasting audit log | 4 | P2 |
| 5 | D | Mail bombing through `POST /forgot-password` (found in the review of a teammate) | 4 | P2 |
| 6 | E | A non-ADMIN or a revoked ADMIN reaches the role operations | 3 | P2 |
| 7 | T | `X-Forwarded-For` chosen by a caller that skips the gateway | 2 | P3 |

Repudiation is mitigated with an audit log that carries the `trace_id`, which is why the story of row S also asks to record every lock.

## 3. The security backlog

Three stories, one per risk of 6, with the labels `security` and `P1`, in the team board and linked to the pull request. Each criterion has an automated test.

| Story | Threat | What it asks | Test |
|---|---|---|---|
| [#280](https://github.com/code-corhuila/csp-docs/issues/280) | S | The 11th failed login from one address answers `429` even with the right password; every lock writes an audit event; a second counter per account slows the answers | Negative test of the login, unit and Redis integration tests |
| [#281](https://github.com/code-corhuila/csp-docs/issues/281) | D | Own stricter limit for `login`, `register` and `forgot-password` at the gateway, a bound on concurrent bcrypt hashes, and an alert on the latency | Smoke check of the gateway, test with a hasher that blocks |
| [#282](https://github.com/code-corhuila/csp-docs/issues/282) | I | `POST /register` answers the same for an existing and a new email, and the owner gets a notice | Contract test; it changes the contract, so it needs the agreement of the team |

The order is by effort and dependency: #280 first (the lock already exists), #281 second, #282 third, once the team approves the change of the contract.

## 4. Security testing in the pipeline: the plan for MVP 3

What the workflows run today, read in `.github/workflows/` of `csp-auth-api`, `csp-auth-portal`, `csp-api-gateway` and `csp-auth-db`: unit and integration
tests, and the contract and component tests of the portal. None of them runs a dependency scan, a secret scan or static analysis (I searched for `govulncheck`,
`gitleaks`, `npm audit`, `gosec`, CodeQL, `trivy` and `semgrep`). GitHub Actions is locked by billing since 2026-10-14, so nothing of this can be seen running yet.

| Gate | `csp-auth-api` (Go) | Portals and `csp-front` (npm) | `csp-api-gateway` (NGINX) | `-db` repositories (Flyway) |
|---|---|---|---|---|
| Dependency scan | `govulncheck ./...` | `npm audit --audit-level=high` | not applicable (image only) | not applicable |
| Secret scan | `gitleaks` on every pull request | `gitleaks` | `gitleaks` | `gitleaks` |
| Static analysis (SAST) | `gosec` | ESLint with a security rule set | `nginx -t` on the rendered configuration | SQL review in the pull request |
| Authorization tests | `go test -run Authz` (the 401, 403 and 429 cases) | contract tests against `auth-service.yaml` | `tests/smoke.sh`: `404` on `internal`, `401` without a token | not applicable |

Rules: a finding fails the build, like a failing unit test; an exception is written in the pull request with an owner and a date; the first run only reports,
so the existing findings are triaged before the gate blocks. The proof of concept for the first column is in
[04-security-ci](../workshop/04-security-ci/README.md) (workflow `security.yml` and the negative test of the login, with the real output in `HOW-TO-RUN.md`).

## 5. What this plan does not prove

- The gates have not run in GitHub; the commands of the first column were run by hand on my machine, and `govulncheck` reported a vulnerability of the standard library of the installed Go version (`GO-2026-6613`), not of a dependency.
- The plan is written for the domain I own and for the shared repositories I read; the other domains have their own threat models in the pull requests of my teammates.
- Choosing `gosec`, the ESLint rule set and where the exceptions are recorded needs the agreement of the team and an ADR before anything is installed. No story is created for it: the team does not open stories without approval.

## 6. Self-check (answers)

1. **Threats concentrate at...** the trust boundaries, where data crosses trust levels.
2. **STRIDE's "E" stands for...** Elevation of privilege.
3. **A security finding should become...** a prioritized story with a testable acceptance criterion; #280, #281 and #282 are examples.
4. **Repudiation is mitigated by...** audit logs with a `trace_id`.
5. **Security testing in CI includes...** dependency and secret scanning, static analysis and authorization tests, and they fail the build.
6. **"We'll review security before launch" tends to...** surface expensive, cross-cutting fixes too late, as the review of a late IDOR in the page shows.
