package dataworks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
)

// dataWorksTestClient stands up a credential-free AliyunClient whose
// dataworks-public RPC transport is pointed at an in-process httptest server,
// mirroring the rds-account unit-test construction. This avoids needing real
// cloud credentials (and the ACC account gap) while exercising the CRUD logic.
func dataWorksTestClient(t *testing.T, handler http.HandlerFunc) *connectivity.AliyunClient {
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
	endpoints.Store("dataworks_public", endpoint)
	return client
}

// writeJSON writes a 200 JSON success body.
func writeJSON(t *testing.T, w http.ResponseWriter, body interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Error(err)
	}
}

// writeGone writes a 404 with the DataWorks resource-gone error code, the shape
// the real API returns when a resource no longer exists.
func writeGone(t *testing.T, w http.ResponseWriter) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"Code":    errorCodeResourceGone,
		"Message": "the resource does not exist",
	})
}

// resourceObject is a canned GetResource Resource body shared across tests.
func resourceObject(id, projectID, name, spec string) map[string]interface{} {
	return map[string]interface{}{
		"Id":         id,
		"ProjectId":  projectID,
		"Name":       name,
		"Spec":       spec,
		"CreateTime": 1700000000,
		"ModifyTime": 1700000001,
		"Owner":      "owner-1",
		"Type":       "file",
		"SourceType": "local",
		"SourcePath": "/src",
		"TargetType": "oss",
		"TargetPath": "/tgt",
		"DataSource": map[string]interface{}{"Name": "ds-1", "Type": "mysql"},
		"Script": map[string]interface{}{
			"Path":    "/script",
			"Id":      "sc-1",
			"Runtime": map[string]interface{}{"Command": "python"},
		},
	}
}

// TestUnitDataWorksResourceCreateAppliesName (F3): resource_name is not a
// CreateResource parameter; it must be applied via a separate RenameResource
// call right after creation. The previous execute skipped this.
func TestUnitDataWorksResourceCreateAppliesName(t *testing.T) {
	var (
		createBody  map[string]string
		renameCalls int32
	)
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "CreateResource":
			createBody = formValues(r)
			writeJSON(t, w, map[string]interface{}{"Id": "res-1"})
		case "RenameResource":
			atomic.AddInt32(&renameCalls, 1)
			assert.Equal(t, "myname", r.Form.Get("Name"), "RenameResource must carry the planned name")
			writeJSON(t, w, map[string]interface{}{})
		case "GetResource":
			assert.Equal(t, "proj-1", r.Form.Get("ProjectId"), "GetResource must use the project id")
			assert.Equal(t, "res-1", r.Form.Get("Id"))
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-1", "proj-1", "myname", "{}")})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	plan := resourceModel{
		ProjectId:    types.StringValue("proj-1"),
		Spec:         types.StringValue("{}"),
		ResourceName: types.StringValue("myname"),
		ResourceFile: types.StringValue("file-content"),
	}
	out, err := resourceCreate(client, plan)
	assert.NoError(t, err)
	assert.Equal(t, "proj-1:res-1", out.Id.ValueString())
	assert.Equal(t, "res-1", out.ResourceId.ValueString())
	assert.Equal(t, "myname", out.ResourceName.ValueString())
	assert.Equal(t, "file-content", out.ResourceFile.ValueString(), "write-only resource_file must carry from plan to state")
	assert.Equal(t, int64(1700000000), out.CreateTime.ValueInt64())
	assert.Equal(t, "ds-1", out.DataSourceName.ValueString())
	assert.Equal(t, "sc-1", out.ScriptId.ValueString())
	assert.Equal(t, "python", out.ScriptRuntimeCmd.ValueString())
	assert.Equal(t, int32(1), atomic.LoadInt32(&renameCalls), "RenameResource must fire once (F3)")
	assert.Equal(t, "{}", createBody["Spec"])
	assert.Equal(t, "file-content", createBody["ResourceFile"])
}

// TestUnitDataWorksResourceCreateOmitsRename: when no resource_name is planned,
// RenameResource must not be called.
func TestUnitDataWorksResourceCreateOmitsRename(t *testing.T) {
	var renameCalls int32
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "CreateResource":
			writeJSON(t, w, map[string]interface{}{"Id": "res-2"})
		case "RenameResource":
			atomic.AddInt32(&renameCalls, 1)
			writeJSON(t, w, map[string]interface{}{})
		case "GetResource":
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-2", "proj-1", "", "{}")})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	plan := resourceModel{
		ProjectId: types.StringValue("proj-1"),
		Spec:      types.StringValue("{}"),
		// ResourceName intentionally null.
	}
	out, err := resourceCreate(client, plan)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), atomic.LoadInt32(&renameCalls), "RenameResource must not fire when no name is planned")
	assert.True(t, out.ResourceName.IsNull(), "resource_name reads null when the API has none")
}

// TestUnitDataWorksResourceReadRecoversProjectId (F2): the composite id carries
// the project id, so Read recovers it without a config-only lookup.
func TestUnitDataWorksResourceReadRecoversProjectId(t *testing.T) {
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "GetResource":
			assert.Equal(t, "proj-1", r.Form.Get("ProjectId"), "Read must split the project id out of the composite id (F2)")
			assert.Equal(t, "res-1", r.Form.Get("Id"))
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-1", "proj-1", "n1", "{}")})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	state := resourceModel{Id: types.StringValue("proj-1:res-1")}
	out, gone, err := resourceRead(client, state)
	assert.NoError(t, err)
	assert.False(t, gone)
	assert.Equal(t, "proj-1:res-1", out.Id.ValueString())
	assert.Equal(t, "res-1", out.ResourceId.ValueString())
}

// TestUnitDataWorksResourceReadGone: a 404 resource-gone response drops the
// resource out of state rather than hard-failing.
func TestUnitDataWorksResourceReadGone(t *testing.T) {
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "GetResource":
			writeGone(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	state := resourceModel{Id: types.StringValue("proj-1:res-1")}
	out, gone, err := resourceRead(client, state)
	assert.NoError(t, err, "Read must not hard-fail when the resource is gone")
	assert.True(t, gone, "the resource must be reported gone so the framework removes it")
	assert.Equal(t, resourceModel{}, out)
}

// TestUnitDataWorksResourceUpdateWriteOnlyConverges (F4): removing write-only
// resource_file and path from config must NOT re-fire UpdateResource/MoveResource
// every plan. The new value is empty, so the send guard suppresses the call.
func TestUnitDataWorksResourceUpdateWriteOnlyConverges(t *testing.T) {
	var (
		updateCalls int32
		moveCalls   int32
		renameCalls int32
	)
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "UpdateResource":
			atomic.AddInt32(&updateCalls, 1)
			writeJSON(t, w, map[string]interface{}{})
		case "MoveResource":
			atomic.AddInt32(&moveCalls, 1)
			writeJSON(t, w, map[string]interface{}{})
		case "RenameResource":
			atomic.AddInt32(&renameCalls, 1)
			writeJSON(t, w, map[string]interface{}{})
		case "GetResource":
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-1", "proj-1", "n1", "{}")})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	state := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue("{}"),
		ResourceFile: types.StringValue("old-file"),
		Path:         types.StringValue("/old/path"),
		ResourceName: types.StringValue("n1"),
	}
	// Plan removes resource_file and path; spec and name unchanged.
	plan := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue("{}"),
		ResourceFile: types.StringNull(),
		Path:         types.StringNull(),
		ResourceName: types.StringValue("n1"),
	}
	out, err := resourceUpdate(client, plan, state)
	assert.NoError(t, err)
	assert.Equal(t, int32(0), atomic.LoadInt32(&updateCalls), "UpdateResource must not fire when only write-only fields were removed (F4)")
	assert.Equal(t, int32(0), atomic.LoadInt32(&moveCalls), "MoveResource must not fire when path was removed (F4)")
	assert.Equal(t, int32(0), atomic.LoadInt32(&renameCalls), "RenameResource must not fire when name is unchanged")
	assert.True(t, out.ResourceFile.IsNull(), "write-only resource_file is null in the refreshed state after removal")
	assert.True(t, out.Path.IsNull(), "write-only path is null in the refreshed state after removal")
}

// TestUnitDataWorksResourceUpdateSendsSpecAndFile: when spec and a non-empty
// resource_file change, UpdateResource fires once carrying both.
func TestUnitDataWorksResourceUpdateSendsSpecAndFile(t *testing.T) {
	var updateBody map[string]string
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "UpdateResource":
			updateBody = formValues(r)
			writeJSON(t, w, map[string]interface{}{})
		case "GetResource":
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-1", "proj-1", "n1", `{"x":1}`)})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	state := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue("{}"),
		ResourceFile: types.StringValue("old-file"),
	}
	plan := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue(`{"x":1}`),
		ResourceFile: types.StringValue("new-file"),
	}
	out, err := resourceUpdate(client, plan, state)
	assert.NoError(t, err)
	assert.Equal(t, `{"x":1}`, updateBody["Spec"], "UpdateResource must carry the new spec")
	assert.Equal(t, "new-file", updateBody["ResourceFile"], "UpdateResource must carry the new non-empty resource_file")
	assert.Equal(t, "new-file", out.ResourceFile.ValueString(), "write-only resource_file carries from plan into refreshed state")
	assert.Equal(t, `{"x":1}`, out.Spec.ValueString())
}

// TestUnitDataWorksResourceUpdateRenamesAndMoves: changing resource_name and
// path fires RenameResource and MoveResource, but not UpdateResource.
func TestUnitDataWorksResourceUpdateRenamesAndMoves(t *testing.T) {
	var (
		renameName string
		movePath   string
		updateSeen int32
	)
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "UpdateResource":
			atomic.AddInt32(&updateSeen, 1)
			writeJSON(t, w, map[string]interface{}{})
		case "RenameResource":
			renameName = r.Form.Get("Name")
			writeJSON(t, w, map[string]interface{}{})
		case "MoveResource":
			movePath = r.Form.Get("Path")
			writeJSON(t, w, map[string]interface{}{})
		case "GetResource":
			writeJSON(t, w, map[string]interface{}{"Resource": resourceObject("res-1", "proj-1", "n2", "{}")})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})

	state := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue("{}"),
		ResourceName: types.StringValue("n1"),
		Path:         types.StringValue("/old"),
	}
	plan := resourceModel{
		Id:           types.StringValue("proj-1:res-1"),
		Spec:         types.StringValue("{}"),
		ResourceName: types.StringValue("n2"),
		Path:         types.StringValue("/new/path"),
	}
	_, err := resourceUpdate(client, plan, state)
	assert.NoError(t, err)
	assert.Equal(t, "n2", renameName, "RenameResource must carry the new name")
	assert.Equal(t, "/new/path", movePath, "MoveResource must carry the new path")
	assert.Equal(t, int32(0), atomic.LoadInt32(&updateSeen), "UpdateResource must not fire when only name/path changed")
}

// TestUnitDataWorksResourceDeleteIdempotent: a 404 resource-gone on delete is
// success, so destroy converges.
func TestUnitDataWorksResourceDeleteIdempotent(t *testing.T) {
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DeleteResource":
			writeGone(t, w)
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	state := resourceModel{Id: types.StringValue("proj-1:res-1")}
	assert.NoError(t, resourceDelete(client, state), "delete must be idempotent when the resource is already gone")
}

// TestUnitDataWorksResourceDeleteOK: a clean delete returns no error.
func TestUnitDataWorksResourceDeleteOK(t *testing.T) {
	var deleteBody map[string]string
	client := dataWorksTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Form.Get("Action") {
		case "DeleteResource":
			deleteBody = formValues(r)
			writeJSON(t, w, map[string]interface{}{})
		default:
			t.Errorf("unexpected action %q", r.Form.Get("Action"))
		}
	})
	state := resourceModel{Id: types.StringValue("proj-1:res-1")}
	assert.NoError(t, resourceDelete(client, state))
	assert.Equal(t, "res-1", deleteBody["Id"], "DeleteResource must target the resource id half of the composite id")
	assert.Equal(t, "proj-1", deleteBody["ProjectId"])
}

// formValues copies the parsed form into a string map for assertions.
func formValues(r *http.Request) map[string]string {
	out := map[string]string{}
	for k, vs := range r.Form {
		if len(vs) > 0 {
			out[k] = vs[0]
		}
	}
	return out
}
