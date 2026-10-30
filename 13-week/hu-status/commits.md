# Git Commit Log - CineSync Platform

## Week 13

Period: 2026-10-26 to 2026-10-30. Author: Oscar Sierra (`Oscarsl10`). Only commits authored by me are listed, including those in the shared repositories (`csp-api-gateway`, `csp-infra-*`) and in my fork. Same change promoted across several branches is listed once; automatic merge commits are omitted. Newest first. The list was taken on 2026-10-30; the commit that publishes this file is not in it.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 12 |
| csp-auth-api | 26 |
| csp-auth-db | 2 |
| csp-auth-portal | 1 |
| csp-api-gateway | 1 |
| csp-infra-postgres | 8 |
| sistemas-distribuidos-2026-b-g1 (my fork) | 14 |
| **Total** | **64** |

### csp-docs

`code-corhuila/csp-docs` - 12 commits

| Commit | Message |
|---|---|
| `b5c7dff0547c` | docs(auth): add ADR-061 on the auth http metrics and the latency alert |
| `c0e93bc8f1ef` | docs(auth): add ADR-060 on the credential rate limit and the bound on concurrent password hashes |
| `5525f75f5797` | docs(auth): add ADR-059 on the login slowdown per account and the audit of the lock |
| `8d2e1564dbeb` | docs(auth): mark the claims of ADR-052 as superseded by ADR-056 |
| `f7a1e71a45a5` | docs(auth): require issuer and audience on service credentials |
| `ea0e718e0a73` | docs(auth): accept ADR-055 and use TOO_MANY_REQUESTS for the 429 |
| `f9fdc94d2c38` | docs(auth): specify users newest first, email change idempotency key and audit contract fixes |
| `820506b6253d` | docs(auth): require idempotency-key on forgot-password |
| `1bf6fd96cc26` | docs(security): add book pages to the references of the STRIDE |
| `b7bbc56db930` | docs(security): address review on STRIDE for csp-api-gateway to csp-auth-api |
| `b38bdb1d2616` | docs(security): complete STRIDE and risk ranking for csp-api-gateway to csp-auth-api |
| `68a9fb88fad9` | docs(security): STRIDE context and first threats for csp-api-gateway to csp-auth-api |

### csp-auth-api

`code-corhuila/csp-auth-api` - 26 commits

| Commit | Message |
|---|---|
| `975b3aad2222` | test(auth): keep two timing tests from failing on a loaded machine |
| `732b1ce75f39` | feat(auth): export the latency and the outcome of every route on /metrics |
| `cfdab917f932` | feat(auth): bound the concurrent password hashes of the service |
| `428ad1ab1bb1` | feat(auth): slow the failed logins of an email without locking the account |
| `347f4b742a15` | feat(auth): record every lock of the login as an audit event |
| `76c070c46c4b` | docs(auth): say what the mail relay does when the provider is unreachable or slow |
| `0cf76ffae1f2` | fix(auth): give the mail-relay variables defaults so compose loads without them |
| `bd927ff90639` | feat(auth): add a mail-relay profile that relays Mailpit to a real provider |
| `5af79f150db5` | fix(auth): expect 401 for a user's access token in the purge test |
| `6613111646d8` | fix(auth): sign the service token of the purge test with issuer and audience |
| `c23ebac9fd2b` | fix(auth): require issuer and audience on service credentials |
| `9640bf054ad0` | fix(auth): answer 429 with the code TOO_MANY_REQUESTS of the norm |
| `db51e643d7d2` | fix(auth): log structured JSON with the trace id and one access line per request |
| `6bcc07239cb7` | feat(auth): store the idempotency key of the email change request |
| `e112b2789622` | fix(auth): declare the request limits and bound register and login |
| `d66b41645735` | fix(auth): list users from newest to oldest |
| `c3de53e63ee3` | docs(auth): record the delivery test fix in the changelog |
| `1e8abb93f839` | test(auth): make the password reset sender delivery test independent of the speed of the machine |
| `707606a7b61e` | test(auth): check invalid user tokens over http on the protected routes |
| `2147db1b0410` | test(auth): expect the redirect status the mux sends for a path to clean |
| `4a11c03a7a1d` | fix(auth): add a root wrapper with correlation id, enveloped 404 and 405 and panic recovery |
| `1a754d3cd5c7` | feat(auth): require an idempotency key on forgot-password |
| `4007ebaa145e` | feat(auth): add the offline command that issues the service tokens |
| `0d49dc17676d` | feat(auth): issue the concessions:inventory:read permission to ADMIN |
| `1bed06b0ae7c` | fix(auth): issue the access token for 15 minutes by default |
| `4573cf23f482` | fix(auth): add a caching rule to GET /jwks |

### csp-auth-db

`code-corhuila/csp-auth-db` - 2 commits

| Commit | Message |
|---|---|
| `30e3d8dbf356` | feat(auth-db): add email change idempotency key table |
| `5232eaeafe55` | feat(auth-db): add password reset idempotency key table |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 1 commits

| Commit | Message |
|---|---|
| `f4a4ea140398` | feat(auth-portal): send an idempotency key on forgot-password |

### csp-api-gateway

`code-corhuila/csp-api-gateway` - 1 commits

| Commit | Message |
|---|---|
| `49aa0f8beb3b` | feat(gateway): stricter rate limit for login, register and forgot password |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 8 commits

| Commit | Message |
|---|---|
| `e3307e12b600` | feat(infra): scrape csp-auth-api and alert on the latency of login and register |
| `afaf1d005ed8` | docs(infra): say what the mail relay does when the provider is unreachable or slow |
| `bf8b0bd1b6ba` | feat(infra): add the script and variables of the development mail relay |
| `d7b875c261bc` | chore(env): pass the SMTP relay of csp-auth-api from the platform env files |
| `5b964c9120ef` | docs(infra): document the steps of an existing environment for the migrators |
| `c3d66c06d6b9` | feat(infra): create the auth_migrator login that owns the auth schema |
| `f2e91b52c68d` | feat(infra): create the booking_migrator login that owns the booking schema |
| `7aa2b682cf91` | fix(infra): add issuer and audience to the dev service token |

### sistemas-distribuidos-2026-b-g1 (my fork)

`Oscarsl10/sistemas-distribuidos-2026-b-g1` - 14 commits (workshop and session activities of week 13)

| Commit | Message |
|---|---|
| `84cb9d8` | docs(w13): add the session 2 threat model and security testing plan |
| `f652109` | docs(w13): add the session 1 report on security in distributed systems |
| `520fefb` | test(w13): add the session 1 authentication and authorization rehearsal |
| `ff01038` | docs(w13): close threat modeling workshop |
| `81a7292` | test(w13): run negative authorization test and record output |
| `d3f80bd` | docs(w13): refine security backlog stories |
| `18d965a` | docs(w13): correct the book section and pages |
| `794dd0b` | docs(w13): update hu-status for threat modeling workshop |
| `879d19f` | docs(w13): document the security gates and how to run the test |
| `5f559a4` | test(w13): add the negative authorization test |
| `f14c995` | ci(w13): add security gates and negative authorization test |
| `eb628a0` | docs(w13): link security backlog stories |
| `e0176bd` | docs(w13): map attack surface and trust boundaries |
| `d28ab59` | docs(w13): start threat modeling workshop |
