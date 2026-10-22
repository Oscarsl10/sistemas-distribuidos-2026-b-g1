# Git Commit Log - CineSync Platform

## Week 12

Period: 2026-10-19 to 2026-10-22. Author: Oscar Sierra (`Oscarsl10`). Only commits authored by me are listed, including those in the shared repositories (`csp-front`, `csp-api-gateway`, `csp-worker`, `csp-infra-*`). Same change promoted across several branches is listed once; automatic merge commits are omitted, and changes already reported in weeks 9, 10 and 11 are not repeated. Newest first.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 8 |
| csp-auth-api | 7 |
| csp-auth-db | 3 |
| csp-auth-portal | 2 |
| csp-front | 2 |
| csp-infra-postgres | 3 |
| csp-infra-rabbitmq | 6 |
| csp-worker | 6 |
| **Total** | **37** |

### csp-docs

`code-corhuila/csp-docs` - 8 commits

| Commit | Message |
|---|---|
| `8540dbedc505` | docs(architecture): accept adr-027 and adr-036 |
| `ff6f0c0955e2` | docs(auth): align adr-040 with the implemented login attempt limit |
| `369b731fdd08` | docs(auth): record the login attempt limit by client ip and align the auth docs (hu-auth-005) |
| `cf69d33df778` | docs(auth): record the smtp relay provider and its limits in the runbook and deployment |
| `2239b22a031e` | docs(infra): amend adr-036 with the permissions verified against a real broker |
| `69053db225cf` | docs(auth): contract 1.6.1, adr-030 amendment for email_verification and the profile route in the shell |
| `d9fe21d49ba3` | docs(infra): add adr-036 for the rabbitmq instance csp-infra-rabbitmq |
| `525313b32f43` | docs(adr): amend ADR-026 so the CORS regex is anchored and the mechanism stays a template per repository |

### csp-auth-api

`code-corhuila/csp-auth-api` - 7 commits

| Commit | Message |
|---|---|
| `89a4817a22e6` | feat(auth): answer 429 rate_limited on login after repeated failed attempts |
| `b2f2b437c767` | feat(auth): add the redis adapter of the login attempt limit |
| `77c27f3078c8` | feat(auth): add the domain and use case of the login attempt limit |
| `2b049204d0d1` | docs(purge): document the email verification retention |
| `910a5e3a126f` | feat(purge): delete email verification requests expired more than 7 days ago |
| `c158fc79294b` | docs(auth): document where the smtp variables go per environment |
| `a928435a7900` | fix(auth): answer the exact remaining seconds in the retry-after of the email change limit |

### csp-auth-db

`code-corhuila/csp-auth-db` - 3 commits

| Commit | Message |
|---|---|
| `6b352d72ad0d` | test(auth-db): check the email verification purge grant and index |
| `37c643721f99` | chore(auth-db): add complete expires_at index on email_verification |
| `fc87476378ce` | chore(auth-db): grant delete on email_verification to auth_writer |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 2 commits

| Commit | Message |
|---|---|
| `a507c4c78fb5` | feat(auth): explain the failed-login rate limit on the login form |
| `ac68d342f911` | fix(auth-portal): refuse an unanchored CORS regex and serve only remoteEntry.json as no-store json |

### csp-front

`code-corhuila/csp-front` - 2 commits

| Commit | Message |
|---|---|
| `df75283bae1c` | feat(front): mount the profile page at /profile and link it from the navigation |
| `f28a26ab488a` | fix(front): refuse an unanchored CORS regex and check the CORS rule of the image in the CI |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 3 commits

| Commit | Message |
|---|---|
| `00d808a0e3b7` | docs(infra): document the steps of an existing environment for the migrators |
| `cbc9e4f9adf0` | feat(infra): create the auth_migrator login that owns the auth schema |
| `e540807aed93` | feat(infra): create the booking_migrator login that owns the booking schema |

### csp-infra-rabbitmq

`code-corhuila/csp-infra-rabbitmq` - 6 commits

| Commit | Message |
|---|---|
| `e635297d16d8` | docs: describe the instance, its logins and its limits |
| `fe0d877b8070` | docs: add the changelog |
| `ee98ed91c578` | chore: keep every file with LF endings |
| `f03c98f8b4fe` | ci: start the instance as dev and qa and run the checks |
| `d884d7b6043a` | test: check the exchanges, logins, permissions and limits against a real broker |
| `817d1e68c39e` | chore: add the rabbitmq instance with its definitions and scripts |

### csp-worker

`code-corhuila/csp-worker` - 6 commits

| Commit | Message |
|---|---|
| `35fa7db7794a` | fix(worker): make the auth relay end-to-end test independent of the publication order |
| `a0947ca65e37` | chore(worker): pass the auth relay variables in the compose file |
| `8976a442efc9` | docs(worker): document the auth outbox relay and its flag |
| `bf3de6a6d096` | refactor(worker): parameterize the outbox relay adapters by domain |
| `3a32a98dc49d` | feat(worker): relay the auth outbox behind its own flag |
| `ee8660919f3d` | feat(worker): take an instance lock before each run of the outbox relay |
