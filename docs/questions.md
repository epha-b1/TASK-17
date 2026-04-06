# Required Document Description: Business Logic Questions Log

This file records business-level ambiguities from the prompt and the implementation decisions taken.
Each entry follows: **Question + My Understanding/Hypothesis + Solution**.

## 1) [Business Rule] Inventory rollback after cancellation
Question: The prompt allows cancellation/expiration but does not explicitly define the rollback behavior for reserved stalls.
My Understanding/Hypothesis: Any cancellation or expiration must release held/consumed capacity immediately to prevent artificial shortage.
Solution: Implemented capacity release as part of reservation lifecycle (`cancelled`/`expired`) and reconciliation compensating events for missed late cases.

## 2) [Boundary Condition] Hold timeout scope and default
Question: The prompt says default hold timeout is 15 minutes, but does not state whether this is global or zone-specific.
My Understanding/Hypothesis: 15 minutes is a default value and zones may override it when needed.
Solution: Added/used `hold_timeout_minutes` on zones with default `15`; hold expiration logic reads zone configuration.

## 3) [Boundary Condition] Out-of-order events beyond 10-minute window
Question: Prompt defines correction for out-of-order arrivals within 10 minutes but does not define behavior after the window.
My Understanding/Hypothesis: Events arriving after the reorder window should still be recorded, then corrected through reconciliation, not dropped.
Solution: Persist late arrivals, avoid unsafe reorder beyond window, and rely on reconciliation run + compensating release to restore consistency.

## 4) [Business Rule] DND policy versus critical alerts
Question: The prompt supports DND hours but does not clarify whether DND blocks every notification type.
My Understanding/Hypothesis: DND suppresses normal reminders but should not block critical operational alerts.
Solution: Notification dispatch distinguishes critical vs non-critical topics; non-critical respects DND/frequency cap, critical delivers immediately.

## 5) [Business Process] Frequency cap counting window
Question: The prompt states “no more than 3 reminders per booking per day” but does not define the counting key.
My Understanding/Hypothesis: Count should be keyed by `(booking_id, topic, local-day)` to avoid cross-booking interference.
Solution: Added dedupe/counter logic in reminder scheduling and dispatch checks before enqueue/send.

## 6) [Data Relationship] Campaign vs task entity boundary
Question: The prompt mentions a campaign/task area but does not clearly separate responsibilities.
My Understanding/Hypothesis: Campaign is the parent container; task is the executable/remindable work item.
Solution: Keep separate campaign and task models; campaign stores context/targeting window, task stores status/deadline/reminder behavior.

## 7) [Data Relationship] Message rules vs notification topics
Question: Prompt includes “message rules” and subscribable topics but does not define the relation.
My Understanding/Hypothesis: Message rule maps domain event to topic + template used by notification pipeline.
Solution: Model message rules as trigger definitions and evaluate them on business events before creating notification jobs.

## 8) [Boundary Condition] Segment execution schedule
Question: The prompt requires on-demand and nightly segment runs but does not define exact nightly timing.
My Understanding/Hypothesis: One configurable nightly schedule is sufficient for MVP offline deployment.
Solution: Scheduler runs nightly segment evaluation at configured system time; supports on-demand execution from operator action.

## 9) [Business Rule] Export sharing restrictions
Question: Prompt says exports are restricted by role and segment membership but does not define precedence.
My Understanding/Hypothesis: Access must pass both checks (role authorization AND segment scope), not either/or.
Solution: Enforced conjunctive authorization in export endpoints and generated data scope filters accordingly.

## 10) [Boundary Condition] Replay without double counting
Question: Prompt requires controlled replay and idempotency but does not define stable dedupe key.
My Understanding/Hypothesis: Dedupe must use immutable event identity (`device_id + event_key`) independent of replay attempts.
Solution: Enforced unique event key per device and idempotent ingest path so replays do not alter counts twice.

## 11) [Business Process] Offline buffering ownership
Question: Prompt requires offline buffering/retransmission but does not specify whether buffering occurs server-side or on device.
My Understanding/Hypothesis: Device/client side performs buffering; server focuses on idempotent acceptance and reconciliation.
Solution: Kept ingest API tolerant to delayed retransmissions; server stores accepted events and deduplicates reliably.

## 12) [Business Rule] Signed device time trust model
Question: Prompt requires trusted timestamp evidence using server time plus signed device time when available, but does not define trust fallback.
My Understanding/Hypothesis: Server time is always authoritative baseline; signed device time is secondary evidence when signature verification succeeds.
Solution: Store both timestamps and a trust indicator; verification failure keeps event usable with server timestamp only.

## 13) [Business Rule] Session expiration and concurrent sessions
Question: Prompt specifies 30-minute inactivity timeout but does not explicitly define concurrent login behavior.
My Understanding/Hypothesis: Multiple sessions may exist, each independently expired by inactivity; admin can revoke sessions when needed.
Solution: Session model tracks `last_active_at` and per-session expiry; middleware enforces inactivity timeout on every authenticated request.

## 14) [Data Relationship] Multi-role users
Question: Prompt defines four roles but does not state whether users are single-role or multi-role.
My Understanding/Hypothesis: Multi-role support is needed for real operations and admin flexibility.
Solution: Use many-to-many `user_roles`; permission checks aggregate user capabilities across assigned roles.

## 15) [Boundary Condition] Password reset without external email
Question: System is offline-first/local network, but prompt does not define forgotten-password recovery flow without external mail/SMS.
My Understanding/Hypothesis: Reset must be admin-initiated locally with forced password change on next login.
Solution: Admin reset path sets temporary credentials and `force_password_change=true`, then user updates password on first authenticated session.

## 16) [Business Rule] Encryption scope by field type
Question: Prompt says “sensitive fields encrypted at rest” while also requiring password hashes, which are not reversible encryption.
My Understanding/Hypothesis: Passwords are hashed only; reversible encryption applies to API tokens and contact notes.
Solution: Use Argon2id for password hashes and AES-256-GCM at rest for token/notes fields with environment-managed key.

## 17) [Business Process] Audit log coverage for tamper-evident history
Question: Prompt names some audited actions (tag changes, exports, replay) but does not define complete minimum scope.
My Understanding/Hypothesis: All security-sensitive and business-critical mutations should be audited for forensic completeness.
Solution: Audit events include auth events, role/permission changes, reservation lifecycle mutations, tag changes, exports/imports, and replay/reconciliation actions.
