# ADR-C3-01: Run PostgreSQL as RDS primary/standby in two zones with cross-region backup copies

## Status

Proposed

## Context

All relational schemas of CSP (`auth`, `booking`, `concessions`, `ticketing_fulfillment`, `worker`, `workflow`) live in one PostgreSQL instance per environment (`05-architecture/decisions/records/ADR-006-environment-level-database-instances.md`). ADR-001 already lists the risk: if that instance fails, all services lose persistence at the same time, and it suggests a managed instance, backups and a documented restore procedure. The `booking` schema is the most critical one, because it keeps the rule that a seat is never sold twice. `13-operations/backup-and-recovery.md` proposes for Booking an RTO of 30 minutes and an RPO of 5 minutes, and states that no production environment exists and no restore test has been run. At the moment the data lives in Docker volumes, which that page says are not backups.

## Decision

Host the PostgreSQL instance on Huawei Cloud RDS for PostgreSQL as a primary in AZ1 and a standby in AZ2 with semi-synchronous replication and automatic failover. Take incremental log backups every 5 minutes plus a daily full backup, and copy them to a second region with Cloud Backup and Recovery. Keep the targets RTO 30 minutes and RPO 5 minutes, and run the quarterly restore test that the backup page requires. The `-db` repositories remain the source of truth of the schemas, so a restored schema must match the released tag.

## Alternatives considered

- **Daily snapshot only.** Rejected: it gives an RPO of up to 24 hours, far from the 5 minutes needed to protect confirmed bookings.
- **Self-managed PostgreSQL with streaming replication on ECS.** Rejected: the four-person team would own patching, failover scripts and backup tooling, and nobody has run a restore even once.
- **Active-active databases in two regions.** Rejected: it reaches an RPO near zero but multiplies the cost and adds conflict handling that the seat-locking transactions do not need for this MVP.

## Consequences

- Positive: losing one zone costs no committed data and recovery needs no manual work, which fits the 30-minute RTO.
- Positive: the backup policy becomes concrete and testable instead of a proposal.
- Negative: a standby and cross-region backup storage roughly double the database cost compared with a single node, which weighs on the cost principle.
- Negative: semi-synchronous replication adds a little commit latency, so the P95 target of 300 ms for seat holds has to be checked in the qa load test.
