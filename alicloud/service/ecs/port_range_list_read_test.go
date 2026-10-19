package ecs

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/aliyun/credentials-go/credentials"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// portRangeListTestClient builds an *connectivity.AliyunClient whose ECS RPC
// calls are served by a local httptest server returning caller-crafted JSON,
// so the Read-side response-shape handling can be exercised without a real
// cloud account. It mirrors the cloud-firewall unit-test harness: fake
// credentials, the endpoint stored under the lowercase product key rpcRequest
// looks up ("ecs"), and NO_PROXY pinned to the server host so the SDK HTTP
// client never routes the loopback call through a corporate proxy.
func portRangeListTestClient(t *testing.T, describeBody string) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := `{}`
		switch req.FormValue("Action") {
		case "DescribePortRangeLists":
			body = describeBody
		case "DescribePortRangeListEntries":
			// A confirmed-existing list with no entries is a legitimate empty
			// set; the happy-path Read maps it to a null entry set.
			body = `{"Entries":[]}`
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("create credential: %v", err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	endpoints.Store("ecs", endpoint)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("init client: %v", err)
	}
	return client
}

// TestUnitPortRangeListReadShapeHandling asserts the contract of
// describePortRangeList for every DescribePortRangeLists response shape:
// only a structurally-correct empty array drops the resource (nil, nil);
// a missing root, a nil root, a wrong-typed field, a non-object element, or a
// non-empty list that carries only sibling ids all surface as an error so Read
// keeps the prior id in state instead of silently clearing it.
func TestUnitPortRangeListReadShapeHandling(t *testing.T) {
	const requestedID = "pl-abc"
	cases := []struct {
		name       string
		body       string
		wantErr    bool
		wantObject bool
	}{
		{
			name:    "missing_root_field",
			body:    `{"RequestId":"x"}`,
			wantErr: true,
		},
		{
			name:    "nil_root_field",
			body:    `{"PortRangeLists":null,"RequestId":"x"}`,
			wantErr: true,
		},
		{
			name:    "empty_list_is_gone",
			body:    `{"PortRangeLists":[],"RequestId":"x"}`,
			wantErr: false,
		},
		{
			name:    "non_empty_only_other_id",
			body:    `{"PortRangeLists":[{"PortRangeListId":"pl-other","MaxEntries":10}],"RequestId":"x"}`,
			wantErr: true,
		},
		{
			name:       "matching_id",
			body:       `{"PortRangeLists":[{"PortRangeListId":"pl-abc","MaxEntries":10,"PortRangeListName":"n"}],"RequestId":"x"}`,
			wantErr:    false,
			wantObject: true,
		},
		{
			name:    "malformed_field_type",
			body:    `{"PortRangeLists":"not-a-list","RequestId":"x"}`,
			wantErr: true,
		},
		{
			name:    "non_object_element",
			body:    `{"PortRangeLists":["not-an-object"],"RequestId":"x"}`,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := portRangeListTestClient(t, tc.body)
			object, err := describePortRangeList(client, requestedID)
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error so Read keeps the id, got object=%v err=%v", object, err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got err=%v", err)
			}
			if tc.wantObject && object == nil {
				t.Fatalf("expected the matching PortRangeList object, got nil")
			}
			if !tc.wantObject && object != nil {
				t.Fatalf("expected nil object, got %v", object)
			}
			if tc.wantObject && fmt.Sprint(object["PortRangeListId"]) != requestedID {
				t.Fatalf("returned object id = %v, want %s", object["PortRangeListId"], requestedID)
			}
		})
	}
}

// TestUnitPortRangeListBuildStateKeepsIDOnMissingRoot and ...DropsOnEmptyList
// verify the Read entry point (buildPortRangeListState) end-to-end: a missing
// PortRangeLists root surfaces as error diagnostics with a nil model (Read
// returns the error and keeps the prior id), while a structurally-correct
// empty array returns a nil model with no diagnostics (Read removes the
// resource from state).
func TestUnitPortRangeListBuildStateKeepsIDOnMissingRoot(t *testing.T) {
	client := portRangeListTestClient(t, `{"RequestId":"x"}`)
	state, diags := buildPortRangeListState(context.Background(), client, "pl-abc", types.MapNull(types.StringType))
	if state != nil {
		t.Fatalf("expected nil state, got %+v", state)
	}
	if !diags.HasError() {
		t.Fatalf("expected error diagnostics so Read keeps the prior id, got %v", diags)
	}
}

func TestUnitPortRangeListBuildStateDropsOnEmptyList(t *testing.T) {
	client := portRangeListTestClient(t, `{"PortRangeLists":[],"RequestId":"x"}`)
	state, diags := buildPortRangeListState(context.Background(), client, "pl-abc", types.MapNull(types.StringType))
	if state != nil {
		t.Fatalf("expected nil state for a gone list, got %+v", state)
	}
	if diags.HasError() {
		t.Fatalf("expected no diagnostics for a clean empty array, got %v", diags)
	}
}

// TestUnitPortRangeListStateTagsNullVsEmpty reproduces the D1 consistency defect
// at the readback that Create/Update/Read all overwrite state with: when the API
// legitimately reports no tags (Tag is an optional 0~20-length array, so tags={}
// is a valid config) the old code unconditionally wrote types.MapNull. With an
// explicit tags = {} plan that made the known empty map drift to null. The fix
// preserves the configured null-vs-explicit-empty distinction instead of
// swapping it, and still refreshes real API tags so an out-of-band change is not
// masked. These are no-cloud cases driven through the real describe call chain.
func TestUnitPortRangeListStateTagsNullVsEmpty(t *testing.T) {
	const id = "pl-abc"
	const matchingNoTags = `{"PortRangeLists":[{"PortRangeListId":"pl-abc","MaxEntries":10,"PortRangeListName":"n"}],"RequestId":"x"}`
	const matchingWithTags = `{"PortRangeLists":[{"PortRangeListId":"pl-abc","MaxEntries":10,"PortRangeListName":"n","Tags":[{"TagKey":"env","TagValue":"prod"}]}],"RequestId":"x"}`

	// An explicit empty map (tags = {}) is a known, non-null value with zero
	// elements, distinct from types.MapNull (tags omitted).
	emptyTags, d := types.MapValueFrom(context.Background(), types.StringType, map[string]string{})
	if d.HasError() {
		t.Fatalf("build explicit empty map: %v", d)
	}

	cases := []struct {
		name      string
		body      string
		priorTags types.Map
		wantNull  bool
		wantLen   int
	}{
		{
			name:      "api_empty_config_omitted_keeps_null",
			body:      matchingNoTags,
			priorTags: types.MapNull(types.StringType),
			wantNull:  true,
		},
		{
			// D1: the old code returned null here, drifting an explicit tags = {}.
			name:      "api_empty_config_explicit_empty_keeps_empty_not_null",
			body:      matchingNoTags,
			priorTags: emptyTags,
			wantNull:  false,
			wantLen:   0,
		},
		{
			name:      "api_nonempty_refreshes_actual_values",
			body:      matchingWithTags,
			priorTags: types.MapNull(types.StringType),
			wantNull:  false,
			wantLen:   1,
		},
		{
			// Real API tags override an explicit empty plan; the prior value is
			// not copied back, so drift is surfaced, not masked.
			name:      "api_nonempty_overrides_explicit_empty",
			body:      matchingWithTags,
			priorTags: emptyTags,
			wantNull:  false,
			wantLen:   1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := portRangeListTestClient(t, tc.body)
			state, diags := buildPortRangeListState(context.Background(), client, id, tc.priorTags)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if state == nil {
				t.Fatalf("expected a state model, got nil")
			}
			if tc.wantNull && !state.Tags.IsNull() {
				t.Fatalf("expected null tags, got %v", state.Tags)
			}
			if !tc.wantNull && state.Tags.IsNull() {
				t.Fatalf("expected a non-null tags value, got null")
			}
			if !tc.wantNull && len(state.Tags.Elements()) != tc.wantLen {
				t.Fatalf("expected %d tag elements, got %d", tc.wantLen, len(state.Tags.Elements()))
			}
		})
	}
}

// TestUnitPortRangeListDestroyDecision is the no-cloud positive/negative pair for
// the acceptance CheckDestroy's actual call chain (describePortRangeList): a
// structurally-correct empty array is the only "confirmed gone" signal (nil,
// nil); a still-present object fails; anything malformed fails closed.
func TestUnitPortRangeListDestroyDecision(t *testing.T) {
	const id = "pl-abc"
	cases := []struct {
		name     string
		body     string
		wantGone bool
		wantErr  bool
	}{
		{name: "empty_array_confirmed_gone", body: `{"PortRangeLists":[],"RequestId":"x"}`, wantGone: true},
		{name: "still_exists_fails", body: `{"PortRangeLists":[{"PortRangeListId":"pl-abc","MaxEntries":10}],"RequestId":"x"}`, wantGone: false},
		{name: "missing_root_fails_closed", body: `{"RequestId":"x"}`, wantGone: false, wantErr: true},
		{name: "malformed_field_fails_closed", body: `{"PortRangeLists":"not-a-list","RequestId":"x"}`, wantGone: false, wantErr: true},
		{name: "non_object_element_fails_closed", body: `{"PortRangeLists":["x"],"RequestId":"x"}`, wantGone: false, wantErr: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := portRangeListTestClient(t, tc.body)
			object, err := describePortRangeList(client, id)
			if (object == nil && err == nil) != tc.wantGone {
				t.Fatalf("destroy decision: gone=%v wantGone=%v (object=%v err=%v)", object == nil && err == nil, tc.wantGone, object, err)
			}
			if tc.wantErr && err == nil {
				t.Fatalf("expected a fail-closed error, got nil")
			}
		})
	}
}

// TestUnitPortRangeListEntriesShapeHandling keeps the entries-side type guards
// regression-protected: an absent/nil/empty Entries field is the legitimate
// empty set (no entries, no error), while a wrong-typed Entries field or a
// non-object Entry element is an error so Read never silently replaces a
// populated entry set with an empty one.
func TestUnitPortRangeListEntriesShapeHandling(t *testing.T) {
	const requestedID = "pl-abc"
	cases := []struct {
		name    string
		body    string
		wantErr bool
		wantLen int
	}{
		{
			name:    "absent_entries_is_empty",
			body:    `{"RequestId":"x"}`,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "nil_entries_is_empty",
			body:    `{"Entries":null,"RequestId":"x"}`,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "empty_entries_array",
			body:    `{"Entries":[],"RequestId":"x"}`,
			wantErr: false,
			wantLen: 0,
		},
		{
			name:    "malformed_entries_field",
			body:    `{"Entries":"not-a-list","RequestId":"x"}`,
			wantErr: true,
		},
		{
			name:    "non_object_entry_element",
			body:    `{"Entries":["not-an-object"],"RequestId":"x"}`,
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := portRangeListEntriesTestClient(t, tc.body)
			entries, err := describePortRangeListEntries(client, requestedID)
			if tc.wantErr && err == nil {
				t.Fatalf("expected an error, got entries=%v err=%v", entries, err)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got err=%v", err)
			}
			if !tc.wantErr && len(entries) != tc.wantLen {
				t.Fatalf("entries len = %d, want %d", len(entries), tc.wantLen)
			}
		})
	}
}

// portRangeListEntriesTestClient serves DescribePortRangeListEntries with a
// caller-crafted body so describePortRangeListEntries' own response-shape guards
// are exercised in isolation (DescribePortRangeLists is not reached).
func portRangeListEntriesTestClient(t *testing.T, entriesBody string) *connectivity.AliyunClient {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := `{}`
		if req.FormValue("Action") == "DescribePortRangeListEntries" {
			body = entriesBody
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId("test-key").
		SetAccessKeySecret("test-secret"))
	if err != nil {
		t.Fatalf("create credential: %v", err)
	}
	endpoint := strings.TrimPrefix(server.URL, "http://")
	t.Setenv("NO_PROXY", endpoint)
	endpoints := new(sync.Map)
	endpoints.Store("ecs", endpoint)
	config := &connectivity.Config{
		AccessKey:            "test-key",
		SecretKey:            "test-secret",
		Credential:           credential,
		RegionId:             "cn-hangzhou",
		AccountType:          "test",
		Protocol:             "http",
		Endpoints:            endpoints,
		SignVersion:          new(sync.Map),
		SkipRegionValidation: true,
	}
	client, err := config.Client()
	if err != nil {
		t.Fatalf("init client: %v", err)
	}
	return client
}
