<!-- HU-STATUS TEMPLATE - do NOT remove the <!-- ... --> markers or the table headers.
     Your weekly grade is read AUTOMATICALLY from this file:
       13-week/hu-status/README.md  (inside YOUR fork). English. -->

# Weekly Status - Week 13

<!-- CONFIG-START - must match your profile repo (username/username) CONFIG -->
- FULL_NAME: Oscar Guillermo Sierra Lozano
- GITHUB_USER: Oscarsl10
- TEAM: CineSync Platform
- SPRINT_GOAL: Threat model the boundary from `csp-api-gateway` to `csp-auth-api` and turn the three highest risks into security stories and a security gate for the pipeline (week 13 workshop, progress delivery).
<!-- CONFIG-END -->

## 1. User stories worked this week
| HU ID | Title | Status (todo/doing/done) | Evidence (PR or commit URL) |
|---|---|---|---|
| W13-THREAT-MODEL | Threat model of the boundary `csp-api-gateway` to `csp-auth-api` (workshop, no HU) | doing | [csp-docs#278](https://github.com/code-corhuila/csp-docs/pull/278) |
| W13-SECURITY-BACKLOG | Three security stories: [#280](https://github.com/code-corhuila/csp-docs/issues/280) (S), [#281](https://github.com/code-corhuila/csp-docs/issues/281) (D), [#282](https://github.com/code-corhuila/csp-docs/issues/282) (I) | doing | [#280](https://github.com/code-corhuila/csp-docs/issues/280), [#281](https://github.com/code-corhuila/csp-docs/issues/281), [#282](https://github.com/code-corhuila/csp-docs/issues/282) |

## 2. My individual contribution

Workshop of today (2026-10-28), individual, boundary claimed in [csp-docs#275](https://github.com/code-corhuila/csp-docs/issues/275).

- **Attack surface map** ([01-attack-surface.md](../workshop/01-attack-surface.md)): data flow of the platform with seven trust boundaries (mine: TB2a, gateway to `csp-auth-api`), nine assets with their classification and the entry points of the boundary with how they authenticate today.
- **STRIDE** in the team `csp-docs` ([pull request #278](https://github.com/code-corhuila/csp-docs/pull/278), `05-architecture/security/stride-Oscarsl10.md`, commits [68a9fb8](https://github.com/code-corhuila/csp-docs/commit/68a9fb8), [b38bdb1](https://github.com/code-corhuila/csp-docs/commit/b38bdb1)): six threats with likelihood, impact, risk, existing control and missing mitigation, a ranking and the residual risk. The pull request is open and I do not merge it. The row S and its mitigation come from the library book *Hacking y seguridad en Internet* (García-Morán et al., 2nd ed., Ediciones de la U, 2013, chapter 3, section 3.4.1, page 132, and chapter 4, section 4.2.1.3, page 149).
- **Security backlog**: [#280](https://github.com/code-corhuila/csp-docs/issues/280) failed logins limited, recorded and reviewed (S, risk 6), [#281](https://github.com/code-corhuila/csp-docs/issues/281) stricter rate limit for registration and login (D, risk 6) and [#282](https://github.com/code-corhuila/csp-docs/issues/282) registration that does not reveal the account (I, risk 6), with the labels `security` and `P1`, in the team board.
- **Security gates** ([04-security-ci](../workshop/04-security-ci/README.md)): workflow `security.yml` with gitleaks, `govulncheck` and the authorization tests, checked with `actionlint`, and a negative test of the login of `csp-auth-api` (the 11th failed login of one address answers `429` even with valid credentials).
- **Peer review:** my review comment on the STRIDE pull request of a teammate: https://github.com/code-corhuila/csp-docs/pull/277#pullrequestreview-5691705056.
- Commits of the workshop: CP0 [d28ab59](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/d28ab59), CP1 [e0176bd](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/e0176bd), CP3 [eb628a0](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/eb628a0), CP4 [f14c995](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/f14c995), CP4 [5f559a4](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/5f559a4), CP4 [879d19f](https://github.com/Oscarsl10/sistemas-distribuidos-2026-b-g1/commit/879d19f).

## 3. Blockers and risks

- GitHub Actions is locked by billing, so `security.yml` has not run in GitHub; its commands were run on my machine and `govulncheck` reports a vulnerability of the Go standard library (`GO-2026-6613`) in the installed toolchain.
- Stories #280 to #282 change the Auth contract or the gateway and need the agreement of the team before they are built.
- The negative test runs against a minimal model of the login, not against the real service yet.

## 4. Plan for next week

- Answer the review of my STRIDE pull request, refine the three stories and run the negative test against the real `csp-auth-api` (final delivery of the workshop, until 2026-10-30).

## 5. Compliance self-check

- [x] Conventional Commits - `type(scope): summary`
- [x] Per-environment HU branch + PR to that environment (code repos: child branch -> `develop`; `csp-docs`: `docs/* -> main`). The STRIDE branch follows the workshop name `security/stride-<user>`.
- [x] Testable acceptance criteria
- [x] Tests added/updated (unit / integration). Four negative authorization tests of the login limit.
- [x] DDD / hexagonal boundaries respected (domain has no I/O)
- [x] No secrets; config via environment variables

## 6. Evidence links

- [Workshop folder](../workshop/README.md)
- [Pull request csp-docs#278](https://github.com/code-corhuila/csp-docs/pull/278)
- [Claim issue csp-docs#275](https://github.com/code-corhuila/csp-docs/issues/275)
- Security stories: [#280](https://github.com/code-corhuila/csp-docs/issues/280), [#281](https://github.com/code-corhuila/csp-docs/issues/281), [#282](https://github.com/code-corhuila/csp-docs/issues/282)
