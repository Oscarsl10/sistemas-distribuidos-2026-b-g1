# Git Commit Log - CineSync Platform

## Week 14

Period: 2026-11-02 to 2026-11-07. Author: Oscar Sierra (`Oscarsl10`). Only commits authored by me are listed, including those in the shared repositories (`csp-api-gateway`, `csp-front`, `csp-infra-*`, `csp-workflow`, `csp-worker`) and in my fork. Same change promoted across several branches is listed once; automatic merge commits are omitted. Newest first. The list was taken on 2026-11-07.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 20 |
| csp-auth-api | 4 |
| csp-auth-db | 2 |
| csp-auth-portal | 4 |
| csp-api-gateway | 6 |
| csp-booking-api | 4 |
| csp-booking-portal | 0 |
| csp-front | 1 |
| csp-infra-postgres | 11 |
| csp-infra-rabbitmq | 10 |
| csp-infra-rustfs | 2 |
| csp-ticketing-api | 16 |
| csp-ticketing-db | 6 |
| csp-ticketing-portal | 6 |
| csp-worker | 4 |
| csp-workflow | 31 |
| sistemas-distribuidos-2026-b-g1 (my fork) | 11 |
| **Total** | **138** |

### csp-docs

`code-corhuila/csp-docs` - 20 commits

| Commit | Date | Message |
|---|---|---|
| `649fec7` | 2026-11-07 | docs(booking): specify administrator reservation listing |
| `ecf3e97` | 2026-11-07 | docs(project-control): answer auth event topology question |
| `7acea36` | 2026-11-06 | docs(ticketing): write the stories hu-ticketing-006, hu-fe-ticketing-002 and hu-fe-ticketing-003 |
| `ccd8bc8` | 2026-11-06 | docs(ticketing): amend adr-068 for event order, message text and hold length |
| `278a8b4` | 2026-11-06 | docs(ticketing): design the hold expiry reminder of hu-ticketing-005 (adr-068) |
| `1999096` | 2026-11-06 | docs(ticketing): #348 correct showtimestartsat and describe get me tickets as implemented |
| `9925b28` | 2026-11-06 | docs(ticketing): #346 document the notifications and the consumer of reservationexpired |
| `4a1c27a` | 2026-11-06 | docs(adr): accept tls between booking and catalog (ADR-067) |
| `5358cf6` | 2026-11-06 | docs(governance): add oscar as co-owner of the ticketing domain |
| `754f437` | 2026-11-06 | docs(adr): date the acceptance of the workflow adr on the day it is merged |
| `2353bdc` | 2026-11-06 | docs(project-control): close td-040 now that the hold keeps the start time |
| `5cda25e` | 2026-11-05 | docs(adr): accept the adr of the workflow and its booking operations |
| `680ed63` | 2026-11-05 | docs(project-control): record that the hold takes its snapshot from catalog |
| `4e22141` | 2026-11-05 | docs(governance): add oscar as owner of the booking domain |
| `619742a` | 2026-11-05 | docs(workflow): keep the steps of a saga in its row in adr-065 |
| `fbb74db` | 2026-11-05 | docs(workflow): propose the book-seats saga and the contract of the sagas api |
| `bcd6f58` | 2026-11-05 | docs(workflow): propose postgresql schema workflow for the state of the sagas |
| `f4eacbb` | 2026-11-05 | docs(workflow): propose java and spring boot as the stack of csp-workflow |
| `53b23d1` | 2026-11-03 | docs(auth): answer the registration the same for every email and confirm it by link |
| `8165623` | 2026-11-03 | docs(project-control): record the state of the broker and the ticketing queues in q-001 |

### csp-auth-api

`code-corhuila/csp-auth-api` - 4 commits

| Commit | Date | Message |
|---|---|---|
| `a998029` | 2026-11-03 | chore(auth): remove the register user use case and document the registration flow |
| `01df7e5` | 2026-11-03 | feat(auth): answer the registration the same for every email and confirm it by link |
| `46d2ce5` | 2026-11-03 | feat(auth): add the confirm registration use case and its storage |
| `321e414` | 2026-11-03 | feat(auth): add the request registration use case |

### csp-auth-db

`code-corhuila/csp-auth-db` - 2 commits

| Commit | Date | Message |
|---|---|---|
| `d324aa1` | 2026-11-02 | feat(db): add registration request table for email confirmation |
| `fb1804a` | 2026-11-02 | fix(deploy): retry the connection of the migration executor while the instance starts |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 4 commits

| Commit | Date | Message |
|---|---|---|
| `79399dc` | 2026-11-03 | feat(auth): confirm the registration by link and show a neutral message |
| `eb5478b` | 2026-11-02 | chore(auth-portal): leave a guarded screen when the api refuses the session |
| `e61c758` | 2026-11-02 | chore(auth-portal): keep the sign-in and session copy in one typed file |
| `76df75c` | 2026-11-02 | chore(auth-portal): move the sign-in rules out of the login and register pages |

### csp-api-gateway

`code-corhuila/csp-api-gateway` - 6 commits

| Commit | Date | Message |
|---|---|---|
| `d38b917` | 2026-11-07 | fix(gateway): route get /me/notifications to csp-ticketing-api |
| `f711abd` | 2026-11-06 | feat(gateway): route the client operations of ticketing |
| `81a9f7b` | 2026-11-06 | feat(gateway): reach catalog over verified tls (HU-289) |
| `5a1cf62` | 2026-11-05 | feat(gateway): route the sagas of the workflow (HU-CONCESSIONS-004) |
| `67be52a` | 2026-11-05 | test(gateway): specify the route of the sagas of the workflow |
| `601ef0b` | 2026-11-03 | feat(gateway): route the confirmation of the registration to the auth service |

### csp-booking-api

`code-corhuila/csp-booking-api` - 4 commits

| Commit | Date | Message |
|---|---|---|
| `727d9b5` | 2026-11-06 | feat(booking): #88 add expiresAt to the reservation held event payload |
| `0d99dbe` | 2026-11-06 | feat(booking): trust only the internal ca when calling catalog |
| `ba90ec0` | 2026-11-06 | fix(hold): keep the start time of the showtime in the reservation and its confirmed event |
| `3b5fe14` | 2026-11-05 | fix(hold): take the title, the room and the amount of a hold from catalog |

### csp-booking-portal

`code-corhuila/csp-booking-portal` - 0 commits

| Commit | Date | Message |
|---|---|---|

### csp-front

`code-corhuila/csp-front` - 1 commits

| Commit | Date | Message |
|---|---|---|
| `62ececc` | 2026-11-07 | fix(front): prevent stale shell after deployment |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 11 commits

| Commit | Date | Message |
|---|---|---|
| `380dd3e` | 2026-11-07 | chore(infra): enable auth relay in development example |
| `eabac11` | 2026-11-07 | chore(env): list the worker variables missing from the env examples |
| `3b76821` | 2026-11-07 | fix(infra): run every migration executor from up.sh |
| `7da8925` | 2026-11-06 | chore(infra): #58 compose the ticket pdf storage and the mail catcher |
| `b80bf01` | 2026-11-06 | chore(infra): #41 compose ticketing in the platform |
| `aa3ee36` | 2026-11-06 | fix(infra): start csp-workflow after its migration creates workflow.saga |
| `4f6f5e8` | 2026-11-06 | feat(infra): add internal CA and Catalog TLS certificate to the platform |
| `2da7583` | 2026-11-05 | feat(infra): compose csp-workflow in the platform |
| `aa21fbb` | 2026-11-03 | feat(env): add the page that the registration link opens |
| `0cc2bc2` | 2026-11-03 | chore(infra): compose the rabbitmq broker in the platform |
| `3c232b2` | 2026-11-02 | chore(infra): add the password reset url to the environment examples |

### csp-infra-rabbitmq

`code-corhuila/csp-infra-rabbitmq` - 10 commits

| Commit | Date | Message |
|---|---|---|
| `1ed0f50` | 2026-11-06 | fix(rabbitmq): run the health check as the rabbitmq user so the erlang cookie is not created by root |
| `757ef8a` | 2026-11-03 | docs(changelog): record the 3.0.0 release |
| `21ac6bc` | 2026-10-23 | test(rabbitmq): validate the vhost of the AMQP URLs of the env examples |
| `f550050` | 2026-10-23 | fix(rabbitmq): end the documented AMQP URLs in the default vhost as %2F |
| `ce2557b` | 2026-10-20 | docs: describe the instance, its logins and its limits |
| `feb86c2` | 2026-10-20 | docs: add the changelog |
| `c70e192` | 2026-10-20 | chore: keep every file with LF endings |
| `3004d8b` | 2026-10-20 | ci: start the instance as dev and qa and run the checks |
| `540cf21` | 2026-10-20 | test: check the exchanges, logins, permissions and limits against a real broker |
| `22dba17` | 2026-10-20 | chore: add the rabbitmq instance with its definitions and scripts |

### csp-infra-rustfs

`code-corhuila/csp-infra-rustfs` - 2 commits

| Commit | Date | Message |
|---|---|---|
| `1b9d54e` | 2026-11-07 | docs(changelog): record the 3.0.0 release |
| `a932e02` | 2026-11-06 | docs(storage): #5 say how the platform composes the storage |

### csp-ticketing-api

`code-corhuila/csp-ticketing-api` - 16 commits

| Commit | Date | Message |
|---|---|---|
| `e5a1bba` | 2026-11-07 | docs(changelog): record the 3.0.0 release |
| `063c423` | 2026-11-07 | fix(api): make the delivery test independent of the database clock |
| `489697c` | 2026-11-07 | fix(api): skip a delivery that cannot be processed and keep dispatching |
| `f16b5a0` | 2026-11-06 | feat(api): #70 create the hold expiry warning notifications |
| `49450d4` | 2026-11-06 | feat(api): #70 keep the held reservations that expire |
| `56a8a9d` | 2026-11-06 | chore(api): #71 make smtp auth and starttls configurable for the dev catcher |
| `365f617` | 2026-11-06 | docs(api): #68 say why a ticket can have no showtime start |
| `7eea93c` | 2026-11-06 | feat(api): #68 list the tickets of the user |
| `aa550c9` | 2026-11-06 | feat(api): #66 notify the user when a reservation expires |
| `0ec8fe4` | 2026-11-06 | feat(api): list the notifications of the user and create them from booking and ticket events |
| `bc21570` | 2026-11-06 | fix(api): accept a pdf reference only on the object storage of the platform |
| `1d9f276` | 2026-11-06 | chore(api): add the ci pipeline, board tracking and pull request template |
| `15035d7` | 2026-11-06 | chore(api): #291 add the negative authorization test of the token |
| `f349e34` | 2026-11-06 | chore(api): read domain variables in the platform composition |
| `5b06e0a` | 2026-11-06 | chore(api): #54 make the api runnable in the platform composition |
| `805a577` | 2026-11-06 | chore(api): #54 stop tracking the target folders |

### csp-ticketing-db

`code-corhuila/csp-ticketing-db` - 6 commits

| Commit | Date | Message |
|---|---|---|
| `9658acf` | 2026-11-07 | docs(changelog): record the 3.0.0 release |
| `cceb4d8` | 2026-11-06 | feat(db): #20 add the held_reservation table for the hold expiry reminder |
| `ce464fc` | 2026-11-06 | feat(db): add the notification table for hu-ticketing-005 |
| `759727a` | 2026-11-06 | fix(db): grant the writer role to the login of the service |
| `516fcdc` | 2026-11-06 | chore(db): #12 add the migration executor to the platform composition |
| `231c749` | 2026-11-06 | fix(db): create the schema and its default privileges as the user that migrates |

### csp-ticketing-portal

`code-corhuila/csp-ticketing-portal` - 6 commits

| Commit | Date | Message |
|---|---|---|
| `eb0d1d3` | 2026-11-07 | docs(changelog): record the 3.0.0 release |
| `d65b32b` | 2026-11-07 | feat(notifications): add notifications page for HU-FE-TICKETING-003 |
| `bf41288` | 2026-11-07 | feat(portal): show the real tickets of the client on the my tickets page |
| `9ea8d55` | 2026-11-06 | fix(portal): validate the shape of the shell error and use forward slashes in tsconfig.federation.json |
| `64e007f` | 2026-11-06 | chore(portal): serve on port 80 and align the container with the other portals |
| `6682c86` | 2026-11-06 | fix(portal): show the user-visible text in spanish |

### csp-worker

`code-corhuila/csp-worker` - 4 commits

| Commit | Date | Message |
|---|---|---|
| `a69d1bb` | 2026-11-07 | fix(storage): give the pdf job its own writer storage variables |
| `ae35192` | 2026-11-06 | feat(worker): #68 add the hold-expiry-warnings job |
| `1dde2d4` | 2026-11-03 | docs(worker): name the source of the ticketing relay in the replay commands |
| `bd3a53f` | 2026-11-03 | docs(worker): document the start-up order of the relays and the replay of failed events |

### csp-workflow

`code-corhuila/csp-workflow` - 31 commits

| Commit | Date | Message |
|---|---|---|
| `d8b0b7e` | 2026-11-06 | docs(changelog): record the 3.0.0 release |
| `6070f87` | 2026-11-06 | test(e2e): check the original correlation id after a restart and the start time of the showtime |
| `4455463` | 2026-11-06 | docs(saga): record the stored correlation id in the readme and the changelog |
| `203488c` | 2026-11-06 | feat(saga): keep the correlation id of the request with the saga so a restart resumes with it |
| `695c14b` | 2026-11-05 | test(saga): specify that the saga keeps the correlation id of its request |
| `903e001` | 2026-11-06 | test(e2e): hold the confirms for less than the gateway timeout and document the failed compensation |
| `4a5cf38` | 2026-11-05 | test(e2e): prove a failing compensation leaves the saga failed |
| `72d1806` | 2026-11-05 | docs(e2e): document the booking ref of the end-to-end test |
| `6756170` | 2026-11-05 | test(e2e): tamper the signature deterministically and make the probe repeatable |
| `6dcf691` | 2026-11-05 | test(e2e): run against a booking ref and check the snapshot comes from catalog |
| `f6511e2` | 2026-11-05 | chore(e2e): keep LF in the shell scripts |
| `5c3b30e` | 2026-11-05 | docs(e2e): explain how to run the end-to-end test and record it in the changelog |
| `769b195` | 2026-11-05 | test(e2e): prove the saga book-seats against the real booking and catalog |
| `9ac07a6` | 2026-11-05 | fix(deploy): read the variables of the platform environment and grant the writer role to the service login |
| `577e048` | 2026-11-05 | test(state): specify the grant of the writer role to the login of the service |
| `d3235e0` | 2026-11-05 | feat(saga): recover the sagas a restart cut in the middle when the service starts |
| `6e29248` | 2026-11-05 | test(saga): specify the recovery of the sagas a restart cut in the middle |
| `e222c87` | 2026-11-05 | feat(state): add the postgresql saga store with its migrations and the unique key race |
| `d456c7a` | 2026-11-05 | test(state): specify the postgresql saga store, its migrations and the race of two starts |
| `8c5c241` | 2026-11-05 | feat(saga): keep the request and the version in the saga and make it restorable |
| `d3e67d5` | 2026-11-05 | test(saga): specify the request, the version and the snapshot of a saga |
| `4645ba9` | 2026-11-05 | feat(http): add the saga endpoints and compose the service |
| `cf064ab` | 2026-11-05 | test(http): specify the saga endpoints, the request validation and the end to end flow |
| `c7e3635` | 2026-11-05 | feat(participants): add the booking client and the service token |
| `124c542` | 2026-11-05 | test(participants): specify the booking client and the service token |
| `b4c5423` | 2026-11-05 | feat(http): add the token verifier, the filters and the error envelope |
| `b8ee12e` | 2026-11-05 | test(http): specify the token verifier, the filters and the error envelope |
| `3eb81df` | 2026-11-05 | feat(saga): add the saga domain and the book-seats orchestrator |
| `0ac50fb` | 2026-11-05 | test(saga): specify the saga domain and the book-seats orchestrator |
| `be9e249` | 2026-11-05 | chore(scaffold): add the base structure of the workflow service |
| `8dd13ea` | 2026-11-05 | chore(governance): replace the template identity and add the changelog |

### sistemas-distribuidos-2026-b-g1 (my fork)

`Oscarsl10/sistemas-distribuidos-2026-b-g1` - 11 commits (week 14 sessions and status)

| Commit | Date | Message |
|---|---|---|
| `7db030c` | 2026-11-07 | docs(hu-status): add the week 14 commit log and summary, and the approval of the registration contract |
| `497af35` | 2026-11-07 | docs(w14): record the approval of the registration contract in the session 2 plan |
| `cbf11d8` | 2026-11-07 | docs(hu-status): update week 14 with the work of the whole week |
| `9a66220` | 2026-11-07 | docs(w14): add the session 2 plan for the pipeline and branch protection |
| `6004b87` | 2026-11-07 | docs(w14): add the session 1 report on CI/CD and quality gates |
| `73749ab` | 2026-11-07 | test(w14): add the session 1 quality gates rehearsal |
| `e66ac52` | 2026-11-07 | docs(certification): split chapter 7 study hours across two days |
| `d187432` | 2026-11-07 | docs(certification): add chapter 7 evidence and study log row |
| `db56e70` | 2026-11-05 | docs(certification): add chapter 6 evidence and study log row |
| `f19c2e7` | 2026-11-04 | docs(certification): add chapter 5 evidence and study log row |
| `b16d2f5` | 2026-11-03 | docs(certification): add chapter 4 evidence and study log row |
