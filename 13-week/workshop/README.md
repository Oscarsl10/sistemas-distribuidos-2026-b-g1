# Week 13 Workshop - Threat Modeling and the Security Backlog

<!-- CONFIG-START -->
- FULL_NAME: Oscar Guillermo Sierra Lozano
- GITHUB_USER: Oscarsl10
- TEAM: CineSync Platform
- DOCS_REPO: https://github.com/code-corhuila/csp-docs
- BOUNDARY: csp-api-gateway to csp-auth-api (requests under /api/v1/auth, inbound to Auth)
- CLAIM_ISSUE: https://github.com/code-corhuila/csp-docs/issues/275
- REFLECTION_DOC: https://docs.google.com/document/d/1WgrVs7Jf5x53A961HDXYZA3ljsBY6otg4e1RXYTfnwQ/edit?tab=t.0
- LIBRARY_BOOK: Jean Paul Garcia-Moran, Yago Fernandez Hansen, Ruben Martinez Sanchez, Angel Ochoa Martin and Antonio Angel Ramos Varon. Hacking y seguridad en internet. 2nd edition. Year and call number: to be completed from the physical copy. Library registration number 0011099.
- LIBRARY_RECORD_TIME: 2026-10-28 19:37
<!-- CONFIG-END -->

## The boundary I claimed

I own the Auth domain, so I chose the boundary where the requests of every client reach it: from `csp-api-gateway` (NGINX, port 8000) to
`csp-auth-api` (Go) under `/api/v1/auth`. It carries credentials (`POST /login`, `POST /register`), access and refresh tokens, the contact
data stored in `app_user` and the role changes of `PATCH /users/{userId}/role`. The boundary of the gateway to the Booking service belongs to
a teammate (issue #274), so the two analyses do not overlap.
