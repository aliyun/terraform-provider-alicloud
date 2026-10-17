package connectivity

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUnitRoaPaiPluginGroupWire captures the real HTTP wire the PAIPlugin
// 2022-01-12 Group CRUD sends through the provider's roaRequest path, against
// an httptest server. It is the cloud-free regression for the SDK contract:
// the official Go SDK (alibabacloud-go/paiplugin-20220112) emits
// Action=CreateGroup/GetGroup/DeleteGroup, Version=2022-01-12 over
// POST/GET/DELETE /api/v2/groups[/<Id>] with a JSON body.
//
// It proves RoaPostWithApiName/RoaGetWithApiName/RoaDeleteWithApiName drive
// DoRequestWithAction and therefore carry the x-acs-action header, while the
// generic RoaPost (apiName="") routes through DoRequest and sends an EMPTY
// x-acs-action — the exact difference that left the backend answering at the
// auth/routing layer (MissingAccessKeyId / InvalidApi.NotFound) rather than
// the business route. That auth-layer response alone does NOT prove the
// account/backend is provisioned; provisioning stays a product-confirmed
// assumption, not an excluded code factor.
//
// No cloud is touched: a fake access_key credential signs locally and the
// endpoint map is pre-seeded so loadApiEndpoint never calls the location
// service. getTeaRoaDslSdkConfig only assembles the roa.Config.
func TestUnitRoaPaiPluginGroupWire(t *testing.T) {
	type captured struct {
		method      string
		path        string
		action      string
		version     string
		contentType string
		body        []byte
	}
	var (
		mu  sync.Mutex
		got []captured
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		got = append(got, captured{
			method:      r.Method,
			path:        r.URL.EscapedPath(),
			action:      r.Header.Get("x-acs-action"),
			version:     r.Header.Get("x-acs-version"),
			contentType: r.Header.Get("content-type"),
			body:        b,
		})
		mu.Unlock()
		// GetGroup/DeleteGroup need a 200 with a Data object; CreateGroup
		// returns the created id under Data.Id.
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_, _ = io.WriteString(w, `{"Data":{"Id":123,"Name":"g1"}}`)
	}))
	defer srv.Close()

	// Bare host:port — the roa SDK composes <protocol>://<endpoint><path>.
	host := strings.TrimPrefix(srv.URL, "http://")
	endpoints := &sync.Map{}
	endpoints.Store("paiplugin", host)

	cred, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("FAKEPAIPLUGINAK").
		SetAccessKeySecret("FAKEPAIPLUGINSK"))
	require.NoError(t, err)

	cfg := &Config{
		RegionId:             "cn-hangzhou",
		AccessKey:            "FAKEPAIPLUGINAK",
		SecretKey:            "FAKEPAIPLUGINSK",
		Protocol:             "HTTP",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		Credential:           cred,
		SkipRegionValidation: true,
	}
	roaCfg, err := cfg.getTeaRoaDslSdkConfig(true)
	require.NoError(t, err)
	client := &AliyunClient{config: cfg, teaRoaSdkConfig: roaCfg}

	// CreateGroup — POST /api/v2/groups with a JSON body carrying Name/Source.
	body := map[string]interface{}{"Name": "g1", "Source": 0}
	if _, err := client.RoaPostWithApiName("PAIPlugin", "2022-01-12", "CreateGroup", "/api/v2/groups", nil, nil, body, false); err != nil {
		t.Fatalf("RoaPostWithApiName CreateGroup: %v", err)
	}
	// GetGroup — GET /api/v2/groups/<Id>.
	if _, err := client.RoaGetWithApiName("PAIPlugin", "2022-01-12", "GetGroup", "/api/v2/groups/123", nil, nil, nil); err != nil {
		t.Fatalf("RoaGetWithApiName GetGroup: %v", err)
	}
	// DeleteGroup — DELETE /api/v2/groups/<Id>.
	if _, err := client.RoaDeleteWithApiName("PAIPlugin", "2022-01-12", "DeleteGroup", "/api/v2/groups/123", nil, nil, nil, false); err != nil {
		t.Fatalf("RoaDeleteWithApiName DeleteGroup: %v", err)
	}
	// Generic RoaPost (apiName="") — the contrast: no x-acs-action is sent.
	if _, err := client.RoaPost("PAIPlugin", "2022-01-12", "/api/v2/groups", nil, nil, body, false); err != nil {
		// A backend-free mock returns 200, so no error is expected; log only.
		t.Logf("generic RoaPost returned err (acceptable for the contrast): %v", err)
	}

	require.Equal(t, 4, len(got), "expected Create/Get/Delete + one generic contrast request")
	c, g, d := got[0], got[1], got[2]
	contrast := got[3]

	// --- CreateGroup ---
	assert.Equal(t, http.MethodPost, c.method, "CreateGroup method")
	assert.Equal(t, "/api/v2/groups", c.path, "CreateGroup path")
	assert.Equal(t, "CreateGroup", c.action, "CreateGroup x-acs-action sent via DoRequestWithAction")
	assert.Equal(t, "2022-01-12", c.version, "CreateGroup x-acs-version")
	assert.NotEmpty(t, c.contentType, "CreateGroup carries a content-type header")
	t.Logf("CreateGroup content-type=%q body=%q", c.contentType, string(c.body))
	var jb map[string]interface{}
	if assert.NoError(t, json.Unmarshal(c.body, &jb), "CreateGroup body is JSON") {
		assert.Equal(t, "g1", jb["Name"], "CreateGroup JSON body carries Name")
	}

	// --- GetGroup ---
	assert.Equal(t, http.MethodGet, g.method, "GetGroup method")
	assert.Equal(t, "/api/v2/groups/123", g.path, "GetGroup encoded path /api/v2/groups/<Id>")
	assert.Equal(t, "GetGroup", g.action, "GetGroup x-acs-action")
	assert.Equal(t, "2022-01-12", g.version, "GetGroup x-acs-version")

	// --- DeleteGroup ---
	assert.Equal(t, http.MethodDelete, d.method, "DeleteGroup method")
	assert.Equal(t, "/api/v2/groups/123", d.path, "DeleteGroup encoded path /api/v2/groups/<Id>")
	assert.Equal(t, "DeleteGroup", d.action, "DeleteGroup x-acs-action")
	assert.Equal(t, "2022-01-12", d.version, "DeleteGroup x-acs-version")

	// --- Contrast: generic RoaPost sends an EMPTY x-acs-action ---
	assert.Equal(t, http.MethodPost, contrast.method, "generic RoaPost method")
	assert.Equal(t, "/api/v2/groups", contrast.path, "generic RoaPost path")
	assert.Equal(t, "", contrast.action, "generic RoaPost sends NO x-acs-action (the original defect)")
	assert.Equal(t, "2022-01-12", contrast.version, "generic RoaPost still carries x-acs-version")
}
