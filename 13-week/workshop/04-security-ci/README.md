# Security gates in the pipeline (CP4)

Proposal for the pull request pipeline of `csp-auth-api` (Go). It lives in my fork; it is not installed in the team repositories.

| Gate | Job of `security.yml` | Tool | Blocks the merge when |
|---|---|---|---|
| Secret scanning | `secret-scan` | gitleaks | a key, token or password is found in the history of the pull request |
| Dependency audit | `dependency-audit` | `govulncheck` | the code calls a known vulnerability of a dependency or of the standard library |
| Negative authorization tests | `authz-tests` | `go test -run Authz` | a request that must be refused (for example the 11th failed login of one address) is answered with success |

The three jobs run on `pull_request`; a red step blocks the pull request and the security story stays open.

**Status of the evidence.** The workflow was not seen running in GitHub because the Actions of the organization are locked by billing since 2026-10-14. The
YAML is valid (`actionlint`), and the command of each job was run on my machine, see [HOW-TO-RUN.md](HOW-TO-RUN.md). The URL of a real run will be added here when the Actions work again.
