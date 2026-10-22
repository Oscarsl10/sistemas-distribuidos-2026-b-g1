<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       12-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 12

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Guillermo Sierra Lozano
- GITHUB_USER: Oscarsl10
- TEAM: CineSync Platform
- SPRINT_GOAL: Finish HU-AUTH-004 (change and verify the email) and implement HU-AUTH-005 (protection after repeated failed logins) in the Auth domain, add the RabbitMQ instance and the Auth outbox relay that the asynchronous journey needs, and keep the 3.0.0 releases ready for approval.
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| HU-AUTH-003 | Admin Manages User Roles | done | Release 3.0.0 pull requests, open: [csp-auth-api#83](https://github.com/code-corhuila/csp-auth-api/pull/83), [csp-auth-db#33](https://github.com/code-corhuila/csp-auth-db/pull/33), [csp-auth-portal#49](https://github.com/code-corhuila/csp-auth-portal/pull/49), [csp-front#61](https://github.com/code-corhuila/csp-front/pull/61), [csp-api-gateway#6](https://github.com/code-corhuila/csp-api-gateway/pull/6) |
| HU-AUTH-004 | User Changes and Verifies Email | doing | [csp-auth-api#105](https://github.com/code-corhuila/csp-auth-api/pull/105), [csp-auth-api#107](https://github.com/code-corhuila/csp-auth-api/pull/107), [csp-auth-api#110](https://github.com/code-corhuila/csp-auth-api/pull/110), [csp-auth-db#36](https://github.com/code-corhuila/csp-auth-db/pull/36), [csp-auth-portal#57](https://github.com/code-corhuila/csp-auth-portal/pull/57), [csp-front#60](https://github.com/code-corhuila/csp-front/pull/60) (in review on the board) |
| HU-AUTH-005 | Account Protection After Repeated Failed Logins | doing | [csp-auth-api#129](https://github.com/code-corhuila/csp-auth-api/pull/129), [csp-auth-api#130](https://github.com/code-corhuila/csp-auth-api/pull/130), [csp-auth-api#131](https://github.com/code-corhuila/csp-auth-api/pull/131), [csp-auth-portal#64](https://github.com/code-corhuila/csp-auth-portal/pull/64), [csp-docs#178](https://github.com/code-corhuila/csp-docs/pull/178) |

## 2. My individual contribution

Only my own work is listed; in the repositories shared with the team (`csp-front`, `csp-worker`, `csp-infra-*`) only the commits I authored are counted.

- **HU-AUTH-005, login protection (in progress)**: domain and use case of the failed-login limit ([csp-auth-api#129](https://github.com/code-corhuila/csp-auth-api/pull/129)), the Redis adapter of the attempt counter and lock ([csp-auth-api#130](https://github.com/code-corhuila/csp-auth-api/pull/130)), `429 rate_limited` on `POST /login` with the client address and the wiring ([csp-auth-api#131](https://github.com/code-corhuila/csp-auth-api/pull/131)), and the message that explains the limit on the login form ([csp-auth-portal#64](https://github.com/code-corhuila/csp-auth-portal/pull/64)). The decision is recorded in ADR-040 and the Auth docs ([csp-docs#178](https://github.com/code-corhuila/csp-docs/pull/178), [csp-docs#180](https://github.com/code-corhuila/csp-docs/pull/180)).
- **HU-AUTH-004, email change (finishing)**: the use cases, the adapters, `PATCH /me/email` and `POST /me/email/verify`, and the SMTP delivery of the code ([csp-auth-api#105](https://github.com/code-corhuila/csp-auth-api/pull/105) to [csp-auth-api#112](https://github.com/code-corhuila/csp-auth-api/pull/112)), the pending email and attempt counter in `email_verification` ([csp-auth-db#36](https://github.com/code-corhuila/csp-auth-db/pull/36)), the profile page of the portal ([csp-auth-portal#57](https://github.com/code-corhuila/csp-auth-portal/pull/57)) mounted at `/profile` in the shell ([csp-front#60](https://github.com/code-corhuila/csp-front/pull/60)), the exact `Retry-After` seconds of the limit ([csp-auth-api#118](https://github.com/code-corhuila/csp-auth-api/pull/118)), and the purge of verification requests expired more than 7 days ago ([csp-auth-api#117](https://github.com/code-corhuila/csp-auth-api/pull/117), [csp-auth-db#37](https://github.com/code-corhuila/csp-auth-db/pull/37)). The SMTP relay provider and its limits are documented in the runbook ([csp-docs#169](https://github.com/code-corhuila/csp-docs/pull/169), [csp-auth-api#113](https://github.com/code-corhuila/csp-auth-api/pull/113)).
- **RabbitMQ for the domain events**: ADR-036 for the instance `csp-infra-rabbitmq` ([csp-docs#160](https://github.com/code-corhuila/csp-docs/pull/160), amended after checking the permissions against a real broker in [csp-docs#164](https://github.com/code-corhuila/csp-docs/pull/164), accepted with ADR-027 in [csp-docs#179](https://github.com/code-corhuila/csp-docs/pull/179)), and the repository itself with its definitions, scripts and tests that run against a real broker ([csp-infra-rabbitmq#2](https://github.com/code-corhuila/csp-infra-rabbitmq/pull/2), [csp-infra-rabbitmq#3](https://github.com/code-corhuila/csp-infra-rabbitmq/pull/3)).
- **Auth outbox relay in `csp-worker`**: instance lock before each run, the relay of the Auth outbox behind its own flag (off by default until the consumers exist), the variables in the compose file and an end-to-end test that does not depend on the publication order ([csp-worker#27](https://github.com/code-corhuila/csp-worker/pull/27), [csp-worker#29](https://github.com/code-corhuila/csp-worker/pull/29), [csp-worker#31](https://github.com/code-corhuila/csp-worker/pull/31), [csp-worker#33](https://github.com/code-corhuila/csp-worker/pull/33)).
- **Platform**: logins `auth_migrator` and `booking_migrator` that own their schemas in the shared PostgreSQL instance, and the steps for an existing environment ([csp-infra-postgres#17](https://github.com/code-corhuila/csp-infra-postgres/pull/17), [csp-infra-postgres#19](https://github.com/code-corhuila/csp-infra-postgres/pull/19), [csp-infra-postgres#21](https://github.com/code-corhuila/csp-infra-postgres/pull/21)). The shell and the Auth portal now refuse an unanchored CORS regex (ADR-026 amended in [csp-docs#156](https://github.com/code-corhuila/csp-docs/pull/156); [csp-front#58](https://github.com/code-corhuila/csp-front/pull/58), [csp-auth-portal#59](https://github.com/code-corhuila/csp-auth-portal/pull/59)).
- **Releases 3.0.0 prepared**: eight pull requests to `main`, open and waiting for the approval: [csp-auth-api#83](https://github.com/code-corhuila/csp-auth-api/pull/83), [csp-auth-db#33](https://github.com/code-corhuila/csp-auth-db/pull/33), [csp-auth-portal#49](https://github.com/code-corhuila/csp-auth-portal/pull/49), [csp-front#61](https://github.com/code-corhuila/csp-front/pull/61), [csp-api-gateway#6](https://github.com/code-corhuila/csp-api-gateway/pull/6), [csp-infra-redis#9](https://github.com/code-corhuila/csp-infra-redis/pull/9), [csp-infra-postgres#15](https://github.com/code-corhuila/csp-infra-postgres/pull/15), [csp-worker#24](https://github.com/code-corhuila/csp-worker/pull/24).
- **Session activities (not applied to the project)**: [session 1](../01-session/README.md), a timed blue-green rehearsal on my machine with a stand-in application (image built once and promoted by digest, health checks, 0 failed requests of 54 during the release and the rollback), and [session 2](../02-session/README.md), the deployment plan for MVP 3 (artifact promotion, Docker Compose as IaC, blue-green through the gateway, expand-and-contract migrations, timed rollback) with six story drafts that are not created on the board.
- **Observability workshop of week 11** (delivery 2 closed on 2026-10-21): reported in [the week 11 status](../../11-week/hu-status/README.md).
- Full commit log of the week: [commits.md](commits.md).

## 3. Blockers and risks

- GitHub Actions is still locked by billing since 2026-10-14: the pull requests of the week were validated locally and merged with the explicit authorization of the owner, with the bot review done. CI evidence is missing for them.
- The eight 3.0.0 release pull requests need the approval of the professor and of the leader; HU-AUTH-001 to HU-AUTH-003 reach `main` only with them.
- The Auth relay stays off: `csp-ticketing-api` has no consumer and no queue yet, so the asynchronous journey cannot run end to end.
- HU-AUTH-004 is in review and needs a real SMTP server for the delivery of the code outside development.
- The Compose files of `csp-auth-api` and `csp-api-gateway` build the image from the source on the host and pull no tagged image, so the build-once, promote-the-tagged-image rule of the sessions of this week is not met yet, and there is no cloud environment; the rehearsal ran on one local host with a stand-in application.

## 4. Plan for next week

- Close HU-AUTH-004 and HU-AUTH-005 (review, promotion to `qa`) and start HU-AUTH-006, spec and contract first, then a failing test, then the code.
- Take the deployment plan for MVP 3 (written this week, not applied) to the team: decide the registry and the IaC tool in an ADR, get the approval of the stories before creating them, and rehearse and time the rollback on `qa` when it is applied.
- Keep the 3.0.0 releases ready and the documentation in sync with every decision through `docs/* -> main` pull requests.

## 5. Compliance self-check

- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (code repos: child branch -> `develop`; `csp-docs`: `docs/* -> main`)
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration). Unit tests of the login limit, integration tests of the Redis adapter, tests of `csp-infra-rabbitmq` against a real broker and an end-to-end test of the relay.
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links

- [Week 12 Full Commit Log](commits.md)
- [Repository: csp-docs](https://github.com/code-corhuila/csp-docs)
- [Repository: csp-auth-api](https://github.com/code-corhuila/csp-auth-api)
- [Repository: csp-auth-db](https://github.com/code-corhuila/csp-auth-db)
- [Repository: csp-auth-portal](https://github.com/code-corhuila/csp-auth-portal)
- [Repository: csp-front](https://github.com/code-corhuila/csp-front)
- [Repository: csp-worker](https://github.com/code-corhuila/csp-worker)
- [Repository: csp-infra-rabbitmq](https://github.com/code-corhuila/csp-infra-rabbitmq)
- [Repository: csp-infra-postgres](https://github.com/code-corhuila/csp-infra-postgres)
- Week 12 summary session 1-2:

  ![Week 12 summary session 1-2](summary_week_12.png)
