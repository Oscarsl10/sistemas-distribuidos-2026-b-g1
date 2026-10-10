# Git Commit Log - CineSync Platform

## Week 10

Period: 2026-10-05 to 2026-10-10. Author: Oscar Sierra (`Oscarsl10`). Same change promoted across several branches is listed once; automatic merge commits are omitted, and promotions of changes already reported in week 9 are not repeated. Newest first.

### Summary

| Repository | Commits |
|---|---|
| csp-docs | 19 |
| csp-front | 26 |
| csp-auth-api | 25 |
| csp-auth-db | 5 |
| csp-auth-portal | 15 |
| csp-concessions-portal | 22 |
| csp-catalog-portal | 1 |
| csp-ticketing-portal | 1 |
| csp-infra-mongo | 5 |
| csp-infra-postgres | 3 |
| **Total** | **122** |

### csp-docs

`code-corhuila/csp-docs` - 19 commits

| Commit | Message |
|---|---|
| `ead8cb7041e5` | docs(api): name the error field of the 409 codes as the shared envelope does |
| `610201c5aff8` | docs(api): name the 409 codes and the Location of POST /register |
| `3f0b706cba6a` | docs(auth): describe the retention of the idempotency keys in the runbook |
| `83d1558c42a3` | docs(api): align the Idempotency-Key length with the norm (8 to 128 characters) |
| `7be049c510ea` | docs(auth): document the Idempotency-Key of POST /register |
| `713cdb239f44` | docs(api): remove the pull request body file added by mistake |
| `99ca4764ba6f` | docs(api): limit the auth password to 72 bytes, the limit of bcrypt |
| `b60aa9b62254` | docs(auth): state that expired refresh tokens are not purged in the MVP |
| `170f31189a9f` | docs(adr): link the open issues and state the partial adoption in ADR-026 |
| `5cd2e8bebda7` | docs(adr): amend ADR-026 so the remote portals answer CORS from the environment |
| `b6f7507a3bd8` | docs(architecture): add ADR-027 one federated entry per role in each domain portal |
| `05aeb0308559` | docs(architecture): fix the encoding of ADR-026 and add traceability references |
| `5400debbbd00` | docs(architecture): add ADR-026 front-end runtime configuration and align remote ports |
| `896cb69544cc` | docs(architecture): track the concessions and ticketing remote alignment |
| `95be5d01c3de` | docs(architecture): address review of the gateway port change |
| `a4ceaf9639f5` | docs(architecture): publish the api-gateway at port 8000 and document the front shell and remotes (ADR-025) |
| `1ac2598d7609` | docs(auth): record that the portal joins first and last name into name |
| `459ab6a46869` | docs(architecture): add ADR-024 to the decision index |
| `771e00a06df1` | docs(auth): add phone and address to the register contract and the auth data model |

### csp-front

`code-corhuila/csp-front` - 26 commits

| Commit | Message |
|---|---|
| `f14af36e2636` | docs(changelog): add the administration mounts, the real spec count and two known limits |
| `d0dde0c95553` | docs(changelog): describe the development token as accepted only in development builds |
| `f9222b58f762` | docs(changelog): add the Unreleased section and the client-side limit of the guards |
| `c0cb62bb1d64` | docs(changelog): state the release process as the norm defines it |
| `1ed5b36d6b84` | docs(changelog): add the changelog with the 2.0.0 release |
| `89975125eaae` | feat(front): mount the booking and catalog administration entries |
| `dc56a0675a70` | fix(front): accept the development token only in development builds |
| `65df7a22dfb5` | test(front): cover the portal mounts, sign-in, session loader and unavailable notice |
| `250d44db3390` | feat(front): role guard and the customer snack mount |
| `ba677ab0b84d` | feat(front): open protected routes with the auth portal session |
| `45c35e2e4a10` | fix(front): keep /movies pointing to the start address |
| `cc5fc8547188` | feat(front): mount the catalog portal at the start address |
| `44be88441311` | fix(front): keep Angular out of main.ts so federation starts before it |
| `7558af332c39` | fix(front): keep the shell navigation on one row on narrow screens |
| `da07100a00ed` | feat(front): style the shell layout and its pages with the design system |
| `5fe05e35ed9c` | feat(front): add design tokens, buttons and logos from the design system |
| `b5dbc9206824` | chore(deploy): render environment config from container variables |
| `82ca35bc342e` | feat(front): read the gateway URL from a runtime config.json |
| `f61ef6cdeea6` | fix(front): use a single separator in the federation tsconfig paths |
| `6eb022652787` | chore(deploy): keep the healthcheck only in the Dockerfile |
| `24c81ad2df76` | chore(config): normalize federation paths and final newlines |
| `0aaab2cb082f` | chore(front): mount the token sign-in only in development |
| `c3da2fdfbfeb` | chore(front): add deploy files and CI workflow |
| `3896742b276a` | chore(front): add package-lock.json |
| `3d74f1225eff` | chore(front): add shell runtime with single http client and session |
| `9728af9b2389` | chore(front): add shell scaffold configuration |

### csp-auth-api

`code-corhuila/csp-auth-api` - 25 commits

| Commit | Message |
|---|---|
| `9b22f42778f9` | refactor(auth): drop reuse detection from the refresh endpoint |
| `46eec7e4c735` | feat(auth): add POST /refresh with refresh token rotation |
| `316140dc6149` | feat(auth): add POST /login and wire it in the composition root |
| `fb8d2fcd770d` | refactor(auth): extract the JSON body decoding from the register handler |
| `e75e265fd4d0` | chore(deploy): pass the service settings from the platform environment |
| `ae7e5af5550f` | feat(auth): wire POST /register to the database in the composition root |
| `46b63c0d4f5c` | feat(auth): add the register endpoint to the http adapter |
| `41c3989a4970` | feat(auth): make the registration idempotent and return the session |
| `20624090ba10` | feat(auth): add the login use case and its adapters |
| `f89516ae822a` | fix(auth): use a UUID v4 refresh token and seconds for its lifetime |
| `4d6bdf3a410e` | fix(auth): read APP_AUTH_JWT_EXPIRY as integer seconds |
| `c3f9d293d710` | fix(auth): state the 72-byte limit in the weak password error |
| `9fc9a46ebfc0` | feat(auth): store and issue refresh tokens |
| `8f85baff515d` | fix(auth): limit the password to 72 bytes |
| `f29e008d1a9a` | feat(auth): publish the signing key at GET /jwks |
| `eada83b15817` | feat(auth): add RS256 access token issuer with RFC 7638 key id |
| `4ae31b98d35a` | ci(auth): run the integration tests against the csp-auth-db schema |
| `826759128dcc` | feat(auth): read database pool, timeout and bcrypt settings from the environment |
| `c70c7ad96848` | feat(auth): add PostgreSQL and bcrypt adapters for user registration |
| `4f1827e5ce05` | feat(application): add register user use case and its ports |
| `327cd201f9a5` | feat(domain): add password policy, name and UserRegistered event |
| `799ce5968900` | test(auth): load the test configuration before starting the server goroutine |
| `da505dfe2f83` | test(auth): serve on a listener the test owns and assert the listen error |
| `b2533bf2c55f` | test(auth): cover the composition root and the user accessors |
| `a2c9daca50e6` | feat(auth): add phone and address contact data to the user aggregate |

### csp-auth-db

`code-corhuila/csp-auth-db` - 5 commits

| Commit | Message |
|---|---|
| `88a39212daed` | feat(auth): add the idempotency_key table for POST /register |
| `e99505ce36bb` | test(auth): check seed, roles, grants and constraints of the schema in CI |
| `1764f1b0ec81` | feat(auth): add phone and address columns to app_user |
| `a74f0837b1e9` | fix(schema): grant auth_writer the columns the outbox purge reads |
| `1fb174e5af9b` | fix(governance): name the executor credentials after the domain |

### csp-auth-portal

`code-corhuila/csp-auth-portal` - 15 commits

| Commit | Message |
|---|---|
| `6a565e29feb5` | docs(changelog): add the Unreleased section and the client-side limit of the guards |
| `2946e10659ee` | docs(changelog): state the release process as the norm defines it |
| `962895bcab5d` | docs(changelog): add the changelog with the 2.0.0 release |
| `01b251f4afa5` | chore(auth): keep LF in the shell scripts of the container |
| `62b6347cf8c8` | fix(auth): answer CORS for the shell from the portal container |
| `d8b39c1f7957` | test(auth): cover the lazy routes, the toast container and the recovery guard |
| `9ea407596882` | feat(auth): expose the in-memory session to the shell |
| `e09a3eeedc01` | feat(auth): open login and register as modals with toasts, as in the mockup |
| `6776d44e55e7` | fix(auth): apply the mockup background to the standalone portal |
| `b231e888fa93` | fix(auth): mount the standalone portal under /auth and stand in for /movies and /admin |
| `eef1f95c9aa6` | feat(auth): ask for phone and address in the register form |
| `efd8c30804c2` | feat(auth): add the register form and document the synthetic users |
| `ea3d49e6a6af` | feat(auth): apply the mockup design and Spanish copy to login and recovery |
| `6494b7c49751` | feat(auth): add login and forgot-password forms |
| `2f63dcce6150` | feat(auth): add in-memory session and role guards |

### csp-concessions-portal

`code-corhuila/csp-concessions-portal` - 22 commits

| Commit | Message |
|---|---|
| `ee6aefca3aee` | docs(changelog): add the Unreleased section and the client-side limit of the guards |
| `5bd7949c8281` | docs(changelog): state the release process as the norm defines it |
| `076166dd3bc9` | docs(changelog): add the changelog with the 2.0.0 release |
| `b069ac4d7d0b` | fix(concessions): answer CORS for the shell from the portal container |
| `9300a7a207e4` | test(concessions): cover price rounding, empty price, admin routes and the HTTP service |
| `45d43bdd0a16` | feat(concessions): admin combos screen with synthetic data |
| `238565effb44` | fix(concessions): show the seeded categories and stock reasons in Spanish |
| `8a69f5f8e9ce` | docs(concessions): document the administration screens |
| `17e725a1490f` | feat(concessions): admin inventory screen with synthetic data |
| `d5b47f53ed63` | feat(concessions): admin products screen with synthetic data |
| `d7d2d5eae465` | feat(concessions): admin area layout, navigation and shared look |
| `ba4e448f8a2b` | docs(concessions): document the synthetic catalog and the routes like csp-auth-portal |
| `281196265aee` | feat(concessions): show the snack selection in Spanish, as the mockup does |
| `06127e5c021a` | feat(concessions): snack selection screen with the mockup look |
| `696dd32b6fce` | fix(concessions): assert that the admin area rejects the customer route |
| `bd3bd2b12feb` | feat(concessions): order draft with quantities and total |
| `9b835f75b049` | feat(concessions): synthetic catalog and published filter |
| `0ab39348890c` | fix(concessions): show the title of the snack selection in Spanish |
| `4b77fdcc4484` | fix(concessions): export the snack routes from the federation barrel |
| `de463e8f6e50` | refactor(concessions): split customer and administration routes into two federated entries |
| `4fdb39396303` | chore(concessions): standalone runtime parity and mockup base look |
| `155f418efa26` | chore(concessions): move the portal dev port to 4204 |

### csp-catalog-portal

`code-corhuila/csp-catalog-portal` - 1 commits

| Commit | Message |
|---|---|
| `1a0ce64bcb20` | chore(catalog): move the portal dev port to 4203 |

### csp-ticketing-portal

`code-corhuila/csp-ticketing-portal` - 1 commits

| Commit | Message |
|---|---|
| `f4fada27b966` | chore(ticketing): align the remote name, exported routes and dev port with the shell |

### csp-infra-mongo

`code-corhuila/csp-infra-mongo` - 5 commits

| Commit | Message |
|---|---|
| `5963aa21cc29` | fix(infra): leave an existing user as it is and print the key file hint on one line |
| `e77887c6d79a` | fix(infra): read an unpublished port as :0 in the CI check |
| `7b7b6656fbc4` | ci(infra): start the instance from an empty volume and check users and replica set |
| `f7858a2cc478` | feat(infra): run the single MongoDB instance as a single-node replica set |
| `68ac65b7acbe` | chore(infra): lay the repository structure, env examples and README |

### csp-infra-postgres

`code-corhuila/csp-infra-postgres` - 3 commits

| Commit | Message |
|---|---|
| `ab51c47fb6d1` | chore(infra): compose the auth service in the platform |
| `22a4d2a425d2` | fix(infra): print the key file hint on one line |
| `568dda1af459` | feat(infra): compose the MongoDB instance in the root composition |
