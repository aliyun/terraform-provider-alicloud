package cms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// cmsTestClient stands up a credential-free AliyunClient whose CMS ROA transport
// is pointed at an in-process httptest server, mirroring the rds-account and
// ecs-snapshot unit-test construction. It runs under a plain `go test` and
// exercises the real CRUD/state/failure paths without touching the cloud.
func cmsTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
		}
		handler(w, r)
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").SetAccessKeyId("test-key").SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	endpoints := new(sync.Map)
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	config := &connectivity.Config{
		AccessKey: "test-key", SecretKey: "test-secret", Credential: credential,
		RegionId: "cn-hangzhou", AccountType: "test", Protocol: "http",
		Endpoints: endpoints, SignVersion: new(sync.Map), SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatal(err)
	}
	endpoints.Store("cms", endpoint)
	return client
}

func writeJson(t *testing.T, w http.ResponseWriter, v interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		t.Error(err)
	}
}

func writeCmsError(t *testing.T, w http.ResponseWriter, status int, code string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	writeJson(t, w, map[string]interface{}{"Code": code, "Message": code})
}

// robotJSON builds a List-API robot object. webhook is intentionally omitted to
// model the real List API, which never returns the write-only URL.
func robotJSON(id, name, rtype, lang, workspace, employee, signKey string) map[string]interface{} {
	return map[string]interface{}{
		"robotId":             id,
		"name":                name,
		"type":                rtype,
		"lang":                lang,
		"workspace":           workspace,
		"digitalEmployeeName": employee,
		"robotSignKey":        signKey,
	}
}

func robotFromCreatePlan(plan alertRobotResourceModel) map[string]interface{} {
	signKey := ""
	if !plan.RobotSignKey.IsNull() {
		signKey = plan.RobotSignKey.ValueString()
	}
	employee := ""
	if !plan.DigitalEmployeeName.IsNull() {
		employee = plan.DigitalEmployeeName.ValueString()
	}
	workspace := ""
	if !plan.Workspace.IsNull() {
		workspace = plan.Workspace.ValueString()
	}
	return robotJSON("rid-1", plan.AlertRobotName.ValueString(), plan.Type.ValueString(),
		resolveLang(plan.Lang), workspace, employee, signKey)
}

func TestUnitCmsAlertRobotCreateDefaultsAndPersistsWriteOnly(t *testing.T) {
	var getCalls, postCalls int
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/robot":
			postCalls++
			writeJson(t, w, map[string]interface{}{"alertRobotId": "rid-1"})
		case r.Method == http.MethodGet && r.URL.Path == "/robots":
			getCalls++
			// Echo back the created robot without webhook (write-only).
			writeJson(t, w, map[string]interface{}{
				"robots": []interface{}{robotJSON("rid-1", "robot-1", "DING", "zh_CN", "ws-1", "emp-1", "key-1")},
				"total":  1,
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	r := &alertRobotResource{}
	plan := alertRobotResourceModel{
		AlertRobotName:      types.StringValue("robot-1"),
		Type:                types.StringValue("DING"),
		Url:                 types.StringValue("https://oapi.dingtalk.com/robot/send?access_token=acc"),
		Workspace:           types.StringValue("ws-1"),
		DigitalEmployeeName: types.StringValue("emp-1"),
		RobotSignKey:        types.StringValue("key-1"),
		// lang omitted: must default to zh_CN
	}

	state, err := r.create(context.Background(), client, plan)
	assert.NoError(t, err)
	assert.Equal(t, 1, postCalls, "Create must POST /robot once")
	assert.Equal(t, 1, getCalls, "Create must re-read via GET /robots once")
	assert.Equal(t, "rid-1", state.Id.ValueString())
	assert.Equal(t, "rid-1", state.AlertRobotId.ValueString())
	assert.Equal(t, "zh_CN", state.Lang.ValueString(), "omitted lang must default to zh_CN and persist")
	assert.Equal(t, plan.Url.ValueString(), state.Url.ValueString(), "write-only url must be preserved from plan")
	assert.Equal(t, "key-1", state.RobotSignKey.ValueString())
	assert.Equal(t, "ws-1", state.Workspace.ValueString())
}

func TestUnitCmsAlertRobotCreateHonorsExplicitLang(t *testing.T) {
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/robot":
			writeJson(t, w, map[string]interface{}{"alertRobotId": "rid-2"})
		case r.Method == http.MethodGet && r.URL.Path == "/robots":
			writeJson(t, w, map[string]interface{}{
				"robots": []interface{}{robotJSON("rid-2", "robot-2", "FEISHU", "en_US", "", "", "")},
				"total":  1,
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	r := &alertRobotResource{}
	plan := alertRobotResourceModel{
		AlertRobotName: types.StringValue("robot-2"),
		Type:           types.StringValue("FEISHU"),
		Url:            types.StringValue("https://example.com/feishu"),
		Lang:           types.StringValue("en_US"),
	}
	state, err := r.create(context.Background(), client, plan)
	assert.NoError(t, err)
	assert.Equal(t, "en_US", state.Lang.ValueString(), "explicit lang must be honored")
	assert.Equal(t, "https://example.com/feishu", state.Url.ValueString())
}

func TestUnitCmsAlertRobotReadGoneRemovesResource(t *testing.T) {
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/robots" {
			writeJson(t, w, map[string]interface{}{"robots": []interface{}{}, "total": 0})
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})

	r := &alertRobotResource{}
	state := alertRobotResourceModel{Id: types.StringValue("rid-gone")}
	updated, gone, err := r.read(context.Background(), client, state)
	assert.NoError(t, err)
	assert.True(t, gone, "a missing robot must signal removal, not error")
	assert.Equal(t, "", updated.Id.ValueString())
}

func TestUnitCmsAlertRobotReadPersistsWriteOnlyUrl(t *testing.T) {
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/robots" {
			writeJson(t, w, map[string]interface{}{
				"robots": []interface{}{robotJSON("rid-3", "robot-3", "DING", "zh_CN", "ws", "emp", "key")},
				"total":  1,
			})
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})

	r := &alertRobotResource{}
	state := alertRobotResourceModel{
		Id:  types.StringValue("rid-3"),
		Url: types.StringValue("https://preserved.example.com/ro"),
	}
	updated, gone, err := r.read(context.Background(), client, state)
	assert.NoError(t, err)
	assert.False(t, gone)
	assert.Equal(t, "rid-3", updated.Id.ValueString())
	assert.Equal(t, "https://preserved.example.com/ro", updated.Url.ValueString(), "write-only url must be preserved across read")
	assert.Equal(t, "key", updated.RobotSignKey.ValueString())
}

func TestUnitCmsAlertRobotUpdatePatchesChangedFields(t *testing.T) {
	var patchBody map[string]interface{}
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPatch && r.URL.Path == "/robot/rid-4":
			_ = json.NewDecoder(r.Body).Decode(&patchBody)
			writeJson(t, w, map[string]interface{}{})
		case r.Method == http.MethodGet && r.URL.Path == "/robots":
			writeJson(t, w, map[string]interface{}{
				"robots": []interface{}{robotJSON("rid-4", "robot-4-upd", "DING", "en_US", "ws-4", "emp-upd", "key-upd")},
				"total":  1,
			})
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusBadRequest)
		}
	})

	r := &alertRobotResource{}
	state := alertRobotResourceModel{
		Id:             types.StringValue("rid-4"),
		AlertRobotName: types.StringValue("robot-4"),
		Type:           types.StringValue("DING"),
		Url:            types.StringValue("https://oapi.dingtalk.com/robot/send?access_token=old"),
		Lang:           types.StringValue("zh_CN"),
		Workspace:      types.StringValue("ws-4"),
	}
	plan := alertRobotResourceModel{
		Id:                  types.StringValue("rid-4"),
		AlertRobotName:      types.StringValue("robot-4-upd"),
		Type:                types.StringValue("DING"),
		Url:                 types.StringValue("https://oapi.dingtalk.com/robot/send?access_token=new"),
		Lang:                types.StringValue("en_US"),
		Workspace:           types.StringValue("ws-4"),
		DigitalEmployeeName: types.StringValue("emp-upd"),
		RobotSignKey:        types.StringValue("key-upd"),
	}
	updated, err := r.update(context.Background(), client, plan, state)
	assert.NoError(t, err)
	assert.Equal(t, "robot-4-upd", updated.AlertRobotName.ValueString())
	assert.Equal(t, "en_US", updated.Lang.ValueString())
	assert.Equal(t, "https://oapi.dingtalk.com/robot/send?access_token=new", updated.Url.ValueString())
	assert.Equal(t, "key-upd", updated.RobotSignKey.ValueString())
	assert.Equal(t, "emp-upd", updated.DigitalEmployeeName.ValueString())
	// workspace did not change, so it must not be sent on the patch body.
	_, workspaceSent := patchBody["workspace"]
	assert.False(t, workspaceSent, "unchanged workspace must not be sent on Update")
}

func TestUnitCmsAlertRobotDeleteIdempotentOnGone(t *testing.T) {
	var deleteCalls int
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && r.URL.Path == "/robots" {
			deleteCalls++
			writeCmsError(t, w, http.StatusNotFound, "ResourceNotFound")
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})

	r := &alertRobotResource{}
	state := alertRobotResourceModel{Id: types.StringValue("rid-gone"), Type: types.StringValue("DING")}
	err := r.delete(context.Background(), client, state)
	assert.NoError(t, err, "Delete on a missing robot must be idempotent, not error")
	assert.Equal(t, 1, deleteCalls)
}

func TestUnitCmsAlertRobotDeletePropagatesUnexpectedError(t *testing.T) {
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && r.URL.Path == "/robots" {
			writeCmsError(t, w, http.StatusForbidden, "AccessDenied")
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})

	r := &alertRobotResource{}
	state := alertRobotResourceModel{Id: types.StringValue("rid-5"), Type: types.StringValue("DING")}
	err := r.delete(context.Background(), client, state)
	assert.Error(t, err, "a non-not-found Delete error must propagate, not be swallowed")
}

func TestUnitCmsAlertRobotsDataSourceFiltersAndMaps(t *testing.T) {
	client := cmsTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/robots" {
			robotIds := r.URL.Query().Get("robotIds")
			robots := []interface{}{
				robotJSON("rid-1", "alpha", "DING", "zh_CN", "ws-1", "e1", "k1"),
				robotJSON("rid-2", "beta", "FEISHU", "en_US", "ws-2", "e2", ""),
			}
			// The ids filter narrows to rid-1.
			if strings.Contains(robotIds, "rid-1") && !strings.Contains(robotIds, "rid-2") {
				robots = robots[:1]
			}
			writeJson(t, w, map[string]interface{}{"robots": robots, "total": len(robots)})
			return
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
	})

	d := &alertRobotsDataSource{}
	ids, _ := types.ListValueFrom(context.Background(), types.StringType, []string{"rid-1"})
	config := alertRobotsDataSourceModel{
		Ids:       ids,
		NameRegex: types.StringNull(),
		Type:      types.StringNull(),
		Workspace: types.StringNull(),
	}
	state, err := d.read(context.Background(), client, config)
	assert.NoError(t, err)
	assert.Equal(t, "rid-1", state.Ids.Elements()[0].(types.String).ValueString())
	// robots list must contain exactly the filtered robot.
	robotsLen := len(state.Robots.Elements())
	assert.Equal(t, 1, robotsLen, "ids filter must narrow to one robot")
	assert.NotEqual(t, "", state.Id.ValueString(), "data source id must be a stable hash")
}
