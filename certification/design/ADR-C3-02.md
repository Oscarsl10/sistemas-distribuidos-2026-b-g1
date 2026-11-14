# ADR-C3-02: Deploy the services on CCE with autoscaling and a CodeArts pipeline with rolling updates

## Status

Proposed

## Context

`10-devops/environments.md` says that only the local environment exists and that no workflow deploys. The target design in `10-devops/ci-cd.md` already assumes Kubernetes namespaces (`csp-dev`, `csp-qa`, `csp-prod`), a rolling deploy with `kubectl rollout` and a rollback with `kubectl rollout undo`. Each API repository builds one container (`deploy/Dockerfile`) and `NFR-003` in `04-requirements/non-functional.md` expects Booking, Ticketing and Concessions to scale horizontally. The worker is the exception: it must stay at one instance because of its advisory lock. GitHub Actions of the organization is locked by billing, so the pipeline that exists cannot run.

## Decision

Run every service of CSP as a CCE deployment, one namespace per environment, with horizontal pod autoscaling and node pool autoscaling. Build and scan the images in CodeArts, store them in SWR and deploy with a CodeArts pipeline: automatic for `dev` and `qa`, and a manual approval with a rolling update for `prod`. Roll back with `kubectl rollout undo` when `/health/ready` fails, and run the Flyway job of the `-db` repository before each API deploy. Scale Booking, API Gateway and Ticketing between the limits of section 3 of the design note.

## Alternatives considered

- **ECS virtual machines with Docker Compose.** Rejected: it copies the local setup, but gives no rolling update, no autoscaling per service and no clean rollback, and each VM becomes a pet to patch.
- **Serverless pods on Cloud Container Instance.** Rejected: the services are steady and hold database and AMQP connections, so per-second billing saves little and it gives less control over node pools and network placement.
- **Keep GitHub Actions for deployment.** Rejected: the organization account is locked by billing, and a deploy that depends on it would stop with it.

## Consequences

- Positive: one mechanism gives autoscaling, rolling updates and rollback for all services, and it matches the Kubernetes commands that the runbooks already document.
- Positive: qa and production run the same images and the same manifests, so a release is promoted unchanged.
- Negative: the team must learn CCE and keep the cluster upgraded, which is extra operational work for four people.
- Negative: the worker needs an exception to the scaling rule, and a wrong setting could start two relays and publish events twice.
