# Week 14 - Session 2: Planning the pipeline, the Definition of Done as gates and branch protection (MVP 3)

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: map each item of the Definition of Done (DoD) to an automated gate, define the branch protection of each environment,
design the pipeline for a system of many services, and slice the CI/CD work into stories for MVP 3.

**Scope of this work.** This is a plan. It is **not applied to the project**: no `csp-*` repository is changed, and the stories of section 4 are
drafts that are not created on the board (the team does not open stories without approval). Every statement about the project was read in
the repositories (branch `develop`) and in the workflow files; the checks of the first session are in
[01-session](../01-session/README.md).

## 1. The DoD as a checklist of machines

A DoD item with no gate that can reject it is self-certified. The table maps each item of the page to a gate and says what the project runs today.

| DoD item | Gate that proves it | What the project runs today | Gap |
|---|---|---|---|
| Code style is consistent | Linter, blocking | Go: `gofmt -l` in `csp-auth-api` only. Portal: `lint`. Gateway, workflow and ticketing have none. | Go repositories without a linter; a format check in the others |
| Types are correct | Strict type check | Go: `go vet` and `go build`. Portal: the Angular build. | No type check beyond the compiler in the Go repositories |
| Behaviour is tested | Unit, integration and contract tests | Go: `go test`, with an integration job in `csp-auth-api`. Portal: `test` and the contract check against `auth-service.yaml`. | Not every repository has an integration job |
| Enough coverage | Coverage check, declared at most measured | Nothing | No coverage floor in any of the ten workflows I read |
| No secrets and no known CVEs | Secret scan and dependency scan | Nothing in the workflows (week 13 wrote a proposal in `security.yml`, not installed in any repository) | Both scans are missing |
| The system runs for real | Smoke test: compose up and `/health` | Gateway: `smoke.sh` and integration checks. Workflow: the image answers its health check. Infra: the compose files render. | No smoke test of the whole platform after deploy |

Two rows are the most useful to start with: the coverage floor and the dependency scan, because each one is a single job in a repository that already runs `go test`.

## 2. Branch protection per environment

The page asks for protection on `develop`, `qa` and `main`: no direct pushes, required status checks, required review and "up to date with base".

**The constraint is the plan, not the rule.** On the Free plan of the organization, a private repository cannot enable protected branches. I checked it on `csp-docs`: the branch protection of `main` answers `403` with "Upgrade to GitHub Pro or make this repository public to enable this feature". The same answer is expected for every private repository of the organization on this plan, but I checked only `csp-docs`.

The required status checks have a second blocker: GitHub Actions is locked by billing since 2026-10-14, so no check runs on a pull request, and a required check that never reports blocks every merge.

Options to decide with the team (each needs its own ADR, not a decision made in a pull request):

| Option | What it fixes | Cost or risk |
|---|---|---|
| Make the repositories public | Protection and required checks become available | The repositories expose their code and their documentation; the team must agree |
| Move the organization to a paid plan | Protection and required checks in private repositories | Cost, decided by the team and the professor |
| Keep the plan and use a detective check | A person or a script finds the merges with no review | It detects, it does not prevent |

Until one of them is decided, the rule "no direct pushes and one review" stays a team agreement, not a control.

## 3. Pipeline design for many services

The platform has separate repositories for each service (API, database and portal), so each service keeps its own pipeline: build, test, scan and package its image.

| Element | Design | State today |
|---|---|---|
| Per-service pipeline | Each repository runs its own jobs on its pull requests | Yes, in every repository I read |
| Path filters | A change to one service runs only that service's jobs | No `paths:` filter in the five main workflows I read (`csp-auth-api`, `csp-ticketing-api`, `csp-workflow`, `csp-api-gateway`, `csp-auth-portal`) |
| Contract tests across services | The contract of the API is checked by the consumers | Partly: the portal checks its requests against `auth-service.yaml`, and `csp-auth-api` validates its responses against the same contract |
| Release step | The tagged image of a green build is promoted to the next environment, not rebuilt | Not yet: the Compose files of `csp-auth-api`, `csp-api-gateway`, `csp-auth-portal`, `csp-booking-api` and `csp-catalog-api` build from the source (`build:`) |

The release step is the largest gap. Week 12 found the same rule: a rebuild is a different artifact, so what passed in QA is not what runs in production. The fix is the registry and the image variable of each service, decided in an ADR before anything is applied.

## 4. The CI/CD stories for MVP 3

These are drafts with testable criteria. They are not created on the board: the team does not open stories without approval, and the project rule does not allow new HUs without the owner's order.

| Draft story | Testable criterion | Gate it adds |
|---|---|---|
| Pull requests to `develop` need green build, tests and lint | A pull request with a failing test cannot be merged | Required checks, after the billing lock is solved or the plan changes |
| Coverage gate | A change that drops the coverage under the declared floor fails | Coverage check in the Go repositories |
| Secret and dependency scan block on findings | A pull request that adds a key, or uses a vulnerable module, fails | `gitleaks` and `govulncheck` in each repository |
| Merge to `main` builds and pushes a versioned image | A tagged release produces one image with its digest | Image build and push on the tag |
| A smoke test verifies the deployed system responds | After deploy, `/health` of the gateway and of each service answers `200` | Smoke job after the deploy |

Order: the coverage gate and the scan are the first two, because they need no decision outside the repositories. Branch protection comes after the decision of section 2.

## 5. Common mistakes, with evidence from our own repositories

- **A gate bypassed under pressure.** The page describes a test disabled to force a merge. I did not find that case, but the merge record shows the same risk: on 2026-10-30 I checked the last three merged pull requests of `csp-docs` (#313, #314, #315), and each was merged by its author with no approval recorded in GitHub. With the billing lock, the project rule replaces the bot review with the local tests and with my explicit authorization for each merge. Nothing in GitHub enforces either one.
- **A change merged without the agreed review.** The contract change of `POST /register` (ADR-062, `auth-service.yaml` 1.17.0) was merged on 2026-11-03 with my authorization recorded in the pull request. I did not find a review of the team in the pull request or in the story (#282) before the merge. The rule to fix this is not in the page, and it belongs to the team agreement of section 2.
- **Gates that exist only in a document.** Week 13 wrote `security.yml` and checked it only as YAML. A gate that has never run in the pipeline has not been tested.
- **Long-lived branches.** Not a problem this week by the measure of section 5 of the first session: the pull requests were open for minutes, median 0.1 hours.

## 6. Self-check (answers)

1. **A DoD item with no automated gate is...** self-certified, which is a risk. Make it a gate.
2. **Branch protection ensures...** no direct pushes, and a green check and a review before each merge.
3. **In a multi-service product, pipelines should...** run per service with path filters, and contract tests should check the services against each other. One giant pipeline is slow.
4. **A failing gate under deadline pressure should be...** fixed at its cause and never bypassed. If the gate is wrong, change it through review.
5. **The release step in CD...** promotes the green, tagged images to the next environment, without rebuilding them.
6. **A good final smoke test checks...** that the deployed system answers, for example `compose up` and `/health`. Compiling is not enough.

## 7. Next week (from the page)

The page says that the next session covers testing at scale, and then the MVP 3 release. The stories of section 4 are the input for that work, once they are approved.
