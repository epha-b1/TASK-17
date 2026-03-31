package API_tests

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestManualReconciliationRunCreatesCompensatingEventsAndAuditLog(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")
	fx := createReservationFixture(t, env, admin, 10, 15)

	_, err := env.pool.Exec(context.Background(), `
		INSERT INTO devices(zone_id, device_type, device_key, status)
		VALUES ($1::uuid, 'camera', 'reconcile-dev', 'online')
	`, fx.zoneID)
	if err != nil {
		t.Fatalf("seed device: %v", err)
	}

	var deviceID string
	err = env.pool.QueryRow(context.Background(), `SELECT id::text FROM devices WHERE device_key='reconcile-dev'`).Scan(&deviceID)
	if err != nil {
		t.Fatalf("read seeded device: %v", err)
	}

	_, err = env.pool.Exec(context.Background(), `
		INSERT INTO device_events(device_id, event_key, sequence_number, event_type, received_at, processed)
		VALUES
		($1::uuid, $2, 1, 'stall_freed', $3, true),
		($1::uuid, $4, 2, 'stall_freed', $3, true)
	`, deviceID, "rc-ev-1", time.Now().UTC(), "rc-ev-2")
	if err != nil {
		t.Fatalf("seed device events: %v", err)
	}

	_, err = env.pool.Exec(context.Background(), `
		INSERT INTO capacity_snapshots(zone_id, snapshot_at, authoritative_stalls)
		VALUES ($1::uuid, now(), 0)
	`, fx.zoneID)
	if err != nil {
		t.Fatalf("seed snapshot: %v", err)
	}

	run := apiRequest(t, env.r, http.MethodPost, "/api/reconciliation/run", nil, admin)
	logStep(t, "POST", "/api/reconciliation/run", run.Code, run.Body.String())
	if run.Code != http.StatusOK || !strings.Contains(run.Body.String(), `"discrepancies_found":1`) {
		t.Fatalf("expected reconciliation discrepancy, got %d %s", run.Code, run.Body.String())
	}

	var eventType string
	var stallCount int
	err = env.pool.QueryRow(context.Background(), `
		SELECT event_type, stall_count
		FROM compensating_events
		ORDER BY created_at DESC
		LIMIT 1
	`).Scan(&eventType, &stallCount)
	if err != nil {
		t.Fatalf("read compensating event: %v", err)
	}
	if eventType != "hold" || stallCount != 2 {
		t.Fatalf("expected compensating hold delta=2, got type=%s count=%d", eventType, stallCount)
	}

	auditor := loginAs(t, env, "auditor", "UserPass1234")
	audit := apiRequest(t, env.r, http.MethodGet, "/api/admin/audit-logs?page=1&limit=20", nil, auditor)
	logStep(t, "GET", "/api/admin/audit-logs", audit.Code, audit.Body.String())
	if audit.Code != http.StatusOK || !strings.Contains(audit.Body.String(), `"action":"reconciliation.run"`) {
		t.Fatalf("expected reconciliation audit log entry, got %d %s", audit.Code, audit.Body.String())
	}
}
