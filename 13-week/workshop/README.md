# Week 13 Workshop - Threat Modeling and the Security Backlog

<!-- CONFIG-START -->
- FULL_NAME: Oscar Guillermo Sierra Lozano
- GITHUB_USER: Oscarsl10
- TEAM: CineSync Platform
- DOCS_REPO: https://github.com/code-corhuila/csp-docs
- BOUNDARY: csp-api-gateway to csp-auth-api (requests under /api/v1/auth, inbound to Auth)
- CLAIM_ISSUE: https://github.com/code-corhuila/csp-docs/issues/275
- REFLECTION_DOC: https://docs.google.com/document/d/1WgrVs7Jf5x53A961HDXYZA3ljsBY6otg4e1RXYTfnwQ/edit?tab=t.0
- LIBRARY_BOOK: Jean Paul García-Morán, Yago Fernández Hansen, Rubén Martínez Sánchez, Ángel Ochoa Martín and Antonio Ángel Ramos Varón. Hacking y seguridad en Internet. 2nd edition. Bogotá: Ediciones de la U, 2013. Call number 0011099. Used: chapter 3, section 3.4.1, page 132 (brute force against credentials) and chapter 4, section 4.2.1.3, page 149 (account lockout policy)
- LIBRARY_RECORD_TIME: 2026-10-28 19:37
- PULL_REQUEST: https://github.com/code-corhuila/csp-docs/pull/278
- STORIES: https://github.com/code-corhuila/csp-docs/issues/280, https://github.com/code-corhuila/csp-docs/issues/281, https://github.com/code-corhuila/csp-docs/issues/282
<!-- CONFIG-END -->

## The boundary I claimed

I own the Auth domain, so I chose the boundary where the requests of every client reach it: from `csp-api-gateway` (NGINX, port 8000) to
`csp-auth-api` (Go) under `/api/v1/auth`. It carries credentials (`POST /login`, `POST /register`), access and refresh tokens, the contact
data stored in `app_user` and the role changes of `PATCH /users/{userId}/role`. The boundary of the gateway to the Booking service belongs to
a teammate (issue #274), so the two analyses do not overlap.

## Security backlog

Each story comes from one threat of the STRIDE table in [pull request #278](https://github.com/code-corhuila/csp-docs/pull/278) and has a label `security`, a priority label, testable criteria with one automated test and a place in the team board.

| Story | STRIDE | Risk | Priority | Related to |
|---|---|---|---|---|
| [#280](https://github.com/code-corhuila/csp-docs/issues/280) Failed logins are limited, recorded and reviewed | S | 6 | P1 | [#281](https://github.com/code-corhuila/csp-docs/issues/281) (same route `POST /login`), [#288](https://github.com/code-corhuila/csp-docs/issues/288) (record of rejected calls), [#284](https://github.com/code-corhuila/csp-docs/issues/284) (stolen or revoked access token) |
| [#281](https://github.com/code-corhuila/csp-docs/issues/281) Registration and login have their own stricter rate limit | D | 6 | P1 | [#280](https://github.com/code-corhuila/csp-docs/issues/280), [#282](https://github.com/code-corhuila/csp-docs/issues/282) (same route `POST /register`) |
| [#282](https://github.com/code-corhuila/csp-docs/issues/282) Registering an existing email does not reveal the account | I | 6 | P1 | [#281](https://github.com/code-corhuila/csp-docs/issues/281) |

After the review of my pull request, #281 also limits `POST /forgot-password`, the mail bombing threat that a teammate found.
