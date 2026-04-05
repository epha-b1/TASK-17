# Required Document Description: Business Logic Questions Log

## 1) Capacity hold timeout
Question: The prompt sets a 15-minute hold timeout, but does not define whether timeout is global or configurable by zone.
My Understanding/Hypothesis: Timeout should be configurable per zone, with 15 minutes as the default.
Solution: Add `hold_timeout_minutes` on zones with default `15`; reservation hold logic reads from zone-level config.

## 2) Out-of-order device events after reorder window
Question: The prompt defines a 10-minute reorder window but does not define behavior for events arriving later.
My Understanding/Hypothesis: Events arriving after 10 minutes should still be accepted, flagged as late, and not reordered.
Solution: Persist late events with `is_late=true`; apply in arrival order; trigger reconciliation check for impacted zone.

## 3) DND scope for notifications
Question: The prompt allows DND hours but does not specify whether DND suppresses all notification topics or only reminders.
My Understanding/Hypothesis: DND should suppress all user-facing topics except critical system alerts.
Solution: Add topic priority; notification worker delays non-critical notifications during DND and sends critical alerts immediately.

## 4) Tag version export format
Question: The prompt requires tag version export for rollback, but format and fields are not specified.
My Understanding/Hypothesis: Export should be a structured JSON snapshot with member, tag, and timestamp data.
Solution: Implement JSON export/import schema including `member_id`, `tags[]`, `exported_at`; enforce validation and audit log on import.

## 5) Analytics export row limits
Question: The prompt requests CSV/Excel/PDF exports but does not define size limits.
My Understanding/Hypothesis: CSV can be uncapped; Excel and PDF should have practical safety limits.
Solution: Keep CSV unlimited; cap Excel at 1,048,576 rows and PDF at 10,000 rows; return `truncated=true` metadata when capped.

## 6) Fleet manager organization boundary
Question: The prompt does not explicitly define whether Fleet Managers can view vehicles across organizations.
My Understanding/Hypothesis: Fleet Managers should only access records within their own organization.
Solution: Enforce `organization_id` scoping in all fleet queries and mutations; return 403 for cross-org access attempts.

## 7) Reconciliation tolerance
Question: The prompt requires compensating actions for snapshot mismatch, but does not define tolerance threshold.
My Understanding/Hypothesis: Any non-zero mismatch should be corrected to keep capacity accuracy strict.
Solution: Reconciliation generates compensating release for all non-zero deltas and records each correction in audit log.

## 8) GPS drift smoothing rule
Question: The prompt requires smoothing sudden GPS jumps but does not define algorithm.
My Understanding/Hypothesis: A threshold-based two-point confirmation approach is sufficient for MVP.
Solution: Mark jump as suspect if distance/time exceeds threshold; accept only when next point confirms, otherwise discard suspect point.

## 9) Rate plan application timing
Question: The prompt lists rate plans but does not define whether pricing is calculated at creation or confirmation.
My Understanding/Hypothesis: Price should be calculated at reservation creation and stored as a snapshot for consistency.
Solution: Compute and persist `applied_rate` and `estimated_total` at create-time; use stored snapshot for downstream payment/display.

## 10) Arrears source of truth
Question: The prompt mentions arrears reminders but does not define how overdue balances are tracked.
My Understanding/Hypothesis: Arrears should come from a ledger-backed member balance, not manual notes only.
Solution: Add member balance ledger entries (`charge`, `payment`, `adjustment`); arrears reminder rule triggers when balance exceeds threshold.

## 11) Meaning of message rules
Question: The prompt says admins manage message rules but does not distinguish them from notification triggers.
My Understanding/Hypothesis: Message rules are trigger definitions mapping events to template + topic.
Solution: Create `message_rules` with `event_type`, `topic`, `template_id`, `active`; dispatcher evaluates rules per domain event.

## 12) Campaign versus task model
Question: The prompt says campaign/task area, but entity boundary and lifecycle are unclear.
My Understanding/Hypothesis: Campaign is a container; tasks are executable items under a campaign.
Solution: Keep separate `campaigns` and `tasks` tables; campaign handles targeting/window, task handles status/deadline/reminders.

## 13) Incremental UI update mechanism
Question: The prompt asks for incremental UI updates but does not require SSE, WebSocket, or polling explicitly.
My Understanding/Hypothesis: Polling is acceptable for MVP and lower complexity in offline-first environment.
Solution: Implement 10-15s polling endpoints first; design event stream abstraction so SSE can be added later without API breakage.

## 14) SMS/email export package structure
Question: The prompt requires exportable SMS/email packages for manual handling, but expected structure is unspecified.
My Understanding/Hypothesis: CSV and JSON both should be supported for operational portability.
Solution: Provide package export with `recipient`, `channel`, `subject`, `body`, `created_at`, `reference_id`; include checksum and export metadata.

## 15) Audit log action coverage
Question: Prompt lists only a subset of actions that must be audited.
My Understanding/Hypothesis: Security and business-critical lifecycle events should all be audit logged.
Solution: Audit login, auth failures, role changes, reservation lifecycle, tag operations, exports/imports, event replay, reconciliation, admin actions.

## 16) Encryption scope for sensitive fields
Question: Prompt says some fields are encrypted at rest, but exact field set and method boundaries are unclear.
My Understanding/Hypothesis: Passwords should be hashed only; tokens and sensitive notes should be encrypted.
Solution: Use Argon2id for passwords; AES-256-GCM encryption for API tokens and contact notes; key from environment-managed secret.

## 17) Multi-role user support
Question: Prompt defines four roles but does not state whether users can hold multiple roles.
My Understanding/Hypothesis: Multi-role assignment is needed for operational flexibility.
Solution: Implement `user_roles` join table with many-to-many mapping; authorization checks aggregate permissions across assigned roles.

## 18) Offline buffering ownership
Question: Prompt mentions offline buffering/retransmission but does not specify device-side versus server-side buffering responsibilities.
My Understanding/Hypothesis: Buffering should be device-side; server handles idempotent ingest.
Solution: Keep single ingest endpoint with idempotency key; accept retransmitted events and deduplicate on `(device_id, event_id)`.

## 19) Signed device time definition
Question: Prompt references signed device time but does not define signature scheme.
My Understanding/Hypothesis: HMAC payload signing with per-device shared secret is intended.
Solution: Verify HMAC on inbound payload; store both server timestamp and verified device timestamp with trust flag.

## 20) Nightly segment run time
Question: Prompt says segments run nightly but does not specify schedule granularity.
My Understanding/Hypothesis: System-wide default schedule is sufficient for MVP.
Solution: Configure nightly segment job at `02:00` local time via scheduler config; allow global override via environment setting.

## 21) Password reset in offline deployment
Question: Prompt does not define reset flow when email-based recovery is unavailable.
My Understanding/Hypothesis: Reset should be admin-initiated with forced password change at next login.
Solution: Add admin reset action generating temporary password + `force_password_change=true`; require user update on first authenticated session.

## 22) Concurrent session policy
Question: Prompt defines inactivity timeout but not whether multiple active sessions are allowed.
My Understanding/Hypothesis: Multiple sessions should be allowed, each with independent inactivity expiration.
Solution: Track sessions by token/session ID; expire per session; provide admin endpoint to revoke all sessions for a target user.

## 23) Exception acknowledgement state flow
Question: Prompt says operators monitor exceptions but does not define acknowledgement semantics.
My Understanding/Hypothesis: Acknowledgement should move exception from active queue to tracked state, with optional reopen.
Solution: Use states `open -> acknowledged -> resolved` and allow `acknowledged -> open` on recurrence; log actor and note.

## 24) Editable reservation fields after creation
Question: Prompt includes reservation update capability but does not define editable fields and guard rules.
My Understanding/Hypothesis: Time window, quantity, and notes are editable until reservation is completed/cancelled.
Solution: Permit updates only in `hold` and `confirmed`; on capacity-impacting edits, rerun availability and refresh hold atomically.
