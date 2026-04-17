# Delivery Acceptance / Project Architecture Review (v2)

Project: `repo/`  
Date: 2026-04-06  
Inspector role: Delivery Acceptance / Project Architecture Review

## Plan Checklist (Executed In Order)
- [x] 1) Mandatory Thresholds (runnable + theme alignment)
- [x] 2) Delivery Completeness (core feature coverage + 0->1 form)
- [x] 3) Engineering & Architecture Quality
- [x] 4) Engineering Details & Professionalism (incl. security priority checks)
- [x] 5) Prompt Understanding & Fitness
- [x] 6) Aesthetics (full-stack UI applicability)
- [x] 7) Test Coverage Assessment (Static Audit)

## Environment Restriction Notes / Verification Boundary
- Runtime/test command attempted in this environment failed because Go toolchain is unavailable (`go: command not found`), so dynamic verification is bounded by static audit + provided tests.
- Command attempted:
```bash
cd repo
TEST_DATABASE_URL='postgres://parkops:parkops@127.0.0.1:5432/parkops?sslmode=disable' go test -mod=mod ./unit_tests/... ./API_tests/... -v -count=1
```
- Local reproducible commands (user side):
```bash
cd repo
go version
TEST_DATABASE_URL='postgres://parkops:parkops@127.0.0.1:5432/parkops?sslmode=disable' go test -mod=mod ./unit_tests/... ./API_tests/... -v -count=1
DATABASE_URL='postgres://parkops:parkops@127.0.0.1:5432/parkops?sslmode=disable' SESSION_SECRET='dev-session-secret' ENCRYPTION_KEY='00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff' go run ./cmd/server
```
- Evidence for runnable docs: `repo/README.md:5`, `repo/README.md:9`, `repo/README.md:17`.

---

## 1) Mandatory Thresholds

### 1.1 Deliverable can run and be verified
- **Conclusion**: Partially Pass
- **Reason (basis)**: startup/test instructions exist and are concrete, but runtime was not executable in this environment due missing `go` binary (environment boundary, not project defect).
- **Evidence**: `repo/README.md:5`, `repo/README.md:9`, `repo/README.md:17`, `repo/internal/config/config.go:26`, `repo/internal/config/config.go:29`, `repo/internal/config/config.go:35`.
- **Reproducible verification**:
  1) Run README commands above.
  2) Expected: `/api/health` returns 200; tests run if local PostgreSQL + Go are available.

### 1.3 Prompt-theme deviation check
- **Conclusion**: Pass
- **Reason (basis)**: implementation remains centered on ParkOps domain (reservations/capacity/devices/tracking/notifications/campaigns/segments/analytics/security/audit).
- **Evidence**: `repo/internal/server/router.go:27`, `repo/internal/server/router.go:35`, `repo/internal/server/reservation_handlers.go:29`, `repo/internal/server/device_handlers.go:28`, `repo/internal/server/tracking_handlers.go:22`, `repo/internal/server/campaign_handlers.go:23`, `repo/internal/server/segment_handlers.go:24`, `repo/internal/server/analytics_handlers.go:24`.
- **Reproducible verification**: inspect route groups and call representative endpoints after login.

---

## 2) Delivery Completeness

### 2.1 Core Prompt requirements coverage

1) Reservations/capacity/oversell/hold lifecycle  
- **Conclusion**: Pass  
- **Reason**: atomic lock + overlap usage computation + expiration release + confirm/cancel release are implemented.
- **Evidence**: `repo/internal/server/reservation_handlers.go:94`, `repo/internal/server/reservation_handlers.go:101`, `repo/internal/server/reservation_handlers.go:110`, `repo/internal/server/reservation_handlers.go:115`, `repo/internal/server/reservation_handlers.go:227`, `repo/internal/server/reservation_handlers.go:1059`.
- **Verification**: `go test -run TestCreateHoldOversellAndConcurrentOversell ./API_tests -v`.

2) Device integration/idempotency/out-of-order/replay  
- **Conclusion**: Pass  
- **Reason**: unique event_key dedupe, reorder classification window, replay counters, late-event handling all present.
- **Evidence**: `repo/internal/server/device_handlers.go:277`, `repo/internal/server/device_handlers.go:314`, `repo/internal/server/device_handlers.go:332`, `repo/internal/server/device_handlers.go:461`, `repo/internal/devices/logic.go:5`.
- **Verification**: `go test -run 'TestIngestDuplicateEventKey|TestLateEventFlag|TestReplayEventAndDuplicateReplay' ./API_tests -v`.

3) Tracking smoothing/stop detection/trusted timestamps  
- **Conclusion**: Pass  
- **Reason**: suspect-jump detection, confirmation/discard flow, 3-minute stop threshold, optional signed device-time trust indicator.
- **Evidence**: `repo/internal/server/tracking_handlers.go:154`, `repo/internal/server/tracking_handlers.go:125`, `repo/internal/server/tracking_handlers.go:148`, `repo/internal/tracking/logic.go:17`, `repo/internal/tracking/logic.go:47`.
- **Verification**: `go test -run 'TestTrackingSuspectPositionHandling|TestTrackingStopEventsEndpoint|TestTrackingInvalidSignatureNotTrusted' ./API_tests -v`.

4) Notifications topics/DND/frequency cap/retries/export package  
- **Conclusion**: Pass  
- **Reason**: topic subscriptions, DND settings, per-booking daily cap, retry/backoff, exportable package endpoints implemented.
- **Evidence**: `repo/internal/server/notification_handlers.go:29`, `repo/internal/server/notification_handlers.go:143`, `repo/internal/notifications/service.go:97`, `repo/internal/notifications/logic.go:35`, `repo/internal/notifications/logic.go:39`, `repo/internal/server/notification_handlers.go:332`.
- **Verification**: `go test -run 'TestNotificationDNDSettings|TestNotificationFrequencyCapAndPersistence' ./API_tests -v`.

5) Campaign/task + reminders until complete  
- **Conclusion**: Pass  
- **Evidence**: `repo/internal/server/campaign_handlers.go:34`, `repo/internal/server/campaign_handlers.go:39`, `repo/internal/campaigns/service.go:31`, `repo/internal/campaigns/service.go:94`.
- **Verification**: `go test -run TestCampaignTaskReminderStopsAfterComplete ./API_tests -v`.

6) Tags/segments preview/run/nightly/import-export rollback  
- **Conclusion**: Pass  
- **Evidence**: `repo/internal/server/segment_handlers.go:44`, `repo/internal/server/segment_handlers.go:45`, `repo/internal/server/segment_handlers.go:49`, `repo/internal/server/segment_handlers.go:50`, `repo/internal/segments/service.go:203`, `repo/internal/segments/service.go:229`.
- **Verification**: `go test -run 'TestSegmentCRUDAndPreviewRun|TestTagExportImport' ./API_tests -v`.

7) Analytics pivots + CSV/Excel/PDF exports + role+segment restrictions  
- **Conclusion**: **Partially Pass (issue exists)**  
- **Reason**: code only accepts `csv`, while Prompt requires CSV/Excel/PDF; segment restriction is role-only (admin-only when segment_id given), not role AND segment membership.
- **Evidence**: `repo/migrations/000014_analytics_exports.up.sql:4`, `repo/internal/server/analytics_handlers.go:272`, `repo/internal/server/analytics_handlers.go:355`, `repo/internal/server/analytics_handlers.go:289`.
- **Verification**:
  - POST `/api/exports` with `{"format":"pdf","scope":"bookings"}` -> current expected `400 invalid format`.
  - Compare with Prompt export requirement.

8) Security controls (password/lockout/session/RBAC/audit)  
- **Conclusion**: Pass (with one crypto-model caveat listed in issues)
- **Evidence**: `repo/internal/platform/security/password.go:61`, `repo/internal/auth/types.go:7`, `repo/internal/auth/types.go:8`, `repo/internal/auth/service.go:134`, `repo/internal/server/auth_middleware.go:42`, `repo/migrations/000016_audit_log_hash_chain.up.sql:1`.
- **Verification**: `go test -run 'TestLoginLockoutAfterFiveFails|TestSessionTimeout|TestDispatchRoleForbiddenEndpointsAndAuditLog' ./API_tests ./unit_tests -v`.

### 2.2 Basic delivery form (0->1), hardcode/mock, documentation
- **Conclusion**: Pass
- **Reason**: full project structure, migrations, tests, docs are present; no evidence of placeholder-only scaffold.
- **Evidence**: `repo/README.md:1`, `repo/go.mod:1`, `repo/migrations/000001_initial_schema.up.sql:1`, `repo/API_tests/auth_api_test.go:29`, `repo/unit_tests/security_test.go:10`.
- **Verification**: repo tree inspection + test command from README.

---

## 3) Engineering & Architecture Quality

### 3.1 Structure and module division
- **Conclusion**: Pass
- **Reason**: domain-separated packages with route registration per business domain; migrations/test suites are organized.
- **Evidence**: `repo/internal/server/router.go:24`, `repo/internal/server/router.go:35`, `repo/internal/auth/service.go:23`, `repo/internal/notifications/service.go:13`, `repo/internal/segments/service.go:15`.
- **Verification**: inspect package boundaries and imports.

### 3.2 Maintainability/extensibility awareness
- **Conclusion**: Partially Pass
- **Reason**: architecture is serviceable, but some handler files are very large and combine API/validation/business + SQL concerns, raising extension cost.
- **Evidence**: `repo/internal/server/reservation_handlers.go:1` (1284 lines), `repo/internal/server/master_handlers.go:1` (820 lines), `repo/internal/server/device_handlers.go:1` (604 lines).
- **Verification**: review file lengths and mixed responsibilities.

---

## 4) Engineering Details & Professionalism

### 4.1 Error handling/logging/validation/interface design
- **Conclusion**: Pass
- **Reason**: consistent API error shape and status use, middleware recovery, request logging without body dump, key input validations.
- **Evidence**: `repo/internal/server/errors.go:3`, `repo/internal/server/middleware.go:16`, `repo/internal/server/middleware.go:28`, `repo/internal/server/reservation_handlers.go:79`, `repo/internal/server/device_handlers.go:253`.
- **Verification**: call invalid input endpoints and observe standardized error payload.

### 4.2 Real product form vs demo form
- **Conclusion**: Pass
- **Reason**: includes auth, RBAC, multiple bounded contexts, schedulers, audit, exports, test suites and UI pages.
- **Evidence**: `repo/cmd/server/main.go:56`, `repo/cmd/server/main.go:58`, `repo/cmd/server/main.go:64`, `repo/internal/server/router.go:37`.
- **Verification**: start server and navigate `/login`, `/dashboard`, `/reservations`, `/analytics`.

### Security Priority Checks (Authentication/Authorization/Data Isolation)

1) **Authentication entry points**  
- **Conclusion**: Pass  
- **Evidence**: `repo/internal/server/auth_handlers.go:49`, `repo/internal/server/auth_handlers.go:55`, `repo/internal/auth/service.go:78`, `repo/internal/auth/service.go:134`.
- **Repro idea**: call `/api/me` without cookie -> expect 401.

2) **Route-level authorization**  
- **Conclusion**: Pass  
- **Evidence**: `repo/internal/server/auth_middleware.go:42`, `repo/internal/server/auth_handlers.go:62`, `repo/internal/server/analytics_handlers.go:40`.
- **Repro idea**: dispatch calling `/api/admin/users` -> 403.

3) **Object-level authorization (resource ownership/scope)**  
- **Conclusion**: Partially Pass  
- **Reason**: many object-scope checks exist for fleet/org resources; however export+segment policy required by prompt is weaker than specified (role-only gate for segment exports).
- **Evidence**: `repo/internal/server/master_handlers.go:107`, `repo/internal/server/reservation_handlers.go:1208`, `repo/internal/server/analytics_handlers.go:289`.
- **Repro idea**: non-admin segment export is rejected, but no segment-membership enforcement path exists.

4) **Data isolation (cross-org)**  
- **Conclusion**: Pass  
- **Evidence**: `repo/internal/server/master_handlers.go:125`, `repo/internal/server/device_handlers.go:67`, `repo/internal/server/tracking_handlers.go:102`, `repo/internal/server/reservation_handlers.go:1185`.
- **Repro idea**: fleet user reads cross-org resource -> 403 or filtered out.

---

## 5) Prompt Requirement Understanding & Fitness

- **Overall Conclusion**: Partially Pass
- **Reason (basis)**: major ParkOps workflows and constraints are understood and implemented; key gaps remain where explicit Prompt statements were narrowed (export format set and segment-sharing semantics).
- **Evidence**: `repo/internal/server/analytics_handlers.go:272`, `repo/internal/server/analytics_handlers.go:355`, `repo/internal/server/analytics_handlers.go:289`, `repo/migrations/000014_analytics_exports.up.sql:4`.
- **Verification**: request `excel`/`pdf` export and segment-scoped sharing scenarios.

---

## 6) Aesthetics (Applicable to full-stack/frontend)

- **Conclusion**: Pass (basic product-grade)
- **Reason**: pages render with clear sections/navigation and functional interactions; styling is utilitarian rather than polished but coherent.
- **Evidence**: `repo/internal/web/layout.go:49`, `repo/internal/web/dashboard.go:14`, `repo/internal/web/reservations.go:19`, `repo/internal/web/analytics_page.go:15`.
- **Verification**: run server and review `/dashboard`, `/reservations`, `/capacity`, `/analytics` on desktop/mobile viewport.

---

## Issues Found (Prioritized)

### [High] Export format requirement mismatch (CSV-only)
- **Impact**: Prompt explicitly requires CSV/Excel/PDF exports; current delivery rejects Excel/PDF, causing functional acceptance gap.
- **Evidence**: `repo/migrations/000014_analytics_exports.up.sql:4`, `repo/internal/server/analytics_handlers.go:272`, `repo/internal/server/analytics_handlers.go:355`, `repo/internal/server/analytics_handlers.go:493`.
- **Minimal actionable fix**: add `excel` and `pdf` generators (or adapter layer), keep `csv` path unchanged, and add API tests for both success and download content-type.

### [High] Export-sharing policy weaker than prompt (missing segment-membership enforcement)
- **Impact**: Prompt states sharing restricted by role **and** segment membership; current code validates role and segment existence only, not membership scope.
- **Evidence**: `repo/internal/server/analytics_handlers.go:281`, `repo/internal/server/analytics_handlers.go:289`.
- **Minimal actionable fix**: introduce export access predicate `role_allowed && actor_in_segment_scope` for create/list/download paths; add negative tests for out-of-segment users.

### [Medium] Trusted timestamp secret model is weakly isolated
- **Impact**: device-time trust can be forged more easily by privileged insiders because signing secrets are operational identifiers (`device_key`, `plate_number`) that are also used/displayed as normal data.
- **Evidence**: `repo/internal/server/device_handlers.go:295`, `repo/internal/server/device_handlers.go:323`, `repo/internal/server/device_handlers.go:63`, `repo/internal/server/tracking_handlers.go:90`, `repo/internal/server/tracking_handlers.go:111`.
- **Minimal actionable fix**: store per-device/per-vehicle dedicated HMAC secrets (encrypted at rest), never return secrets in API responses.

### [Low] Nightly segment schedule not configurable in code path
- **Impact**: Prompt/assumption expects configurable nightly run time; scheduler currently hardcodes `02:00 UTC`.
- **Evidence**: `repo/internal/segments/service.go:228`, `repo/internal/segments/service.go:238`, `repo/internal/config/config.go:10`.
- **Minimal actionable fix**: add config/env for nightly hour/minute/timezone and validate at startup.

---

## Unit/API/Logging Separate Conclusions

- **Unit tests**: Present and meaningful for security/capacity/device/tracking/reconciliation/notification logic; good baseline but does not fully cover all prompt risks (see coverage section). Evidence: `repo/unit_tests/security_test.go:10`, `repo/unit_tests/capacity_test.go:120`, `repo/unit_tests/device_test.go:10`, `repo/unit_tests/tracking_test.go:13`.
- **API/integration tests**: Extensive endpoint coverage including auth/RBAC/object-scope and major flows. Evidence: `repo/API_tests/auth_api_test.go:239`, `repo/API_tests/rbac_api_test.go:24`, `repo/API_tests/reservations_api_test.go:74`, `repo/API_tests/authorization_scope_api_test.go:10`.
- **Log categorization & sensitive leakage risk**: structured `slog` request logging exists and does not log request bodies by default; no direct password/token logging found in server middleware. Evidence: `repo/internal/server/middleware.go:16`, `repo/internal/server/logger.go:8`, `repo/API_tests/auth_api_test.go:378`.

---

## Test Coverage Assessment (Static Audit)

### Test Overview
- Unit tests exist under `repo/unit_tests/` and API/integration tests under `repo/API_tests/`.
- Framework/entry: Go `testing` with `go test`; README provides executable command. Evidence: `repo/README.md:17`, `repo/go.mod:3`, `repo/API_tests/auth_api_test.go:29`.
- DB dependency behavior: tests skip when DB unavailable unless `TEST_DATABASE_URL` explicitly set. Evidence: `repo/API_tests/auth_api_test.go:32`, `repo/API_tests/auth_api_test.go:45`.

### Coverage Mapping Table

| Requirement / Risk Point | Corresponding Test Case | Key Assertion / Fixture | Coverage Judgment | Gap | Minimal Test Addition Suggestion |
|---|---|---|---|---|---|
| Auth login/lockout/session timeout | `repo/API_tests/auth_api_test.go:239`; `repo/API_tests/auth_api_test.go:322`; `repo/unit_tests/auth_test.go:214` | 200 login cookie flags, 429 on 5th fail, timeout behavior | Sufficient | None major | Add refresh-token absence regression test (if future feature added) |
| Route-level RBAC (403) | `repo/API_tests/rbac_api_test.go:24` | multiple forbidden endpoint assertions + denied audit logs | Sufficient | None major | Keep matrix updated when new routes added |
| Object-level auth / cross-org scope | `repo/API_tests/authorization_scope_api_test.go:10`; `repo/API_tests/reservations_api_test.go:310` | fleet user blocked from cross-org reservations/devices/tracking/tags | Sufficient | Export segment membership not tested | Add cross-segment export access tests |
| Data isolation list/get filtering | `repo/API_tests/master_data_api_test.go:169`; `repo/API_tests/analytics_api_test.go:97` | cross-org 403 + ownership filtering for exports | Basic Coverage | Segment membership dimension absent | Add list/download filtering tests by segment |
| Capacity hold happy path + oversell/concurrency | `repo/API_tests/reservations_api_test.go:74`; `repo/API_tests/reservations_api_test.go:185`; `repo/unit_tests/capacity_test.go:120` | availability deltas and conflict checks in race | Sufficient | None major | Add transactional rollback assertion for partial failure path |
| Hold expiry / confirm expired | `repo/API_tests/reservations_api_test.go:151`; `repo/unit_tests/capacity_test.go:151` | expired hold cannot confirm; capacity restored | Sufficient | None major | Add delayed-expiry reconciliation interaction case |
| Device idempotency/out-of-order/replay | `repo/API_tests/devices_api_test.go:116`; `repo/API_tests/devices_api_test.go:143`; `repo/unit_tests/device_test.go:17` | duplicate event handling, replay count logic, late flag | Sufficient | None major | Add mixed-device same `event_key` negative case |
| Tracking drift/stop/trusted timestamp | `repo/API_tests/tracking_api_test.go:66`; `repo/API_tests/tracking_api_test.go:109`; `repo/unit_tests/tracking_test.go:27` | suspect confirmation flow, stop events, signature trust | Sufficient | Secret-robustness not tested | Add test asserting dedicated non-public signing secret |
| Notifications DND/frequency/retry | `repo/API_tests/notifications_api_test.go:139`; `repo/API_tests/notifications_api_test.go:160`; `repo/unit_tests/notifications_test.go:33` | DND config, cap behavior, backoff logic | Basic Coverage | Critical-vs-noncritical policy not explicitly tested | Add notification priority bypass tests |
| Analytics export formats CSV/Excel/PDF | `repo/API_tests/analytics_api_test.go:46`; `repo/API_tests/analytics_api_test.go:155` | CSV create/download works | **Insufficient** | Excel/PDF required by prompt are untested and currently unsupported | Add `excel`/`pdf` create+download tests with content type checks |
| 401/404/409 exception paths | `repo/API_tests/auth_api_test.go:371`; `repo/API_tests/router_api_test.go:34`; `repo/API_tests/reservations_api_test.go:215` | unauthorized after logout, not found error shape, capacity conflict | Basic Coverage | 404/409 not systematic for all major resources | Add table-driven error-code contract tests |
| Logs & sensitive info leakage | `repo/API_tests/auth_api_test.go:378`; `repo/internal/server/middleware.go:16` | reset response excludes token; request logs omit body | Basic Coverage | no dedicated test for log redaction under errors | Add logger redaction test harness |

### Security Coverage Audit (Mandatory Focus)
- **Authentication**: Covered; login, lockout, logout, forced password change, session timeout have direct tests. (`repo/API_tests/auth_api_test.go:239`, `repo/unit_tests/auth_test.go:214`)
- **Route Authorization**: Covered; role-based 403 matrix exists with audit verification. (`repo/API_tests/rbac_api_test.go:24`, `repo/API_tests/rbac_api_test.go:155`)
- **Object-level Authorization**: Mostly covered for org-scoped resources; export segment-membership policy missing. (`repo/API_tests/authorization_scope_api_test.go:10`, `repo/internal/server/analytics_handlers.go:289`)
- **Data Isolation**: Covered for fleet cross-org reads/writes in core domains; analytics segment scope gap remains. (`repo/API_tests/master_data_api_test.go:169`, `repo/API_tests/authorization_scope_api_test.go:50`)

### Mock/Stub Handling
- No payment integration in scope; no payment mock issue applicable.
- Tests use real DB-backed API flow (not pure mocked HTTP business logic), with skip behavior if DB unavailable. Evidence: `repo/API_tests/auth_api_test.go:40`, `repo/API_tests/auth_api_test.go:55`.

### Overall Judgment on “Can tests catch the vast majority of problems?”
- **Conclusion**: Partially Pass
- **Boundary**:
  - Covered well: auth/RBAC/capacity/device/tracking/notification base flows and several high-risk exception paths.
  - Not sufficiently covered: explicit Prompt export requirements (Excel/PDF), segment-membership sharing control, and cryptographic trust-secret robustness.
  - Risk implication: tests can pass while severe acceptance defects still exist in export capability/access policy.
