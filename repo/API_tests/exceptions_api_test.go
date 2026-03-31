package API_tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func createExceptionFromDeviceEvent(t *testing.T, env *apiTestEnv, admin *http.Cookie, eventType string) string {
	t.Helper()
	deviceID, _ := createDeviceWithZone(t, env, admin)

	ingest := apiRequest(t, env.r, http.MethodPost, "/api/device-events", map[string]any{
		"device_id":       deviceID,
		"event_key":       "ev-ex-" + time.Now().UTC().Format("150405.000000"),
		"sequence_number": 1,
		"event_type":      eventType,
	}, admin)
	logStep(t, "POST", "/api/device-events", ingest.Code, ingest.Body.String())
	if ingest.Code != http.StatusCreated {
		t.Fatalf("create exception event failed: %d %s", ingest.Code, ingest.Body.String())
	}

	list := apiRequest(t, env.r, http.MethodGet, "/api/exceptions", nil, admin)
	logStep(t, "GET", "/api/exceptions", list.Code, list.Body.String())
	if list.Code != http.StatusOK {
		t.Fatalf("list exceptions failed: %d %s", list.Code, list.Body.String())
	}

	var payload struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal(list.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode exceptions payload: %v", err)
	}
	for _, item := range payload.Items {
		typ, _ := item["exception_type"].(string)
		if typ != eventType {
			continue
		}
		id, _ := item["id"].(string)
		if id != "" {
			return id
		}
	}
	t.Fatalf("expected %s exception in list, got %s", eventType, list.Body.String())
	return ""
}

func TestExceptionsOpenListAndGet(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")
	exceptionID := createExceptionFromDeviceEvent(t, env, admin, "gate_stuck")

	open := apiRequest(t, env.r, http.MethodGet, "/api/exceptions", nil, admin)
	logStep(t, "GET", "/api/exceptions", open.Code, open.Body.String())
	if open.Code != http.StatusOK || !strings.Contains(open.Body.String(), `"exception_type":"gate_stuck"`) || !strings.Contains(open.Body.String(), `"status":"open"`) {
		t.Fatalf("expected open gate_stuck exception, got %d %s", open.Code, open.Body.String())
	}

	getOne := apiRequest(t, env.r, http.MethodGet, "/api/exceptions/"+exceptionID, nil, admin)
	logStep(t, "GET", "/api/exceptions/:id", getOne.Code, getOne.Body.String())
	if getOne.Code != http.StatusOK || !strings.Contains(getOne.Body.String(), exceptionID) {
		t.Fatalf("expected get exception success, got %d %s", getOne.Code, getOne.Body.String())
	}
}

func TestAcknowledgeExceptionAsDispatchAndHistory(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")
	dispatch := loginAs(t, env, "operator", "UserPass1234")
	exceptionID := createExceptionFromDeviceEvent(t, env, admin, "sensor_offline")

	ack := apiRequest(t, env.r, http.MethodPost, "/api/exceptions/"+exceptionID+"/acknowledge", map[string]any{"note": "checked sensor"}, dispatch)
	logStep(t, "POST", "/api/exceptions/:id/acknowledge", ack.Code, ack.Body.String())
	if ack.Code != http.StatusOK || !strings.Contains(ack.Body.String(), `"status":"acknowledged"`) {
		t.Fatalf("expected acknowledge success, got %d %s", ack.Code, ack.Body.String())
	}

	open := apiRequest(t, env.r, http.MethodGet, "/api/exceptions", nil, dispatch)
	logStep(t, "GET", "/api/exceptions", open.Code, open.Body.String())
	if open.Code != http.StatusOK || strings.Contains(open.Body.String(), exceptionID) {
		t.Fatalf("expected acknowledged exception removed from open list, got %d %s", open.Code, open.Body.String())
	}

	history := apiRequest(t, env.r, http.MethodGet, "/api/exceptions/history?page=1", nil, dispatch)
	logStep(t, "GET", "/api/exceptions/history", history.Code, history.Body.String())
	if history.Code != http.StatusOK || !strings.Contains(history.Body.String(), exceptionID) {
		t.Fatalf("expected history to include acknowledged exception, got %d %s", history.Code, history.Body.String())
	}
}

func TestAcknowledgeExceptionWrongRoleForbiddenAndAuditLogged(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")
	auditor := loginAs(t, env, "auditor", "UserPass1234")
	exceptionID := createExceptionFromDeviceEvent(t, env, admin, "camera_error")

	forbidden := apiRequest(t, env.r, http.MethodPost, "/api/exceptions/"+exceptionID+"/acknowledge", map[string]any{"note": "admin cannot ack"}, admin)
	logStep(t, "POST", "/api/exceptions/:id/acknowledge", forbidden.Code, forbidden.Body.String())
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for wrong role, got %d %s", forbidden.Code, forbidden.Body.String())
	}

	audit := apiRequest(t, env.r, http.MethodGet, "/api/admin/audit-logs?page=1&limit=20", nil, auditor)
	logStep(t, "GET", "/api/admin/audit-logs", audit.Code, audit.Body.String())
	if audit.Code != http.StatusOK || !strings.Contains(audit.Body.String(), `"action":"rbac_denied"`) || !strings.Contains(audit.Body.String(), `/api/exceptions/`) {
		t.Fatalf("expected rbac_denied audit log for exception acknowledge, got %d %s", audit.Code, audit.Body.String())
	}
}
