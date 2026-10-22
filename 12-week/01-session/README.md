# Week 12 - Session 1: Cloud deployment, models, strategies and zero-downtime

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: deploy the system to a cloud environment or a production-like host, using the **tagged image**, with a
**zero-downtime strategy** (rolling or blue-green) and **health checks**.

**Scope of this work.** It is a rehearsal, not a deployment of the project. By decision of the owner nothing is applied to the
`csp-*` repositories yet: the rehearsal runs on my own machine with a small stand-in application (NGINX with `/health` and
`/version`), kept in [rehearsal/](rehearsal/). What it proves is the mechanism; what it does not prove is listed in section 5.
The application of this to the real services is planned in [session 2](../02-session/README.md).

## 1. Which deployment model fits CineSync Platform

| Model | What I manage | Fit for the platform |
|---|---|---|
| IaaS (virtual machines) | The VM, the operating system, Docker, updates | Most control and most work. Too much operation for a team of four students. |
| Containers on one host with Docker Compose | The host and the Compose files | **What the project does today** (`csp-infra-postgres/compose.yml` composes every service). Simple, reproducible, but one host is one point of failure. |
| Containers on a managed service or an orchestrator (Kubernetes) | Images and configuration | The natural next step: the services are already containers and stateless. It adds health-based replacement and rolling updates. |
| PaaS | Only the code | Possible for the portals, but the platform already has its own gateway, broker and databases. |
| Serverless | Only functions | Does not fit: the services are long-running and some jobs (`csp-worker`) keep state in PostgreSQL. |

My recommendation for MVP 3: keep Docker Compose as the base of every environment and add a production-like host with the images promoted by tag.
Moving to an orchestrator is a later decision, once the build-once and promotion pipeline works.

## 2. Build once, deploy the artifact

The rehearsal follows build, release, run:

1. **Build.** `docker build --build-arg VERSION=v1.1.0` runs once per version and the image is pushed to a registry (a local `registry:2`).
2. **Release.** The environment is the pair "image + configuration". The environment name (`ENVIRONMENT`) is read when the container starts, it is not baked in the image.
3. **Run.** `qa` and `prod` run the **same digest** (`sha256:f953ad...` for v1.0.0). Only the variable changes, as the output of section 4 shows.

## 3. The rehearsal (blue-green)

Files in [rehearsal/](rehearsal/): `Dockerfile` (the artifact, with a `HEALTHCHECK`), `app/default.conf.template`
(`/health` and `/version`), `compose.yml` (blue, green and the edge that sends the traffic), `edge/nginx.conf` and `run.sh` (the whole rehearsal).

- `blue` runs the current version (v1.0.0) and `green` the new one (v1.1.0), side by side. The old one is never stopped before the new one answers.
- The edge starts only when both are `healthy` (`depends_on: condition: service_healthy`).
- The switch rewrites one line (`set $active_upstream green:8080;`) and runs `nginx -s reload`: the reload keeps the connections that are in progress, so there is no restart.
- Rolling back is the same switch in the other direction, because blue is still running.

Run it with Docker started: `cd 12-week/01-session/rehearsal && ./run.sh`.

## 4. Real output

Docker 27.4.0 on Windows 11, 2026-10-22. The first lines (download of the registry image) are left out.

```
== 2. Build each version ONCE and push it
v1.0.0 -> 127.0.0.1:5000/deploy-demo@sha256:f953ad1b64a98bf26015b97953c7bb7e1a2f00fbc545b1c79bf0c85f45f8fea3
v1.1.0 -> 127.0.0.1:5000/deploy-demo@sha256:86187efdaef9e88fd317a3069fe1f27da0b5f7a5dca486511ae283ca6e972e77

== 3. Promotion: the same digests run in qa and in prod, only the environment changes
qa: version=v1.0.0 environment=qa host=d3e3d7d8c9d7
prod: version=v1.0.0 environment=prod host=0a9dd1b2f4b3

== 4. Health checks: both colours healthy before any traffic moves
blue  Up 9 seconds (healthy)
green  Up 9 seconds (healthy)

== 5. Blue-green release under load (a request every 50 ms)
release to v1.1.0 done in 2245 ms: version=v1.1.0 environment=prod host=846c8589734e

== 6. Rollback to the previous version
rollback to v1.0.0 done in 2250 ms: version=v1.0.0 environment=prod host=0a9dd1b2f4b3

== 7. Result of the load during the release and the rollback
requests: 54, answered 200: 54, failed: 0
```

The times include the `docker compose exec` that writes the line and reloads the edge, which is most of the 2.2 seconds; the switch itself is the reload of NGINX.

## 5. What the rehearsal does not prove

- It is one host on my machine, not a cloud environment, and the stand-in application is stateless: no database, no Redis, no broker.
- 54 requests are a light load. It shows that no request failed in this run, not how the strategy behaves at the real traffic.
- A database migration is not covered. The migrations of `csp-auth-db` (Flyway) only go forward, so a rollback of the service must work with the schema that is already applied. That needs the expand-and-contract rule of the plan in session 2.
- The registry is local. A real one (for example GitHub Packages of the organization) and its permissions have to be decided with the team.

## 6. Where the project stands against these ideas

Read in the repositories (branch `develop`), without changing anything:

- The Compose files of `csp-auth-api` and `csp-api-gateway` build the service image **from the source** (`build:`), and the platform composition includes them, so the host that deploys also builds. The rule "never rebuild for production" is not met yet.
- The Dockerfiles of `csp-auth-api` and `csp-api-gateway` have no `HEALTHCHECK` line. The services do have health routes (`/api/v1/auth/health` and `/health/ready`).
- The gateway already resolves its upstreams on every request through Docker DNS and takes the address of each domain from a variable (`APP_GATEWAY_AUTH_UPSTREAM`). It is the same mechanism as the edge of the rehearsal, so a blue-green switch of a domain would be a change of that variable and a reload.
- Releases are already tagged (`v2.0.0` for the shell and the portals, `v0.1.x` for the others), but the tag points to the source, not to a published image.
- The environments already exist as files (`env/dev.env`, `qa.env`, `main.env`, each with its own Compose project), which helps with promotion.

## 7. Self-check (answers)

- **Why never rebuild for production:** a rebuild is a different artifact; what was tested in QA is no longer what runs.
- **Why not stop the old version first:** between the stop and the start nobody answers.
- **Blue-green versus rolling:** blue-green switches all the traffic at once and keeps the old version for an instant rollback; rolling replaces instances a few at a time and needs more than one instance of the service.
