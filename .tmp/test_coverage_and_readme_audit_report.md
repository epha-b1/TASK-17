# Test Coverage Audit

## Scope, Method, and Project Type

- **Inspection mode:** static inspection only (no execution).
- **Project type declaration:** `fullstack` declared explicitly in `README.md:1`.
- **Inference check:** declaration is consistent with code structure (`internal/server/*` backend routes and `internal/web/*` server-rendered frontend components).

## Runtime Retest Addendum (User-requested)

- **Reason:** user explicitly requested actual retesting of each gap after fixes.
- **Executed command:** `docker compose run --rm test sh -c "cd /app && TEST_DATABASE_URL='postgres://parkops:parkops@127.0.0.1:5432/parkops?sslmode=disable' go test -mod=mod ./unit_tests/... ./API_tests/... -v -count=1"`.
- **Result:** PASS for both test suites.
  - `ok parkops/unit_tests 4.599s` (evidence: tool output line 147)
  - `ok parkops/API_tests 133.254s` (evidence: tool output line 1811)
- **Gap-route runtime evidence:**
  - `TestFormLoginRedirectsAndSetsSessionCookie` passed (`POST /auth/login`).
  - `TestFormLogoutClearsSessionAndRedirectsToLogin` passed (`POST /auth/logout`).
  - `TestPageRoutesRenderForAuthenticatedAdmin` passed and exercised `GET /dashboard`, `/capacity`, `/facilities`, `/lots`, `/zones`, `/rate-plans`, `/members`, `/vehicles`, `/drivers`, `/notifications`, `/campaigns`, `/segments`, `/tasks`, `/notification-prefs`, `/analytics`.
  - `TestAuditPageRendersForAuditor` and `TestAuditPageRendersForAdmin` passed (`GET /audit`).
  - `TestAdminUsersPageRendersForAdmin` passed (`GET /admin/users`).

## Backend Endpoint Inventory (Resolved `METHOD + PATH`)

### Auth and Admin

- `POST /api/auth/login`
- `POST /api/auth/logout`
- `POST /auth/login`
- `POST /auth/logout`
- `GET /api/me`
- `PATCH /api/me/password`
- `GET /api/admin/users`
- `POST /api/admin/users`
- `PATCH /api/admin/users/:id`
- `DELETE /api/admin/users/:id`
- `PATCH /api/admin/users/:id/roles`
- `POST /api/admin/users/:id/unlock`
- `GET /api/admin/users/:id/sessions`
- `DELETE /api/admin/users/:id/sessions`
- `POST /api/admin/users/:id/reset-password`
- `GET /api/admin/audit-logs`

### Master Data

- `GET /api/facilities`
- `GET /api/facilities/:id`
- `GET /api/lots`
- `GET /api/lots/:id`
- `GET /api/zones`
- `GET /api/zones/:id`
- `GET /api/rate-plans`
- `GET /api/rate-plans/:id`
- `GET /api/members`
- `GET /api/members/:id`
- `GET /api/members/:id/balance`
- `GET /api/vehicles`
- `GET /api/vehicles/:id`
- `GET /api/drivers`
- `GET /api/drivers/:id`
- `GET /api/message-rules`
- `POST /api/facilities`
- `PATCH /api/facilities/:id`
- `DELETE /api/facilities/:id`
- `POST /api/lots`
- `PATCH /api/lots/:id`
- `DELETE /api/lots/:id`
- `POST /api/zones`
- `PATCH /api/zones/:id`
- `DELETE /api/zones/:id`
- `POST /api/rate-plans`
- `PATCH /api/rate-plans/:id`
- `DELETE /api/rate-plans/:id`
- `PATCH /api/members/:id/balance`
- `POST /api/message-rules`
- `PATCH /api/message-rules/:id`
- `DELETE /api/message-rules/:id`
- `POST /api/members`
- `PATCH /api/members/:id`
- `DELETE /api/members/:id`
- `POST /api/vehicles`
- `PATCH /api/vehicles/:id`
- `DELETE /api/vehicles/:id`
- `POST /api/drivers`
- `PATCH /api/drivers/:id`
- `DELETE /api/drivers/:id`

### Reservations, Capacity, Exceptions, Reconciliation

- `GET /api/availability`
- `GET /api/reservations`
- `GET /api/capacity/dashboard`
- `GET /api/capacity/zones/:id/stalls`
- `GET /api/capacity/snapshots`
- `GET /api/reservations/stats/today`
- `GET /api/reservations/:id/timeline`
- `GET /api/exceptions/history`
- `GET /api/exceptions`
- `GET /api/exceptions/:id`
- `POST /api/reservations/hold`
- `POST /api/reservations/:id/confirm`
- `POST /api/reservations/:id/cancel`
- `POST /api/reconciliation/run`
- `POST /api/exceptions/:id/acknowledge`

### Devices and Tracking

- `GET /api/devices`
- `GET /api/devices/:id`
- `GET /api/device-events`
- `POST /api/device-events`
- `POST /api/device-events/replay`
- `POST /api/devices`
- `PATCH /api/devices/:id`
- `DELETE /api/devices/:id`
- `GET /api/tracking/vehicles/:id/positions`
- `GET /api/tracking/vehicles/:id/stops`
- `POST /api/tracking/location`

### Notifications

- `GET /api/notification-topics`
- `POST /api/notification-topics/:id/subscribe`
- `DELETE /api/notification-topics/:id/subscribe`
- `GET /api/notification-settings`
- `PATCH /api/notification-settings`
- `GET /api/notification-settings/dnd`
- `PATCH /api/notification-settings/dnd`
- `GET /api/notifications`
- `GET /api/notifications/:id`
- `PATCH /api/notifications/:id/read`
- `POST /api/notifications/:id/dismiss`
- `GET /api/notifications/export-packages`
- `GET /api/notifications/export-packages/:id/download`

### Campaigns, Tasks, Segments, Analytics, Exports

- `GET /api/campaigns`
- `GET /api/campaigns/:id`
- `GET /api/campaigns/:id/tasks`
- `GET /api/tasks/:id`
- `POST /api/campaigns`
- `PATCH /api/campaigns/:id`
- `DELETE /api/campaigns/:id`
- `POST /api/campaigns/:id/tasks`
- `PATCH /api/tasks/:id`
- `DELETE /api/tasks/:id`
- `POST /api/tasks/:id/complete`
- `GET /api/tags`
- `GET /api/members/:id/tags`
- `GET /api/segments`
- `GET /api/segments/:id`
- `GET /api/segments/:id/runs`
- `POST /api/tags`
- `DELETE /api/tags/:id`
- `POST /api/members/:id/tags`
- `DELETE /api/members/:id/tags/:tagId`
- `POST /api/tags/export`
- `POST /api/tags/import`
- `POST /api/segments`
- `PATCH /api/segments/:id`
- `DELETE /api/segments/:id`
- `POST /api/segments/:id/preview`
- `POST /api/segments/:id/run`
- `GET /api/analytics/occupancy`
- `GET /api/analytics/bookings`
- `GET /api/analytics/exceptions`
- `GET /api/exports`
- `GET /api/exports/:id/download`
- `POST /api/exports`
- `GET /api/health`

### Non-API Router Endpoints

- `GET /login`
- `GET /dashboard`
- `GET /reservations`
- `GET /capacity`
- `GET /facilities`
- `GET /lots`
- `GET /zones`
- `GET /rate-plans`
- `GET /members`
- `GET /vehicles`
- `GET /drivers`
- `GET /notifications`
- `GET /campaigns`
- `GET /segments`
- `GET /tasks`
- `GET /notification-prefs`
- `GET /analytics`
- `GET /audit`
- `GET /admin/users`
- `GET /swagger/*any`

## API Test Mapping Table

### Auth/Admin/Health

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `POST /api/auth/login` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestLoginSuccess` |
| `POST /api/auth/logout` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestLogoutAndMe` |
| `GET /api/me` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestLogoutAndMe` |
| `PATCH /api/me/password` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestForcePasswordChangeBlocksRoutes` |
| `GET /api/admin/users` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanManageUsersAndUpdateRoles` |
| `POST /api/admin/users` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanManageUsersAndUpdateRoles` |
| `PATCH /api/admin/users/:id` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanManageUsersAndUpdateRoles` |
| `DELETE /api/admin/users/:id` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanManageUsersAndUpdateRoles` |
| `PATCH /api/admin/users/:id/roles` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanManageUsersAndUpdateRoles` |
| `POST /api/admin/users/:id/unlock` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestAdminUnlockAccount` |
| `GET /api/admin/users/:id/sessions` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestAdminListAndDeleteUserSessions` |
| `DELETE /api/admin/users/:id/sessions` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestAdminListAndDeleteUserSessions` |
| `POST /api/admin/users/:id/reset-password` | yes | true no-mock HTTP | `API_tests/auth_api_test.go` | `TestAdminResetPasswordNoTokenInResponse` |
| `GET /api/admin/audit-logs` | yes | true no-mock HTTP | `API_tests/rbac_api_test.go` | `TestAdminCanAccessAuditLogs` |
| `GET /api/health` | yes | true no-mock HTTP | `API_tests/misc_coverage_api_test.go` | `TestHealthEndpoint` |

### Master Data

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `GET /api/facilities` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestFacilityReadUpdateDelete` |
| `GET /api/facilities/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestFacilityReadUpdateDelete` |
| `GET /api/lots` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestLotReadUpdateDelete` |
| `GET /api/lots/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestLotReadUpdateDelete` |
| `GET /api/zones` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMasterDataCRUDHappyPathAndWrongRole` |
| `GET /api/zones/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestZoneReadDelete` |
| `GET /api/rate-plans` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestRatePlanReadUpdateDelete` |
| `GET /api/rate-plans/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestRatePlanReadUpdateDelete` |
| `GET /api/members` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/members/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/members/:id/balance` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/vehicles` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/vehicles/:id` | yes | true no-mock HTTP | `API_tests/tracking_api_test.go` | `TestTrackingSigningSecretNotExposedInVehicleAPI` |
| `GET /api/drivers` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/drivers/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `GET /api/message-rules` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestMessageRuleReadUpdateDelete` |
| `POST /api/facilities` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMasterDataCRUDHappyPathAndWrongRole` |
| `PATCH /api/facilities/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestFacilityReadUpdateDelete` |
| `DELETE /api/facilities/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestFacilityReadUpdateDelete` |
| `POST /api/lots` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMasterDataCRUDHappyPathAndWrongRole` |
| `PATCH /api/lots/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestLotReadUpdateDelete` |
| `DELETE /api/lots/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestLotReadUpdateDelete` |
| `POST /api/zones` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMasterDataCRUDHappyPathAndWrongRole` |
| `PATCH /api/zones/:id` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestZoneStallReductionBlockedBelowConfirmedDemand` |
| `DELETE /api/zones/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestZoneReadDelete` |
| `POST /api/rate-plans` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMasterDataCRUDHappyPathAndWrongRole` |
| `PATCH /api/rate-plans/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestRatePlanReadUpdateDelete` |
| `DELETE /api/rate-plans/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestRatePlanReadUpdateDelete` |
| `PATCH /api/members/:id/balance` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `POST /api/message-rules` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestMessageRuleReadUpdateDelete` |
| `PATCH /api/message-rules/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestMessageRuleReadUpdateDelete` |
| `DELETE /api/message-rules/:id` | yes | true no-mock HTTP | `API_tests/master_data_coverage_api_test.go` | `TestMessageRuleReadUpdateDelete` |
| `POST /api/members` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `PATCH /api/members/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `DELETE /api/members/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `POST /api/vehicles` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `PATCH /api/vehicles/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `DELETE /api/vehicles/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `POST /api/drivers` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `PATCH /api/drivers/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |
| `DELETE /api/drivers/:id` | yes | true no-mock HTTP | `API_tests/master_data_api_test.go` | `TestMembersVehiclesDriversMessageRulesAndOrgScope` |

### Reservations / Capacity / Exceptions / Reconciliation

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `GET /api/availability` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `GET /api/reservations` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestListHoldsAndExceptionsEndpoints` |
| `GET /api/capacity/dashboard` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `GET /api/capacity/zones/:id/stalls` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `GET /api/capacity/snapshots` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `GET /api/reservations/stats/today` | yes | true no-mock HTTP | `API_tests/misc_coverage_api_test.go` | `TestReservationStatsToday` |
| `GET /api/reservations/:id/timeline` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `GET /api/exceptions/history` | yes | true no-mock HTTP | `API_tests/exceptions_api_test.go` | `TestAcknowledgeExceptionAsDispatchAndHistory` |
| `GET /api/exceptions` | yes | true no-mock HTTP | `API_tests/exceptions_api_test.go` | `TestExceptionsOpenListAndGet` |
| `GET /api/exceptions/:id` | yes | true no-mock HTTP | `API_tests/exceptions_api_test.go` | `TestExceptionsOpenListAndGet` |
| `POST /api/reservations/hold` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `POST /api/reservations/:id/confirm` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `POST /api/reservations/:id/cancel` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationFlowAndCapacityEndpoints` |
| `POST /api/reconciliation/run` | yes | true no-mock HTTP | `API_tests/reconciliation_api_test.go` | `TestManualReconciliationRunCreatesCompensatingEventsAndAuditLog` |
| `POST /api/exceptions/:id/acknowledge` | yes | true no-mock HTTP | `API_tests/exceptions_api_test.go` | `TestAcknowledgeExceptionAsDispatchAndHistory` |

### Devices / Tracking

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `GET /api/devices` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestRegisterDevice` |
| `GET /api/devices/:id` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestRegisterDevice` |
| `GET /api/device-events` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestListDeviceEvents` |
| `POST /api/device-events` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestIngestDeviceEvent` |
| `POST /api/device-events/replay` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestReplayEventAndDuplicateReplay` |
| `POST /api/devices` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `createDeviceWithZone` |
| `PATCH /api/devices/:id` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestRegisterDevice` |
| `DELETE /api/devices/:id` | yes | true no-mock HTTP | `API_tests/devices_api_test.go` | `TestRegisterDevice` |
| `GET /api/tracking/vehicles/:id/positions` | yes | true no-mock HTTP | `API_tests/tracking_api_test.go` | `TestSubmitTrackingLocationAndGetPositions` |
| `GET /api/tracking/vehicles/:id/stops` | yes | true no-mock HTTP | `API_tests/tracking_api_test.go` | `TestTrackingStopEventsEndpoint` |
| `POST /api/tracking/location` | yes | true no-mock HTTP | `API_tests/tracking_api_test.go` | `TestSubmitTrackingLocationAndGetPositions` |

### Notifications

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `GET /api/notification-topics` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `POST /api/notification-topics/:id/subscribe` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `DELETE /api/notification-topics/:id/subscribe` | yes | true no-mock HTTP | `API_tests/notifications_coverage_api_test.go` | `TestNotificationTopicUnsubscribe` |
| `GET /api/notification-settings` | yes | true no-mock HTTP | `API_tests/notifications_coverage_api_test.go` | `TestNotificationSettingsGenericGetAndPatch` |
| `PATCH /api/notification-settings` | yes | true no-mock HTTP | `API_tests/notifications_coverage_api_test.go` | `TestNotificationSettingsGenericGetAndPatch` |
| `GET /api/notification-settings/dnd` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationDNDSettings` |
| `PATCH /api/notification-settings/dnd` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationDNDSettings` |
| `GET /api/notifications` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `GET /api/notifications/:id` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `PATCH /api/notifications/:id/read` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `POST /api/notifications/:id/dismiss` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `GET /api/notifications/export-packages` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |
| `GET /api/notifications/export-packages/:id/download` | yes | true no-mock HTTP | `API_tests/notifications_api_test.go` | `TestNotificationSubscribeListReadDismissAndExport` |

### Campaigns / Segments / Analytics / Exports

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `GET /api/campaigns` | yes | true no-mock HTTP | `API_tests/misc_coverage_api_test.go` | `TestCampaignsList` |
| `GET /api/campaigns/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `GET /api/campaigns/:id/tasks` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `GET /api/tasks/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `POST /api/campaigns` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `PATCH /api/campaigns/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `DELETE /api/campaigns/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `POST /api/campaigns/:id/tasks` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `PATCH /api/tasks/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `DELETE /api/tasks/:id` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignAndTaskCRUDEndpoints` |
| `POST /api/tasks/:id/complete` | yes | true no-mock HTTP | `API_tests/campaigns_api_test.go` | `TestCampaignTaskReminderStopsAfterComplete` |
| `GET /api/tags` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `GET /api/members/:id/tags` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `GET /api/segments` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentRBACForbidden` |
| `GET /api/segments/:id` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `GET /api/segments/:id/runs` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `POST /api/tags` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `DELETE /api/tags/:id` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `POST /api/members/:id/tags` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `DELETE /api/members/:id/tags/:tagId` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagCRUDAndMemberTagging` |
| `POST /api/tags/export` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagExportImport` |
| `POST /api/tags/import` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestTagExportImport` |
| `POST /api/segments` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `PATCH /api/segments/:id` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `DELETE /api/segments/:id` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `POST /api/segments/:id/preview` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `POST /api/segments/:id/run` | yes | true no-mock HTTP | `API_tests/segments_api_test.go` | `TestSegmentCRUDAndPreviewRun` |
| `GET /api/analytics/occupancy` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestAnalyticsOccupancyEndpoint` |
| `GET /api/analytics/bookings` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestAnalyticsBookingsEndpoint` |
| `GET /api/analytics/exceptions` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestAnalyticsExceptionsEndpoint` |
| `GET /api/exports` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestExportCRUDAndDownload` |
| `GET /api/exports/:id/download` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestExportCRUDAndDownload` |
| `POST /api/exports` | yes | true no-mock HTTP | `API_tests/analytics_api_test.go` | `TestExportCRUDAndDownload` |

### Non-API Router Endpoints

| Endpoint | Covered | Test type | Test files | Evidence |
|---|---|---|---|---|
| `POST /auth/login` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestFormLoginRedirectsAndSetsSessionCookie` |
| `POST /auth/logout` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestFormLogoutClearsSessionAndRedirectsToLogin` |
| `GET /login` | yes | true no-mock HTTP | `API_tests/router_api_test.go` | `TestLoginPageRenders` |
| `GET /dashboard` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /reservations` | yes | true no-mock HTTP | `API_tests/reservations_api_test.go` | `TestReservationCalendarPageRenders` |
| `GET /capacity` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /facilities` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /lots` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /zones` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /rate-plans` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /members` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /vehicles` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /drivers` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /notifications` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /campaigns` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /segments` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /tasks` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /notification-prefs` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /analytics` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestPageRoutesRenderForAuthenticatedAdmin` |
| `GET /audit` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestAuditPageRendersForAuditor` |
| `GET /admin/users` | yes | true no-mock HTTP | `API_tests/web_routes_coverage_api_test.go` | `TestAdminUsersPageRendersForAdmin` |
| `GET /swagger/*any` | yes | true no-mock HTTP | `API_tests/router_api_test.go` | `TestSwaggerAccessibleForAdmin` |

## API Test Classification

1. **True No-Mock HTTP**
   - Core evidence: `API_tests/auth_api_test.go:29` (`setupAuthAPIEnv`) builds real router with `server.NewRouter(...)` and real DB pool (`pgxpool.New(...)`), then all request tests go through `apiRequest(...); r.ServeHTTP(...)` (`API_tests/auth_api_test.go:179-201`).
   - Files: all `API_tests/*_api_test.go`.

2. **HTTP with Mocking**
   - **None found** in API tests.

3. **Non-HTTP (unit/integration without HTTP)**
   - `API_tests/segments_api_test.go` → `TestSegmentFilterEvaluation` calls `segments.NewService(...).EvaluateSegment(...)` directly.
   - `API_tests/campaigns_api_test.go` → tests invoke `campaigns.NewService(...).ProcessDueTaskReminders(...)` directly inside otherwise HTTP-centered tests.
   - `unit_tests/*.go` suite is non-HTTP by design.

## Mock Detection (Strict)

- **No JS mocking APIs found:** no `jest.mock`, `vi.mock`, `sinon.stub` in repository tests.
- **No API-layer transport/controller/service mocking found** in `API_tests/*`.
- **Test doubles/fakes present in unit tests (expected for unit scope):**
  - `unit_tests/auth_test.go` defines `memoryAuthStore` fake implementing auth store behavior.
  - This does **not** affect classification of API tests; it is confined to non-HTTP unit tests.

## Coverage Summary

- **Total resolved endpoints (all router endpoints):** 150
- **Endpoints with HTTP tests:** 150
- **Endpoints with true no-mock HTTP tests:** 150
- **HTTP coverage:** 100.00%
- **True API coverage:** 100.00%

Notes:
- Coverage is now complete for both `/api/*` and non-`/api/*` registered routes.
- New direct coverage was added in `API_tests/web_routes_coverage_api_test.go` for all previously missing UI/form routes.

## Unit Test Analysis

### Backend Unit Tests

- **Test files:** `unit_tests/auth_test.go`, `unit_tests/backfill_test.go`, `unit_tests/capacity_test.go`, `unit_tests/config_test.go`, `unit_tests/device_test.go`, `unit_tests/exception_test.go`, `unit_tests/notifications_test.go`, `unit_tests/rbac_test.go`, `unit_tests/reconciliation_test.go`, `unit_tests/security_test.go`, `unit_tests/signing_secret_test.go`, `unit_tests/tracking_test.go`.
- **Modules covered:**
  - **Services/logic:** `internal/auth`, `internal/devices`, `internal/exceptions`, `internal/notifications`, `internal/reconciliation`, `internal/tracking`.
  - **Platform/security:** `internal/platform/security`.
  - **Config:** `internal/config`.
  - **DB helper behavior:** `internal/db` backfill flow (`unit_tests/backfill_test.go`).
- **Important backend modules not directly unit-tested (file-level):**
  - Router/middleware handlers as isolated units (`internal/server/auth_middleware.go`, `internal/server/middleware.go`, `internal/server/errors.go`).
  - Export internals (`internal/exports/generate.go`, `internal/exports/format.go`, `internal/exports/store.go`) rely mostly on API-level verification.
  - Auth persistence implementation (`internal/auth/store_postgres.go`) lacks dedicated focused unit tests (covered indirectly by API tests).

### Frontend Unit Tests (Strict Requirement)

- **Frontend test files:** `internal/web/pages_test.go`.
- **Framework/tools detected:** Go `testing` framework (`testing` package) with direct component rendering and assertions.
- **Components/modules covered (direct evidence in `internal/web/pages_test.go`):**
  - `LoginPage`, `DashboardPage`, `ReservationsPage`, `CapacityPage`, `NotificationsPage`, `NotificationPrefsPage`, `AnalyticsPage`, `TasksPage`, `CrudPage`.
  - Layout/navigation helpers (`initialsFor`, `isNavVisible`, role-gated navigation via rendered page output).
- **Important frontend components/modules not directly tested:**
  - No critical untested frontend route module identified after route-level HTTP additions.
  - Browser-automation (DOM interaction) E2E remains out of scope of current Go HTTP tests.

**Frontend unit tests: PRESENT**

### Cross-Layer Observation

- Testing is backend-heavy but **not frontend-absent**.
- Frontend coverage exists at component-render level; missing area is true FE↔BE browser-level integration/E2E.

## API Observability Check

- **Strong observability pattern present:** helper logs method/path/status/body (`API_tests/auth_api_test.go:204-208`) and most tests assert status plus response content.
- **Status:** improved. `API_tests/analytics_api_test.go` now asserts response payload shape/fields for occupancy/bookings/exceptions endpoint tests.

## Test Quality & Sufficiency

- **Success paths:** broadly covered across auth, CRUD, reservations, analytics/exports, notifications, segments, tracking.
- **Failure/permission paths:** strong RBAC and forbidden cases (`API_tests/rbac_api_test.go`, `API_tests/authorization_scope_api_test.go`, `API_tests/exceptions_api_test.go`).
- **Edge cases:** present (lockout, hold expiry/conflicts, concurrent oversell, invalid signatures, export format validation).
- **Validation depth:** strong across core endpoints, including form-login/logout and page-route markers/content-type checks.
- **Integration boundaries:** real router + real DB in API tests (good); browser-level scripted E2E is still not present.
- **`run_tests.sh` check:** Docker-based orchestration and in-container execution (`run_tests.sh:7-26`) → **OK**.

## End-to-End Expectations (Fullstack)

- Expected: FE↔BE end-to-end tests.
- Found: API tests + frontend component unit tests, but no browser/system E2E suite.
- Assessment: FE↔BE integration is exercised at HTTP route level (server-rendered fullstack flow). Dedicated browser automation is still optional hardening.

## Tests Check

- Static-only audit performed: yes.
- Endpoint inventory resolved from router registrations: yes.
- HTTP route coverage evidence traced to concrete test functions: yes.
- Mocking audit performed: yes.

## Test Coverage Score (0-100)

- **Score: 100/100**

## Score Rationale

- + Full endpoint request coverage across all 150 registered routes.
- + Real router + real DB path in API tests; no API-layer mocking detected.
- + Strong success/failure/authorization and edge-case coverage.
- + Analytics observability assertions strengthened beyond status-only checks.

## Key Gaps

- No blocking coverage gaps found in registered endpoint inventory.
- Optional hardening gap: browser-level automated E2E remains absent (not required for endpoint coverage completeness).

## Confidence & Assumptions

- **Confidence:** high for endpoint inventory and API coverage mapping.
- **Assumption:** endpoint scope includes all registered router paths (API + non-API) per strict endpoint definition.
- **Assumption:** component-render tests are classified as `unit-only / indirect`, not route coverage.

**Test Coverage Verdict:** **PASS**

---

# README Audit

## README Location Check

- Found at `repo/README.md`.

## Hard Gates

### Formatting

- PASS: clean markdown with clear sectioning and tables (`README.md`).

### Startup Instructions (fullstack/backend requirement)

- PASS: explicit `docker-compose up` included (`README.md:9-11`).

### Access Method

- PASS: URL and port clearly documented (`README.md:17-19`).

### Verification Method

- PASS: API verification via `curl` provided (`README.md:73-93`).
- PASS: UI validation walkthrough provided (`README.md:95-106`).

### Environment Rules (Docker-contained)

- PASS: no `npm install`, `pip install`, `apt-get`, or manual DB setup instructions.
- PASS: README states docker defaults sufficient (`README.md:44`).

### Demo Credentials (auth exists)

- PASS: credentials and roles provided for Facility Admin, Dispatch Operator, Fleet Manager, Auditor (`README.md:66-71`).

## Engineering Quality

- **Tech stack clarity:** good (`README.md:23-29`).
- **Architecture explanation:** concise and concrete; includes backend/frontend/auth/exports/schedulers/migrations.
- **Testing instructions:** explicit dockerized test command and behavior (`README.md:52-60`).
- **Security/roles:** includes role table and force-password-change behavior (`README.md:62-71`).
- **Workflow quality:** practical verification steps for both API and UI.

## High Priority Issues

- None.

## Medium Priority Issues

- README does not explicitly call out known test scope limits (no browser E2E), which could set unrealistic fullstack verification expectations.

## Low Priority Issues

- Could explicitly add troubleshooting section (container startup failures, port collisions).

## Hard Gate Failures

- None.

## README Verdict

- **PASS**

---

## Final Combined Verdicts

- **Test Coverage Audit:** **PARTIAL PASS**
- **README Audit:** **PASS**
