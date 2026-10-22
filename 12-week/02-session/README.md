# Week 12 - Session 2: Deployment plan, Infrastructure as Code and rollback (MVP 3)

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

This is a plan. It is **not applied to the project yet**: no `csp-*` repository is changed and the stories of section 6 are drafts that are
not created on the board (the team does not open stories without approval). The rehearsal that backs the numbers is in
[session 1](../01-session/README.md). Every statement about the project was read in the repositories (branch `develop`).

## 1. Artifact promotion across environments

One image per service, built once from a release tag and promoted by its digest. Only the configuration changes between environments.

| Step | What happens | Where it is defined |
|---|---|---|
| Release | The release `release/X.Y.Z` is merged to `main` and tagged `vX.Y.Z` (norm 6.3.2, 11). | Exists today. |
| Build once | CI builds the image of the service from the tagged commit and pushes it to the registry as `vX.Y.Z`, keeping its digest. | To build. The workflows (`ci.yml`) push no image to any registry, and the CI is locked by billing since 2026-10-14. |
| `dev` | Runs the image of `develop` (`develop-<short sha>`). | To build. |
| `qa` | Runs the digest of the release candidate, with `env/qa.env`. | The example of the env file and the Compose project name (`csp-qa`) exist. |
| `main` | Runs the **same digest** that passed `qa`, with `env/main.env`. | The example of the env file and the Compose project name (`csp-main`) exist. |

Rules: an image is never rebuilt for an environment; the configuration comes from the environment file and the secrets of the environment
(norm 5.9), never from the image; the registry keeps the previous digests so a rollback has something to go back to.

## 2. Infrastructure as Code

What already works as code: the Compose file of each repository, the platform composition of `csp-infra-postgres` (`up.sh dev|qa|main`),
the environment examples (`env/*.env.example`), the definitions of the broker (`csp-infra-rabbitmq`), the Flyway executor of each `-db`
repository and the observability configuration. The environments are created from files, not by hand.

What is missing for MVP 3:
- The services take their image from a variable (`AUTH_API_IMAGE`, ...) instead of `build:` (today the Compose files of `csp-auth-api` and `csp-api-gateway` build from the source).
- A registry and its credentials, loaded by the person as secrets.
- A check of the Compose files in CI (`docker compose config`) so a change of the environment is reviewed like code.

Decision proposed: **Docker Compose is the IaC of MVP 3**; Terraform is not needed while the platform runs on one host. Choosing the registry
and the tool needs an ADR before anything is applied.

## 3. Zero-downtime strategy

**Blue-green per domain**, with the API gateway as the switch. The gateway already takes the address of each domain from a variable
(`APP_GATEWAY_AUTH_UPSTREAM`) and resolves it on every request; the rehearsal of session 1 used the same mechanism. The switch is a change
of that variable and an NGINX reload.

| Component | Strategy | Why |
|---|---|---|
| `csp-auth-api`, `csp-booking-api`, other services | Blue-green: the new version starts next to the old one, passes its health check, the gateway switches. | Stateless, one instance per environment today, so rolling has nothing to roll. |
| `csp-worker` | Blue-green is safe | The relay takes a PostgreSQL lock and an instance that does not get it skips the run, so two workers do not publish twice. |
| Portals and `csp-front` | Blue-green behind the shell | Static files; the rollback is the previous tag. |
| PostgreSQL, Redis, RabbitMQ, MongoDB | No blue-green: stateful | Updated rarely, with the operations plan of each instance. |

Canary is not needed: there is no real traffic to split.

**Database rule (expand and contract).** The Flyway migrations only go forward, so a release must work with the schema of the previous
one. Migrations are applied first, they only add (columns, tables, indexes), and the destructive step goes to a later release. Without
this, a rollback of the service would meet a schema it does not understand.

**Health checks.** Every Dockerfile gets a `HEALTHCHECK` on the health route of the service (today `csp-auth-api` and `csp-api-gateway` have none),
and the deploy script does not switch until the new colour is `healthy`.

## 4. Rollback, documented and timed

1. Before each release the deploy script writes the digest that is running (`releases.log`).
2. To go back: run the switch to the previous colour (still running), or start the previous digest if it was already stopped.
3. Verify recovery: the health route answers `200`, a registration with a test account through the gateway answers `201`, and the error ratio of the alert of the workshop (week 11) goes down.

| Measure | Value | Status |
|---|---|---|
| Rollback by switching to the previous colour, rehearsal on one host | 2.25 s (includes the `docker compose exec`); 0 failed requests of 54 | Measured in session 1, not on the project. |
| Rollback of the real platform | Target: under 5 minutes, from the decision to the verification | **To rehearse and time** when the plan is applied. |

An untested rollback is not trusted: the first story of the rollback is to rehearse it on `qa`, with the time written in the PR.

## 5. Risks and dependencies

- The CI is locked, so the first builds and pushes would be manual until it works again.
- The registry, the tool and the move from `build:` to images need an ADR and the agreement of the other owners (their services, `csp-front`, `csp-api-gateway` and `csp-infra-*` are shared).
- The gateway is one instance: it cannot be blue-green with itself; it would be updated by stop and start in a quiet moment, or with a second edge.
- Stateful instances stay a single point of failure; backups are a separate topic (production readiness, next sessions).

## 6. Deployment stories (drafts, not created)

Each one would be a technical issue under an approved story.

**D1 - The image of a service is built once on the release tag**
```gherkin
Scenario 1: Tag publishes the image
  Given the release vX.Y.Z is tagged in main
  When  the pipeline runs
  Then  one image is pushed to the registry tagged vX.Y.Z and its digest is written in the job summary

Scenario 2: No rebuild for an environment
  Given the image of vX.Y.Z exists
  When  it is deployed to qa and then to main
  Then  both environments run the same digest
```

**D2 - Every environment is created from files**
```gherkin
Scenario: Compose takes the image from a variable
  Given the Compose file of a service has no "build:" key
  When  the environment starts with the variable of the image and env/<environment>.env
  Then  the service runs that image and no secret is read from the image
```

**D3 - Every service reports its health to the orchestrator**
```gherkin
Scenario: The new version is not switched before it is healthy
  Given a new version starts next to the old one
  When  its health check does not pass
  Then  the switch does not happen, the old version keeps serving and the deploy fails
```

**D4 - A domain is released with blue-green and no failed request**
```gherkin
Scenario: Switch under load
  Given a load of one request every 50 ms on the gateway
  When  the domain switches from blue to green
  Then  every request answers 200 and the version of the answer changes
```

**D5 - Rollback is rehearsed and timed**
```gherkin
Scenario: Rollback on qa
  Given a release is running on qa
  When  the previous version is restored following the runbook
  Then  the health route and the test registration pass in under 5 minutes, and the time is recorded in the PR
```

**D6 - A migration can be rolled back with the service**
```gherkin
Scenario: Release with a migration
  Given a release carries a new Flyway migration
  When  the previous version of the service runs against the migrated schema
  Then  its tests pass, and a migration that drops or renames is rejected by the check
```
