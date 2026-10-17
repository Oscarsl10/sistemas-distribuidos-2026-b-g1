# Git Commit Log - CineSync Platform

## Week 11

Period: 2026-10-12 to 2026-10-17. Author: Oscar Sierra (`Oscarsl10`). Only commits authored by me are listed, including those in the shared repositories (`csp-front`, `csp-api-gateway`, `csp-worker`, `csp-infra-*`). Same change promoted across several branches is listed once; automatic merge commits are omitted, and changes already reported in weeks 9 and 10 are not repeated. Newest first.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 33 |
| csp-auth-api | 47 |
| csp-auth-db | 8 |
| csp-auth-portal | 14 |
| csp-front | 13 |
| csp-api-gateway | 4 |
| csp-infra-postgres | 7 |
| csp-infra-redis | 2 |
| csp-worker | 2 |
| **Total** | **130** |

### csp-docs

`code-corhuila/csp-docs` - 33 commits

| Commit | Message |
|---|---|
| `b1a1a4b5de88` | docs(auth): record the review decisions of the email change and the profile route (HU-AUTH-004) |
| `6f5bee7983c5` | docs(auth): the email-change code is sent by SMTP from auth and the log adapter is removed (HU-AUTH-004) |
| `6238c80bfcb4` | docs(auth): the real delivery of the email-change code is an SMTP adapter inside auth (HU-AUTH-004) |
| `27455a6c6ec5` | docs(auth): the unique index of the email change is V025 |
| `46841edf8b45` | docs(auth): ADR-034 for the email change and its code delivery, contract 1.6.0 and data model (HU-AUTH-004) |
| `51f9c69eb2ac` | docs(auth): update the auth service runbook for the roles, the logout and the web session |
| `05f756e806e0` | docs(requirements): mark HU-AUTH-001 and HU-AUTH-002 as Done |
| `9dc57319a2c4` | docs(auth): declare the 400 of logout and refine the Redis outage behaviour in ADR-033 |
| `042375760879` | docs(auth): specify logout as scenario 3 of HU-AUTH-003 and record ADR-033 |
| `ff49253a0ca3` | docs(auth): record the persistence of the web session in ADR-032 |
| `a251c1bec92f` | docs(ux): align the user roles screen text with the final Spanish copy |
| `715e088d7b5f` | docs(ux): fix two Spanish words of the user roles screen text |
| `0e776043e24f` | docs(ux): add the Spanish accents to the user roles screen text |
| `b72b3ef2e614` | docs(ux): write the user roles screen text in Spanish |
| `af540b961723` | docs(auth): reject revoking the only role of a user with USER_WITHOUT_ROLES |
| `162575849e38` | docs(changelog): record the user role management documentation |
| `9251450b7662` | docs(ux): add the user roles screen wireframe and the admin users route |
| `598607fc302e` | docs(auth): add UserRoleChanged and the last-admin rule to the domain and service docs |
| `777bb918459e` | docs(adr): add ADR-031 administrator management of user roles |
| `2f6028b29e87` | docs(auth): specify GET /users and the role change of the auth contract (1.5.0) |
| `f2e3df27abae` | docs(auth): declare the error responses of the contract and correct the worker facts |
| `e4bc0e355b39` | docs(adr): add ADR-030 auth retention and purge |
| `62e54a304c79` | docs(changelog): record ADR-029 |
| `bef100ed26a3` | docs(adr): add ADR-029 refresh token reuse detection |
| `f036e8a47387` | docs(changelog): record the promotional video |
| `94a9687bbe44` | docs(changelog): record the alignment of the auth and platform documents |
| `c02ba2d905ac` | docs(changelog): record ADR-028 and the alignment of the gateway documents |
| `3cbaeedd3319` | docs(auth): align the auth and platform documents with ADR-028 |
| `3ef78f8e1218` | docs(gateway): align the gateway documents with ADR-028 |
| `0591d6f9810f` | docs(adr): add ADR-028 Redis instance for the auth state |
| `0bb342d8754f` | docs(changelog): record the ADR-026 amendment of PR #122 |
| `b4a1eb0ac69c` | docs(changelog): add the changelog of the documentation |
| `fb21500b2566` | docs(adr): amend ADR-026 so the shell and the portals serve on port 80 and check health on 127.0.0.1 |

### csp-auth-api

`code-corhuila/csp-auth-api` - 47 commits

| Commit | Message |
|---|---|
| `c5ca92ab19c0` | fix(auth): take the account lock first when verifying and reject an unknown SMTP TLS mode (HU-AUTH-004) |
| `0a4ca289f15b` | feat(auth): send the email-change code by SMTP and remove the log adapter (HU-AUTH-004) |
| `9f00162e34f0` | feat(auth): endpoints to request and verify an email change (HU-AUTH-004) |
| `e8feb9759e9c` | feat(auth): adapters to store and deliver the email change (HU-AUTH-004) |
| `51aa7dcd0399` | feat(auth): domain and use cases to request and verify an email change (HU-AUTH-004) |
| `912d39ef94a1` | test(auth): expire the token of the end-to-end refresh test with a margin the clock skew cannot cancel |
| `d92649b9663d` | test(auth): make the persistence integration tests independent of the clock and the table statistics |
| `6346d0e67fb7` | fix(auth): revoke the refresh tokens on logout while Redis is down |
| `8d71ec238d2b` | feat(auth): reject revoked access tokens on the protected routes |
| `50bf6ccd92ba` | test(auth): check the 400 of logout against the contract |
| `0acbd6b724e1` | feat(auth): add POST /logout |
| `3b7706a26bb1` | feat(auth): add the revoked access token list in Redis |
| `cf7dc0da0527` | fix(auth): restart the service when the database is not ready at start |
| `2d13465e73fc` | fix(auth): copy go.sum in the container image build |
| `3ca40132b8cb` | test(auth): every test that creates an administrator takes the administrators lock |
| `a7706139bceb` | feat(auth): plan the user search per execution so it can use the email pattern index |
| `ed942c5e9569` | docs(auth): document the user list endpoint |
| `f150ccc2d049` | test(auth): contract, wiring and end-to-end tests of GET /users |
| `55152e951f71` | feat(auth): GET /users handler wired in the composition root |
| `6c03da6ab2f4` | feat(auth): read a page of users by email prefix from PostgreSQL |
| `2c855cfefe04` | feat(auth): ListUsers use case with the email prefix and page value objects |
| `fbe7f4fff40a` | test(auth): the persistence tests that create an administrator take the administrators lock |
| `bdd6f511f9f6` | test(auth): USER_WITHOUT_ROLES in unit, contract and end-to-end tests, docs |
| `e1180788ecb0` | feat(auth): answer 409 USER_WITHOUT_ROLES when a revoke would leave no role |
| `aa6cc36edf77` | docs(auth): document the role change endpoint |
| `6120c1635f13` | test(auth): end-to-end test of the role change and shared lock of the administrator tests |
| `9b4da2a2689a` | test(auth): contract tests of updateUserRole against auth-service.yaml 1.5.0 |
| `30d5fc36a236` | feat(auth): wire the role change route in the composition root |
| `0f5f474af846` | feat(auth): PATCH /users/{userId}/role handler and error mapping |
| `fbeac3699d07` | feat(auth): requirePermission guard answering 401 and 403 |
| `29556243e5e9` | feat(auth): AuthorizeUser use case |
| `3d13ab53e033` | feat(auth): verify the RS256 access token of a user |
| `4bee2b465a01` | test(auth): integration tests of the role change against PostgreSQL |
| `8767e7cd8153` | feat(auth): store role changes and the UserRoleChanged outbox row |
| `eb96d8518217` | feat(auth): ChangeUserRole use case and its ports |
| `1f6d2b94918f` | feat(auth): grant and revoke the ADMIN role in the user aggregate |
| `3fb3114f2dcb` | docs(changelog): record the 3.0.0 release |
| `d9ab0786e672` | fix(auth): delete the purged rows by primary key (HU-AUTH-002) |
| `d088090bffb5` | test(auth): validate every contract response in full |
| `4aa4a3c4c33d` | feat(auth): add POST /internal/maintenance/purge |
| `6d7d5b5a29cc` | feat(auth): add the purge use case and its PostgreSQL adapters |
| `b863fcd574fb` | feat(auth): verify the service credential of the internal operations |
| `03dded92b93c` | feat(auth): revoke the sessions of a user when a spent refresh token is reused |
| `27a4f2317cdf` | test(auth): validate the HTTP responses against auth-service.yaml |
| `8a307f6a96af` | docs(changelog): add the changelog of the auth api |
| `35f183b8b097` | docs(auth): fix the JWT key comment in .env.example and state how compose gets the key |
| `202ca27f735c` | fix(auth): log rejected refresh attempts and rename the token fields of the session use cases |

### csp-auth-db

`code-corhuila/csp-auth-db` - 8 commits

| Commit | Message |
|---|---|
| `1fb805371a08` | test(auth-db): cover the negative attempt count and the slot of a verified request (HU-AUTH-004) |
| `1320e3e00d91` | feat(auth-db): pending email and attempt counter in email_verification (HU-AUTH-004) |
| `ac15a4b48a03` | docs(auth): record V023 in the changelog |
| `812f0ed7fe17` | feat(auth): add the email prefix index for the user listing |
| `c43c56de68cf` | docs(changelog): record the 3.0.0 release |
| `ba624e61cbb4` | feat(auth-db): indexes for the retention purge of refresh_token and idempotency_key (HU-AUTH-002) |
| `25ca130e1b71` | feat(auth-db): grant DELETE on refresh_token and idempotency_key for the purge |
| `4fcdd0c0a803` | docs(changelog): add the changelog of the auth database |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 14 commits

| Commit | Message |
|---|---|
| `0f286c5ec57e` | fix(auth-portal): save the profile route title as UTF-8 and tidy the code step (HU-AUTH-004) |
| `6988d372611a` | feat(auth-portal): profile page to change and verify the email (HU-AUTH-004) |
| `f1389c19fb15` | fix(auth-portal): keep a failed session restore silent |
| `f979f925c829` | feat(auth-portal): restore and renew the session with the refresh token |
| `6410c404e557` | fix(auth-portal): show the user roles screen in Spanish |
| `2d75c6c610dc` | feat(auth-portal): handle the USER_WITHOUT_ROLES conflict on role revoke |
| `83a0e036792e` | docs(auth-portal): document the role change |
| `d5ff2a956fcb` | feat(auth-portal): grant and revoke the ADMIN role from the admin page |
| `ec383a8def1f` | docs(auth-portal): document the admin entry point |
| `eb7d2fac9b77` | feat(auth-portal): add the admin page to search users by email |
| `03d528b6f955` | docs(auth-portal): document the real auth API integration |
| `f89db6e861b5` | feat(auth-portal): call csp-auth-api for login, register and password recovery |
| `65c12ce5fd19` | docs(changelog): record the 3.0.0 release |
| `7d59b0aebfaa` | fix(deploy): use 127.0.0.1 in the container healthcheck |

### csp-front

`code-corhuila/csp-front` - 13 commits

| Commit | Message |
|---|---|
| `d9ab8ab85aa9` | fix(front): leave a protected page when the session ends |
| `46df5141c31a` | docs(front): document the explicit sign out |
| `22e14e4204fa` | test(front): cover the explicit sign out and the 401 path |
| `ba18a46ad48d` | feat(front): close the session on the server when the person signs out |
| `0641c6587390` | docs(front): document the session restore on reload |
| `a66d102ba82d` | test(front): cover the session restore wait, its timeout and the guards after it |
| `34ac770aa35b` | feat(front): wait for the session restore of the auth portal before the first navigation |
| `2210ef03dff0` | fix(front): save the Spanish copy of the shell as UTF-8 |
| `5fee370109fb` | fix(front): show the whole shell in Spanish |
| `821580d03947` | fix(front): show the user roles administration in Spanish |
| `3ac3c6804c58` | feat(front): mount the user roles administration of the auth portal |
| `7e0cc671d858` | feat(front): send the access token of the auth portal session |
| `d1751cc8e3d5` | fix(deploy): use 127.0.0.1 in the container healthcheck |

### csp-api-gateway

`code-corhuila/csp-api-gateway` - 4 commits

| Commit | Message |
|---|---|
| `5836213472a0` | feat(gateway): route GET /api/v1/auth/users as protected |
| `f3b38843ce74` | feat(gateway): route /api/v1/auth/ to csp-auth-api |
| `e6a5648c1d23` | docs(changelog): record the 3.0.0 release |
| `d9a069807217` | chore(scaffold): add base structure of the NGINX api gateway |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 7 commits

| Commit | Message |
|---|---|
| `84fb76766afd` | docs(changelog): record the 3.0.0 release |
| `0d65d4d641c5` | chore(infra): pass the identity purge settings through the platform environment |
| `941df1a96fe2` | docs(infra): say that the gateway is the only application service published to the host |
| `51d44f11b536` | chore(infra): compose the API gateway in the platform |
| `8a8f6ee7d582` | chore(infra): compose the Redis instance of csp-infra-redis |
| `35f992b7612e` | docs(infra): add the changelog |
| `5a1a96a6a4d7` | chore(infra): move the URL-safe note to AUTH_APP_PASSWORD and document the auth signing key |

### csp-infra-redis

`code-corhuila/csp-infra-redis` - 2 commits

| Commit | Message |
|---|---|
| `b4ce93a4b6d7` | docs(changelog): record the 3.0.0 release |
| `cb0665331054` | feat(infra): add the single Redis instance for the auth state |

### csp-worker

`code-corhuila/csp-worker` - 2 commits

| Commit | Message |
|---|---|
| `edeb3589f42b` | docs(changelog): record the 3.0.0 release |
| `822452f23e37` | feat(worker): add the purge-auth job |
