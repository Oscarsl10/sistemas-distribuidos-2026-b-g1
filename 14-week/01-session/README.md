# Week 14 - Session 1: CI/CD, pipelines and quality gates

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: tell CI, Continuous Delivery and Continuous Deployment apart, design a pipeline whose quality gates block the merge,
protect `main` with PR checks, and measure the delivery health with the DORA metrics.

**Scope of this work.** It is a rehearsal, not a change to the project. Nothing is applied to the `csp-*` repositories from this session. The
rehearsal is in [rehearsal/](rehearsal/): a small Go module with one function, its tests, and a script that runs the gates in the order of the
pipeline of the session. Section 6 reads the workflows of the project and compares them with the gates; it changes nothing.

## 1. CI, Continuous Delivery and Continuous Deployment

| Term | What is automatic | What a person still does |
|---|---|---|
| Continuous Integration | Every change is built and tested when it is merged, in small batches | Merges the change |
| Continuous Delivery | Every green build is releasable, and one click deploys it | Decides when to deploy |
| Continuous Deployment | Every green build reaches production by itself | Nothing, once the change is merged |

All three share the same pipeline. The difference is how far the automation goes.

## 2. The gates of the rehearsal

Each gate is a command with an exit code. A red gate stops the run, so the change cannot reach the next step.

| Gate | Command in `gates.sh` | What it blocks |
|---|---|---|
| Lint | `gofmt -l .` must print nothing | Code that is not formatted |
| Typecheck | `go vet ./...` and `go build ./...` | Code that does not compile or has vet errors |
| Test | `go test -v ./...` | A failing unit test |
| Coverage check | the lowest coverage of the packages must be at least the floor in `coverage.floor` (declared 80, measured 100) | Coverage below the declared floor |
| No placeholders | `no-todo.sh` fails on `TODO`, `FIXME` or `PLACEHOLDER` in `*.go` | Code that still has a marker |

The page names two more gates, dependency and secret scanning (and SAST). They are not in this rehearsal: they need tools that are not installed here, and Section 6 says what the project runs instead of them.

`make` is not installed on this machine, so the steps are in a shell script with the same names and the same order as the pipeline snippet of the page.

## 3. Real output

Go 1.27.0 on Windows 11, 2026-11-07. Command: `cd 14-week/01-session/rehearsal && ./gates.sh`.

**Run 1: all gates green.**

```text
== lint: formatting
lint: OK, gofmt reports 0 files
== typecheck: vet and compile
typecheck: OK
== test: unit tests
--- PASS: TestDivideReturnsTheQuotient (0.00s)
--- PASS: TestDivideRejectsZeroDivisor (0.00s)
test: OK, 2 PASS lines, 0 FAIL lines
== coverage-check: measured must be >= declared floor
coverage-check: OK, measured 100.0% >= floor 80%
== no-todo: no placeholders left in the code (the dependency, secret and SAST scans are not part of this rehearsal)
no-todo: OK, no placeholder markers in *.go files
== all gates green: the change is releasable
EXIT=0
```

**Run 2: a `TODO` left in the code (done in a copy, the rehearsal is not changed).** The gate `no-todo` stops the run.

```text
== no-todo: no placeholders left in the code (the dependency, secret and SAST scans are not part of this rehearsal)
./calc/calc.go:8:// TODO: handle negative numbers
no-todo: FAIL, placeholder markers found above
EXIT=1
```

**Run 3: the test of the zero divisor removed (done in a copy).** The tests still pass, but coverage falls under the floor and the gate `coverage-check` stops the run.

```text
== test: unit tests
--- PASS: TestDivideReturnsTheQuotient (0.00s)
test: OK, 1 PASS lines, 0 FAIL lines
== coverage-check: measured must be >= declared floor
coverage-check: FAIL, measured 66.7% < floor 80%
EXIT=1
```

Two more red runs happened while I built the copies, and both were correct stops: a left-over blank line made `gofmt` fail, and a removed test left an unused import that `go vet` rejects. I kept them out of the list because the three runs above are the cases that matter.

## 4. What the rehearsal does not prove

- It is one small package. A real service has many packages and its gates take longer, and coverage depends on which code the tests reach.
- The dependency scan (`govulncheck`), the secret scan (`gitleaks`) and the static analysis are not run here. Week 13 ran `govulncheck` on a module and found a vulnerability of the standard library of the installed Go toolchain (`GO-2026-6613`), so those gates matter.
- The gates run on my machine. Nothing here shows that GitHub runs them on a pull request, because GitHub Actions is locked by billing since 2026-10-14.
- A green run proves the change passes the gates, not that it is correct. The tests are the only check of behaviour.

## 5. DORA metrics from this week (a proxy, not a measurement of the team)

DORA has four metrics: deployment frequency, lead time for changes, change failure rate and time to restore service. The data I have is the merged pull requests of my account in the organization from 2026-11-02 to 2026-11-07, read with `gh`. The table uses it as a proxy.

| Metric | Proxy used here | Value this week | What it really says |
|---|---|---|---|
| Deployment frequency | Merges of my pull requests per day | 144 merges on 4 days | It counts merges, not deployments; 16 repositories were touched |
| Lead time for changes | Time from pull request creation to merge | median 0.1 hours, 90th percentile 0.4 hours | Very short, but not a sign of speed: the merges happen with no check and no review (see Section 6) |
| Change failure rate | Share of merged titles that start with `fix` or `revert` | 22 of 144, about 15 % | An upper estimate: a `fix` title is not always a production failure |
| Time to restore | Not measured | not available | The data has no incident or deployment time |

The lead time is the number to read with care. A median of six minutes from creation to merge means the pull requests were merged right after they were opened, with no pipeline and no reviewer in between. A fast lead time without gates is the same risk the session describes: a change reaches the main branch before anyone checks it.

## 6. Where the project stands against these gates

Read in the workflow files (`.github/workflows/`) of the repositories, without changing anything:

| Repository | Jobs that run on a pull request today | Gate missing against the page |
|---|---|---|
| `csp-auth-api` | `build-and-test`: format check (`gofmt -l`), `go vet`, `go build`, `go test`; an `integration-tests` job with a database and Redis | No coverage floor, no dependency or secret scan, no SAST |
| `csp-auth-portal` | Install, contract check against `auth-service.yaml`, build, lint, test, container image and CORS check | No coverage floor, no dependency or secret scan |
| `csp-api-gateway` | Build, configuration check, smoke checks and integration checks | No coverage floor, no lint, no scan |
| `csp-workflow` | Build and test, migrations, compose render, image build and health check | No coverage floor, no lint, no scan |
| `csp-ticketing-api` | Build and test, and a board tracking job (`env-tracking.yml`) | No coverage floor, no scan |
| `csp-ticketing-db` | `migrate`: Flyway migration, tests and rollback (`db-ci.yml`) | No coverage floor, no scan |
| `csp-booking-api`, `csp-worker` | `build-and-test` | No coverage floor, no scan |

I searched the workflows of ten repositories for `govulncheck`, `gitleaks`, `npm audit`, `gosec`, CodeQL, `trivy` and `semgrep`, and for a coverage floor: none of them is there. The only lint gate is in the auth portal.

Two more facts:

- **No path filters.** The five main workflows have no `paths:` filter, so every pull request runs the whole pipeline of its repository.
- **Branch protection is not available.** On the Free plan of the organization, a private repository cannot enable protected branches. I checked it on `csp-docs`: the API for the branch protection of `main` answers `403` with "Upgrade to GitHub Pro or make this repository public". So required checks and required reviews cannot be enforced by GitHub today, and the rules of the team depend on each person.

## 7. Self-check (answers)

1. **Continuous Integration means...** frequent merges that are built and tested automatically, in small batches.
2. **A quality gate in the pipeline should...** block the merge or the deploy when it fails. A gate that only prints a warning is ignored.
3. **To keep main always releasable you...** protect it: required checks and review, and no direct pushes.
4. **Continuous Deployment means...** every green build goes to production by itself. Delivery stops one step earlier: every green build is only releasable.
5. **DORA metrics include...** deployment frequency, lead time for changes, change failure rate and time to restore service.
6. **"Broken main blocks the team" is prevented by...** required CI gates on the pull requests, so a red build cannot merge.

The page has no answer key; these answers come from its own explanations.
