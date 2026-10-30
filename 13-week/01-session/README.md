# Week 13 - Session 1: Security in distributed systems (authentication, authorization and zero trust)

Author: Oscar Guillermo Sierra Lozano (`Oscarsl10`). Team CineSync Platform, Group 1.

The task of the session: authenticate at the edge, validate the token and authorize in **each** service, secure the service-to-service calls,
hash the secrets and design against the usual vulnerability classes.

**Scope of this work.** It is a rehearsal, not a change to the project. By decision of the owner nothing is applied to the `csp-*` repositories
from this session: the rehearsal runs on my own machine with a small stand-in of the gateway and of `csp-auth-api`, kept in [rehearsal/](rehearsal/).
It uses the real names of the project (`iss` `csp-auth`, `aud` `csp-api`, the permission `user:role:update`, the route `/api/v1/auth/users/{userId}`),
so what it shows can be compared with the real service. Section 5 lists what it does not prove, and
the threat model of the same boundary is in the [workshop](../workshop/README.md).

## 1. The scenario of the session, in the platform

The session describes an endpoint such as `DELETE /admin/users`, called "internal only" and without authentication, that a compromised service
or an SSRF can reach. In CineSync Platform the nearest case is the administrator route of `csp-auth-api`: `PATCH /users/{userId}/role` and
`GET /users`, which can give the ADMIN role to anybody. The rehearsal uses `DELETE /api/v1/auth/users/{userId}` with the same permission as a stand-in,
because deleting is the clearest example of an action that must never depend on the network.

## 2. Authentication and authorization are two checks

| Question | Where it is answered | In the rehearsal |
|---|---|---|
| Who are you? (authentication) | The token: signature, issuer, audience, expiry | `Verify` in `token.go`: `401` when any of them fails |
| What may you do? (authorization) | The permission inside the token, on **every** protected action | `protect` in `service.go`: `403` for a valid token without `user:role:update` |

Rules that the code applies:

1. The algorithm is fixed to `RS256` by the service, never read from the token. A token with `alg` `none` or `HS256` is rejected.
2. The token is checked without calling the database: the service only needs the public key (in the project, `GET /jwks`).
3. The answer is the same for every cause of `401`, so an attacker learns nothing about why it failed.
4. The edge (the gateway) only checks that `Authorization` is not empty and removes the headers `X-User-Id` and `X-User-Roles` that a client may send; the service never trusts it.

## 3. Zero trust: the service checks even when the caller skips the gateway

The rehearsal has two entrances to the same service: through the edge and directly. The two tests that matter are:

- `TestAuthzTheEdgeAloneDoesNotValidateTheToken`: a made-up token passes the edge, because it only looks for the header, and the service answers `401`.
- `TestAuthzCallingTheServiceDirectlyIsStillRejected`: a request that skips the edge and writes `X-User-Roles: ADMIN` by hand is still `401`.

This is the same reason why `csp-auth-api` validates the token itself and does not read the headers of the gateway.

## 4. Real output

Go go1.27.0 on Windows 11, 2026-10-30. Command: `cd 13-week/01-session/rehearsal && go test ./... -count=1 -v`.

```text
=== RUN   TestAuthzNoTokenAnswers401
--- PASS: TestAuthzNoTokenAnswers401 (0.00s)
=== RUN   TestAuthzValidTokenWithoutThePermissionAnswers403
--- PASS: TestAuthzValidTokenWithoutThePermissionAnswers403 (0.00s)
=== RUN   TestAuthzTokenWithThePermissionIsAccepted
--- PASS: TestAuthzTokenWithThePermissionIsAccepted (0.00s)
=== RUN   TestAuthzTokenSignedByAnotherKeyAnswers401
--- PASS: TestAuthzTokenSignedByAnotherKeyAnswers401 (0.40s)
=== RUN   TestAuthzTokenThatNamesAnotherAlgorithmAnswers401
--- PASS: TestAuthzTokenThatNamesAnotherAlgorithmAnswers401 (0.01s)
=== RUN   TestAuthzExpiredTokenAnswers401
--- PASS: TestAuthzExpiredTokenAnswers401 (0.00s)
=== RUN   TestAuthzTokenForAnotherAudienceAnswers401
--- PASS: TestAuthzTokenForAnotherAudienceAnswers401 (0.00s)
=== RUN   TestAuthzTheEdgeAloneDoesNotValidateTheToken
--- PASS: TestAuthzTheEdgeAloneDoesNotValidateTheToken (0.00s)
=== RUN   TestAuthzCallingTheServiceDirectlyIsStillRejected
--- PASS: TestAuthzCallingTheServiceDirectlyIsStillRejected (0.00s)
=== RUN   TestAuthzTheEdgeDropsTheIdentityHeadersOfTheClient
--- PASS: TestAuthzTheEdgeDropsTheIdentityHeadersOfTheClient (0.00s)
=== RUN   TestPasswordIsStoredHashedAndStillVerifies
--- PASS: TestPasswordIsStoredHashedAndStillVerifies (0.28s)
PASS
ok  	cinesync/w13-session1-rehearsal	2.648s
```

The tests catch a weak verifier: I removed the fixed-algorithm check from `Verify` for one run, and `TestAuthzTokenThatNamesAnotherAlgorithmAnswers401` failed; I put the check back and the whole run passed again.

## 5. What the rehearsal does not prove

- It is a stand-in of two components in one process, not the real gateway or the real `csp-auth-api`, and the key is generated in the test.
- Revocation is not covered here: a revoked access token is accepted by this stand-in until it expires. In the project the revoked-token list is in the Redis of Auth and only `csp-auth-api` reads it (ADR-033).
- Service-to-service trust is only argued, not run: there is no mTLS in the rehearsal and no service token. The project uses short-lived service tokens (section 6).
- Injection, SSRF and the other classes of OWASP are named in the self-check but not tested here; the threat model of the workshop covers the boundary `csp-api-gateway` to `csp-auth-api` with STRIDE.

## 6. Where the project stands against these ideas

Read in the repositories (branch `develop`), without changing anything for this session:

- **Gateway.** In `csp-api-gateway/nginx/routes/` the protected routes only check that `Authorization` is present ("the presence of Authorization is checked here, the token by the service"); the token is checked by each service. The `internal` routes of Auth and Booking answer `404` at the gateway.
- **Token validation per service.** `csp-auth-api` verifies the token with `RS256` only, issuer `csp-auth` and audience `csp-api` (`internal/adapter/out/token/access_verifier.go`), and checks the permission of each route with `requirePermission`. The revoked-token list is checked too.
- **Service-to-service.** Service credentials are short-lived tokens issued by an offline command; this week they were made to require issuer and audience as well (ADR-056, `csp-auth-api#198`). There is no mTLS in the gateway, in `csp-auth-api` or in the platform composition: the trust between services is the token, not the network.
- **Secrets.** Passwords are hashed with bcrypt (`internal/adapter/out/security/bcrypt_hasher.go`); keys and credentials come from the environment, not from the repositories (norm 5.9).
- **Open points.** The gateway does not validate tokens, so a stolen token reaches the service until it checks it; the lock of the login fails open when Redis is down (ADR-040), which is the row S of my STRIDE table. I did not audit the other domains: their own checks are in the stories of the teammates (for example `csp-docs#287` and `csp-docs#291`).

## 7. Self-check (answers)

1. **Authorization is checked...** on every protected action, not only at login: a token proves who the caller is, but each action asks for its own permission.
2. **A service validates a JWT by...** checking its signature and expiry, without a database call. The rehearsal also checks issuer, audience and the algorithm.
3. **Zero trust means...** not trusting the network and verifying every request, even the internal ones. The test that calls the service directly shows it.
4. **Service-to-service calls should use...** mTLS or short-lived tokens with the least privilege. A shared admin password or no authentication are the two mistakes of the session.
5. **The first class of web vulnerability is...** broken access control, that is, a missing authorization check. The test of the `403` is the guard against it.
6. **An "internal only" endpoint without authentication is...** a serious risk: a compromised service or an SSRF can reach it. The same checks must run behind the gateway.
