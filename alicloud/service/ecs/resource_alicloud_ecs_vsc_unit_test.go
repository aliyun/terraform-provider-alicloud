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

	credential "github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// tagMap builds a tags map from the given key/value pairs the way the plan
// decoder would, so the pure helpers can be exercised without the cloud.
func tagMap(t *testing.T, tags ...[2]string) types.Map {
	t.Helper()
	m := map[string]string{}
	for _, kv := range tags {
		m[kv[0]] = kv[1]
	}
	v, d := types.MapValueFrom(context.Background(), types.StringType, m)
	ensureNoDiag(t, d)
	return v
}

func ensureNoDiag(t *testing.T, d diag.Diagnostics) {
	t.Helper()
	if d.HasError() {
		t.Fatalf("unexpected diagnostics: %v", d.Errors())
	}
}

// TestUnitParseEcsVscResponse covers the Read robustness rules: a malformed or
// non-matching DescribeVscs response is an error, never a silent NotFound that
// would RemoveResource a live entry. Empty array stays NotFound (the Ecs
// 2014-05-26 contract's NotFound shape); a matched entry missing Status is an
// error so Read preserves state instead of silently forgetting a live, ID'd
// resource; a sibling entry is never mapped onto the queried id.
func TestUnitParseEcsVscResponse(t *testing.T) {
	const id = "vsc-1"

	cases := []struct {
		name    string
		resp    map[string]interface{}
		wantObj map[string]interface{}
		wantErr bool
	}{
		{
			name: "matching entry returns the object",
			resp: map[string]interface{}{
				"Vscs": []interface{}{
					map[string]interface{}{"VscId": id, "Status": "Available", "InstanceId": "i-1"},
				},
			},
			wantObj: map[string]interface{}{"VscId": id, "Status": "Available", "InstanceId": "i-1"},
		},
		{
			name: "matching entry found after skipping a sibling",
			resp: map[string]interface{}{
				"Vscs": []interface{}{
					map[string]interface{}{"VscId": "vsc-other", "Status": "Available"},
					map[string]interface{}{"VscId": id, "Status": "Available"},
				},
			},
			wantObj: map[string]interface{}{"VscId": id, "Status": "Available"},
		},
		{
			name: "empty array is NotFound",
			resp: map[string]interface{}{
				"Vscs": []interface{}{},
			},
		},
		{
			name: "matching entry without Status is an error (preserve state)",
			resp: map[string]interface{}{
				"Vscs": []interface{}{
					map[string]interface{}{"VscId": id},
				},
			},
			wantErr: true,
		},
		{
			name:    "missing $.Vscs root is an error",
			resp:    map[string]interface{}{"RequestId": "r-1"},
			wantErr: true,
		},
		{
			name:    "null $.Vscs root is an error",
			resp:    map[string]interface{}{"Vscs": nil},
			wantErr: true,
		},
		{
			name:    "non-array $.Vscs is an error",
			resp:    map[string]interface{}{"Vscs": "not-an-array"},
			wantErr: true,
		},
		{
			name: "non-empty result without the queried id is an error",
			resp: map[string]interface{}{
				"Vscs": []interface{}{
					map[string]interface{}{"VscId": "vsc-other", "Status": "Available"},
				},
			},
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			obj, err := parseEcsVscResponse(c.resp, id)
			if c.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got obj=%v err=nil", obj)
				}
				if obj != nil {
					t.Fatalf("expected nil object on error, got %v", obj)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if c.wantObj == nil {
				if obj != nil {
					t.Fatalf("expected NotFound (nil), got %v", obj)
				}
				return
			}
			if obj == nil {
				t.Fatalf("expected object %v, got nil", c.wantObj)
			}
			if got := obj["VscId"]; got != c.wantObj["VscId"] {
				t.Fatalf("VscId = %v, want %v", got, c.wantObj["VscId"])
			}
		})
	}
}

// TestUnitVscCreateTagParams verifies the CreateVsc request shape built from a
// decoded tags map: 1-based Tag.<i>.Key / Tag.<i>.Value, exactly one pair per
// tag. This is the no-cloud HTTP-request positive case. Map iteration order is
// unspecified, so the test asserts the pair invariant — each tag lands as
// exactly one matching Tag.<i>.Key/Tag.<i>.Value pair — rather than exact
// indices.
func TestUnitVscCreateTagParams(t *testing.T) {
	tags := map[string]string{
		"env":   "dev",
		"team":  "qa",
		"empty": "",
	}
	params := vscCreateTagParams(tags)
	if len(params) != len(tags)*2 {
		t.Fatalf("params has %d entries, want %d (%v)", len(params), len(tags)*2, params)
	}
	seen := map[string]bool{}
	for k, want := range tags {
		idx := 0
		for i := 1; i <= len(tags); i++ {
			if params[fmt.Sprintf("Tag.%d.Key", i)] == k {
				idx = i
				break
			}
		}
		if idx == 0 {
			t.Fatalf("tag %q has no Tag.<i>.Key entry (params=%v)", k, params)
		}
		if got := params[fmt.Sprintf("Tag.%d.Value", idx)]; got != want {
			t.Fatalf("Tag.%d.Value = %q, want %q (for key %q)", idx, got, want, k)
		}
		if seen[k] {
			t.Fatalf("tag %q appeared at more than one index", k)
		}
		seen[k] = true
	}
}

// TestUnitTagsToMap covers the decoder that backs the CreateVsc Tag.<i> shape and
// the Update diff: a known map decodes to a key->value map with the empty
// string preserved (so an unchanged empty value is not seen as a change); a
// null or unknown map yields an empty map so Create/Update skip the tag path.
func TestUnitTagsToMap(t *testing.T) {
	ctx := context.Background()

	t.Run("known map decodes with empty value preserved", func(t *testing.T) {
		m := tagMap(t, [2]string{"env", "dev"}, [2]string{"note", ""})
		var d diag.Diagnostics
		got := tagsToMap(ctx, &d, m)
		ensureNoDiag(t, d)
		if got["env"] != "dev" {
			t.Fatalf("got[env] = %q, want dev", got["env"])
		}
		if _, ok := got["note"]; !ok {
			t.Fatalf("got[note] missing; empty value must be preserved, not dropped")
		}
		if got["note"] != "" {
			t.Fatalf("got[note] = %q, want empty string", got["note"])
		}
		if len(got) != 2 {
			t.Fatalf("map has %d entries, want 2 (%v)", len(got), got)
		}
	})

	t.Run("null map yields empty", func(t *testing.T) {
		var d diag.Diagnostics
		got := tagsToMap(ctx, &d, types.MapNull(types.StringType))
		ensureNoDiag(t, d)
		if len(got) != 0 {
			t.Fatalf("expected empty map for null, got %v", got)
		}
	})

	t.Run("unknown map yields empty", func(t *testing.T) {
		var d diag.Diagnostics
		got := tagsToMap(ctx, &d, types.MapUnknown(types.StringType))
		ensureNoDiag(t, d)
		if len(got) != 0 {
			t.Fatalf("expected empty map for unknown, got %v", got)
		}
	})
}

// TestUnitEcsVscStateTagValueFidelity is the F5 regression: a Read tag whose
// TagValue is the legitimate empty string (or absent) must round-trip as a
// known value, never null, because the schema marks value Required. A null
// would drift against a config value="" on the next plan. This asserts the
// state-fidelity property without running the apply protocol or guessing an
// apply-time error message.
func TestUnitEcsVscStateTagValueFidelity(t *testing.T) {
	resp := map[string]interface{}{
		"VscId":      "vsc-1",
		"Status":     "Available",
		"InstanceId": "i-1",
		"Tags": []interface{}{
			map[string]interface{}{"TagKey": "env", "TagValue": ""},    // empty string
			map[string]interface{}{"TagKey": "team", "TagValue": "qa"}, // normal value
			map[string]interface{}{"TagKey": "note"},                   // TagValue absent
		},
	}
	model, diags := ecsVscState(context.Background(), "vsc-1", resp)
	ensureNoDiag(t, diags)

	got := tagsToMap(context.Background(), &diag.Diagnostics{}, model.Tags)

	want := map[string]string{"env": "", "team": "qa", "note": ""}
	if len(got) != len(want) {
		t.Fatalf("got %d tags, want %d (%v)", len(got), len(want), got)
	}
	for key, w := range want {
		v, ok := got[key]
		if !ok {
			t.Fatalf("tag %q dropped; empty/absent values must be preserved as \"\", not omitted", key)
		}
		// A map value is a string, never null — the prior "value is Required,
		// never null in state" invariant holds structurally under MapAttribute
		// (there is no nullable value to drift against a config value="").
		if v != w {
			t.Fatalf("tag %q value = %q, want %q", key, v, w)
		}
	}
}

// TestUnitApplyReadResult exercises the Read decision end-to-end through the
// parseEcsVscResponse -> applyReadResult chain (the exact logic Read delegates
// to), without a live client. The true Read-method regression — driving Read
// -> getEcsVsc(RpcPost) -> HTTP -> parseEcsVscResponse -> applyReadResult ->
// resp.State.Set/RemoveResource against a local httptest.Server — is
// TestUnitEcsVscReadCloudless below in this file. This test stays at the pure-helper
// seam so the state-preservation rule is exercisable without a Framework
// ReadRequest: a matched DescribeVscs entry that omits Status is a malformed
// response, surfaced as a diagnostic with remove=false so the ID'd resource is
// NOT silently RemoveResource-d, while a legitimate empty-array NotFound still
// removes.
func TestUnitApplyReadResult(t *testing.T) {
	ctx := context.Background()
	const id = "vsc-1"
	prior := ecsVscResourceModel{
		Id:         types.StringValue(id),
		InstanceId: types.StringValue("i-1"),
		Tags:       types.MapNull(types.StringType),
	}

	t.Run("matched entry missing Status preserves state with a diagnostic", func(t *testing.T) {
		resp := map[string]interface{}{
			"Vscs": []interface{}{
				map[string]interface{}{"VscId": id},
			},
		}
		obj, err := parseEcsVscResponse(resp, id)
		if err == nil {
			t.Fatalf("parseEcsVscResponse: expected an error for a Status-less match, got obj=%v", obj)
		}
		got, diags, remove := applyReadResult(ctx, prior, obj, err)
		if remove {
			t.Fatalf("remove=true; a matched-ID entry missing Status must preserve state, not RemoveResource")
		}
		if !diags.HasError() {
			t.Fatalf("expected an error diagnostic so the malformed response is visible, got diags=%v", diags)
		}
		// The ID'd resource must not be silently forgotten: the prior identity is retained.
		if got.Id.ValueString() != id {
			t.Fatalf("prior Id lost: got %q, want %q (state must be preserved on a malformed response)", got.Id.ValueString(), id)
		}
	})

	t.Run("legitimate empty-array NotFound removes the resource", func(t *testing.T) {
		resp := map[string]interface{}{
			"Vscs": []interface{}{},
		}
		obj, err := parseEcsVscResponse(resp, id)
		if err != nil {
			t.Fatalf("parseEcsVscResponse: unexpected error for an empty array: %v", err)
		}
		if obj != nil {
			t.Fatalf("parseEcsVscResponse: expected nil object for an empty array, got %v", obj)
		}
		_, diags, remove := applyReadResult(ctx, prior, obj, err)
		if !remove {
			t.Fatalf("remove=false; an empty Vscs array is the legitimate NotFound and must RemoveResource")
		}
		if diags.HasError() {
			t.Fatalf("legitimate NotFound must not raise a diagnostic, got %v", diags)
		}
	})

	t.Run("well-formed response refreshes state without removing", func(t *testing.T) {
		resp := map[string]interface{}{
			"Vscs": []interface{}{
				map[string]interface{}{
					"VscId":      id,
					"Status":     "Available",
					"InstanceId": "i-2",
				},
			},
		}
		obj, err := parseEcsVscResponse(resp, id)
		if err != nil {
			t.Fatalf("parseEcsVscResponse: unexpected error: %v", err)
		}
		got, diags, remove := applyReadResult(ctx, prior, obj, err)
		ensureNoDiag(t, diags)
		if remove {
			t.Fatalf("remove=true; a matched entry must refresh, not remove")
		}
		if got.Id.ValueString() != id {
			t.Fatalf("refreshed Id = %q, want %q", got.Id.ValueString(), id)
		}
		if got.Status.ValueString() != "Available" {
			t.Fatalf("refreshed Status = %q, want Available", got.Status.ValueString())
		}
		if got.InstanceId.ValueString() != "i-2" {
			t.Fatalf("refreshed InstanceId = %q, want i-2", got.InstanceId.ValueString())
		}
	})
}

// vscDescribeResponse builds a DescribeVscs 200 body whose $.Vscs is the given
// slice, so each Read scenario can express the Ecs 2014-05-26 response shape
// (an empty call is the legitimate NotFound: an empty Vscs array).
func vscDescribeResponse(vscs ...map[string]interface{}) map[string]interface{} {
	arr := make([]interface{}, len(vscs))
	for i, v := range vscs {
		arr[i] = v
	}
	return map[string]interface{}{"Vscs": arr}
}

// vscTestResource wires an httptest.Server behind a real AliyunClient and
// injects the client into a fresh ecsVscResource through the fwadapt Configure
// seam, so Read and Create run end-to-end against canned responses without
// touching the cloud. handler serves the canned per-Action responses. No
// production code is changed: Configure is the injection point the production
// resource already exposes.
func vscTestResource(t *testing.T, handler http.HandlerFunc) *ecsVscResource {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	endpoints := &sync.Map{}
	// The tea-rpc client builds the request URL as {protocol}://{endpoint}/,
	// so the endpoint must be a bare host:port (the way production stores
	// "ecs.cn-hangzhou.aliyuncs.com"), not a full http:// URL — otherwise the
	// scheme doubles up and the host resolves to "http".
	endpoints.Store("ecs", strings.TrimPrefix(srv.URL, "http://"))

	// Mirror the provider.Configure access_key path: a credentials.Config with
	// Type=access_key, then NewCredential. SkipRegionValidation + AccountType
	// keep Client() off the network (no loadAndValidate, no getAccountType), and
	// Endpoints["ecs"] routes DescribeVscs to the test server.
	credConfig := new(credential.Config).
		SetType("access_key").
		SetAccessKeyId("vsc-test-ak").
		SetAccessKeySecret("vsc-test-sk")
	cred, err := credential.NewCredential(credConfig)
	if err != nil {
		t.Fatalf("NewCredential: %v", err)
	}

	cfg := &connectivity.Config{
		AccessKey:            "vsc-test-ak",
		SecretKey:            "vsc-test-sk",
		RegionId:             "cn-hangzhou",
		Protocol:             "HTTP",
		AccountType:          "Domestic",
		SkipRegionValidation: true,
		Endpoints:            endpoints,
		Credential:           cred,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("cfg.Client: %v", err)
	}

	r := NewEcsVscResource().(*ecsVscResource)
	var cfgResp resource.ConfigureResponse
	r.Configure(context.Background(), resource.ConfigureRequest{ProviderData: client}, &cfgResp)
	if cfgResp.Diagnostics.HasError() {
		t.Fatalf("Configure diagnostics: %v", cfgResp.Diagnostics.Errors())
	}
	if r.Client() == nil {
		t.Fatal("client not injected after Configure")
	}
	return r
}

// vscReadEnv wires a Read-only server: DescribeVscs returns describe (or a 4xx
// API error when apiError is non-nil), and any other Action is rejected. It is
// the shared Read environment TestUnitEcsVscReadCloudless uses.
func vscReadEnv(t *testing.T, describe map[string]interface{}, apiError map[string]interface{}) *ecsVscResource {
	t.Helper()
	return vscTestResource(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("Action") != "DescribeVscs" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if apiError != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(apiError)
			return
		}
		_ = json.NewEncoder(w).Encode(describe)
	})
}

// vscPriorState builds a minimal prior model: an ID'd resource with an instance
// id and null tags, the shape Read receives when refreshing an existing entry.
func vscPriorState(id string) ecsVscResourceModel {
	return ecsVscResourceModel{
		Id:         types.StringValue(id),
		InstanceId: types.StringValue("i-prior"),
		Tags:       types.MapNull(types.StringType),
	}
}

// vscReadCall drives a Read against the given prior model and returns the
// response for assertions. The request/state schemas come from the resource's
// own Schema, and the prior model is seeded into ReadRequest.State the way the
// Framework runtime would before invoking Read.
func vscReadCall(t *testing.T, r *ecsVscResource, prior ecsVscResourceModel) resource.ReadResponse {
	t.Helper()
	ctx := context.Background()

	schemaResp := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	readReq := resource.ReadRequest{State: tfsdk.State{Schema: schema}}
	if diags := readReq.State.Set(ctx, &prior); diags.HasError() {
		t.Fatalf("seed prior state: %v", diags.Errors())
	}

	readResp := resource.ReadResponse{State: tfsdk.State{Schema: schema}}
	r.Read(ctx, readReq, &readResp)
	return readResp
}

// readRetainedId gets the Id out of a (non-removed) ReadResponse.State, failing
// the test if the state was removed or unreadable. It is the shared assertion
// for every "preserve state" scenario.
func readRetainedId(t *testing.T, resp resource.ReadResponse, wantId string) {
	t.Helper()
	if resp.State.Raw.IsNull() {
		t.Fatal("state was removed; this scenario must preserve state, not RemoveResource")
	}
	var got ecsVscResourceModel
	if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
		t.Fatalf("state Get: %v", diags.Errors())
	}
	if got.Id.ValueString() != wantId {
		t.Fatalf("prior Id lost: got %q, want %q (state must be preserved on this response)", got.Id.ValueString(), wantId)
	}
}

// TestUnitEcsVscReadCloudless is the F-TEST-GAP regression: it exercises the real
// Framework Read method end-to-end — Read -> getEcsVsc(RpcPost) ->
// parseEcsVscResponse -> applyReadResult -> resp.State.Set/RemoveResource —
// against a local httptest.Server, asserting on the actual ReadResponse.State
// (Raw.IsNull for RemoveResource; Get for retained/refreshed models) and
// Diagnostics. It proves the state-preservation rule the QA verify round flagged
// as a test gap: a matched DescribeVscs entry missing Status, a bad root, a
// wrong id, and a 4xx API error all return an error diagnostic and PRESERVE the
// prior state, while a legitimate empty-array NotFound removes the resource and
// a well-formed response refreshes it.
//
// The prior TestUnitApplyReadResult only tested the pure parseEcsVscResponse ->
// applyReadResult helper chain; this test closes the gap by driving the real
// Read method, the RpcPost HTTP call, and the resp.State.Set/RemoveResource
// outcome. No production code is changed.
func TestUnitEcsVscReadCloudless(t *testing.T) {
	const id = "vsc-1"

	t.Run("matched entry missing Status errors and preserves state", func(t *testing.T) {
		r := vscReadEnv(t, vscDescribeResponse(map[string]interface{}{"VscId": id}), nil)
		resp := vscReadCall(t, r, vscPriorState(id))

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected an error diagnostic for a Status-less match, got %v", resp.Diagnostics)
		}
		readRetainedId(t, resp, id)
	})

	t.Run("legitimate empty-array NotFound removes the resource", func(t *testing.T) {
		r := vscReadEnv(t, vscDescribeResponse(), nil) // empty Vscs array
		resp := vscReadCall(t, r, vscPriorState(id))

		if resp.Diagnostics.HasError() {
			t.Fatalf("legitimate NotFound must not raise a diagnostic, got %v", resp.Diagnostics)
		}
		if !resp.State.Raw.IsNull() {
			t.Fatal("state was not removed; an empty Vscs array is the legitimate NotFound and must RemoveResource")
		}
	})

	t.Run("missing $.Vscs root errors and preserves state", func(t *testing.T) {
		r := vscReadEnv(t, map[string]interface{}{"RequestId": "r-1"}, nil)
		resp := vscReadCall(t, r, vscPriorState(id))

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected an error diagnostic for a missing $.Vscs root, got %v", resp.Diagnostics)
		}
		readRetainedId(t, resp, id)
	})

	t.Run("non-array $.Vscs errors and preserves state", func(t *testing.T) {
		r := vscReadEnv(t, map[string]interface{}{"Vscs": "not-an-array"}, nil)
		resp := vscReadCall(t, r, vscPriorState(id))

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected an error diagnostic for a non-array $.Vscs, got %v", resp.Diagnostics)
		}
		readRetainedId(t, resp, id)
	})

	t.Run("wrong id in a non-empty result errors and preserves state", func(t *testing.T) {
		r := vscReadEnv(t, vscDescribeResponse(
			map[string]interface{}{"VscId": "vsc-other", "Status": "Available"},
		), nil)
		resp := vscReadCall(t, r, vscPriorState(id))

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected an error diagnostic for a non-matching id, got %v", resp.Diagnostics)
		}
		readRetainedId(t, resp, id)
	})

	t.Run("4xx API error errors and preserves state", func(t *testing.T) {
		r := vscReadEnv(t, nil, map[string]interface{}{
			"Code":      "InternalError",
			"Message":   "boom",
			"RequestId": "r-1",
		})
		resp := vscReadCall(t, r, vscPriorState(id))

		if !resp.Diagnostics.HasError() {
			t.Fatalf("expected an error diagnostic for a 4xx API error, got %v", resp.Diagnostics)
		}
		readRetainedId(t, resp, id)
	})

	t.Run("well-formed response refreshes state without removing", func(t *testing.T) {
		r := vscReadEnv(t, vscDescribeResponse(map[string]interface{}{
			"VscId":           id,
			"Status":          "Available",
			"InstanceId":      "i-refreshed",
			"ResourceGroupId": "rg-1",
			"VscName":         "name-1",
			"VscType":         "type-1",
			"Description":     "desc-1",
			"Tags": []interface{}{
				map[string]interface{}{"TagKey": "env", "TagValue": "dev"},
			},
		}), nil)
		resp := vscReadCall(t, r, vscPriorState(id))

		if resp.Diagnostics.HasError() {
			t.Fatalf("a well-formed response must not raise a diagnostic, got %v", resp.Diagnostics)
		}
		if resp.State.Raw.IsNull() {
			t.Fatal("state was removed; a well-formed response must refresh, not RemoveResource")
		}
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if got.Id.ValueString() != id {
			t.Fatalf("refreshed Id = %q, want %q", got.Id.ValueString(), id)
		}
		if got.Status.ValueString() != "Available" {
			t.Fatalf("refreshed Status = %q, want Available", got.Status.ValueString())
		}
		if got.InstanceId.ValueString() != "i-refreshed" {
			t.Fatalf("refreshed InstanceId = %q, want i-refreshed", got.InstanceId.ValueString())
		}
		if got.ResourceGroupId.ValueString() != "rg-1" {
			t.Fatalf("refreshed ResourceGroupId = %q, want rg-1", got.ResourceGroupId.ValueString())
		}
		if got.VscName.ValueString() != "name-1" {
			t.Fatalf("refreshed VscName = %q, want name-1", got.VscName.ValueString())
		}
		// Read tags round-trip through the same TagKey/TagValue -> key/value
		// decoder the Create/Update paths use.
		gotTags := tagsToMap(context.Background(), &diag.Diagnostics{}, got.Tags)
		if len(gotTags) != 1 || gotTags["env"] != "dev" {
			t.Fatalf("refreshed tags = %v, want {env:dev}", gotTags)
		}
	})

	// noUserTagsDescribe is a well-formed DescribeVscs entry that carries no
	// user tag (and no description/vsc_name): ecsVscState renders all of those
	// as null, the shape that exercises the Read-path empty-value retention.
	noUserTagsDescribe := vscDescribeResponse(map[string]interface{}{
		"VscId":      id,
		"Status":     "Available",
		"InstanceId": "i-refreshed",
	})

	t.Run("tags={} prior preserved through refresh when API has no user tags", func(t *testing.T) {
		// Prior state holds a known empty tags map (tags={}). DescribeVscs returns
		// a well-formed entry with no user tag, so ecsVscState normalizes tags to
		// null. The Read-path empty-value retention must carry the prior empty map
		// forward — else a bare refresh drifts: config tags={} plans +tags={}.
		prior := vscPriorState(id)
		prior.Tags = tagMap(t) // known empty map
		r := vscReadEnv(t, noUserTagsDescribe, nil)
		resp := vscReadCall(t, r, prior)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if got.Tags.IsNull() {
			t.Fatal("tags={} prior was normalized to null on refresh; a known empty map must be carried forward")
		}
		if len(got.Tags.Elements()) != 0 {
			t.Fatalf("tags has %d elements, want 0 (empty map preserved)", len(got.Tags.Elements()))
		}
	})

	t.Run("non-empty tags prior reflects API removal (drift surfaced, not masked)", func(t *testing.T) {
		// Prior state holds a non-empty tags map; the API no longer reports any
		// user tag (external removal). Read must surface that as null, NOT mask it
		// by carrying the stale non-empty prior forward — only known EMPTY values
		// are retained, so real drift stays visible.
		prior := vscPriorState(id)
		prior.Tags = tagMap(t, [2]string{"env", "dev"})
		r := vscReadEnv(t, noUserTagsDescribe, nil)
		resp := vscReadCall(t, r, prior)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if !got.Tags.IsNull() {
			t.Fatalf("non-empty tags prior must reflect the API removal as null (drift surfaced), got %v", got.Tags)
		}
	})

	t.Run("null tags prior stays null when API has no user tags", func(t *testing.T) {
		// Omitted (null) tags must stay null — the empty-value retention only
		// carries forward known EMPTY maps, not null.
		r := vscReadEnv(t, noUserTagsDescribe, nil)
		resp := vscReadCall(t, r, vscPriorState(id)) // null tags
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if !got.Tags.IsNull() {
			t.Fatalf("null tags prior must stay null, got %v", got.Tags)
		}
	})
}

// vscCreateEnv wires a server that handles CreateVsc (returns a fixed VscId)
// and DescribeVscs (returns describe), so the Create method runs end-to-end —
// Create -> CreateVsc(RpcPost) -> getEcsVsc(DescribeVscs) -> ecsVscState ->
// plan-known retention -> resp.State.Set — against canned responses, without
// touching the cloud.
func vscCreateEnv(t *testing.T, describe map[string]interface{}) *ecsVscResource {
	t.Helper()
	return vscTestResource(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Query().Get("Action") {
		case "CreateVsc":
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"VscId":     "vsc-new",
				"RequestId": "r-create",
			})
		case "DescribeVscs":
			_ = json.NewEncoder(w).Encode(describe)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	})
}

// vscCreateCall drives Create against the given plan model and returns the
// response for assertions. The plan schema comes from the resource's own
// Schema, and the plan model is seeded into CreateRequest.Plan the way the
// Framework runtime would before invoking Create.
func vscCreateCall(t *testing.T, r *ecsVscResource, plan ecsVscResourceModel) resource.CreateResponse {
	t.Helper()
	ctx := context.Background()

	schemaResp := resource.SchemaResponse{}
	r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
	schema := schemaResp.Schema

	createReq := resource.CreateRequest{Plan: tfsdk.Plan{Schema: schema}}
	if diags := createReq.Plan.Set(ctx, &plan); diags.HasError() {
		t.Fatalf("seed plan: %v", diags.Errors())
	}

	createResp := resource.CreateResponse{State: tfsdk.State{Schema: schema}}
	r.Create(ctx, createReq, &createResp)
	return createResp
}

// TestUnitEcsVscCreateCloudless is the Create-path regression for the plan-vs-
// applied consistency defect. Create first sets the plan, then overwrites
// state with ecsVscState(Describe), which normalizes a known empty tags map /
// explicit empty-string description / vsc_name to null (MapNull / StringNull
// when Describe carries no user tag or an empty/absent string). That
// divergence trips Framework's "inconsistent result after apply". The fix
// retains plan-known (non-null, non-unknown) Optional values after ecsVscState,
// mirroring Update's plan retention. These subtests assert the applied state
// matches the plan for each empty/omitted/non-empty shape, driving the real
// Create method + RpcPost HTTP against a local httptest.Server.
//
// Before the fix, the tags={} and empty-string subtests fail: state.Tags is
// null (not a known empty map) and state.Description/VscName is null (not an
// empty string). They are the fix-before-pass regression the static review
// flagged: the existing ACC (2->2->1 tags, always non-empty description/vsc_name)
// and unit tests only covered non-empty tag-value fidelity, not the empty
// map / empty string / omitted collection.
func TestUnitEcsVscCreateCloudless(t *testing.T) {
	// DescribeVscs returns a well-formed entry with no user tags and no
	// description/vsc_name — the shape that makes ecsVscState normalize plan
	// values to null before the retention fix.
	noUserValues := vscDescribeResponse(map[string]interface{}{
		"VscId":      "vsc-new",
		"Status":     "Available",
		"InstanceId": "i-1",
	})

	basePlan := func() ecsVscResourceModel {
		// Omitted tags arrive from the real Terraform runtime as a typed null
		// map (MapNull String), not the zero-value Map; the zero value carries
		// no element type and tfsdk rejects it as a DynamicPseudoType conversion
		// when seeding the plan. Mirror vscPriorState's null-tags shape so the
		// omitted-tags / empty-string sub-tests exercise Create, not Plan.Set.
		return ecsVscResourceModel{
			InstanceId: types.StringValue("i-1"),
			Tags:       types.MapNull(types.StringType),
		}
	}

	t.Run("tags={} (known empty map) stays a known empty map, not null", func(t *testing.T) {
		plan := basePlan()
		plan.Tags = tagMap(t) // known empty map (non-nil, zero elements)
		r := vscCreateEnv(t, noUserValues)
		resp := vscCreateCall(t, r, plan)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if got.Tags.IsNull() {
			t.Fatal("tags={} was normalized to null; a known empty map must stay a known empty map (plan-state consistency)")
		}
		if len(got.Tags.Elements()) != 0 {
			t.Fatalf("tags has %d elements, want 0 (empty map)", len(got.Tags.Elements()))
		}
	})

	t.Run("omitted tags (null) stays null", func(t *testing.T) {
		plan := basePlan() // Tags unset = null
		r := vscCreateEnv(t, noUserValues)
		resp := vscCreateCall(t, r, plan)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if !got.Tags.IsNull() {
			t.Fatalf("omitted tags must stay null, got %v", got.Tags)
		}
	})

	t.Run("non-empty tags round-trip", func(t *testing.T) {
		plan := basePlan()
		plan.Tags = tagMap(t, [2]string{"env", "dev"})
		describe := vscDescribeResponse(map[string]interface{}{
			"VscId":      "vsc-new",
			"Status":     "Available",
			"InstanceId": "i-1",
			"Tags": []interface{}{
				map[string]interface{}{"TagKey": "env", "TagValue": "dev"},
			},
		})
		r := vscCreateEnv(t, describe)
		resp := vscCreateCall(t, r, plan)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		gotTags := tagsToMap(context.Background(), &diag.Diagnostics{}, got.Tags)
		if len(gotTags) != 1 || gotTags["env"] != "dev" {
			t.Fatalf("tags = %v, want {env:dev}", gotTags)
		}
	})

	t.Run(`description="" stays an empty string, not null`, func(t *testing.T) {
		plan := basePlan()
		plan.Description = types.StringValue("")
		r := vscCreateEnv(t, noUserValues)
		resp := vscCreateCall(t, r, plan)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if got.Description.IsNull() {
			t.Fatal(`description="" was normalized to null; an explicit empty string must stay an empty string`)
		}
		if got.Description.ValueString() != "" {
			t.Fatalf("description = %q, want empty string", got.Description.ValueString())
		}
	})

	t.Run(`vsc_name="" stays an empty string, not null`, func(t *testing.T) {
		plan := basePlan()
		plan.VscName = types.StringValue("")
		r := vscCreateEnv(t, noUserValues)
		resp := vscCreateCall(t, r, plan)
		ensureNoDiag(t, resp.Diagnostics)
		var got ecsVscResourceModel
		if diags := resp.State.Get(context.Background(), &got); diags.HasError() {
			t.Fatalf("state Get: %v", diags.Errors())
		}
		if got.VscName.IsNull() {
			t.Fatal(`vsc_name="" was normalized to null; an explicit empty string must stay an empty string`)
		}
		if got.VscName.ValueString() != "" {
			t.Fatalf("vsc_name = %q, want empty string", got.VscName.ValueString())
		}
	})
}
