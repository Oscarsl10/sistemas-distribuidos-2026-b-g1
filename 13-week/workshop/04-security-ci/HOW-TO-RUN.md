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
