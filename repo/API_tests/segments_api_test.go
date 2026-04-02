package API_tests

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"parkops/internal/segments"
)

func TestTagCRUDAndMemberTagging(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")

	// Create tag
	createTag := apiRequest(t, env.r, http.MethodPost, "/api/tags", map[string]any{
		"name": "downtown_monthly",
	}, admin)
	logStep(t, "POST", "/api/tags", createTag.Code, createTag.Body.String())
	if createTag.Code != http.StatusCreated {
		t.Fatalf("create tag failed: %d %s", createTag.Code, createTag.Body.String())
	}
	tagID := extractID(t, createTag.Body.String())

	// List tags
	listTags := apiRequest(t, env.r, http.MethodGet, "/api/tags", nil, admin)
	logStep(t, "GET", "/api/tags", listTags.Code, listTags.Body.String())
	if listTags.Code != http.StatusOK || !strings.Contains(listTags.Body.String(), "downtown_monthly") {
		t.Fatalf("list tags failed: %d %s", listTags.Code, listTags.Body.String())
	}

	// Create a member to tag
	_, err := env.pool.Exec(context.Background(), `
		INSERT INTO organizations(id, name) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'OrgA') ON CONFLICT DO NOTHING
	`)
	if err != nil {
		t.Fatalf("seed org: %v", err)
	}
	createMember := apiRequest(t, env.r, http.MethodPost, "/api/members", map[string]any{
		"organization_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"display_name":    "Test Member",
	}, admin)
	logStep(t, "POST", "/api/members", createMember.Code, createMember.Body.String())
	if createMember.Code != http.StatusCreated {
		t.Fatalf("create member failed: %d %s", createMember.Code, createMember.Body.String())
	}
	memberID := extractID(t, createMember.Body.String())

	// Add tag to member
	addTag := apiRequest(t, env.r, http.MethodPost, "/api/members/"+memberID+"/tags", map[string]any{
		"tag_id": tagID,
	}, admin)
	logStep(t, "POST", "/api/members/:id/tags", addTag.Code, addTag.Body.String())
	if addTag.Code != http.StatusOK {
		t.Fatalf("add tag to member failed: %d %s", addTag.Code, addTag.Body.String())
	}

	// Get member tags
	getMemberTags := apiRequest(t, env.r, http.MethodGet, "/api/members/"+memberID+"/tags", nil, admin)
	logStep(t, "GET", "/api/members/:id/tags", getMemberTags.Code, getMemberTags.Body.String())
	if getMemberTags.Code != http.StatusOK || !strings.Contains(getMemberTags.Body.String(), "downtown_monthly") {
		t.Fatalf("get member tags failed: %d %s", getMemberTags.Code, getMemberTags.Body.String())
	}

	// Remove tag from member
	removeTag := apiRequest(t, env.r, http.MethodDelete, "/api/members/"+memberID+"/tags/"+tagID, nil, admin)
	logStep(t, "DELETE", "/api/members/:id/tags/:tagId", removeTag.Code, removeTag.Body.String())
	if removeTag.Code != http.StatusOK {
		t.Fatalf("remove tag failed: %d %s", removeTag.Code, removeTag.Body.String())
	}

	// Verify tag removed
	getMemberTags2 := apiRequest(t, env.r, http.MethodGet, "/api/members/"+memberID+"/tags", nil, admin)
	logStep(t, "GET", "/api/members/:id/tags (after remove)", getMemberTags2.Code, getMemberTags2.Body.String())
	if getMemberTags2.Code != http.StatusOK || strings.Contains(getMemberTags2.Body.String(), "downtown_monthly") {
		t.Fatalf("tag should be removed: %d %s", getMemberTags2.Code, getMemberTags2.Body.String())
	}

	// Delete tag
	deleteTag := apiRequest(t, env.r, http.MethodDelete, "/api/tags/"+tagID, nil, admin)
	logStep(t, "DELETE", "/api/tags/:id", deleteTag.Code, deleteTag.Body.String())
	if deleteTag.Code != http.StatusNoContent {
		t.Fatalf("delete tag failed: %d %s", deleteTag.Code, deleteTag.Body.String())
	}
}

func TestSegmentCRUDAndPreviewRun(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")

	// Seed org and member with arrears
	_, _ = env.pool.Exec(context.Background(), `
		INSERT INTO organizations(id, name) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'OrgA') ON CONFLICT DO NOTHING
	`)
	createMember := apiRequest(t, env.r, http.MethodPost, "/api/members", map[string]any{
		"organization_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"display_name":    "High Arrears Member",
	}, admin)
	logStep(t, "POST", "/api/members", createMember.Code, createMember.Body.String())
	if createMember.Code != http.StatusCreated {
		t.Fatalf("create member failed: %d %s", createMember.Code, createMember.Body.String())
	}
	memberID := extractID(t, createMember.Body.String())

	// Set arrears balance via admin endpoint
	patchBal := apiRequest(t, env.r, http.MethodPatch, "/api/members/"+memberID+"/balance", map[string]any{
		"amount_cents": 10000,
		"reason":       "test seed",
	}, admin)
	logStep(t, "PATCH", "/api/members/:id/balance", patchBal.Code, patchBal.Body.String())
	if patchBal.Code != http.StatusOK {
		t.Fatalf("patch balance failed: %d %s", patchBal.Code, patchBal.Body.String())
	}

	// Create segment with balance filter
	createSeg := apiRequest(t, env.r, http.MethodPost, "/api/segments", map[string]any{
		"name":              "High Arrears",
		"filter_expression": map[string]any{"arrears_balance_cents": map[string]any{"gt": 5000}},
		"schedule":          "manual",
	}, admin)
	logStep(t, "POST", "/api/segments", createSeg.Code, createSeg.Body.String())
	if createSeg.Code != http.StatusCreated {
		t.Fatalf("create segment failed: %d %s", createSeg.Code, createSeg.Body.String())
	}
	segmentID := extractID(t, createSeg.Body.String())

	// Get segment
	getSeg := apiRequest(t, env.r, http.MethodGet, "/api/segments/"+segmentID, nil, admin)
	logStep(t, "GET", "/api/segments/:id", getSeg.Code, getSeg.Body.String())
	if getSeg.Code != http.StatusOK {
		t.Fatalf("get segment failed: %d %s", getSeg.Code, getSeg.Body.String())
	}

	// Preview segment
	preview := apiRequest(t, env.r, http.MethodPost, "/api/segments/"+segmentID+"/preview", nil, admin)
	logStep(t, "POST", "/api/segments/:id/preview", preview.Code, preview.Body.String())
	if preview.Code != http.StatusOK {
		t.Fatalf("preview segment failed: %d %s", preview.Code, preview.Body.String())
	}
	var previewResp map[string]any
	_ = json.Unmarshal(preview.Body.Bytes(), &previewResp)
	if previewResp["member_count"].(float64) < 1 {
		t.Fatalf("expected at least 1 member in preview, got %v", previewResp["member_count"])
	}

	// Run segment
	run := apiRequest(t, env.r, http.MethodPost, "/api/segments/"+segmentID+"/run", nil, admin)
	logStep(t, "POST", "/api/segments/:id/run", run.Code, run.Body.String())
	if run.Code != http.StatusOK {
		t.Fatalf("run segment failed: %d %s", run.Code, run.Body.String())
	}
	var runResp map[string]any
	_ = json.Unmarshal(run.Body.Bytes(), &runResp)
	if runResp["member_count"].(float64) < 1 {
		t.Fatalf("expected at least 1 member in run, got %v", runResp["member_count"])
	}

	// List run history
	runs := apiRequest(t, env.r, http.MethodGet, "/api/segments/"+segmentID+"/runs", nil, admin)
	logStep(t, "GET", "/api/segments/:id/runs", runs.Code, runs.Body.String())
	if runs.Code != http.StatusOK || !strings.Contains(runs.Body.String(), "manual") {
		t.Fatalf("list runs failed: %d %s", runs.Code, runs.Body.String())
	}

	// Patch segment
	patchSeg := apiRequest(t, env.r, http.MethodPatch, "/api/segments/"+segmentID, map[string]any{
		"name": "Updated High Arrears",
	}, admin)
	logStep(t, "PATCH", "/api/segments/:id", patchSeg.Code, patchSeg.Body.String())
	if patchSeg.Code != http.StatusOK {
		t.Fatalf("patch segment failed: %d %s", patchSeg.Code, patchSeg.Body.String())
	}

	// Delete segment
	deleteSeg := apiRequest(t, env.r, http.MethodDelete, "/api/segments/"+segmentID, nil, admin)
	logStep(t, "DELETE", "/api/segments/:id", deleteSeg.Code, deleteSeg.Body.String())
	if deleteSeg.Code != http.StatusNoContent {
		t.Fatalf("delete segment failed: %d %s", deleteSeg.Code, deleteSeg.Body.String())
	}
}

func TestTagExportImport(t *testing.T) {
	env := setupAuthAPIEnv(t)
	admin := loginAs(t, env, "admin", "AdminPass1234")

	// Seed org, member, and tag
	_, _ = env.pool.Exec(context.Background(), `
		INSERT INTO organizations(id, name) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'OrgA') ON CONFLICT DO NOTHING
	`)
	createMember := apiRequest(t, env.r, http.MethodPost, "/api/members", map[string]any{
		"organization_id": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
		"display_name":    "Export Member",
	}, admin)
	if createMember.Code != http.StatusCreated {
		t.Fatalf("create member: %d %s", createMember.Code, createMember.Body.String())
	}
	memberID := extractID(t, createMember.Body.String())

	createTag := apiRequest(t, env.r, http.MethodPost, "/api/tags", map[string]any{"name": "vip"}, admin)
	if createTag.Code != http.StatusCreated {
		t.Fatalf("create tag: %d %s", createTag.Code, createTag.Body.String())
	}
	tagID := extractID(t, createTag.Body.String())

	addTag := apiRequest(t, env.r, http.MethodPost, "/api/members/"+memberID+"/tags", map[string]any{"tag_id": tagID}, admin)
	if addTag.Code != http.StatusOK {
		t.Fatalf("add tag: %d %s", addTag.Code, addTag.Body.String())
	}

	// Export
	export := apiRequest(t, env.r, http.MethodPost, "/api/tags/export", map[string]any{}, admin)
	logStep(t, "POST", "/api/tags/export", export.Code, export.Body.String())
	if export.Code != http.StatusOK || !strings.Contains(export.Body.String(), "snapshot") {
		t.Fatalf("export failed: %d %s", export.Code, export.Body.String())
	}

	// Parse snapshot from export
	var exportResp map[string]any
	_ = json.Unmarshal(export.Body.Bytes(), &exportResp)
	snapshot := exportResp["snapshot"]

	// Import (restore)
	importReq := apiRequest(t, env.r, http.MethodPost, "/api/tags/import", map[string]any{
		"snapshot": snapshot,
	}, admin)
	logStep(t, "POST", "/api/tags/import", importReq.Code, importReq.Body.String())
	if importReq.Code != http.StatusOK {
		t.Fatalf("import failed: %d %s", importReq.Code, importReq.Body.String())
	}

	// Verify audit log has tag_import entry
	var auditCount int
	err := env.pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM audit_logs WHERE action='tag_import'`).Scan(&auditCount)
	if err != nil {
		t.Fatalf("query audit: %v", err)
	}
	if auditCount == 0 {
		t.Fatal("expected audit log entry for tag_import")
	}
}

func TestSegmentFilterEvaluation(t *testing.T) {
	env := setupAuthAPIEnv(t)

	// Seed org and members
	_, _ = env.pool.Exec(context.Background(), `
		INSERT INTO organizations(id, name) VALUES ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'OrgA') ON CONFLICT DO NOTHING
	`)
	_, _ = env.pool.Exec(context.Background(), `
		INSERT INTO members(organization_id, display_name, arrears_balance_cents) VALUES
		('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'Member Low', 100),
		('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa'::uuid, 'Member High', 20000)
	`)

	svc := segments.NewService(env.pool)

	// Test gt filter
	ids, err := svc.EvaluateSegment(context.Background(), []byte(`{"arrears_balance_cents": {"gt": 5000}}`))
	if err != nil {
		t.Fatalf("evaluate segment: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("expected 1 matching member, got %d", len(ids))
	}

	// Test lt filter
	idsLt, err := svc.EvaluateSegment(context.Background(), []byte(`{"arrears_balance_cents": {"lt": 500}}`))
	if err != nil {
		t.Fatalf("evaluate segment lt: %v", err)
	}
	if len(idsLt) != 1 {
		t.Fatalf("expected 1 matching member for lt, got %d", len(idsLt))
	}
}

func TestSegmentRBACForbidden(t *testing.T) {
	env := setupAuthAPIEnv(t)
	operator := loginAs(t, env, "operator", "UserPass1234")

	// Operator (dispatch_operator) should be forbidden from creating segments
	createSeg := apiRequest(t, env.r, http.MethodPost, "/api/segments", map[string]any{
		"name":              "Forbidden Seg",
		"filter_expression": map[string]any{"arrears_balance_cents": map[string]any{"gt": 0}},
	}, operator)
	logStep(t, "POST", "/api/segments (operator)", createSeg.Code, createSeg.Body.String())
	if createSeg.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for operator creating segment, got %d", createSeg.Code)
	}

	// Operator should be able to read segments
	listSeg := apiRequest(t, env.r, http.MethodGet, "/api/segments", nil, operator)
	logStep(t, "GET", "/api/segments (operator)", listSeg.Code, listSeg.Body.String())
	if listSeg.Code != http.StatusOK {
		t.Fatalf("operator should be able to list segments, got %d", listSeg.Code)
	}
}
