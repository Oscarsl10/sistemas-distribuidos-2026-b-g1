# Git Commit Log - CineSync Platform

## Week 12

Period: 2026-10-19 to 2026-10-24. Author: Oscar Sierra (`Oscarsl10`). Only commits authored by me are listed, including those in the shared repositories (`csp-front`, `csp-api-gateway`, `csp-worker`, `csp-infra-*`). Same change promoted across several branches is listed once; automatic merge commits are omitted, and changes already reported in weeks 9, 10 and 11 are not repeated. Newest first.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 28 |
| csp-auth-api | 48 |
| csp-auth-db | 4 |
| csp-auth-portal | 14 |
| csp-front | 5 |
| csp-api-gateway | 5 |
| csp-infra-postgres | 6 |
| csp-infra-rabbitmq | 8 |
| csp-worker | 7 |
| **Total** | **125** |

### csp-docs

`code-corhuila/csp-docs` - 28 commits

| Commit | Message |
|---|---|
| `bec2a5b1cc05` | docs(auth): correct the status of the login lock by client IP |
| `c4f9bb370bf9` | docs(auth): describe the service tokens issued with auth-token |
| `9f1f7d7f5a0a` | docs(auth): record that ADMIN holds concessions:inventory:read |
| `707b0bf2fe6b` | docs(auth): shorten the access token to 15 minutes |
| `95eaef6be8ed` | docs(adr): mark the offline command as the chosen option of ADR-052 |
| `9ac7f2e6dc48` | docs(adr): accept ADR-052 and ADR-053 with the decisions of the tech lead |
| `e1742791cc99` | docs(auth): align the auth documents and contracts with the norm |
| `3b781f6660e0` | docs(ux): update the grant and revoke administrator dialog texts |
| `99b735eba2ce` | docs(auth): end the sessions of a user whose ADMIN role is revoked for HU-AUTH-003 |
| `46ad6da84223` | docs(auth): specify HU-AUTH-007 profile edition and password change |
| `3bae7d6e8d01` | docs(auth): use the service message in the weak password example of reset-password |
| `f8c133b3ec07` | docs(auth): reject the current password on reset for HU-AUTH-006 |
| `42990de73020` | docs(auth): add the reset password route to the navigation map |
| `d28b97e6b52c` | docs(auth): supersede the pending reset token by soft delete |
| `56508e7afbfc` | docs(auth): specify password reset for HU-AUTH-006 |
| `2928133329bc` | docs(auth): describe the service token of csp-worker in the runbook |
| `60982a349e18` | docs(auth): say that the failing call is answered as paused in adr-048 |
| `b6e405634ace` | docs(auth): record the bounded latency of the login attempt limit in adr-048 |
| `000e513383a1` | docs(auth): record the decisions of the integration campaign of HU-AUTH-004 |
| `f749826ee223` | docs(auth): the idempotency key is required on the email change calls (HU-AUTH-004) |
| `8540dbedc505` | docs(architecture): accept adr-027 and adr-036 |
| `ff6f0c0955e2` | docs(auth): align adr-040 with the implemented login attempt limit |
| `369b731fdd08` | docs(auth): record the login attempt limit by client ip and align the auth docs (hu-auth-005) |
| `cf69d33df778` | docs(auth): record the smtp relay provider and its limits in the runbook and deployment |
| `2239b22a031e` | docs(infra): amend adr-036 with the permissions verified against a real broker |
| `69053db225cf` | docs(auth): contract 1.6.1, adr-030 amendment for email_verification and the profile route in the shell |
| `d9fe21d49ba3` | docs(infra): add adr-036 for the rabbitmq instance csp-infra-rabbitmq |
| `525313b32f43` | docs(adr): amend ADR-026 so the CORS regex is anchored and the mechanism stays a template per repository |

### csp-auth-api

`code-corhuila/csp-auth-api` - 48 commits

| Commit | Message |
|---|---|
| `277a8785e0df` | feat(auth): add the offline command that issues the service tokens |
| `3c0ce03f2f27` | feat(auth): issue the concessions:inventory:read permission to ADMIN |
| `0e343600b9df` | fix(auth): issue the access token for 15 minutes by default |
| `9c80026c9d63` | fix(auth): add a caching rule to GET /jwks |
| `7ff3b7d964d5` | fix(auth): keep the removal of admin when the client disconnects |
| `e6cd9be66c23` | fix(auth): keep the logout revocation when the client disconnects |
| `478ecfa78c20` | docs(auth): document the change of the own password |
| `e13fce233666` | feat(auth): wire the password change route |
| `2862bfe54cb5` | feat(auth): add the POST /me/password route |
| `5e1aa5745481` | docs(auth): document the settings of the password change limits |
| `b916188b1140` | feat(auth): add the settings of the password change limits |
| `208b48dfea6c` | feat(auth): lock the account row while the password changes |
| `9727ec11f344` | feat(auth): count the wrong current passwords per user in redis |
| `2ea4343e97ba` | docs(auth): record the change password use case in the changelog |
| `6de384f32a38` | feat(auth): add the change password use case |
| `2274f609dfd7` | feat(auth): add the password change limit and errors to the domain |
| `5e012b96a1c5` | docs(auth): document the own profile routes |
| `ea12f6174716` | feat(auth): wire the own profile routes |
| `4f236671793e` | feat(auth): add the GET and PATCH /me routes |
| `76352c0dabf5` | feat(auth): update the profile fields in one statement |
| `a86ac36254f4` | docs(auth): record the own profile use cases in the changelog |
| `47b9a8330cf6` | feat(auth): add the use cases to read and edit the own profile |
| `ee87debfda8d` | feat(auth): add the profile change value object |
| `561effff21d9` | test(auth): cover the session revocation of a removed admin with PostgreSQL |
| `e0431175a4b8` | docs(auth): document that revoking admin ends the sessions of the user |
| `5d9a0bfdbbf2` | feat(auth): end the sessions of a user whose admin role is revoked |
| `ae81b2435197` | docs(auth): document the rejection of the current password on reset |
| `22521e21ee6c` | feat(auth): reject the current password on reset-password |
| `2870c03cfbbf` | feat(auth): wire the password reset routes and document their settings |
| `619e8a18acb3` | feat(auth): add the forgot-password and reset-password routes |
| `d56bf76a4ffb` | feat(auth): add the password reset sender and its settings |
| `e885f9bd46fd` | feat(auth): add the password reset repository and the purge of password resets |
| `716f5542377c` | docs(auth): record the password reset use cases in the changelog |
| `915bf4838657` | feat(auth): add the password reset use cases and their ports |
| `cc6892636c35` | feat(auth): add the reset token and the password reset to the domain |
| `8b2e55140fff` | fix(auth): make the bulkhead tests independent of the speed of the machine |
| `be2d84c6c2fe` | fix(auth): bound the latency of the login attempt limit when redis does not answer |
| `7f0c0ffc4798` | feat(auth): answer 429 rate_limited on login after repeated failed attempts |
| `25940abd872a` | feat(auth): add the redis adapter of the login attempt limit |
| `857f528c50e2` | feat(auth): add the domain and use case of the login attempt limit |
| `5230e9527f64` | fix(auth): measure the refusal before the contract check in the bulkhead test |
| `6c5fe0989377` | fix(auth): bound the concurrent email change requests with a bulkhead |
| `4f937abe451b` | fix(auth): answer 401 when the account of the token is gone in the email change request |
| `08b7d041af6d` | fix(auth): reject an email longer than 255 characters with 400 |
| `f45d38847119` | docs(purge): document the email verification retention |
| `cb8043e10046` | feat(purge): delete email verification requests expired more than 7 days ago |
| `1173d299742a` | docs(auth): document where the smtp variables go per environment |
| `a2ebb34c96dc` | fix(auth): answer the exact remaining seconds in the retry-after of the email change limit |

### csp-auth-db

`code-corhuila/csp-auth-db` - 4 commits

| Commit | Message |
|---|---|
| `cfdcdb02f32b` | feat(auth-db): add schema for password reset |
| `1a6eecd269e9` | test(auth-db): check the email verification purge grant and index |
| `be4dc9163487` | chore(auth-db): add complete expires_at index on email_verification |
| `6296dbcd8e0a` | chore(auth-db): grant delete on email_verification to auth_writer |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 14 commits

| Commit | Message |
|---|---|
| `6d2c9bd13492` | fix(session): forget the session before the logout answers and keep the request alive |
| `25514087bf8c` | fix(admin): tell that removing admin ends the session of the user |
| `24bd85772b72` | feat(auth-portal): add the password change section to the profile page |
| `f2b99059493c` | feat(auth-portal): add the personal data section to the profile page |
| `ab4b2c5836a8` | docs(auth-portal): document the profile and password requests |
| `90d81a14afa1` | feat(auth-portal): add the own profile and password change requests |
| `a2ea0cb7175f` | feat(auth-portal): reject the current password on the reset password page |
| `64502ae7426d` | docs(auth-portal): document the reset password page |
| `5f77d178e3d8` | feat(auth-portal): add the reset password page for the emailed link |
| `66fe1885c641` | feat(auth-portal): show the wait time of the rate-limited login |
| `553fdc480083` | feat(auth): explain the failed-login rate limit on the login form |
| `939784662389` | test(auth-portal): check the requests of the portal against the auth contract (HU-AUTH-004) |
| `58f543838475` | fix(auth-portal): send the Idempotency-Key on the email change requests (HU-AUTH-004) |
| `152a7600b92e` | fix(auth-portal): refuse an unanchored CORS regex and serve only remoteEntry.json as no-store json |

### csp-front

`code-corhuila/csp-front` - 5 commits

| Commit | Message |
|---|---|
| `e139b7817f88` | fix(shell): send http requests with fetch so keepalive takes effect |
| `cc50e2062acf` | fix(shell): remove my account menu and keep signed-in people off the login page |
| `3b6dee5e4eb4` | feat(front): carry retry-after in api error |
| `466e22e37491` | feat(front): mount the profile page at /profile and link it from the navigation |
| `b581176acc17` | fix(front): refuse an unanchored CORS regex and check the CORS rule of the image in the CI |

### csp-api-gateway

`code-corhuila/csp-api-gateway` - 5 commits

| Commit | Message |
|---|---|
| `025278ddb9b0` | fix(auth): route the exact own-profile path to the auth service |
| `fd5b9922d7de` | fix(gateway): accept only safe correlation ids in responses and errors |
| `e9a5a8bf1910` | fix(gateway): expose retry-after to browsers and emit one correlation id |
| `60bcfd666c0e` | test(gateway): burst the rate limit with parallel curl in the integration checks |
| `7d1d499a7d7e` | fix(gateway): strip client-supplied x-user headers before proxying |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 6 commits

| Commit | Message |
|---|---|
| `33c482410331` | feat(env): document the password reset url of the auth service |
| `affe19ba5301` | chore(infra): enable the auth purge in qa |
| `87f8a85ce8e3` | chore(env): pass the SMTP relay of csp-auth-api from the platform env files |
| `00d808a0e3b7` | docs(infra): document the steps of an existing environment for the migrators |
| `cbc9e4f9adf0` | feat(infra): create the auth_migrator login that owns the auth schema |
| `e540807aed93` | feat(infra): create the booking_migrator login that owns the booking schema |

### csp-infra-rabbitmq

`code-corhuila/csp-infra-rabbitmq` - 8 commits

| Commit | Message |
|---|---|
| `4fc8e3203586` | test(rabbitmq): validate the vhost of the AMQP URLs of the env examples |
| `66e3865cec0b` | fix(rabbitmq): end the documented AMQP URLs in the default vhost as %2F |
| `1464f2b1dde2` | docs: describe the instance, its logins and its limits |
| `282d8373ac94` | docs: add the changelog |
| `9fc634c3366e` | chore: keep every file with LF endings |
| `3c2830c5df27` | ci: start the instance as dev and qa and run the checks |
| `1189592f1d7d` | test: check the exchanges, logins, permissions and limits against a real broker |
| `75a2be450981` | chore: add the rabbitmq instance with its definitions and scripts |

### csp-worker

`code-corhuila/csp-worker` - 7 commits

| Commit | Message |
|---|---|
| `560868c59c59` | fix(worker): check the relay exchange instead of declaring it |
| `886d39031d47` | fix(worker): make the auth relay end-to-end test independent of the publication order |
| `e9c93370ea97` | chore(worker): pass the auth relay variables in the compose file |
| `a8b3a3ae64d4` | docs(worker): document the auth outbox relay and its flag |
| `9bee0e47e1ec` | feat(worker): relay the auth outbox behind its own flag |
| `50cb10b4c081` | refactor(worker): parameterize the outbox relay adapters by domain |
| `d89435086344` | feat(worker): take an instance lock before each run of the outbox relay |
