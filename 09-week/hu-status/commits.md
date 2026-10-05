# Git Commit Log - CineSync Platform

## Week 9

Period: 2026-09-28 to 2026-10-04. Author: Oscar Sierra (`Oscarsl10`). Same change promoted across several branches is listed once; automatic merge commits are omitted. Newest first.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 170 |
| csp-concessions-api | 2 |
| csp-concessions-db | 1 |
| csp-concessions-portal | 2 |
| csp-infra-mongo | 2 |
| csp-booking-api | 15 |
| csp-booking-db | 22 |
| csp-booking-portal | 200 |
| csp-auth-api | 8 |
| csp-auth-db | 9 |
| csp-auth-portal | 6 |
| csp-catalog-api | 8 |
| csp-catalog-db | 12 |
| csp-catalog-portal | 31 |

### csp-docs

`code-corhuila/csp-docs` - 170 commits

| Commit | Message |
|---|---|
| `2a7ecbf8aa6a` | chore(docs): sync auth-schema-annex-a-alignment with auth-scaffold-decisions |
| `5dbeee804423` | docs(data): add password_reset and email_verified to the auth ERD |
| `4c81f7060d73` | chore(docs): sync auth-schema-annex-a-alignment with main |
| `29f13f1511a9` | chore(docs): sync auth-scaffold-decisions with main |
| `cc9161568a05` | docs(auth): align auth schema docs with Annex A and add readiness check |
| `3ae74552a63a` | docs(data): record auth scaffold decisions and reconcile auth tables |
| `4b003e49dd07` | docs(catalog): merge PR #89 - align the Catalog documents with ADR-019 and ADR-021 |
| `b006df2c2fb2` | docs(catalog): align the data model, contract, runbook and env with ADR-019 and ADR-021 |
| `d9624663c66c` | docs(architecture): merge PR #87 - record Angular with Native Federation as the front-end stack |
| `1c383b39c68e` | docs(architecture): record Angular with Native Federation and the csp-front shell name |
| `1a4eaf754097` | docs(architecture): merge PR #86 - record the qa-promote branch prefix for promotions to qa |
| `ef27c1283f04` | chore(docs): sync adr-qa-promotion-branch-prefix with main |
| `07113f64110b` | docs(api): merge PR #84 - declare the internal purge operation in the auth and concessions contracts |
| `d1a9b8241edc` | docs(architecture): record the qa-promote branch prefix for promotions to qa |
| `45970f52868c` | docs(quality): use PostgreSQL 16 in the testing examples |
| `ffee3777c5cd` | docs(architecture): merge PR #83 - record the outbox relay of other domains and the Catalog snapshot |
| `22995b46b49f` | docs(data): merge PR #82 - align booking tables, statuses and PostgreSQL version with the norm |
| `550494b82279` | docs(operations): document the retention and purge in the auth and concessions runbooks |
| `ea098cef4656` | chore(docs): sync internal-purge-auth-concessions with origin/docs/booking-outbox-and-catalog-adrs |
| `1c097f4bfd02` | docs(architecture): fix the Catalog line of the dependency map and the ADR-020 scope |
| `6fb8a3f211e0` | chore(docs): sync booking-outbox-and-catalog-adrs with origin/docs/python-stack-no-api-migrations |
| `7e66cab1c40d` | docs(architecture): record the engine version and the seat index in ADR-013 |
| `f4b0e1332c52` | chore(docs): sync internal-purge-auth-concessions with main |
| `f3eb151a6c37` | chore(docs): sync booking-outbox-and-catalog-adrs with main |
| `3eebc33505d1` | docs(data): use booking tables in the naming examples |
| `95a16edd0588` | docs(architecture): use the per-engine infra names in the ADR-020 scope |
| `3a41053fb39d` | docs(api): record why the purge operation of auth declares its own security |
| `de6f323a6dc5` | chore(docs): sync internal-purge-auth-concessions with booking-outbox-and-catalog-adrs |
| `97b079a32447` | docs(architecture): keep the ADR index in ascending order |
| `20e5ad3535cd` | chore(docs): sync booking-outbox-and-catalog-adrs with python-stack-no-api-migrations |
| `4c6f5e0832cd` | chore(docs): sync python-stack-no-api-migrations with main |
| `ce828139e6a1` | docs(data): finish the singular names in the booking contract and the naming example |
| `7d6120d7c0ce` | chore(docs): sync product-backlog-cut2 with main |
| `7e54adb60282` | chore(docs): sync traceability-cut2 with main |
| `5b4583e84f8d` | chore(docs): sync fr-loose-and-fe-stories with main |
| `949c24155084` | chore(docs): sync hu-fe-auth-001 with main |
| `743d0fdc3caf` | chore(docs): sync hu-fe-ticketing-001 with main |
| `4a995ffd8c58` | chore(docs): sync hu-fe-concessions-001 with main |
| `b06913db596f` | chore(docs): sync hu-fe-booking-001 with main |
| `822f5ea0111c` | chore(docs): sync hu-fe-catalog-001 with main |
| `8fe5dc2d05d4` | chore(docs): sync hu-booking-003 with main |
| `9c3ecfd95af8` | chore(docs): sync hu-booking-002 with main |
| `5b9ac1a423fe` | chore(docs): sync hu-booking-001 with main |
| `d1fd23496690` | docs(booking): align HU-FE-BOOKING-001 with the Cut 2 snapshot modes and the reservations list |
| `65bc5db8443d` | chore(docs): sync hu-fe-booking-001 with hu-booking-001 |
| `7ab67a97c026` | docs(booking): state the Cut 2 snapshot modes and list the reservations endpoint |
| `b8a0c9575aae` | docs(booking): define the Cut 2 transition of the Catalog snapshot |
| `4168b85934ca` | docs(data): merge PR #82 branch |
| `8006abdd24f7` | docs(data): allow a null showtime start while Catalog is not configured |
| `ce20b9a7451e` | docs(api): declare the internal purge operation in the auth and concessions contracts |
| `53a1821d8aff` | docs(data): document the Catalog outbox collection and align the service credential |
| `4f7721eba070` | docs(architecture): record the outbox relay of other domains and the Catalog snapshot |
| `5d74981ab578` | docs(data): align booking tables, statuses and PostgreSQL version with the norm |
| `1a6a9c827081` | docs(stacks): merge main |
| `22a6c9e90d46` | docs(stacks): merge PR #81 - remove the migration library from the Java service guide |
| `196f208e6cc0` | docs(data): merge PR #78 - Liquibase for catalog, Flyway for ticketing and concessions, owned by each db repository |
| `5f7a878f7854` | docs(stacks): remove the migration library from the Python service guide |
| `3a4db58f314c` | docs(ticketing): merge PR #79 - record Flyway owned by csp-ticketing-db |
| `62481cdba067` | docs(concessions): merge PR #80 - record Flyway owned by csp-concessions-db |
| `ec959109f12f` | docs(auth): merge PR #77 - replace golang-migrate with Flyway owned by csp-auth-db |
| `457705b8a784` | docs(booking): merge PR #73 - document health split, read-only outbox relay and expiration |
| `f09acd1a3186` | docs(booking): merge PR #72 - align data model with db-owned migrations and read-only outbox |
| `5e2eb5da6a47` | docs(booking): merge PR #70 - contract 3.0.0 with reservation list, health split and internal maintenance operations |
| `2c3384c7a562` | docs(stacks): remove unnecessary blank line in java-spring.md |
| `37e932e80264` | docs(stacks): add section for tools and minimum versions |
| `1be8dd098688` | docs(devops): merge PR #76 - run migrations from each db executor, never at API startup |
| `f85eaca03538` | docs(data): merge PR #75 - make each db repository the owner of migrations |
| `38435633864e` | docs(adr): merge PR #69 - add booking stack ADR-012, ADR-013 and outbox relay ADR-014 |
| `1d056acd2d84` | chore(docs): sync navigation-map-fe with product-backlog-cut2 |
| `3c7687731325` | docs(requirements): rename Mongock to Liquibase in the story allocation row |
| `1073581ad9f4` | chore(docs): sync product-backlog-cut2 with traceability-cut2 |
| `387a5c2d6364` | chore(docs): sync traceability-cut2 with fr-loose-and-fe-stories |
| `71dc7b7e7583` | chore(docs): sync fr-loose-and-fe-stories with hu-fe-auth-001 |
| `864b5ca78b79` | chore(docs): sync hu-fe-auth-001 with hu-fe-ticketing-001 |
| `f31b952dd289` | chore(docs): sync hu-fe-ticketing-001 with hu-fe-concessions-001 |
| `3197591c4d33` | chore(docs): sync hu-fe-concessions-001 with hu-fe-booking-001 |
| `4c292a14a18a` | chore(docs): sync hu-fe-booking-001 with hu-fe-catalog-001 |
| `b601cad103d3` | fix(booking): align HU-FE-BOOKING-001 with token validation and the 422 error code |
| `10fbc97c8a9d` | chore(docs): sync hu-fe-catalog-001 with hu-booking-003 |
| `44fdaa93cfb6` | chore(docs): sync hu-booking-003 with hu-booking-002 |
| `7cc727319ed6` | fix(booking): align HU-BOOKING-003 errors and token rule with the norm |
| `3d349337450f` | fix(booking): make HU-BOOKING-002 expiration run in the worker as the norm states |
| `ab934b0a76b2` | chore(docs): sync hu-booking-002 with hu-booking-001 |
| `096cff25b4c8` | fix(booking): require JWT validation in HU-BOOKING-001 as the norm states |
| `56675d6459d2` | chore(docs): sync hu-booking-001 with hu-reorganize-sections |
| `03122a621e03` | chore(docs): sync hu-reorganize-sections with main |
| `b57e2493c711` | chore(docs): sync concessions-migrations-ownership with main through its base |
| `8ecd5a2669e7` | chore(docs): sync ticketing-migrations-ownership with main through its base |
| `dace92f8b7a0` | docs(catalog): name csp-infra-mongo and ADR-006 as the Mongo instance source |
| `28f26994fc9d` | chore(docs): sync catalog-migrations-liquibase with main through its base |
| `9371e324d22d` | chore(docs): sync auth-migrations-flyway with main through its base |
| `c9d27c7c9fd5` | chore(docs): sync catalog-migrations-liquibase with its base |
| `11db83c5a2eb` | chore(docs): sync auth-migrations-flyway with its base |
| `0620b0930351` | docs(booking): document outbox retention and purge in the runbook |
| `dd100f6a9ea2` | docs(booking): make the idempotency_key hold index unique |
| `9988b551dab0` | docs(booking): align contract 3.0.0 with the norm and trace internal operations |
| `2b82283606b0` | chore(docs): sync booking-contract-list-and-health with its remote |
| `241ab27b3abd` | chore(docs): sync java-stack-no-api-migrations with main |
| `1bc757ae25d8` | chore(docs): sync migration-commands-devops with main |
| `a2c058b16da6` | chore(docs): sync migration-policy-platform with main |
| `51204db13e0f` | chore(docs): sync booking-data-model-migrations with main |
| `fff7cb839181` | chore(docs): sync booking-health-outbox-docs with main |
| `9683acf74ee0` | chore(docs): sync booking-contract-list-and-health with main |
| `1943a0de85a2` | chore(docs): sync adr-booking-stack with main |
| `73cb22308a3e` | docs(booking): merge PR #74 - add internal maintenance operations for the worker |
| `c728ab2ac3fc` | docs(adr): merge PR #71 - add booking outbox relay ADR-014 |
| `e4a75586e41c` | docs(data): replace MongoDB change units with Liquibase in the models |
| `db3d5ca782e3` | docs(stacks): remove the migration library from the Java service guide |
| `ea9df26ac97f` | docs(concessions): record Flyway owned by csp-concessions-db |
| `5aa6fb4fe8e2` | docs(ticketing): record Flyway owned by csp-ticketing-db |
| `b848bec13d3f` | docs(catalog): rename Mongock to Liquibase in vision and user stories |
| `a17b595a34d1` | docs(catalog): replace Mongock with Liquibase owned by csp-catalog-db |
| `89dd03099013` | docs(auth): replace golang-migrate with Flyway owned by csp-auth-db |
| `d60836ee8a50` | docs(devops): run migrations from each db executor, never at API startup |
| `0d9fcd646b85` | docs(data): make each db repository the owner of migrations |
| `c35175f2730a` | docs(booking): add internal maintenance operations for the worker |
| `be916c6e512f` | docs(booking): document health split, read-only outbox relay and expiration |
| `1b692a7301d9` | docs(booking): align data model with db-owned migrations and read-only outbox |
| `eb19f5160b0d` | docs(adr): add booking outbox relay ADR-014 |
| `5ec68fc95a20` | docs(booking): add reservation list and split liveness from readiness |
| `9e8ff3f7e4bb` | docs(adr): add booking stack ADR-012 and ADR-013 |
| `07a7c85e51e9` | docs(microservices): merge PR #48 - service-catalog corrections for internal health/availability endpoints per ADR-008 |
| `b9b2c29c542a` | docs(microservices): fold H39 note into endpoint, add ADR-010 to ticketing health, remove duplicate refs |
| `94dfc904db9a` | docs(microservices): add H39 traceability to catalog availability, fix exposure consistency (H39, H47) |
| `4d3110b2a4e6` | docs(microservices): fix PR-04d review items - add internal annotation to health endpoints, remove redundant note, fix blank lines |
| `56aea13d3a19` | docs(microservices): remove H-number references from documentation note (internal tracking only) |
| `74624f153482` | docs(microservices): fix PR-04d review items - fix typo, move note to catalog-service, remove redundant internal annotations, add ADR-010 reference |
| `397885374544` | docs(microservices): service-catalog corrections - ticketing health internal, showtimes availability wording, H39/H40 note (H45,H47,H39) |
| `f95923e62b0a` | docs(microservices): service-catalog corrections - ticketing health internal, showtimes availability wording (H45,H47) |
| `d2734f45d0a8` | docs(api - microservices): merge PR #47 - route coverage for auth reset, ticketing notifications, catalog filters |
| `79de74436c30` | docs(microservices): add ADR traceability to catalog availability endpoint in service-catalog |
| `6822ae070727` | docs(api,microservices): add auth password reset, ticketing notifications, and catalog movie filters |
| `9dad05784f97` | docs(concessions): merge PR #46 - add product status and stock history endpoints per HU-CONCESSIONS-006/007 |
| `e971f6b19316` | docs(api): add x-hu-ref and x-adr-ref traceability to concessions endpoints |
| `2a56db1af7f3` | docs(ux-ui): merge PR #45 - align navigation map and UI mockups per governance policies |
| `23cedd0a79bf` | docs(ux): mark /admin/reports as MOCK ONLY in nav-map, Screen Map, and mockup README |
| `a98853e66178` | docs(microservices): update concessions service catalog with new endpoints |
| `4c39d55f8296` | docs(ux): mark route gaps with HU refs and add 9 Screen Map rows (H35-H44) |
| `db04d51eeb56` | docs(api): add reports endpoint and movie filters to catalog contract |
| `aeeea42490e5` | docs(api): add notifications endpoint to ticketing contract |
| `a676f41e77aa` | docs(api): add forgot-password and reset-password endpoints to auth contract |
| `cbeb8e09c1f5` | docs(api): add movie filter params (q, genre, durationMax) to catalog contract |
| `ee521834738e` | docs(api): add PATCH product status and GET stock movements to concessions contract |
| `e6e259317c71` | docs(requirements): add HUs for concessions status, stock history, movie filters, auth reset, notifications, reports |
| `31e0f736478b` | docs(architecture): merge PR #44 - align async boundaries and Catalog sync per ADR-008 |
| `d15a31c44d40` | docs(microservices): finalize ASCII diagram - single nodes, explicit async/sync topology |
| `1448142379c1` | docs(microservices): align ASCII broker topology with C4 - single RabbitMQ Broker |
| `74c23d44cf84` | docs(diagrams): add owner and target to SVG-to-Mermaid follow-up (HU-ARCH-001) |
| `aa75503e8407` | docs(microservices): align Holds/Reservations sync with H25 - add ReservationExpired event |
| `3eeddd54f577` | docs(microservices): fix ASCII async flow direction and single Booking node |
| `12d3c3c468e5` | docs(architecture): unify C4 sync label to 'Synchronous OpenAPI' |
| `7625e0125b25` | docs(microservices): redraw ASCII diagram - single Booking node, clear sync/async paths |
| `c806070df6a2` | docs(diagrams): add HU-ARCH-001 traceability for SVG-to-Mermaid follow-up |
| `e2f998cbf881` | docs(microservices): fix ASCII diagram to show Booking sync validates Catalog |
| `654b8a9ec6bf` | docs(diagrams): document SVG auditability limitation per ADR-009 |
| `8a51279c6768` | docs(arch): reverse sync validation arrow to Booking -> Catalog |
| `678a417bf1f1` | docs(arch): align service catalog holds/reservations with ADR-008 |
| `7e56a2731e8c` | docs(arch): align holds/reservations access with ADR-008 events only |
| `d415c407966e` | docs(arch): fix async boundary - Booking sync validates against Catalog |
| `f95ed7ac4233` | docs(governance-api): merge PR #32 - align policies, ADRs and API contracts per ADR-006 |
| `8b911aec5e4f` | docs(api): add pagination and standard error responses to concessions contract |
| `0220c79de91d` | docs(architecture): add ADR-007 and ADR-008 to architecture README |
| `9cedecd611ef` | docs(architecture): update ADR statuses from Proposed to Accepted |
| `e105a3c756b4` | docs(architecture): update ADR-010 status to Accepted |
| `b7d2be25f5f5` | docs(architecture): update ADR-008 status to Accepted |
| `c88f8e400762` | docs(architecture): update ADR-007 status to Accepted |
| `6f540e8a4760` | docs(context): update scope.md datastore integrations row per ADR-006 |
| `755073a44604` | docs(governance): align microservices doc with ADR-006 topology |
| `fe9b6f4f7697` | docs(governance): align security policy doc with ADR-006 topology |
| `cc002344d134` | docs(governance): align DoD doc with ADR-006 topology |

### csp-concessions-api

`code-corhuila/csp-concessions-api` - 2 commits

| Commit | Message |
|---|---|
| `f47c19d385ec` | chore(governance): seed README and CODEOWNERS |
| `24d4e6e78d69` | Initial commit |

### csp-concessions-db

`code-corhuila/csp-concessions-db` - 1 commits

| Commit | Message |
|---|---|
| `5c1c8ff67813` | Initial commit |

### csp-concessions-portal

`code-corhuila/csp-concessions-portal` - 2 commits

| Commit | Message |
|---|---|
| `789ff2b6e8ba` | chore(governance): seed README and CODEOWNERS |
| `a54b6aa22756` | Initial commit |

### csp-infra-mongo

`code-corhuila/csp-infra-mongo` - 2 commits

| Commit | Message |
|---|---|
| `d4da6d771479` | chore(governance): seed README and CODEOWNERS |
| `549f8b5decec` | Initial commit |

### csp-booking-api

`code-corhuila/csp-booking-api` - 15 commits

| Commit | Message |
|---|---|
| `5310cbff5238` | docs(readme): merge PR #6 - document how to build, test and run the service |
| `3b0b2ac982a2` | docs(readme): document how to build, test and run the service |
| `98311fcd8ba4` | feat(domain): merge PR #5 - add the reservation aggregate and the typed domain errors |
| `64dd0aabb053` | chore(ci): sync feat/booking-domain-model with develop |
| `b9cde437c3ed` | chore(ci): merge PR #4 - add the CI workflow, the Dockerfile and the service compose file |
| `a92293a62eb9` | chore(ci): sync chore/booking-api-ci-deploy with develop |
| `ed62cafad5d2` | chore(github): merge PR #3 - add the pull request template |
| `33ff2a0e5ea3` | chore(ci): sync chore/booking-api-pr-template with develop |
| `b410373890cc` | feat(domain): add the reservation aggregate and the typed domain errors |
| `5d65c083c33d` | chore(ci): add the CI workflow, the Dockerfile and the service compose file |
| `e45f2c612135` | chore(runtime): merge PR #2 - scaffold the three-module Maven build and the composition root |
| `39296cbabdfa` | chore(governance): merge PR #1 - update README and CODEOWNERS identity |
| `77a48e4a5d1d` | chore(github): add the pull request template |
| `2281793a687e` | chore(runtime): scaffold the three-module Maven build and the composition root |
| `6e443c2eaf3d` | chore(governance): update README and CODEOWNERS identity |

### csp-booking-db

`code-corhuila/csp-booking-db` - 22 commits

| Commit | Message |
|---|---|
| `0bf1b48c6708` | docs(github): add the cross-domain check to the pull request template |
| `29521e159ca7` | chore(booking-db): merge PR #10 - promote pull request template to qa |
| `de3095c8163e` | fix(deploy): connect the migration executor with the domain user |
| `c6a86a2afa70` | chore(booking-db): merge PR #9 - promote migration executor to qa |
| `3812d9798fa3` | chore(booking-db): merge PR #8 - promote schema scaffold to qa |
| `54a6dc79946b` | chore(booking-db): merge PR #7 - promote governance to qa |
| `72be3c920fbf` | docs(readme): merge PR #6 - document the layout, the executor and the reversion order |
| `0be833e113bc` | chore(ci): sync chore/booking-db-readme with develop |
| `a01120c64094` | chore(ci): merge PR #5 - add the database rebuild workflow |
| `80b9e49cf219` | chore(ci): sync chore/booking-db-ci with develop |
| `ae66e8dd931c` | chore(github): merge PR #4 - add the pull request template |
| `d07f9d6635c2` | chore(ci): sync chore/booking-db-pr-template with develop |
| `cc9a19986bd9` | docs(readme): document the layout, the executor and the reversion order |
| `b90384054da2` | chore(ci): add the database rebuild workflow |
| `9177e97f1bfc` | chore(deploy): merge PR #3 - add the migration executor, env example and gitignore |
| `494707e9e4af` | chore(db): merge PR #2 - scaffold migration families, booking schema and roles |
| `84c6f016d8fb` | chore(governance): merge PR #1 - update README and CODEOWNERS identity |
| `d14fe8c8219e` | chore(github): add the pull request template |
| `ed044b8c64bf` | fix(db): move the schema settings to the Flyway default environment |
| `650f8d1d0a4e` | chore(deploy): add the migration executor, env example and gitignore |
| `2aa6148508e8` | chore(db): scaffold migration families, booking schema and roles |
| `672b72b5b24b` | chore(governance): update README and CODEOWNERS identity |

### csp-booking-portal

`code-corhuila/csp-booking-portal` - 200 commits

| Commit | Message |
|---|---|
| `9019fb577244` | build(booking-portal): merge PR #14 - promote npm lockfile to qa |
| `9eaf57994564` | test(booking-portal): merge PR #13 - promote domain adapter spec to qa |
| `078917bb6dea` | chore(booking-portal): merge PR #12 - promote domain scaffold to qa |
| `4a14981a1429` | chore(booking-portal): merge PR #11 - promote runtime scaffold to qa |
| `232d7d0dc08a` | chore(booking-portal): merge PR #9 - promote configuration scaffold to qa |
| `936005848afe` | docs(governance): merge PR #8 - promote governance affiliation update to qa |
| `9fde3558a1f7` | feat(booking): merge PR #7 - align the portal with booking-service 3.1.0 |
| `e50d9478eccb` | chore(booking): restore tsconfig.federation.json rewritten by the local build |
| `0caba44a237d` | feat(booking): send the Cut 2 title and room snapshot with the hold request |
| `76df466d9300` | feat(booking): send the createdBefore filter when listing reservations |
| `4c1b4e056734` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 84 |
| `cdc3422e6537` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 84 start |
| `a47d74c86be0` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 83 |
| `48b3e97209cc` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 83 start |
| `c9a9ea020dfc` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 82 |
| `30987ce946a6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 82 start |
| `ba3d39ed5528` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 81 |
| `a1fcb6fc8f74` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 81 start |
| `430337ef1b48` | chore(catalog): merge PR #6 - update governance affiliation |
| `311dc3505cd0` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 80 |
| `5780b0c7cb6d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 80 start |
| `44627ff5bf6e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 79 |
| `f2d4e2eb512a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 79 start |
| `54527425230e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 78 |
| `7e57a32e919f` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 78 start |
| `1377d5153fc3` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 77 start |
| `cbe2faf712e0` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 77 |
| `fe6eaad8620e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 76 |
| `efe1af8e8382` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 76 start |
| `335d5fe1817e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 75 |
| `24fab925598d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 75 start |
| `128e7de2c697` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 74 |
| `85d89033c7db` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 74 start |
| `20177648cab7` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 73 |
| `456e2132cc19` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 73 start |
| `6efdb2d209e5` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 72 start |
| `b455efa41095` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 72 |
| `84ba30f6f3bd` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 71 |
| `40e137918b3d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 71 start |
| `57fee1a17794` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 70 |
| `ae755e4e0935` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 70 start |
| `0d32b5fcbce7` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 69 |
| `37a46ac7cc99` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 69 start |
| `792295f8eae2` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 68 |
| `82e3c026aa8a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 68 start |
| `fd53af6d26af` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 67 |
| `039e56a98ac3` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 67 start |
| `acb0efe712b8` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 66 |
| `edb99f8fd582` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 66 start |
| `5710e13c3236` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 65 |
| `8cd67644f9eb` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 65 start |
| `3e057087a72a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 64 start |
| `a69ce8c1115c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 64 |
| `16865aac8936` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 63 |
| `c52b5c7fc075` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 63 start |
| `fb9c51964d67` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 62 |
| `c8cf15eb4c18` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 62 start |
| `50fe7521d6a6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 61 |
| `f7c130f5de99` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 61 start |
| `340586bb46eb` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 60 start |
| `3f010da656ee` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 60 |
| `5954847c7271` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 59 |
| `aa4413da1d10` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 59 start |
| `1e3161c8bc11` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 58 |
| `45fc301d57f8` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 58 start |
| `030a268c0f87` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 57 start |
| `351f4a9ae8ca` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 57 |
| `32d0b49efe9f` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 56 |
| `224a8d9cc980` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 56 start |
| `e4a3a95a77b6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 55 |
| `54538eaae063` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 55 start |
| `66369926c840` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 54 |
| `1a13798ddff1` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 54 start |
| `3c31041467b9` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 53 |
| `50d398a00f0f` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 53 start |
| `fbfb037237a8` | docs(governance): update project affiliation and doc links to cinesync platform |
| `de101326976e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 52 |
| `7d0f712dc197` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 52 start |
| `094815aedf15` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 51 |
| `37bbcb62fa91` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 51 start |
| `1814f8097112` | chore(booking): merge PR #4 - add CI and deployment scaffold |
| `40a974987351` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 50 |
| `fb43a3fbab5b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 50 start |
| `e257e2afa4b4` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 49 |
| `709fddcc873e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 49 start |
| `0489f9611490` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 48 start |
| `64b19c6abf58` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 48 |
| `2bf465558b73` | chore(booking): merge PR #3 - add booking domain scaffold |
| `dfcedee6a7ef` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 47 |
| `ee5bee3f7202` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 47 start |
| `c670b409fde4` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 46 |
| `dc2e46df8e06` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 46 start |
| `7787a65e00fb` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 45 |
| `93c6a6c8f979` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 45 start |
| `bebfff17e014` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 44 |
| `251804e4c205` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 44 start |
| `5ad7306a275b` | test(booking): validate container deployment scaffold |
| `28cd28536627` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 43 start |
| `812743e6ae7b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 43 |
| `b7aedbb28a30` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 42 |
| `a1ef1141ace2` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 42 start |
| `088ef5deda17` | chore(booking): sync deployment scaffold with domain tests |
| `19150e952bfa` | test(booking): validate domain scaffold adapter |
| `5471af6a8243` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 41 |
| `8598bca3cdf5` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 41 start |
| `7fc3a59de17c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 40 |
| `8b669183811b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 40 start |
| `5bc4b3898cf8` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 39 |
| `401c955ef68f` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 39 start |
| `425dc643ab02` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 38 start |
| `94cc19d417ea` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 38 |
| `99e8de04bfd3` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 37 |
| `f75457c8a650` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 37 start |
| `9fdaf811ef75` | build(booking): add npm lockfile for CI |
| `25ce4cc28f19` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 36 |
| `e4efb5512d35` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 36 start |
| `731cae64088b` | fix(booking): harden CI and container scaffold |
| `a06bbfac10d6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 35 |
| `724289d0e041` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 35 start |
| `8a8f17797d29` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 34 |
| `ac3a02d527c1` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 34 start |
| `b4653011bfb9` | chore(booking): sync deployment scaffold with develop |
| `1259fe7e11b6` | chore(booking): sync domain scaffold with develop |
| `deafe939c94e` | chore(booking): merge PR #2 - add standalone runtime scaffold |
| `7d8492ae2836` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 33 start |
| `863b2aa8ca95` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 33 |
| `5873ab551d1e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 32 |
| `e8d6e4aeabfe` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 32 start |
| `117e8b4da6bc` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 31 start |
| `8984f6c3f858` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 31 |
| `506839d687fd` | fix(booking): align scaffold with booking API contract |
| `a558909c1d06` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 30 |
| `68768a2efa43` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 30 start |
| `9af57e65703c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 29 |
| `ebc1e8ab7acf` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 29 start |
| `a30bb0856231` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 28 |
| `2ee6157260b0` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 28 start |
| `02b96c414880` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 27 |
| `1eff18c4d441` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 27 start |
| `7b0410603016` | test(booking): validate standalone runtime scaffold |
| `11fdabf4e5d6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 26 |
| `7a7869ec387b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 26 start |
| `7f267149d3da` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 25 |
| `a8f9190e5825` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 25 start |
| `69467a8ef424` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 24 |
| `99a531401a48` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 24 start |
| `492cacba9564` | chore(booking-portal): add booking domain scaffold |
| `6bd2ce625c45` | chore(booking-portal): add CI and deployment scaffold |
| `689d7ebc2fc3` | chore(booking-portal): add standalone runtime scaffold |
| `84d21f8790d5` | chore(booking): merge PR #1 - add Angular federation scaffold |
| `28909ea3ab83` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 23 |
| `43a3319769c8` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 23 start |
| `bb167af824ac` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 22 start |
| `c30ab50d6f3d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 22 |
| `04d4f4d77abe` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 21 |
| `ad9bb3b1444d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 21 start |
| `d861f726db9a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 20 |
| `7dce77b5e2ad` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 20 start |
| `d5926deba77d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 19 |
| `b99b1bc3495c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 19 start |
| `7b442aa8ed1d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 18 |
| `d8bf36a03b02` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 18 start |
| `ea6417f2ac91` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 17 |
| `7f73f258ae7c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 17 start |
| `6915c6bc4353` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 16 |
| `bc462b1107f5` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 16 start |
| `e20341d33b76` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 15 |
| `6fc43769589b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 15 start |
| `0d30666d9869` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 14 |
| `498ea8b729ba` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 14 start |
| `096e73bb5bad` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 13 |
| `773c087e8f4d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 13 start |
| `037ee00b22d9` | chore(booking-portal): add Angular federation scaffold configuration |
| `ad97a95f0da7` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 12 |
| `2cebfb25749a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 12 start |
| `c8fda51aa787` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 11 |
| `47f2398b4627` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 11 start |
| `92f4eff5120a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 10 |
| `ca4699a0615b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 10 start |
| `1530d1143dcb` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 9 start |
| `92e4b2211bc5` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 9 |
| `cb711b77755a` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 8 |
| `048a8795300e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 8 start |
| `b6c023b93669` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 7 start |
| `f2a7ee0b37ad` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 7 |
| `16767926348e` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 6 |
| `c01d8e07c1ac` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 6 start |
| `0e49fb981746` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 5 |
| `9397e350ec8d` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 5 start |
| `3a0a1c440831` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 4 start |
| `cdcbe361da2c` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 4 |
| `0dcf18d8b91b` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 3 start |
| `574703a50446` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 3 |
| `4803ce1eaf88` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 2 |
| `a70b3f9d24b1` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 2 start |
| `47d534b83ed6` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 1 |
| `58c098052a75` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - turn 1 start |
| `f332d5d6ea76` | Agent host session 3a3ef06f-80ea-4399-8a4e-1288ccb7df3c - baseline checkpoint |

### csp-auth-api

`code-corhuila/csp-auth-api` - 8 commits

| Commit | Message |
|---|---|
| `0e5f2006bf41` | docs(readme): describe the used and reserved sections of the env example |
| `babb03c4c366` | chore(governance): split env example into used and reserved variables |
| `8ab419b1ab62` | docs(readme): describe purpose, rules and commands of the service |
| `cd76129e5426` | feat(domain): add User, Email and Role with typed errors |
| `dbc9dbe53844` | chore(ci): merge scaffold into ci-deploy |
| `6a8210ebdf39` | chore(ci): add CI workflow, Dockerfile and compose file |
| `59d7d87df1fb` | chore(scaffold): add Go module, composition root and hexagonal packages |
| `44e5ce32c5eb` | chore(governance): add gitignore, env example and PR template |

### csp-auth-db

`code-corhuila/csp-auth-db` - 9 commits

| Commit | Message |
|---|---|
| `a736c6821376` | chore(schema): merge schema into ci |
| `b7021058076e` | fix(schema): let auth_writer purge old rows of the outbox |
| `f8363a9f967a` | docs(readme): describe the auth database repository |
| `c34d6951f7b2` | ci(db): add schema rebuild check |
| `698cc1777bd4` | chore(deploy): add flyway migration executor |
| `942455cc1976` | chore(schema): add auth schema, tables, constraints, seeds and roles |
| `7782705dea7e` | chore(governance): add gitignore, env example and pr template |
| `0f3d6e4116bd` | chore(governance): update project affiliation and doc links to cinesync platform |
| `3b28c53f9ab4` | chore(governance): update project affiliation and doc links to cinesync |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 6 commits

| Commit | Message |
|---|---|
| `43173b00f0f4` | chore(deps): add package-lock.json |
| `4ebee9a210e1` | chore(ci): add CI workflow and container deployment files |
| `ad6cf5d689e8` | chore(domain): add auth routes, typed models, synthetic data service and page placeholders |
| `54da5969e446` | chore(runtime): add federated bootstrap, root component and shell contract |
| `0e9d1900cc55` | chore(config): add Angular 21 and Native Federation project configuration |
| `b45bb0b812eb` | chore(governance): align README, CODEOWNERS and repo hygiene files to Cinesync Platform |

### csp-catalog-api

`code-corhuila/csp-catalog-api` - 8 commits

| Commit | Message |
|---|---|
| `32d7a23cef41` | chore(catalog): merge PR #3 - add CI, deployment files and complete the README |
| `cba672c31305` | chore(catalog): sync chore/catalog-api-ci-deploy-readme with develop |
| `278f3d6eab53` | chore(catalog): merge PR #2 - add the three-module Maven scaffold |
| `9d8b08c92d64` | chore(governance): merge PR #1 - add gitignore, env example and PR template |
| `c226832e0e77` | chore(governance): point the CODEOWNERS comment to csp-docs |
| `34f0744b5ba5` | chore(catalog): add CI, deployment files and complete the README |
| `5181be69b9bd` | chore(catalog): add the three-module Maven scaffold |
| `88f75b23862d` | chore(governance): add gitignore, env example and PR template |

### csp-catalog-db

`code-corhuila/csp-catalog-db` - 12 commits

| Commit | Message |
|---|---|
| `14870e3bdd26` | chore(catalog): merge PR #3 - add the database rebuild workflow and complete the README |
| `c910543645a2` | chore(catalog): sync chore/catalog-db-ci-readme with develop |
| `e2cf26b6bb9b` | chore(catalog): merge PR #2 - add the Liquibase changelog skeleton and the executor |
| `1377540b9870` | chore(governance): merge PR #1 - add gitignore, env example and PR template |
| `2151810d0c55` | fix(catalog): drop the stray test output and document the domain user |
| `b3bf10c30884` | chore(catalog): sync chore/catalog-db-ci-readme with chore/catalog-db-changelog-executor |
| `7bf9a01786e9` | fix(catalog): run the executor as the catalog domain user |
| `3ca62690f1e4` | fix(catalog): use the domain database user in the env example |
| `964ed5a1ee7b` | chore(governance): point the CODEOWNERS comment to csp-docs |
| `0c55c8140273` | chore(catalog): add the database rebuild workflow and complete the README |
| `7c7a8a548635` | chore(catalog): add the Liquibase changelog skeleton and the executor |
| `ebc4431eb9e0` | chore(governance): add gitignore, env example and PR template |

### csp-catalog-portal

`code-corhuila/csp-catalog-portal` - 31 commits

| Commit | Message |
|---|---|
| `6075606ed61a` | chore(catalog): merge PR #7 - align the portal with the Annex H template |
| `510513f275fe` | chore(catalog): align the portal configuration with the Annex H template |
| `c477699c8d47` | feat(catalog): add retry to the billboard error state |
| `8a8e0ebefab3` | chore(catalog): merge PR #5 - add CI and deployment scaffold |
| `3f0278f5dc1d` | chore(catalog): synchronize deployment branch with develop |
| `3fd3b4df7a92` | fix(catalog): harden container deployment |
| `e25c8807da70` | feat(catalog): merge PR #4 - add typed billboard domain scaffold |
| `47a94c8b50ed` | chore(catalog): merge PR #3 - add standalone runtime scaffold |
| `84d5d7c0ddad` | chore(catalog): merge PR #2 - add Angular federation configuration |
| `ba489a72def4` | chore(catalog): merge PR #1 - update governance affiliation |
| `c84b4b9224c6` | chore(catalog): synchronize final scaffold configuration |
| `d0ea52c6e3ee` | chore(catalog): synchronize TypeScript configuration |
| `10827e214025` | fix(catalog): remove unused federation path alias |
| `a7586aa0501a` | fix(catalog): run CI on every pull request |
| `67ac50ec1ead` | fix(catalog): complete CI and deployment scaffold |
| `991acd3bb399` | chore(catalog): synchronize deployment with domain scaffold |
| `fe0cf5a8c4fd` | chore(catalog): synchronize domain configuration |
| `a6493ee2fb52` | chore(catalog): synchronize runtime configuration |
| `19820d11181e` | fix(catalog): align TypeScript path alias |
| `abdfa14342d5` | chore(catalog): keep generated lockfile in ci scaffold |
| `d3d07f1bd5c7` | feat(catalog): add typed billboard adapter |
| `e08b15b61ca5` | chore(catalog): synchronize domain with runtime scaffold |
| `caa9bc7e2a4c` | test(catalog): cover shell error contract |
| `e75ee4428485` | chore(catalog): synchronize runtime with config scaffold |
| `2c0815bda199` | fix(catalog): align portal identity and encoding |
| `219c16fd2545` | chore(catalog): fix CI and deploy scaffold structure |
| `c7c3b0f01ca4` | chore(catalog): add CI and deployment scaffold |
| `cd5afd58feaa` | chore(catalog): add catalog domain scaffold |
| `4e5aa15e5928` | chore(catalog): add standalone runtime scaffold |
| `f382ac33ac66` | chore(catalog): add Angular federation scaffold |
| `dca293b61a5d` | chore(governance): update project affiliation and doc links to cinesync platform |
