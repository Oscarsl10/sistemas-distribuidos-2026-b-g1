# CP1 - Critical journey and trace map

Journey: **a client registers an account** (HU-AUTH-001): `POST /api/v1/auth/register`, from the sign-in portal to the
contact projection of Ticketing & Fulfillment. It crosses five components and one asynchronous boundary.

Sources checked in the team repositories: `csp-docs` (`07-api/contracts/openapi/auth-service.yaml`,
`09-microservices/event-catalog.md`, `09-microservices/services/02-auth-service/events.md`, ADR-014, ADR-019, ADR-036),
`csp-api-gateway` (`nginx/routes/auth.conf`, `nginx/nginx.conf`), `csp-auth-api` (`register_user.go`),
`csp-worker` (`AmqpEventPublisher.java`, section "Auth outbox relay" of the README).

## 1. Hop table

| # | Caller | Callee | Protocol | Operation | Sync/Async | Data |
|---|---|---|---|---|---|---|
| 1 | `csp-auth-portal` (inside the `csp-front` shell) | `csp-api-gateway` | REST | `POST /api/v1/auth/register` with `Idempotency-Key` | sync | none |
| 2 | `csp-api-gateway` | `csp-auth-api` | REST | `POST /register` (public route, rate limit 10 req/s per client) | sync | none |
| 3 | `csp-auth-api` | PostgreSQL schema `auth` | SQL | one transaction: `app_user`, `user_role`, `idempotency_key`, `outbox_event` (`UserRegistered`) | sync | `auth` database |
| 4 | `csp-worker` (job `auth-outbox-relay`, every 5 s) | RabbitMQ exchange `cine.events` | AMQP | publish `UserRegistered`, routing key `cine.auth.user-registered.v1` | async | `auth.outbox_event` (read only), `worker.outbox_relay_state` |
| 5 | `csp-ticketing-api` | RabbitMQ queue `cine.ticketing_fulfillment.user-registered` | AMQP | consume `UserRegistered` | async | `ticketing_fulfillment.user_contacts` |

**Critical hop: hop 3.** It is where the user waits (bcrypt hash of the password plus the transaction on four tables),
it is the only hop that can lose data, and it is the one the incident of CP9 targets.

Facts found while checking the repositories, which shape the trace design:

- The gateway creates or reuses `X-Correlation-Id`, writes it in its JSON access log (`correlationId`) and forwards it to
  `csp-auth-api`, which returns it as `traceId` in every error body. There is no `traceparent` anywhere yet.
- The relay publishes with the header `X-Correlation-Id`, but the Auth envelope carries no `metadata`, so the relay uses
  **the id of its own run**, not the one of the request. The link between hop 3 and hop 4 is broken today: the
  context has to be stored in the outbox row (see "Gap" below).
- The Auth relay is off by default (`AUTH_OUTBOX_RELAY_ENABLED=false`) and `csp-ticketing-api` has only its README in
  `develop`: hops 4 and 5 are specified in `csp-docs` (ADR-019, ADR-036) but not running yet.

## 2. Sequence diagram and trace context

```mermaid
sequenceDiagram
    participant P as csp-auth-portal
    participant G as csp-api-gateway
    participant A as csp-auth-api
    participant D as PostgreSQL (schema auth)
    participant W as csp-worker (auth-outbox-relay)
    participant M as RabbitMQ (cine.events)
    participant T as csp-ticketing-api

    P->>G: POST /api/v1/auth/register (Idempotency-Key, no trace yet)
    Note over G: creates the trace: traceparent 00-trace_id-span_g-01, X-Correlation-Id
    G->>A: POST /register (traceparent, X-Correlation-Id)
    A->>D: BEGIN; insert app_user, user_role, idempotency_key, outbox_event (traceparent saved in the row)
    D-->>A: COMMIT
    A-->>G: 201 (X-Correlation-Id, traceparent)
    G-->>P: 201 (X-Correlation-Id)
    loop every 5 s
        W->>D: SELECT auth.outbox_event (read only)
        W-)M: publish cine.auth.user-registered.v1 (headers: traceparent, X-Correlation-Id, messageId = eventId)
    end
    M-)T: deliver to cine.ticketing_fulfillment.user-registered (same headers)
    Note over T: consumer starts its span with the traceparent of the message and writes user_contacts
```

What travels on each arrow:

| Arrow | Trace context |
|---|---|
| portal -> gateway | nothing (the browser does not start a trace in this design) |
| gateway -> auth-api | HTTP headers `traceparent` (W3C) and `X-Correlation-Id` |
| auth-api -> PostgreSQL | no header: the context is stored as a column value of `outbox_event` |
| worker -> RabbitMQ | AMQP message headers `traceparent` and `X-Correlation-Id`, read from the outbox row |
| RabbitMQ -> ticketing-api | the same AMQP headers; the consumer continues the trace |

## 3. Span tree

```
trace_id = one value for the whole journey
POST /api/v1/auth/register                      [csp-api-gateway]  root span
  proxy to csp-auth-api                         [csp-api-gateway]  child
    POST /register                              [csp-auth-api]     child (server span)
      validate command and idempotency key      [csp-auth-api]
      bcrypt hash of the password               [csp-auth-api]     expected slowest local step
      transaction register user                 [csp-auth-api]     critical hop
        insert app_user and user_role           [PostgreSQL]
        insert idempotency_key                  [PostgreSQL]
        insert outbox_event (UserRegistered)    [PostgreSQL]       traceparent stored here
  ---- asynchronous boundary: the request has already been answered ----
  publish UserRegistered                        [csp-worker]       child of the stored span, same trace_id
    consume UserRegistered                      [csp-ticketing-api] child, same trace_id (AMQP headers)
      insert user_contacts                      [csp-ticketing-api]
```

## 4. Gap to close before instrumenting the real services

The platform documents define `X-Correlation-Id` (ADR-014, ADR-022, the shared error envelope); this workshop uses the
W3C `traceparent`. Mapping used in the next checkpoints: `trace_id` is the 32 hexadecimal characters of `traceparent`,
and `correlation_id` stays the existing `X-Correlation-Id`, so both can be searched in the logs. To keep one trace across
hop 4 the Auth envelope or the `outbox_event` row would have to carry the context; that is a change of the contract and the
data model, so it is only proposed here and not implemented in the project.
