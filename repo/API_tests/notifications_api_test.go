package API_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"
)

func topicIDByName(t *testing.T, body, name string) string {
	t.Helper()
	needle := `"name":"` + name + `"`
	idx := strings.Index(body, needle)
	if idx < 0 {
		t.Fatalf("topic %s not found: %s", name, body)
	}
	idKey := `"id":"`
	start := strings.LastIndex(body[:idx], idKey)
	if start < 0 {
		t.Fatalf("topic id not found near %s: %s", name, body)
	}
	start += len(idKey)
	end := strings.Index(body[start:], `"`)
	if end < 0 {
		t.Fatalf("topic id not terminated: %s", body)
	}
	return body[start : start+end]
}

func firstNotificationID(t *testing.T, body string) string {
	t.Helper()
	var payload struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("parse notifications payload: %v", err)
	}
	if len(payload.Items) == 0 || payload.Items[0].ID == "" {
		t.Fatalf("missing notification id in payload: %s", body)
	}
	return payload.Items[0].ID
}

func firstExportID(t *testing.T, body string) string {
	t.Helper()
	var payload []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		t.Fatalf("parse export payload: %v", err)
	}
	if len(payload) == 0 || payload[0].ID == "" {
		t.Fatalf("missing export id in payload: %s", body)
	}
	return payload[0].ID
}

func TestNotificationSubscribeListReadDismissAndExport(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")

	topics := apiRequest(t, env.r, http.MethodGet, "/api/notification-topics", nil, admin)
	logStep(t, "GET", "/api/notification-topics", topics.Code, topics.Body.String())
	if topics.Code != http.StatusOK {
		t.Fatalf("expected 200 topics, got %d %s", topics.Code, topics.Body.String())
	}
	topicID := topicIDByName(t, topics.Body.String(), "booking_success")

	sub := apiRequest(t, env.r, http.MethodPost, "/api/notification-topics/"+topicID+"/subscribe", nil, admin)
	logStep(t, "POST", "/api/notification-topics/:id/subscribe", sub.Code, sub.Body.String())
	if sub.Code != http.StatusOK {
		t.Fatalf("subscribe failed: %d %s", sub.Code, sub.Body.String())
	}

	fx := createReservationFixture(t, env, admin, 2, 15)
	hold := apiRequest(t, env.r, http.MethodPost, "/api/reservations/hold", map[string]any{
		"zone_id":           fx.zoneID,
		"member_id":         fx.memberID,
		"vehicle_id":        fx.vehicleID,
		"time_window_start": time.Now().UTC().Add(1 * time.Hour).Format(time.RFC3339),
		"time_window_end":   time.Now().UTC().Add(2 * time.Hour).Format(time.RFC3339),
		"stall_count":       1,
	}, admin)
	logStep(t, "POST", "/api/reservations/hold", hold.Code, hold.Body.String())
	if hold.Code != http.StatusCreated {
		t.Fatalf("hold failed: %d %s", hold.Code, hold.Body.String())
	}
	resID := extractID(t, hold.Body.String())

	confirm := apiRequest(t, env.r, http.MethodPost, "/api/reservations/"+resID+"/confirm", nil, admin)
	logStep(t, "POST", "/api/reservations/:id/confirm", confirm.Code, confirm.Body.String())
	if confirm.Code != http.StatusOK {
		t.Fatalf("confirm failed: %d %s", confirm.Code, confirm.Body.String())
	}

	list := apiRequest(t, env.r, http.MethodGet, "/api/notifications?read=false", nil, admin)
	logStep(t, "GET", "/api/notifications", list.Code, list.Body.String())
	if list.Code != http.StatusOK || !strings.Contains(list.Body.String(), "Booking confirmed") {
		t.Fatalf("expected booking notification in list, got %d %s", list.Code, list.Body.String())
	}
	nid := firstNotificationID(t, list.Body.String())

	getOne := apiRequest(t, env.r, http.MethodGet, "/api/notifications/"+nid, nil, admin)
	logStep(t, "GET", "/api/notifications/:id", getOne.Code, getOne.Body.String())
	if getOne.Code != http.StatusOK {
		t.Fatalf("get notification failed: %d %s", getOne.Code, getOne.Body.String())
	}

	markRead := apiRequest(t, env.r, http.MethodPatch, "/api/notifications/"+nid+"/read", nil, admin)
	logStep(t, "PATCH", "/api/notifications/:id/read", markRead.Code, markRead.Body.String())
	if markRead.Code != http.StatusOK {
		t.Fatalf("mark read failed: %d %s", markRead.Code, markRead.Body.String())
	}

	dismiss := apiRequest(t, env.r, http.MethodPost, "/api/notifications/"+nid+"/dismiss", nil, admin)
	logStep(t, "POST", "/api/notifications/:id/dismiss", dismiss.Code, dismiss.Body.String())
	if dismiss.Code != http.StatusOK {
		t.Fatalf("dismiss failed: %d %s", dismiss.Code, dismiss.Body.String())
	}

	exports := apiRequest(t, env.r, http.MethodGet, "/api/notifications/export-packages", nil, admin)
	logStep(t, "GET", "/api/notifications/export-packages", exports.Code, exports.Body.String())
	if exports.Code != http.StatusOK || !strings.Contains(exports.Body.String(), `"id":"`) {
		t.Fatalf("expected export package list, got %d %s", exports.Code, exports.Body.String())
	}
	exportID := firstExportID(t, exports.Body.String())

	download := apiRequest(t, env.r, http.MethodGet, "/api/notifications/export-packages/"+exportID+"/download", nil, admin)
	logStep(t, "GET", "/api/notifications/export-packages/:id/download", download.Code, download.Body.String())
	if download.Code != http.StatusOK {
		t.Fatalf("download export failed: %d %s", download.Code, download.Body.String())
	}
}

func TestNotificationDNDSettings(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")

	patch := apiRequest(t, env.r, http.MethodPatch, "/api/notification-settings/dnd", map[string]any{
		"start_time": "22:00",
		"end_time":   "06:00",
		"enabled":    true,
	}, admin)
	logStep(t, "PATCH", "/api/notification-settings/dnd", patch.Code, patch.Body.String())
	if patch.Code != http.StatusOK || !strings.Contains(patch.Body.String(), `"enabled":true`) {
		t.Fatalf("expected DND patch success, got %d %s", patch.Code, patch.Body.String())
	}

	get := apiRequest(t, env.r, http.MethodGet, "/api/notification-settings/dnd", nil, admin)
	logStep(t, "GET", "/api/notification-settings/dnd", get.Code, get.Body.String())
	if get.Code != http.StatusOK || !strings.Contains(get.Body.String(), `"enabled":true`) {
		t.Fatalf("expected DND get enabled true, got %d %s", get.Code, get.Body.String())
	}
}

func TestNotificationFrequencyCapAndPersistence(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")
	topics := apiRequest(t, env.r, http.MethodGet, "/api/notification-topics", nil, admin)
	logStep(t, "GET", "/api/notification-topics", topics.Code, topics.Body.String())
	topicID := topicIDByName(t, topics.Body.String(), "booking_success")
	sub := apiRequest(t, env.r, http.MethodPost, "/api/notification-topics/"+topicID+"/subscribe", nil, admin)
	logStep(t, "POST", "/api/notification-topics/:id/subscribe", sub.Code, sub.Body.String())

	fx := createReservationFixture(t, env, admin, 2, 15)
	_, err := env.pool.Exec(context.Background(), `
		INSERT INTO notifications(user_id, topic_id, title, body, booking_id)
		VALUES
		('11111111-1111-1111-1111-111111111111'::uuid, $1::uuid, 'r1', 'b1', NULL),
		('11111111-1111-1111-1111-111111111111'::uuid, $1::uuid, 'r2', 'b2', NULL),
		('11111111-1111-1111-1111-111111111111'::uuid, $1::uuid, 'r3', 'b3', NULL)
	`, topicID)
	if err != nil {
		t.Fatalf("seed notifications: %v", err)
	}

	_, err = env.pool.Exec(context.Background(), `
		INSERT INTO notification_jobs(notification_id, user_id, booking_id, topic_id, channel, status, attempt_count)
		SELECT n.id, n.user_id, n.booking_id, n.topic_id, 'in_app', 'delivered', 0
		FROM notifications n
		WHERE n.title IN ('r1','r2','r3')
	`)
	if err != nil {
		t.Fatalf("seed jobs: %v", err)
	}

	hold := apiRequest(t, env.r, http.MethodPost, "/api/reservations/hold", map[string]any{
		"zone_id":           fx.zoneID,
		"member_id":         fx.memberID,
		"vehicle_id":        fx.vehicleID,
		"time_window_start": time.Now().UTC().Add(3 * time.Hour).Format(time.RFC3339),
		"time_window_end":   time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339),
		"stall_count":       1,
	}, admin)
	logStep(t, "POST", "/api/reservations/hold", hold.Code, hold.Body.String())
	resID := extractID(t, hold.Body.String())
	confirm := apiRequest(t, env.r, http.MethodPost, "/api/reservations/"+resID+"/confirm", nil, admin)
	logStep(t, "POST", "/api/reservations/:id/confirm", confirm.Code, confirm.Body.String())
	if confirm.Code != http.StatusOK {
		t.Fatalf("confirm failed: %d %s", confirm.Code, confirm.Body.String())
	}

	var count int
	err = env.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM notification_jobs WHERE booking_id=$1::uuid`, resID).Scan(&count)
	if err != nil {
		t.Fatalf("count jobs: %v", err)
	}
	if count > 3 {
		t.Fatalf("frequency cap violated, got %d jobs", count)
	}
}
