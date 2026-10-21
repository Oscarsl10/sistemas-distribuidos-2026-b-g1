# CP2 - Structured log contract

Scope: the services of the registration journey of [01-journey.md](01-journey.md): `csp-api-gateway`, `csp-auth-api`,
`csp-worker` (job `auth-outbox-relay`) and `csp-ticketing-api`. Real elements used here: `POST /api/v1/auth/register`,
the event `UserRegistered`, the tables `app_user`, `outbox_event` and `user_contacts`, the queue
`cine.ticketing_fulfillment.user-registered`.

Rule: one JSON object per line on stdout, the same field names in every service. Today `csp-auth-api` already writes JSON
with `log/slog` and the gateway writes a JSON access log; this contract adds the missing fields.

## 1. Fields

| Field | Type | Required | Example | Why it exists |
|---|---|---|---|---|
| `timestamp` | string, ISO-8601 UTC with milliseconds | yes | `2026-10-20T19:41:07.318Z` | Orders the lines of different services. |
| `level` | string: `DEBUG`, `INFO`, `WARN`, `ERROR` | yes | `ERROR` | Lets an alert or a query keep only what needs action. |
| `service` | string, repository name | yes | `csp-auth-api` | Says which service wrote the line. |
| `version` | string, release of the service | yes | `3.0.0` | Links a failure with the release that was running. |
| `environment` | string: `dev`, `qa`, `main` | yes | `qa` | The three environments of the platform write to similar logs. |
| `trace_id` | string, 32 hex | yes | `9c1e7a42d5b84f03a6e1b27d38f0c915` | Joins every line of one request, also across RabbitMQ. |
| `span_id` | string, 16 hex | yes | `5d2f0b8c71a94e36` | Says which step of the trace wrote the line. |
| `correlation_id` | string | yes | `0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90` | The existing `X-Correlation-Id` of the gateway; the same value as the `traceId` of an error body. |
| `event` | string, `domain.action.result` | yes | `outbox.event.published` | A stable name that can be counted and alerted on. |
| `message` | string | yes | `UserRegistered published` | Short sentence for a person. No personal data. |
| `duration_ms` | integer | yes | `412` | Latency of the step; `0` when it does not apply. |
| `user_id` | string, UUID of `app_user.id` | when known | `7a3d9e10-52c4-4b8e-a1f6-0d94c2e7b385` | Domain: follows one account without writing the email. |
| `event_id` | string, UUID | for events | `e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13` | Domain: the `eventId` of `outbox_event`, also the AMQP `messageId`; consumers deduplicate by it. |
| `delivery_attempt` | integer | for events | `2` | Domain: how many times the relay or the consumer tried the event (`attemptCount` of the relay). |
| `http_status` | integer | for HTTP | `201` | Domain: the answer that the client saw. |
| `idempotency_key_hash` | string, first 12 hex of SHA-256 | for register | `3f9a1c7e02bd` | Domain: finds a repeated `Idempotency-Key` without writing its value. |

## 2. Never log

| Data | Where it appears | What is written instead |
|---|---|---|
| Password | body of `POST /register` and `/login` | nothing |
| bcrypt hash | `app_user.password_hash` | nothing |
| Access token, refresh token | `Authorization` header, `/login` and `/refresh` answers | nothing; only `user_id` |
| Value of `Idempotency-Key` | header of `POST /register` | `idempotency_key_hash` |
| Email | `app_user.email`, payload of `UserRegistered` | `user_id`; if a person needs it: `j***@example.com` |
| Phone and address | stored by `app_user` (ADR-024), never published | nothing |
| Verification code of the email change | `email_verification` | nothing |
| Service credential, database and broker passwords, Redis URL | environment variables | nothing; only the name of the variable |

## 3. Levels

| Level | When | Example of this project |
|---|---|---|
| `ERROR` | The step failed and a person or a retry has to act. | `csp-ticketing-api` cannot write `user_contacts` because PostgreSQL does not answer in time. |
| `WARN` | Something was refused or degraded, the system keeps working. | A refresh token is rejected; the revoked-token list in Redis is unavailable. |
| `INFO` | A business step that finished normally. | `UserRegistered` published to `cine.events`; `POST /register` answered `201`. |
| `DEBUG` | Detail for a developer, off in `qa` and `main`. | Duration of each SQL statement of the registration transaction, without values. |

## 4. One request in seven lines

The registration of one client. The first two lines are the synchronous part, the other five are the asynchronous part
(the relay runs every 5 seconds). All seven share the same `trace_id`. Two lines show a failure: a timeout of a downstream
dependency (line 5) and a rejection by a rule of the consumer (line 6).

```json
{"timestamp":"2026-10-20T19:41:07.301Z","level":"INFO","service":"csp-auth-api","version":"3.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"5d2f0b8c71a94e36","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"user.registered","message":"user saved with its outbox event","duration_ms":388,"user_id":"7a3d9e10-52c4-4b8e-a1f6-0d94c2e7b385","event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","idempotency_key_hash":"3f9a1c7e02bd","http_status":201}
{"timestamp":"2026-10-20T19:41:07.318Z","level":"INFO","service":"csp-api-gateway","version":"3.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"1a8c4e92b6d70f35","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"http.request.completed","message":"POST /api/v1/auth/register","duration_ms":412,"http_status":201}
{"timestamp":"2026-10-20T19:41:10.044Z","level":"INFO","service":"csp-worker","version":"3.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"c07e3b95a1d24f68","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"outbox.event.published","message":"UserRegistered published to cine.events","duration_ms":21,"event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","delivery_attempt":1}
{"timestamp":"2026-10-20T19:41:10.071Z","level":"INFO","service":"csp-ticketing-api","version":"1.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"e9b4d2a07c3f1586","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"message.received","message":"UserRegistered received from cine.ticketing_fulfillment.user-registered","duration_ms":0,"event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","delivery_attempt":1}
{"timestamp":"2026-10-20T19:41:13.082Z","level":"ERROR","service":"csp-ticketing-api","version":"1.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"e9b4d2a07c3f1586","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"contact.projection.timeout","message":"write of user_contacts did not answer in time, message requeued","duration_ms":3008,"event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","delivery_attempt":1}
{"timestamp":"2026-10-20T19:41:13.540Z","level":"WARN","service":"csp-ticketing-api","version":"1.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"2f6a9c1d84b03e57","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"message.rejected.duplicate","message":"event already processed, rejected without a second write","duration_ms":6,"event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","delivery_attempt":2}
{"timestamp":"2026-10-20T19:41:13.547Z","level":"INFO","service":"csp-ticketing-api","version":"1.0.0","environment":"qa","trace_id":"9c1e7a42d5b84f03a6e1b27d38f0c915","span_id":"2f6a9c1d84b03e57","correlation_id":"0b6f2c1e-48d3-4a77-9e51-3c7d1f2a8b90","event":"message.acknowledged","message":"message acknowledged after the duplicate check","duration_ms":1,"event_id":"e41b6c07-93a2-4d58-b0f1-6a2c8d5e7f13","delivery_attempt":2}
```

How the story reads when the lines are filtered by `trace_id` and sorted by `timestamp` (the gateway writes its line
when the request ends, so it comes after the line of the service it called):

1. The user and the event were saved in one transaction and the registration was answered `201` (lines 1 and 2).
2. Three seconds later the relay published it and the consumer received it (lines 3 and 4).
3. The write of `user_contacts` timed out and the message was requeued (line 5), but the first write had in fact been done.
4. The second delivery was recognised by its `event_id` and rejected as a duplicate, then acknowledged (lines 6 and 7).
