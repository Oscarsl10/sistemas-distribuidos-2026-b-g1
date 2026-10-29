# How to run the negative authorization test (CP4)

**Story covered:** the story of threat S of my STRIDE analysis, [csp-docs pull request](https://github.com/code-corhuila/csp-docs/pulls)
`Threat model: csp-api-gateway to csp-auth-api` (claim issue #275): *failed logins must be limited*. **Acceptance criterion:** the 11th login from one
client address is refused with `429 RATE_LIMITED` and `Retry-After`, even when the credentials are valid (ADR-040 of `csp-docs`).

The test is in `csp-auth-api/login_authz_test.go`. It is named after the real endpoint (`POST /api/v1/auth/login`) and runs against a minimal model of
the login of `csp-auth-api` (`login.go` and `limiter.go`): the limit is kept in memory here, in the real service it is kept in Redis. It does not change the
team code repository.

```bash
cd 13-week/workshop/04-security-ci/csp-auth-api
go vet ./...
go test -run Authz -v ./...
```

Real output of my machine (Go 1.27.0, Windows 11, 2026-10-28):

```
=== RUN   TestAuthzLoginRejectsTheEleventhFailedAttemptFromTheSameAddressEvenWithValidCredentials
--- PASS: TestAuthzLoginRejectsTheEleventhFailedAttemptFromTheSameAddressEvenWithValidCredentials (0.00s)
=== RUN   TestAuthzLockOfOneAddressDoesNotLockAnotherAddress
--- PASS: TestAuthzLockOfOneAddressDoesNotLockAnotherAddress (0.00s)
=== RUN   TestAuthzAddressWrittenByTheClientIsNotUsedToEscapeTheLock
--- PASS: TestAuthzAddressWrittenByTheClientIsNotUsedToEscapeTheLock (0.00s)
=== RUN   TestAuthzLockEndsAfterTheCooldown
--- PASS: TestAuthzLockEndsAfterTheCooldown (0.00s)
PASS
ok  	csp-auth-api-security	1.386s
```

| Test | What it proves |
|---|---|
| `TestAuthzLoginRejectsTheEleventhFailedAttemptFromTheSameAddressEvenWithValidCredentials` | Ten wrong passwords answer `401`; the 11th attempt with the valid password answers `429`, `Retry-After: 300` and `RATE_LIMITED`. This is the criterion of the story. |
| `TestAuthzLockOfOneAddressDoesNotLockAnotherAddress` | The lock is per address. |
| `TestAuthzAddressWrittenByTheClientIsNotUsedToEscapeTheLock` | A first entry written by the client in `X-Forwarded-For` does not free the locked address; only the last entry, the one the gateway appends, is counted (ADR-040 D8). |
| `TestAuthzLockEndsAfterTheCooldown` | After 5 minutes the same address can log in again. |

## The other two gates, run on my machine

The workflow `security.yml` runs three jobs on every pull request. GitHub Actions of the organization is locked by billing, so the workflow could not be seen running
in GitHub; what I ran by hand is the same command of each job.

- **Workflow syntax:** `actionlint` (image `rhysd/actionlint:1.7.7`) reports no problem for `security.yml`.
- **Secret scan:** `gitleaks dir` (image `zricethezav/gitleaks:v8.21.2`) on the workshop folder: `no leaks found`.
- **Dependency audit:** `govulncheck ./...` (v1.1.4) finds no vulnerable dependency, because the module has none, but it **does report one vulnerability of the Go standard library**
  (`GO-2026-6613`, `net/http`, found in go1.27.0 and fixed in go1.27.2). That is the toolchain installed on my machine, not a dependency of the module; the job would be red on this toolchain,
  which is how the gate is meant to work. The CI job uses `go-version: "1.24"`, to be checked when the Actions can run.

## Run against the real service (CP8)

The test above runs against a minimal model of the login. The real `csp-auth-api` already has the same criterion as tests of its use case and of its HTTP adapter, so on 2026-10-29
I ran them there, without changing that repository (commit 17a2fb8 (2026-10-24) of `develop`, clean working tree). Toolchain go1.27.0.

```bash
cd csp-auth-api
go test ./internal/application/usecase -count=1 -v -run 'TestLoginRejectsTheAttemptAfterTenFailuresEvenWithTheRightPassword|TestLoginWorksAgainOnceTheLockHasElapsed|TestLoginGoesOnWhenTheAttemptStoreFails'
```

```text
=== RUN   TestLoginRejectsTheAttemptAfterTenFailuresEvenWithTheRightPassword
--- PASS: TestLoginRejectsTheAttemptAfterTenFailuresEvenWithTheRightPassword (0.00s)
=== RUN   TestLoginWorksAgainOnceTheLockHasElapsed
--- PASS: TestLoginWorksAgainOnceTheLockHasElapsed (0.00s)
=== RUN   TestLoginGoesOnWhenTheAttemptStoreFails
2026/10/29 12:00:25 WARN login attempt limit unavailable, the login goes on error="redis down"
2026/10/29 12:00:25 WARN failed logins could not be cleared error="redis down"
2026/10/29 12:00:25 WARN login attempt limit unavailable, the login goes on error="redis down"
2026/10/29 12:00:25 WARN failed login could not be counted error="redis down"
--- PASS: TestLoginGoesOnWhenTheAttemptStoreFails (0.10s)
PASS
ok  	github.com/code-corhuila/csp-auth-api/internal/application/usecase	1.379s
```

```bash
go test ./internal/adapter/in/httpapi -count=1 -v -run 'TestLoginRateLimitedAnswers429WithRetryAfterAndTheSharedEnvelope'
```

```text
=== RUN   TestLoginRateLimitedAnswers429WithRetryAfterAndTheSharedEnvelope
--- PASS: TestLoginRateLimitedAnswers429WithRetryAfterAndTheSharedEnvelope (0.00s)
PASS
ok  	github.com/code-corhuila/csp-auth-api/internal/adapter/in/httpapi	2.388s
```

What it shows:

- The attempt after ten failures is rejected even with the right password, and the login works again once the lock has elapsed. This is the criterion of story #280 in the real code.
- The HTTP adapter answers `429` with `Retry-After` and the shared error envelope.
- When the store of attempts fails the login goes on (`redis down` in the log). This is the "fails open" behaviour of row S of the STRIDE table, now seen in the real code.

What it does not show: the end-to-end test against PostgreSQL and Redis (`cmd/auth-api/login_attempts_integration_test.go`) was not run, and nothing went through the gateway. The names of these tests do not contain `Authz`,
so the job `authz-tests` of `security.yml` (`go test -run Authz`) would not select them in that repository; to adopt it there the pattern or the names must change.

## Validation of `security.yml` (CP8)

`security.yml` parses as YAML on 2026-10-29: trigger on `pull_request`, jobs `secret-scan`, `dependency-audit`, `authz-tests`. `actionlint` was not run again today because Docker was off; its result of the first run (no problem found) stands for the same file.
