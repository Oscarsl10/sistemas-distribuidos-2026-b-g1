<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Guillermo Sierra Lozano
- GITHUB_USER: Oscarsl10
- TEAM: CineSync Platform
- SPRINT_GOAL: Threat model the boundary from `csp-api-gateway` to `csp-auth-api` and turn the three highest risks into security stories and a security gate for the pipeline (week 13 workshop, final delivery).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W13-THREAT-MODEL | Threat model of the boundary `csp-api-gateway` to `csp-auth-api` (workshop, no HU) | done | [csp-docs#278](https://github.com/code-corhuila/csp-docs/pull/278) |
| W13-SECURITY-BACKLOG | Three security stories: [#280](https://github.com/code-corhuila/csp-docs/issues/280) (S), [#281](https://github.com/code-corhuila/csp-docs/issues/281) (D), [#282](https://github.com/code-corhuila/csp-docs/issues/282) (I) | done | [#280](https://github.com/code-corhuila/csp-docs/issues/280), [#281](https://github.com/code-corhuila/csp-docs/issues/281), [#282](https://github.com/code-corhuila/csp-docs/issues/282) |

## 2. My individual contribution

Workshop of 2026-10-28 (in class) and 2026-10-29 and 2026-10-30 (final delivery), individual, boundary claimed in [csp-docs#275](https://github.com/code-corhuila/csp-docs/issues/275).

- **Attack surface map** ([01-attack-surface.md](../workshop/01-attack-surface.md)): data flow of the platform with seven trust boundaries (mine: TB2a, gateway to `csp-auth-api`), nine assets with their classification and the entry points of the boundary with how they authenticate today.
- **STRIDE** in the team `csp-docs` ([pull request #278](https://github.com/code-corhuila/csp-docs/pull/278), `05-architecture/security/stride-Oscarsl10.md`, commits CP2 [68a9fb8](https://github.com/code-corhuila/csp-docs/commit/68a9fb8), CP2 [b38bdb1](https://github.com/code-corhuila/csp-docs/commit/b38bdb1), CP6 [b7bbc56](https://github.com/code-corhuila/csp-docs/commit/b7bbc56), CP6 [1bf6fd9](https://github.com/code-corhuila/csp-docs/commit/1bf6fd9)): seven threats with likelihood, impact, risk, existing control and missing mitigation, a ranking, the residual risk and the references. The pull request is open and I do not merge it.
- **Review I received:** a teammate found the missing threat of mail bombing in `POST /forgot-password` and questioned a database check of the role. I added the threat as a row D (risk 4), kept the limit of ADR-050 instead of a new Redis counter, and removed the database check because removing ADMIN already ends the sessions of that user. My answer: https://github.com/code-corhuila/csp-docs/pull/278#issuecomment-6498754591.
- **Book:** *Hacking y seguridad en Internet* (García-Morán et al., 2nd ed., Ediciones de la U, 2013, call number 0011099): the threat of brute force (chapter 3, section 3.4.1, page 132) and the lockout after a few failed attempts (chapter 4, section 4.2.1.3, page 149) gave row S. The record of every lock is my own addition.
- **Security backlog**: [#280](https://github.com/code-corhuila/csp-docs/issues/280) failed logins limited, recorded and reviewed (S, risk 6), [#281](https://github.com/code-corhuila/csp-docs/issues/281) stricter rate limit for registration, login and `POST /forgot-password` (D, risk 6) and [#282](https://github.com/code-corhuila/csp-docs/issues/282) registration that does not reveal the account (I, risk 6), with the labels `security` and `P1`, in the team board, linked to each other and to related issues of teammates.
- **Security gates** ([04-security-ci](../workshop/04-security-ci/README.md)): workflow `security.yml` with gitleaks, `govulncheck` and the authorization tests, checked with `actionlint`, and a negative test of the login. I also ran the tests of the login lock in the real `csp-auth-api` and recorded the output in [HOW-TO-RUN.md](../workshop/04-security-ci/HOW-TO-RUN.md).
- **Peer review:** my review comment on the STRIDE pull request of a teammate: https://github.com/code-corhuila/csp-docs/pull/277#pullrequestreview-5691705056.
- Commits of the fork: CP0 [d28ab59](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/d28ab59), CP1 [e0176bd](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/e0176bd), CP3 [eb628a0](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/eb628a0), CP4 [f14c995](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/f14c995), CP4 [5f559a4](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/5f559a4), CP4 [879d19f](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/879d19f), CP5 [794dd0b](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/794dd0b), CP7 [18d965a](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/18d965a), CP7 [d3f80bd](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/d3f80bd), CP8 [81a7292](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/81a7292).

## 3. Blockers and risks

- GitHub Actions is locked by billing, so `security.yml` has not run in GitHub; its commands were run on my machine and `govulncheck` reports a vulnerability of the Go standard library (`GO-2026-6613`) in the installed toolchain.
- Stories #280 to #282 change the Auth contract or the gateway and need the agreement of the team before they are built.
- The end-to-end test of the login lock against PostgreSQL and Redis was not run, and the names of the real tests do not contain `Authz`, so the job `authz-tests` would not select them yet.
- The STRIDE pull request #278 is still open and waits for the team review and merge.

## 4. Plan for next week

- Build #280 first (record of every lock and a counter per account), then #281 and, after the agreement of the team on the contract, #282. Put the negative test in the pipeline of `csp-auth-api` when the Actions run again.

## 5. Compliance self-check

- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (code repos: child branch -> `develop`; `csp-docs`: `docs/* -> main`). The STRIDE branch follows the workshop name `security/stride-<user>`.
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration). Four negative authorization tests of the login limit, and four tests of the real service run.
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links

- [Workshop folder](../workshop/README.md)
- [Pull request csp-docs#278](https://github.com/code-corhuila/csp-docs/pull/278)
- [Claim issue csp-docs#275](https://github.com/code-corhuila/csp-docs/issues/275)
- Security stories: [#280](https://github.com/code-corhuila/csp-docs/issues/280), [#281](https://github.com/code-corhuila/csp-docs/issues/281), [#282](https://github.com/code-corhuila/csp-docs/issues/282)
