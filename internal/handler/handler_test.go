package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"configcenter/internal/config"
	"configcenter/internal/model"
	"configcenter/internal/service"
	"configcenter/internal/store"
	"configcenter/pkg/httpx"
	"configcenter/pkg/logger"
)

func newTestServer() *Server {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	return NewServer(svc, log, cfg)
}

func decodeResp(t *testing.T, r *httptest.ResponseRecorder) httpx.Response {
	var resp httpx.Response
	if err := json.Unmarshal(r.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	return resp
}

func TestApplicationRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1","description":"desc","status":"active"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create status %d", w.Code)
	}
	resp := decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected conflict, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/applications?page=1&size=10", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/applications/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/applications/noexist", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}

	body = `{"name":"app1-up","code":"c1","description":"desc","status":"active"}`
	req = httptest.NewRequest(http.MethodPut, "/api/applications/"+id, strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/applications/noexist", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found on update, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/applications/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/applications/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found on delete, got %d", w.Code)
	}
}

func TestApplicationValidation(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", w.Code)
	}
}

func TestApplicationEnvironmentsCascade(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/applications/"+appID+"/environments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get app envs status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 env, got %d", len(items))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/applications/noexist/environments", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestApplicationConfigItemsCascade(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v1"}`
	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/applications/"+appID+"/config-items", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get app config items status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 config item, got %d", len(items))
	}
}

func TestApplicationReleasesCascade(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","version":"v1","status":"draft"}`
	req = httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/applications/"+appID+"/releases", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get app releases status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 release, got %d", len(items))
	}
}

func TestEnvironmentRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create env status %d", w.Code)
	}
	resp = decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected conflict, got %d", w.Code)
	}

	body = `{"app_id":"noexist","name":"prod","code":"prod"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/environments?app_id="+appID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list env status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/environments/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get env status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/environments/noexist", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}

	body = `{"app_id":"` + appID + `","name":"dev-up","code":"dev"}`
	req = httptest.NewRequest(http.MethodPut, "/api/environments/"+id, strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update env status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/environments/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete env status %d", w.Code)
	}
}

func TestEnvironmentConfigItemsCascade(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v1"}`
	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/environments/"+envID+"/config-items", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get env config items status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 config item, got %d", len(items))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/environments/noexist/config-items", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestEnvironmentReleasesCascade(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","version":"v1","status":"draft"}`
	req = httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/environments/"+envID+"/releases", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get env releases status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 release, got %d", len(items))
	}
}

func TestConfigItemRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v1","value_type":"string","encrypted":true,"status":"enabled"}`
	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create item status %d", w.Code)
	}
	resp = decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)
	val := resp.Data.(map[string]interface{})["value"].(string)
	if !strings.Contains(val, "*") {
		t.Fatalf("expected masked value, got %s", val)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusConflict {
		t.Fatalf("expected conflict, got %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/config-items?app_id="+appID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list item status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/config-items/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get item status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/config-items/noexist", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v2","value_type":"string","encrypted":false,"status":"enabled"}`
	req = httptest.NewRequest(http.MethodPut, "/api/config-items/"+id, strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update item status %d", w.Code)
	}
	resp = decodeResp(t, w)
	val = resp.Data.(map[string]interface{})["value"].(string)
	if val != "v2" {
		t.Fatalf("expected plain v2, got %s", val)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/config-items/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete item status %d", w.Code)
	}
}

func TestConfigItemMissingAppEnv(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"noexist","env_id":"noenv","key":"k1","value":"v1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestConfigVersionRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v1"}`
	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	itemID := resp.Data.(map[string]interface{})["id"].(string)

	req = httptest.NewRequest(http.MethodGet, "/api/config-versions?config_item_id="+itemID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list versions status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.(map[string]interface{})["items"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("expected 1 version, got %d", len(items))
	}
}

func TestReleaseRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","version":"v1","remark":"r1","status":"draft"}`
	req = httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create release status %d", w.Code)
	}
	resp = decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	req = httptest.NewRequest(http.MethodGet, "/api/releases?app_id="+appID, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list release status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/releases/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("get release status %d", w.Code)
	}

	body = `{"status":"released","operator":"u1"}`
	req = httptest.NewRequest(http.MethodPut, "/api/releases/"+id+"/status", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for illegal transition, got %d", w.Code)
	}

	body = `{"status":"review","operator":"u1"}`
	req = httptest.NewRequest(http.MethodPut, "/api/releases/"+id+"/status", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("transition review status %d", w.Code)
	}

	body = `{"status":"released","operator":"u1"}`
	req = httptest.NewRequest(http.MethodPut, "/api/releases/"+id+"/status", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("transition released status %d", w.Code)
	}
	resp = decodeResp(t, w)
	st := resp.Data.(map[string]interface{})["status"].(string)
	if st != "released" {
		t.Fatalf("expected released, got %s", st)
	}

	body = `{"status":"rolled_back","operator":"u1"}`
	req = httptest.NewRequest(http.MethodPut, "/api/releases/"+id+"/status", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("transition rolled_back status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/releases/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("delete release status %d", w.Code)
	}
}

func TestReleaseMissingAppEnv(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"noexist","env_id":"noenv","version":"v1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestAuditLogRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	req = httptest.NewRequest(http.MethodDelete, "/api/applications/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/audit-logs?action=delete", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list audit status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.(map[string]interface{})["items"].([]interface{})
	if len(items) == 0 {
		t.Fatalf("expected audit logs")
	}
}

func TestStatsRoutes(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()

	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	appID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","name":"dev","code":"dev"}`
	req = httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	envID := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","key":"k1","value":"v1"}`
	req = httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	body = `{"app_id":"` + appID + `","env_id":"` + envID + `","version":"v1","status":"draft"}`
	req = httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/stats/config-item-counts", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("stats config item status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/stats/release-status", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("stats release status %d", w.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/stats/audit-actions", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("stats audit status %d", w.Code)
	}
}

func TestInvalidJSON(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1"`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for invalid json, got %d", w.Code)
	}
}

func TestMultipleJSON(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}{"name":"app2","code":"c2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for multiple json, got %d", w.Code)
	}
}

func TestNotFoundRoute(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodGet, "/api/notexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found for unknown route, got %d", w.Code)
	}
}

func TestPaginationParams(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	for i := 0; i < 5; i++ {
		body := `{"name":"app","code":"c` + string(rune('0'+i)) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/applications?page=1&size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 5 {
		t.Fatalf("expected total 5, got %v", p["total"])
	}
	if int(p["size"].(float64)) != 2 {
		t.Fatalf("expected size 2, got %v", p["size"])
	}
}

func TestReleaseGetNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodGet, "/api/releases/noexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestReleaseDeleteNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodDelete, "/api/releases/noexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestEnvironmentDeleteNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodDelete, "/api/environments/noexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestEnvironmentUpdateNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"a1","name":"dev","code":"dev"}`
	req := httptest.NewRequest(http.MethodPut, "/api/environments/noexist", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestConfigItemUpdateNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"a1","env_id":"e1","key":"k1","value":"v1"}`
	req := httptest.NewRequest(http.MethodPut, "/api/config-items/noexist", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestConfigItemDeleteNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodDelete, "/api/config-items/noexist", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestReleaseUpdateStatusNotFound(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"status":"review","operator":"u1"}`
	req := httptest.NewRequest(http.MethodPut, "/api/releases/noexist/status", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected not found, got %d", w.Code)
	}
}

func TestAuditLogFilterCombinations(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)
	req = httptest.NewRequest(http.MethodDelete, "/api/applications/"+id, nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	req = httptest.NewRequest(http.MethodGet, "/api/audit-logs?operator=system&action=delete&resource=application", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list audit status %d", w.Code)
	}
	resp = decodeResp(t, w)
	items := resp.Data.(map[string]interface{})["items"].([]interface{})
	if len(items) == 0 {
		t.Fatalf("expected audit logs with combined filter")
	}

	req = httptest.NewRequest(http.MethodGet, "/api/audit-logs?operator=noexist", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp = decodeResp(t, w)
	items = resp.Data.(map[string]interface{})["items"].([]interface{})
	if len(items) != 0 {
		t.Fatalf("expected no audit logs for noexist operator")
	}
}

func TestStatsEmpty(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	req := httptest.NewRequest(http.MethodGet, "/api/stats/config-item-counts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("stats empty status %d", w.Code)
	}
}

func TestLargeBody(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := bytes.Repeat([]byte("x"), 2<<20)
	req := httptest.NewRequest(http.MethodPost, "/api/applications", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for large body, got %d", w.Code)
	}
}

func TestApplicationCreateValidationBadRequest(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"name":"app1","code":"c1","status":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/applications", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for invalid status, got %d", w.Code)
	}
}

func TestEnvironmentCreateValidationBadRequest(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"a1","name":"","code":"dev"}`
	req := httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for empty name, got %d", w.Code)
	}
}

func TestConfigItemCreateValidationBadRequest(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"a1","env_id":"e1","key":"","value":"v1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for empty key, got %d", w.Code)
	}
}

func TestReleaseCreateValidationBadRequest(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	body := `{"app_id":"a1","env_id":"e1","version":"","status":"draft"}`
	req := httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for empty version, got %d", w.Code)
	}
}

func TestApplicationListKeywordFilter(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	srv.svc.CreateApplication(model.Application{Name: "alpha", Code: "a1"})
	srv.svc.CreateApplication(model.Application{Name: "beta", Code: "b1"})
	req := httptest.NewRequest(http.MethodGet, "/api/applications?keyword=alpha", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Fatalf("expected total 1 for keyword filter, got %v", p["total"])
	}
}

func TestEnvironmentListAppFilter(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	req := httptest.NewRequest(http.MethodGet, "/api/environments?app_id="+app.ID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Fatalf("expected total 1 for app filter, got %v", p["total"])
	}
}

func TestConfigItemListMultiFilter(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	srv.svc.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1", Status: model.ConfigStatusEnabled})
	srv.svc.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k2", Value: "v2", Status: model.ConfigStatusDisabled})
	req := httptest.NewRequest(http.MethodGet, "/api/config-items?app_id="+app.ID+"&status=enabled", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Fatalf("expected total 1 for multi filter, got %v", p["total"])
	}
}

func TestReleaseListStatusFilter(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	srv.svc.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
	srv.svc.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v2", Status: model.ReleaseStatusReleased})
	req := httptest.NewRequest(http.MethodGet, "/api/releases?status=released", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 1 {
		t.Fatalf("expected total 1 for status filter, got %v", p["total"])
	}
}

func TestConfigItemEncryptedUpdate(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	body := `{"app_id":"` + app.ID + `","env_id":"` + env.ID + `","key":"secret","value":"mysecret","encrypted":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	body = `{"app_id":"` + app.ID + `","env_id":"` + env.ID + `","key":"secret","value":"newsecret","encrypted":true}`
	req = httptest.NewRequest(http.MethodPut, "/api/config-items/"+id, strings.NewReader(body))
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("update status %d", w.Code)
	}
	resp = decodeResp(t, w)
	val := resp.Data.(map[string]interface{})["value"].(string)
	if !strings.Contains(val, "*") {
		t.Fatalf("expected masked value after update, got %s", val)
	}
}

func TestConfigVersionListPagination(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := srv.svc.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	for i := 2; i <= 5; i++ {
		srv.svc.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v" + string(rune('0'+i))})
	}
	req := httptest.NewRequest(http.MethodGet, "/api/config-versions?config_item_id="+c.ID+"&page=1&size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list versions status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 5 {
		t.Fatalf("expected total 5 versions, got %v", p["total"])
	}
	if int(p["size"].(float64)) != 2 {
		t.Fatalf("expected size 2, got %v", p["size"])
	}
}

func TestAuditLogListPagination(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	for i := 0; i < 5; i++ {
		srv.svc.CreateApplication(model.Application{Name: "app", Code: "c" + string(rune('0'+i))})
	}
	req := httptest.NewRequest(http.MethodGet, "/api/audit-logs?page=1&size=2", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("list audit status %d", w.Code)
	}
	resp := decodeResp(t, w)
	p := resp.Data.(map[string]interface{})["pagination"].(map[string]interface{})
	if int(p["total"].(float64)) != 5 {
		t.Fatalf("expected total 5 audit logs, got %v", p["total"])
	}
	if int(p["size"].(float64)) != 2 {
		t.Fatalf("expected size 2, got %v", p["size"])
	}
}

func TestReleaseFullLifecycle(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	body := `{"app_id":"` + app.ID + `","env_id":"` + env.ID + `","version":"v1","status":"draft"}`
	req := httptest.NewRequest(http.MethodPost, "/api/releases", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	resp := decodeResp(t, w)
	id := resp.Data.(map[string]interface{})["id"].(string)

	transitions := []string{"review", "released", "rolled_back"}
	for _, st := range transitions {
		body = `{"status":"` + st + `","operator":"u1"}`
		req = httptest.NewRequest(http.MethodPut, "/api/releases/"+id+"/status", strings.NewReader(body))
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("transition to %s status %d", st, w.Code)
		}
	}
}

func TestConfigItemValueTypeValidation(t *testing.T) {
	srv := newTestServer()
	router := srv.Routes()
	app, _ := srv.svc.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := srv.svc.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	body := `{"app_id":"` + app.ID + `","env_id":"` + env.ID + `","key":"k1","value":"v1","value_type":"invalid"}`
	req := httptest.NewRequest(http.MethodPost, "/api/config-items", strings.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request for invalid value_type, got %d", w.Code)
	}
}
