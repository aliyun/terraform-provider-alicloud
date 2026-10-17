package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	credentials "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	tfsdk "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// These are no-cloud (httptest) Framework tests for alicloud_ecs_storage_set.
// They drive the real resource methods (Create/Read/Update/Delete/Import) and
// the describeStorageSet helper against an in-process RPC server, with no cloud
// credentials and no terraform binary. The negative Read cases (malformed
// list type, non-object entry, missing/other id) must report errors and
// preserve state rather than silently removing the resource.

const testStorageSetID = "ss-test-0001"

// storageSetResponder maps an RPC action + its query params to an HTTP status
// and response body. The returned body is written verbatim.
type storageSetResponder func(action string, params map[string]string) (status int, body string)

// storageSetTestServer starts an httptest server that routes by Action and
// records the params of the last request per action. The endpoint resolution
// is bypassed by pinning the "ecs" endpoint to the server URL.
func storageSetTestServer(t *testing.T, respond storageSetResponder) (*httptest.Server, *sync.Map) {
	t.Helper()
	var mu sync.Mutex
	sent := &sync.Map{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// tea-rpc encodes business parameters as application/x-www-form-urlencoded
		// in the POST body (Action/AccessKeyId/Signature stay on the URL query), so
		// merge query and parsed form to capture the full request.
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		params := map[string]string{}
		for k, v := range r.Form {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
		mu.Lock()
		sent.Store(params["Action"], params)
		mu.Unlock()
		status, body := respond(params["Action"], params)
		if status == 0 {
			status = http.StatusOK
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv, sent
}

// storageSetTestClient builds an AliyunClient whose RPC calls hit srv with no
// external network. This is a harness correction: the earlier helper set
// Protocol "HTTP" (uppercase, which the SDK does not normalize) and left
// AccountType empty, so (*Config).Client() ran getAccountType(), which issues
// a REAL BssOpenApi QueryAvailableInstances RPC against business.aliyuncs.com
// with the fake credentials — leaking out of httptest isolation and masking
// the real provider behavior under a 404. The fix mirrors the stable test
// clients elsewhere in the provider (e.g. cloud_firewall_address_book_test.go):
// AccountType "test" short-circuits getAccountType(); lowercase "http" dials
// the local server over plaintext; the ecs AND bssopenapi endpoint keys are
// pinned to the local host so every RPC the resource issues lands on srv; and
// NO_PROXY keeps environment proxies from forwarding the request away.
func storageSetTestClient(t *testing.T, srv *httptest.Server) *connectivity.AliyunClient {
	t.Helper()
	cred, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-access-key").
		SetAccessKeySecret("test-secret-key"))
	if err != nil {
		t.Fatalf("credential: %v", err)
	}
	// Pin the host without a scheme: tea-rpc reconstructs the request URL as
	// <protocol>://<endpoint>, and a value that already carries a scheme
	// (srv.URL) is mangled into "http//host", so requests never reach srv.
	host := strings.TrimPrefix(strings.TrimPrefix(srv.URL, "https://"), "http://")
	// NO_PROXY keeps any environment HTTP proxy from carrying the local
	// request out to the real network.
	t.Setenv("NO_PROXY", host)
	endpoints := &sync.Map{}
	endpoints.Store("ecs", host)
	endpoints.Store("ecs.%s.aliyuncs.com", host)
	// Pin bssopenapi too: getAccountType would otherwise resolve it to
	// business.aliyuncs.com; with AccountType set it is never queried, but
	// pinning keeps the guarantee watertight if a future code path probes it.
	endpoints.Store("bssopenapi", host)
	cfg := &connectivity.Config{
		AccessKey:            "test-access-key",
		SecretKey:            "test-secret-key",
		Region:               connectivity.Hangzhou,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          &sync.Map{},
		Credential:           cred,
		SkipRegionValidation: true,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	return client
}

func newStorageSetResource(t *testing.T, client *connectivity.AliyunClient) (*storageSetResource, schema.Schema) {
	t.Helper()
	r := &storageSetResource{}
	var confResp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, &confResp)
	if confResp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", confResp.Diagnostics)
	}
	if r.Client() == nil {
		t.Fatal("client not injected")
	}
	var sresp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &sresp)
	if sresp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", sresp.Diagnostics)
	}
	return r, sresp.Schema
}

// nullStorageSetModel returns a fully initialized model where every field is
// a typed null. A zero-value storageSetResourceModel leaves types.Map (and
// other collection values) with a nil element type, which reflect.FromValue
// rejects during Plan/State.Set. Building test models on top of this helper
// avoids "value missing type" write errors.
func nullStorageSetModel() storageSetResourceModel {
	return storageSetResourceModel{
		Id:                 types.StringNull(),
		Description:        types.StringNull(),
		MaxPartitionNumber: types.Int64Null(),
		RegionId:           types.StringNull(),
		ResourceGroupId:    types.StringNull(),
		StorageSetId:       types.StringNull(),
		StorageSetName:     types.StringNull(),
		Tags:               types.MapNull(types.StringType),
		ZoneId:             types.StringNull(),
	}
}

func storageSetPlan(t *testing.T, sch schema.Schema, model storageSetResourceModel) tfsdk.Plan {
	plan := tfsdk.Plan{Schema: sch}
	if diags := plan.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	return plan
}

func storageSetStateFromModel(t *testing.T, sch schema.Schema, model storageSetResourceModel) tfsdk.State {
	state := tfsdk.State{Schema: sch}
	if diags := state.Set(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}
	return state
}

// describeStorageSetObjectBody returns a DescribeStorageSets response body
// containing one StorageSet with the requested id and the given fields.
func describeStorageSetObjectBody(id, desc, name, zone, region, rgID string, partition int, tags string) string {
	return fmt.Sprintf(`{"StorageSets":{"StorageSet":[{"StorageSetId":%q,"Description":%q,"StorageSetName":%q,"ZoneId":%q,"RegionId":%q,"ResourceGroupId":%q,"StorageSetPartitionNumber":%d,"Tags":%s}]},"RequestId":"r"}`,
		id, desc, name, zone, region, rgID, partition, tags)
}

func TestUnitDescribeStorageSet_NoCloud(t *testing.T) {
	cases := []struct {
		name         string
		body         string // DescribeStorageSets response body
		status       int    // HTTP status (200 default)
		wantErr      bool
		wantNotFound bool
	}{
		{
			name:    "positive matching id",
			body:    describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`),
			wantErr: false,
		},
		{
			name:         "explicit empty array is not-found",
			body:         `{"StorageSets":{"StorageSet":[]},"RequestId":"r"}`,
			wantNotFound: true,
		},
		{
			// An explicit empty array is authoritative "no matching entries for the
			// requested id" regardless of TotalCount (per the route plan, an empty
			// array is a legitimate not-found even when TotalCount is positive).
			name:         "explicit empty array with positive TotalCount is still not-found",
			body:         `{"TotalCount":1,"StorageSets":{"StorageSet":[]},"RequestId":"r"}`,
			wantNotFound: true,
		},
		{
			// The legitimate not-found shape: TotalCount=0 corroborates an empty
			// result, so the absent StorageSet key is not-found.
			name:         "absent StorageSet key with TotalCount:0 is not-found",
			body:         `{"TotalCount":0,"StorageSets":{},"RequestId":"r"}`,
			wantNotFound: true,
		},
		{
			// Positive TotalCount without the inner StorageSet container: the API
			// claims matches but returned none — inconsistent, must error and
			// preserve state instead of dropping the resource.
			name:    "positive TotalCount without StorageSet key is error (preserve state)",
			body:    `{"TotalCount":1,"StorageSets":{},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			// TotalCount omitted entirely is not a documented DescribeStorageSets
			// shape (the responseDemo always carries TotalCount): absence cannot
			// be established from the count, so the resource is preserved.
			name:    "absent StorageSet key without TotalCount is conservative error (preserve state)",
			body:    `{"StorageSets":{},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:         "not-found error code is not-found",
			body:         `{"Code":"InvalidStorageSetId.NotFound","Message":"The specified parameter StorageSetId was not found.","RequestId":"r"}`,
			wantNotFound: true,
		},
		{
			name:    "StorageSet present but not a list is malformed error",
			body:    `{"StorageSets":{"StorageSet":"not-a-list"},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "non-object entry is malformed error",
			body:    `{"StorageSets":{"StorageSet":["x"]},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "entry missing StorageSetId is malformed error",
			body:    `{"StorageSets":{"StorageSet":[{"ZoneId":"cn-hangzhou-i"}]},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "entry with other id is malformed error (preserve state)",
			body:    `{"StorageSets":{"StorageSet":[{"StorageSetId":"ss-other"}]},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			// Legitimate not-found: TotalCount=0 corroborates absence, so the
			// absent StorageSets wrapper is not-found.
			name:         "no StorageSets wrapper with TotalCount:0 is not-found",
			body:         `{"TotalCount":0,"Foo":"bar","RequestId":"r"}`,
			wantNotFound: true,
		},
		{
			// Positive TotalCount without the top-level StorageSets container:
			// the API claims matches but returned none — inconsistent, must error
			// and preserve state instead of dropping the resource.
			name:    "positive TotalCount without StorageSets wrapper is error (preserve state)",
			body:    `{"TotalCount":1,"RequestId":"r"}`,
			wantErr: true,
		},
		{
			// TotalCount omitted entirely is not a documented DescribeStorageSets
			// shape (the responseDemo always carries TotalCount): absence cannot
			// be established from the count, so the resource is preserved.
			name:    "no StorageSets wrapper without TotalCount is conservative error (preserve state)",
			body:    `{"Foo":"bar","RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "StorageSets present but not an object is malformed error",
			body:    `{"StorageSets":"not-an-object","RequestId":"r"}`,
			wantErr: true,
		},
		{
			// TotalCount present but non-integer (a fractional count is malformed):
			// the count cannot establish absence, so the resource is preserved
			// instead of dropped.
			name:    "non-integer TotalCount without StorageSets is conservative error",
			body:    `{"TotalCount":1.5,"RequestId":"r"}`,
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
				if action == "DescribeStorageSets" {
					return tc.status, tc.body
				}
				return http.StatusOK, `{}`
			})
			client := storageSetTestClient(t, srv)
			obj, err := describeStorageSet(client, testStorageSetID)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got obj=%v err=nil", obj)
				}
				if obj != nil {
					t.Fatalf("expected nil object on error, got %v", obj)
				}
				return
			}
			if tc.wantNotFound {
				if err != nil {
					t.Fatalf("expected not-found (nil,nil), got err=%v", err)
				}
				if obj != nil {
					t.Fatalf("expected nil object for not-found, got %v", obj)
				}
				return
			}
			// positive
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if obj == nil {
				t.Fatal("expected object, got nil")
			}
			if got, _ := obj["StorageSetId"].(string); got != testStorageSetID {
				t.Fatalf("expected id %q, got %q", testStorageSetID, got)
			}
		})
	}
}

// TestUnitStorageSetRead_Framework_NoCloud drives the Read method. Malformed
// responses must produce an error diagnostic and PRESERVE state (the id is
// retained); a legitimate not-found (empty array) must remove the resource.
func TestUnitStorageSetRead_Framework_NoCloud(t *testing.T) {
	type step struct {
		name         string
		describeBody string
		wantDiagErr  bool
		wantRemoved  bool
	}
	steps := []step{
		{
			name:         "malformed non-slice StorageSet errors and preserves state",
			describeBody: `{"StorageSets":{"StorageSet":"not-a-list"},"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			name:         "non-object entry errors and preserves state",
			describeBody: `{"StorageSets":{"StorageSet":["x"]},"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			name:         "other id errors and preserves state",
			describeBody: `{"StorageSets":{"StorageSet":[{"StorageSetId":"ss-other"}]},"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			name:         "missing StorageSetId errors and preserves state",
			describeBody: `{"StorageSets":{"StorageSet":[{"ZoneId":"cn-hangzhou-i"}]},"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			name:         "empty array removes resource (not-found)",
			describeBody: `{"StorageSets":{"StorageSet":[]},"RequestId":"r"}`,
			wantDiagErr:  false,
			wantRemoved:  true,
		},
		{
			// TotalCount=0 corroborates absence, so the absent StorageSets wrapper
			// is a legitimate not-found: Read removes the resource.
			name:         "TotalCount:0 with absent StorageSets removes resource (not-found)",
			describeBody: `{"TotalCount":0,"RequestId":"r"}`,
			wantDiagErr:  false,
			wantRemoved:  true,
		},
		{
			// Positive TotalCount without the top-level StorageSets container: the
			// API claims matches but returned none. Read must error and PRESERVE
			// state (the id is retained) instead of dropping the resource.
			name:         "positive TotalCount without StorageSets errors and preserves state",
			describeBody: `{"TotalCount":1,"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			// Positive TotalCount without the inner StorageSet container: same
			// inconsistency — error and preserve state.
			name:         "positive TotalCount without StorageSet key errors and preserves state",
			describeBody: `{"TotalCount":1,"StorageSets":{},"RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			// TotalCount omitted entirely is not a documented DescribeStorageSets
			// shape: absence cannot be established from the count, so Read errors
			// and preserves state rather than dropping the resource.
			name:         "absent StorageSets without TotalCount errors and preserves state",
			describeBody: `{"Foo":"bar","RequestId":"r"}`,
			wantDiagErr:  true,
			wantRemoved:  false,
		},
		{
			name:         "positive refresh keeps state",
			describeBody: describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`),
			wantDiagErr:  false,
			wantRemoved:  false,
		},
	}

	for _, st := range steps {
		t.Run(st.name, func(t *testing.T) {
			body := st.describeBody
			srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
				if action == "DescribeStorageSets" {
					return http.StatusOK, body
				}
				return http.StatusOK, `{}`
			})
			client := storageSetTestClient(t, srv)
			r, sch := newStorageSetResource(t, client)
			readModel := nullStorageSetModel()
			readModel.Id = types.StringValue(testStorageSetID)
			readModel.ZoneId = types.StringValue("cn-hangzhou-i")
			state := storageSetStateFromModel(t, sch, readModel)
			var resp resource.ReadResponse
			// The framework server seeds response State with the resource schema
			// before calling Read; resp.State.Set later needs that Schema.Type().
			resp.State = tfsdk.State{Schema: sch}
			r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
			if st.wantDiagErr && !resp.Diagnostics.HasError() {
				t.Fatalf("expected error diagnostic, got none; state.Raw null=%v", resp.State.Raw.IsNull())
			}
			if !st.wantDiagErr && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if st.wantRemoved {
				if !resp.State.Raw.IsNull() {
					t.Fatalf("expected state removed, but state still present: %v", resp.State.Raw)
				}
				return
			}
			// state must be preserved: the id is still present.
			if resp.State.Raw.IsNull() {
				t.Fatal("expected state preserved, but state was removed")
			}
			var model storageSetResourceModel
			if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
				t.Fatalf("state.Get: %v", diags)
			}
			if model.Id.ValueString() != testStorageSetID {
				t.Fatalf("expected preserved id %q, got %q", testStorageSetID, model.Id.ValueString())
			}
		})
	}
}

// TestUnitStorageSetCreate_Framework_NoCloud drives Create. Unknown Optional+Computed
// plan values must not be sent as empty/zero to CreateStorageSet; known values
// must be sent. State is persisted after a successful create+refresh.
func TestUnitStorageSetCreate_Framework_NoCloud(t *testing.T) {
	positiveDescribe := describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)

	t.Run("unknown optional values are not sent", func(t *testing.T) {
		srv, sent := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "CreateStorageSet":
				return http.StatusOK, `{"StorageSetId":"` + testStorageSetID + `","RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, positiveDescribe
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		planModel := nullStorageSetModel()
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringUnknown()
		planModel.MaxPartitionNumber = types.Int64Unknown()
		planModel.ResourceGroupId = types.StringUnknown()
		planModel.StorageSetName = types.StringUnknown()
		planModel.Tags = types.MapUnknown(types.StringType)
		plan := storageSetPlan(t, sch, planModel)
		var resp resource.CreateResponse
		// Seed the response State with the resource schema, mirroring the
		// framework server, so Create's resp.State.Set has a typed container.
		resp.State = tfsdk.State{Schema: sch}
		r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected create diagnostics: %v", resp.Diagnostics)
		}
		if resp.State.Raw.IsNull() {
			t.Fatal("expected state set after create, got null state")
		}
		var model storageSetResourceModel
		if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
			t.Fatalf("state.Get: %v", diags)
		}
		if model.Id.ValueString() != testStorageSetID {
			t.Fatalf("expected id %q, got %q", testStorageSetID, model.Id.ValueString())
		}
		loadVal, _ := sent.Load("CreateStorageSet")
		createParams, _ := loadVal.(map[string]string)
		for _, key := range []string{"Description", "MaxPartitionNumber", "ResourceGroupId", "StorageSetName"} {
			if v, ok := createParams[key]; ok {
				t.Fatalf("unknown optional %q must not be sent, but got %q", key, v)
			}
		}
		if _, ok := createParams["Tag.1.Key"]; ok {
			t.Fatal("unknown tags must not be sent to CreateStorageSet")
		}
	})

	t.Run("known optional values are sent", func(t *testing.T) {
		srv, sent := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "CreateStorageSet":
				return http.StatusOK, `{"StorageSetId":"` + testStorageSetID + `","RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, positiveDescribe
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		tags, diags := types.MapValueFrom(context.Background(), types.StringType, map[string]string{"k": "v"})
		if diags.HasError() {
			t.Fatalf("tags: %v", diags)
		}
		planModel := nullStorageSetModel()
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringValue("desc")
		planModel.MaxPartitionNumber = types.Int64Value(3)
		planModel.StorageSetName = types.StringValue("name")
		planModel.Tags = tags
		plan := storageSetPlan(t, sch, planModel)
		var resp resource.CreateResponse
		// Seed the response State with the resource schema, mirroring the
		// framework server, so Create's resp.State.Set has a typed container.
		resp.State = tfsdk.State{Schema: sch}
		r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected create diagnostics: %v", resp.Diagnostics)
		}
		loadVal, _ := sent.Load("CreateStorageSet")
		createParams, _ := loadVal.(map[string]string)
		if createParams["Description"] != "desc" {
			t.Fatalf("expected Description=desc, got %q", createParams["Description"])
		}
		if createParams["MaxPartitionNumber"] != "3" {
			t.Fatalf("expected MaxPartitionNumber=3, got %q", createParams["MaxPartitionNumber"])
		}
		if createParams["StorageSetName"] != "name" {
			t.Fatalf("expected StorageSetName=name, got %q", createParams["StorageSetName"])
		}
		if createParams["Tag.1.Key"] != "k" || createParams["Tag.1.Value"] != "v" {
			t.Fatalf("expected Tag.1.Key=k Value=v, got %q=%q %q=%q", "Tag.1.Key", createParams["Tag.1.Key"], "Tag.1.Value", createParams["Tag.1.Value"])
		}
	})

	t.Run("unknown tag value is diagnosed not sent empty", func(t *testing.T) {
		srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "CreateStorageSet":
				return http.StatusOK, `{"StorageSetId":"` + testStorageSetID + `","RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, positiveDescribe
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		// Build a known tags map whose single element value is unknown.
		tags := types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringUnknown()})
		planModel := nullStorageSetModel()
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Tags = tags
		plan := storageSetPlan(t, sch, planModel)
		var resp resource.CreateResponse
		// Seed the response State with the resource schema, mirroring the
		// framework server, so Create's resp.State.Set has a typed container.
		resp.State = tfsdk.State{Schema: sch}
		r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
		if !resp.Diagnostics.HasError() {
			t.Fatal("expected error diagnostic for unknown tag value, got none")
		}
	})
}

// TestUnitStorageSetUpdate_Framework_NoCloud drives the Update method after the
// ForceNew fix. resource_group_id and tags are RequiresReplace: ECS Storage Set
// is not a supported ResourceType of JoinResourceGroup, TagResources or
// UntagResources (they return InvalidResourceType.NotFound for storageset), so
// neither can be updated in place — Terraform recreates the resource on a
// change instead of routing a diff through Update. Update therefore only
// issues ModifyStorageSetAttribute for description/storage_set_name and must
// NEVER issue JoinResourceGroup, TagResources or UntagResources, even when the
// plan carries a different resource_group_id or tags value (a defensive Update
// must remain a safe no-op for those attributes, since the real recreate happens
// at plan time). Omitting tags (unknown plan) preserves the cloud tags with no
// write.
func TestUnitStorageSetUpdate_Framework_NoCloud(t *testing.T) {
	t.Run("desc/name change calls ModifyStorageSetAttribute only", func(t *testing.T) {
		var modified, joined, tagged, untagged bool
		srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "ModifyStorageSetAttribute":
				modified = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "JoinResourceGroup":
				joined = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "TagResources":
				tagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "UntagResources":
				untagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "newdesc", "newname", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		planModel := nullStorageSetModel()
		planModel.Id = types.StringValue(testStorageSetID)
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringValue("newdesc")
		planModel.StorageSetName = types.StringValue("newname")
		planModel.ResourceGroupId = types.StringValue("rg-1")
		planModel.Tags = types.MapNull(types.StringType)
		stateModel := nullStorageSetModel()
		stateModel.Id = types.StringValue(testStorageSetID)
		stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
		stateModel.Description = types.StringValue("olddesc")
		stateModel.StorageSetName = types.StringValue("oldname")
		stateModel.ResourceGroupId = types.StringValue("rg-1")
		stateModel.Tags = types.MapNull(types.StringType)
		var resp resource.UpdateResponse
		// Seed the response State with the resource schema, mirroring the
		// framework server, so Update's resp.State.Set has a typed container.
		resp.State = tfsdk.State{Schema: sch}
		r.Update(context.Background(), resource.UpdateRequest{
			Plan:  storageSetPlan(t, sch, planModel),
			State: storageSetStateFromModel(t, sch, stateModel),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update diagnostics: %v", resp.Diagnostics)
		}
		if !modified {
			t.Fatal("expected ModifyStorageSetAttribute for desc/name change")
		}
		if joined || tagged || untagged {
			t.Fatalf("Update must not call JoinResourceGroup/Tag/Untag (ForceNew): joined=%v tagged=%v untagged=%v", joined, tagged, untagged)
		}
	})

	t.Run("tags differ but Update issues no Tag/Untag (ForceNew)", func(t *testing.T) {
		var joined, tagged, untagged, modified bool
		srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "ModifyStorageSetAttribute":
				modified = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "JoinResourceGroup":
				joined = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "TagResources":
				tagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "UntagResources":
				untagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[{"Key":"k","Value":"v"}]}`)
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		planTags, _ := types.MapValueFrom(context.Background(), types.StringType, map[string]string{"new": "v"})
		stateTags, _ := types.MapValueFrom(context.Background(), types.StringType, map[string]string{"old": "v"})
		planModel := nullStorageSetModel()
		planModel.Id = types.StringValue(testStorageSetID)
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringValue("d")
		planModel.StorageSetName = types.StringValue("n")
		planModel.ResourceGroupId = types.StringValue("rg-1")
		planModel.Tags = planTags
		stateModel := nullStorageSetModel()
		stateModel.Id = types.StringValue(testStorageSetID)
		stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
		stateModel.Description = types.StringValue("d")
		stateModel.StorageSetName = types.StringValue("n")
		stateModel.ResourceGroupId = types.StringValue("rg-1")
		stateModel.Tags = stateTags
		var resp resource.UpdateResponse
		resp.State = tfsdk.State{Schema: sch}
		r.Update(context.Background(), resource.UpdateRequest{
			Plan:  storageSetPlan(t, sch, planModel),
			State: storageSetStateFromModel(t, sch, stateModel),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update diagnostics: %v", resp.Diagnostics)
		}
		if modified {
			t.Fatal("desc/name unchanged: ModifyStorageSetAttribute must not be called")
		}
		if joined || tagged || untagged {
			t.Fatalf("Update must not issue JoinResourceGroup/Tag/Untag for a tags change (ForceNew recreates): joined=%v tagged=%v untagged=%v", joined, tagged, untagged)
		}
	})

	t.Run("resource_group_id differs but Update issues no JoinResourceGroup (ForceNew)", func(t *testing.T) {
		var joined, tagged, untagged, modified bool
		srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "ModifyStorageSetAttribute":
				modified = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "JoinResourceGroup":
				joined = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "TagResources":
				tagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "UntagResources":
				untagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-2", 3, `{"Tag":[]}`)
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		planModel := nullStorageSetModel()
		planModel.Id = types.StringValue(testStorageSetID)
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringValue("d")
		planModel.StorageSetName = types.StringValue("n")
		planModel.ResourceGroupId = types.StringValue("rg-2")
		planModel.Tags = types.MapNull(types.StringType)
		stateModel := nullStorageSetModel()
		stateModel.Id = types.StringValue(testStorageSetID)
		stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
		stateModel.Description = types.StringValue("d")
		stateModel.StorageSetName = types.StringValue("n")
		stateModel.ResourceGroupId = types.StringValue("rg-1")
		stateModel.Tags = types.MapNull(types.StringType)
		var resp resource.UpdateResponse
		resp.State = tfsdk.State{Schema: sch}
		r.Update(context.Background(), resource.UpdateRequest{
			Plan:  storageSetPlan(t, sch, planModel),
			State: storageSetStateFromModel(t, sch, stateModel),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update diagnostics: %v", resp.Diagnostics)
		}
		if modified {
			t.Fatal("desc/name unchanged: ModifyStorageSetAttribute must not be called")
		}
		if joined {
			t.Fatal("Update must not issue JoinResourceGroup for a resource_group_id change (ForceNew recreates)")
		}
		if tagged || untagged {
			t.Fatalf("Update must not issue Tag/Untag: tagged=%v untagged=%v", tagged, untagged)
		}
	})

	t.Run("omitted (unknown) tags preserves state, no write", func(t *testing.T) {
		var joined, tagged, untagged, modified bool
		srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
			switch action {
			case "ModifyStorageSetAttribute":
				modified = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "JoinResourceGroup":
				joined = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "TagResources":
				tagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "UntagResources":
				untagged = true
				return http.StatusOK, `{"RequestId":"r"}`
			case "DescribeStorageSets":
				return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[{"Key":"managed","Value":"keep"}]}`)
			}
			return http.StatusOK, `{}`
		})
		client := storageSetTestClient(t, srv)
		r, sch := newStorageSetResource(t, client)
		stateTags, _ := types.MapValueFrom(context.Background(), types.StringType, map[string]string{"managed": "keep"})
		planModel := nullStorageSetModel()
		planModel.Id = types.StringValue(testStorageSetID)
		planModel.ZoneId = types.StringValue("cn-hangzhou-i")
		planModel.Description = types.StringValue("d")
		planModel.StorageSetName = types.StringValue("n")
		planModel.ResourceGroupId = types.StringValue("rg-1")
		planModel.Tags = types.MapUnknown(types.StringType)
		stateModel := nullStorageSetModel()
		stateModel.Id = types.StringValue(testStorageSetID)
		stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
		stateModel.Description = types.StringValue("d")
		stateModel.StorageSetName = types.StringValue("n")
		stateModel.ResourceGroupId = types.StringValue("rg-1")
		stateModel.Tags = stateTags
		var resp resource.UpdateResponse
		resp.State = tfsdk.State{Schema: sch}
		r.Update(context.Background(), resource.UpdateRequest{
			Plan:  storageSetPlan(t, sch, planModel),
			State: storageSetStateFromModel(t, sch, stateModel),
		}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected update diagnostics: %v", resp.Diagnostics)
		}
		if modified {
			t.Fatal("desc/name unchanged: ModifyStorageSetAttribute must not be called")
		}
		if joined || tagged || untagged {
			t.Fatalf("omitted (unknown) tags must issue no JoinResourceGroup/Tag/Untag: joined=%v tagged=%v untagged=%v", joined, tagged, untagged)
		}
	})
}

// TestUnitStorageSetCreate_NullTagElement_Framework_NoCloud verifies Create rejects
// a null tag element value before issuing CreateStorageSet. A null element
// (e.g. map[string]*string{"invalid": nil} mapped via MapValueFrom) used to be
// turned into an empty string by ValueString and sent, creating an unintended
// empty tag. An explicit empty string "" remains a valid value and is NOT null.
func TestUnitStorageSetCreate_NullTagElement_Framework_NoCloud(t *testing.T) {
	srv, sent := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		switch action {
		case "CreateStorageSet":
			return http.StatusOK, `{"StorageSetId":"` + testStorageSetID + `","RequestId":"r"}`
		case "DescribeStorageSets":
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	// Known map with a null element value (the shape produced by MapValueFrom
	// with a *string nil). MapValueFrom accepts it (diags empty), so the
	// provider must reject the null element before issuing CreateStorageSet.
	nullTags, nullDiags := types.MapValueFrom(context.Background(), types.StringType, map[string]*string{"invalid": nil})
	if nullDiags.HasError() {
		t.Fatalf("build nullTags: %v", nullDiags)
	}
	planModel := nullStorageSetModel()
	planModel.ZoneId = types.StringValue("cn-hangzhou-i")
	planModel.Tags = nullTags
	plan := storageSetPlan(t, sch, planModel)
	var resp resource.CreateResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostic for null tag element value, got none")
	}
	// CreateStorageSet must not have been issued for an invalid tag element.
	if v, ok := sent.Load("CreateStorageSet"); ok {
		t.Fatalf("CreateStorageSet must not be called for a null tag element, but was called with %v", v)
	}
}

// TestUnitStorageSetImport_Framework_NoCloud drives ImportState passthrough.
func TestUnitStorageSetImport_Framework_NoCloud(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	var resp resource.ImportStateResponse
	// Mirror the framework server, which seeds the response State with an
	// unknown-valued object of the schema type before invoking ImportState, so
	// ImportStatePassthroughID's SetAttribute has a typed container to write
	// into. A bare tfsdk.State{Schema: sch} has a zero-value Raw (no type),
	// which makes SetAttribute fail with "value missing type".
	resp.State = tfsdk.State{
		Schema: sch,
		Raw:    tftypes.NewValue(sch.Type().TerraformType(context.Background()), tftypes.UnknownValue),
	}
	r.ImportState(context.Background(), resource.ImportStateRequest{ID: testStorageSetID}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected import diagnostics: %v", resp.Diagnostics)
	}
	var model storageSetResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Get: %v", diags)
	}
	if model.Id.ValueString() != testStorageSetID {
		t.Fatalf("expected import id %q, got %q", testStorageSetID, model.Id.ValueString())
	}
	// ImportStatePassthroughID targets path.Root("id"); ensure the resource uses it.
	var _ resource.ResourceWithImportState = (*storageSetResource)(nil)
	_ = path.Root("id")
}

// TestUnitStorageSetExists_Framework_NoCloud drives StorageSetExists — the function
// the acceptance-test destroy check (testAccCheckEcsStorageSetDestroy) delegates
// to — over the response-consistency matrix. It reuses the same strict not-found
// logic as Read, so a legitimate not-found (empty array, TotalCount=0 with an
// absent container, the InvalidStorageSetId.NotFound error) yields (false, nil),
// while a positive-count response without a container, an omitted/malformed
// TotalCount, or a mismatched id yields (false, error) so a leaked or
// inconsistent resource surfaces in CheckDestroy instead of masking as absent.
func TestUnitStorageSetExists_Framework_NoCloud(t *testing.T) {
	cases := []struct {
		name       string
		body       string
		status     int
		wantExists bool
		wantErr    bool
	}{
		{
			name:       "matching id exists",
			body:       describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`),
			wantExists: true,
		},
		{
			name:       "explicit empty array is not-found (false, nil)",
			body:       `{"StorageSets":{"StorageSet":[]},"RequestId":"r"}`,
			wantExists: false,
		},
		{
			name:       "TotalCount:0 with absent StorageSets is not-found (false, nil)",
			body:       `{"TotalCount":0,"RequestId":"r"}`,
			wantExists: false,
		},
		{
			name:       "TotalCount:0 with absent StorageSet key is not-found (false, nil)",
			body:       `{"TotalCount":0,"StorageSets":{},"RequestId":"r"}`,
			wantExists: false,
		},
		{
			name:       "not-found error code is not-found (false, nil)",
			body:       `{"Code":"InvalidStorageSetId.NotFound","Message":"The specified parameter StorageSetId was not found.","RequestId":"r"}`,
			wantExists: false,
		},
		{
			// Positive count without the container: CheckDestroy must surface this
			// as an error (the resource may still exist) rather than false.
			name:    "positive TotalCount without StorageSets is error (preserve state)",
			body:    `{"TotalCount":1,"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "positive TotalCount without StorageSet key is error (preserve state)",
			body:    `{"TotalCount":1,"StorageSets":{},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			// Omitted TotalCount: absence cannot be established, so CheckDestroy
			// errors instead of假通过 as destroyed.
			name:    "absent StorageSets without TotalCount is conservative error",
			body:    `{"Foo":"bar","RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "StorageSet present but not a list is error",
			body:    `{"StorageSets":{"StorageSet":"not-a-list"},"RequestId":"r"}`,
			wantErr: true,
		},
		{
			name:    "entry with other id is error (preserve state)",
			body:    `{"StorageSets":{"StorageSet":[{"StorageSetId":"ss-other"}]},"RequestId":"r"}`,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
				if action == "DescribeStorageSets" {
					return tc.status, tc.body
				}
				return http.StatusOK, `{}`
			})
			client := storageSetTestClient(t, srv)
			exists, err := StorageSetExists(client, testStorageSetID)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got exists=%v err=nil", exists)
				}
				if exists {
					t.Fatalf("expected exists=false on error, got true")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got err=%v", err)
			}
			if exists != tc.wantExists {
				t.Fatalf("expected exists=%v, got %v", tc.wantExists, exists)
			}
		})
	}
}

// TestUnitStorageSetState_Mapping verifies storageSetState maps the Describe object
// onto the model, including max_partition_number across numeric encodings and
// tag filtering of system keys.
func TestUnitStorageSetState_Mapping(t *testing.T) {
	cases := []struct {
		name       string
		object     map[string]interface{}
		wantMax    int64
		wantMaxSet bool
		wantTags   map[string]string
	}{
		{
			name: "float64 partition",
			object: map[string]interface{}{
				"StorageSetId":              testStorageSetID,
				"StorageSetPartitionNumber": float64(5),
				"Tags":                      map[string]interface{}{"Tag": []interface{}{map[string]interface{}{"Key": "k", "Value": "v"}}},
			},
			wantMax: 5, wantMaxSet: true, wantTags: map[string]string{"k": "v"},
		},
		{
			// SDK RpcPost decodes responses with json.Decoder.UseNumber(), so a
			// real HTTP StorageSetPartitionNumber:3 arrives as json.Number, not
			// float64. The switch must match it, not fall through to null.
			name: "json.Number partition (SDK UseNumber path)",
			object: map[string]interface{}{
				"StorageSetId":              testStorageSetID,
				"StorageSetPartitionNumber": json.Number("3"),
			},
			wantMax: 3, wantMaxSet: true, wantTags: map[string]string{},
		},
		{
			// A non-integer json.Number must NOT be truncated to 3 or coerced
			// to 0; it stays null (field absent) with accurate semantics.
			name: "non-integer json.Number stays null (not 0)",
			object: map[string]interface{}{
				"StorageSetId":              testStorageSetID,
				"StorageSetPartitionNumber": json.Number("3.5"),
			},
			wantMaxSet: false, wantTags: map[string]string{},
		},
		{
			name: "string partition",
			object: map[string]interface{}{
				"StorageSetId":              testStorageSetID,
				"StorageSetPartitionNumber": "7",
			},
			wantMax: 7, wantMaxSet: true, wantTags: map[string]string{},
		},
		{
			name: "no partition field => null",
			object: map[string]interface{}{
				"StorageSetId": testStorageSetID,
			},
			wantMaxSet: false, wantTags: map[string]string{},
		},
		{
			name: "system tags filtered",
			object: map[string]interface{}{
				"StorageSetId": testStorageSetID,
				"Tags": map[string]interface{}{"Tag": []interface{}{
					map[string]interface{}{"Key": "acs:foo", "Value": "x"},
					map[string]interface{}{"Key": "keep", "Value": "y"},
				}},
			},
			wantMaxSet: false, wantTags: map[string]string{"keep": "y"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := storageSetState(context.Background(), tc.object)
			if model.Id.ValueString() != testStorageSetID {
				t.Fatalf("expected id %q, got %q", testStorageSetID, model.Id.ValueString())
			}
			if tc.wantMaxSet {
				if model.MaxPartitionNumber.ValueInt64() != tc.wantMax {
					t.Fatalf("expected max %d, got %d", tc.wantMax, model.MaxPartitionNumber.ValueInt64())
				}
			} else if !model.MaxPartitionNumber.IsNull() {
				t.Fatalf("expected null max, got %d", model.MaxPartitionNumber.ValueInt64())
			}
			// tags
			var gotTags map[string]types.String
			if diags := model.Tags.ElementsAs(context.Background(), &gotTags, false); diags.HasError() {
				t.Fatalf("tags: %v", diags)
			}
			if len(gotTags) != len(tc.wantTags) {
				t.Fatalf("expected %d tags, got %d (%v)", len(tc.wantTags), len(gotTags), gotTags)
			}
			for k, v := range tc.wantTags {
				got, ok := gotTags[k]
				if !ok {
					t.Fatalf("expected tag %q, missing", k)
				}
				if got.ValueString() != v {
					t.Fatalf("expected tag %q=%q, got %q", k, v, got.ValueString())
				}
			}
		})
	}
}
