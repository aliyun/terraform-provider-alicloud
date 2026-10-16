// Package ecs contains the Framework-native resources of the ECS product.
package ecs

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	tfsdk "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// root_storage_review_test.go mirrors the root independent review gates on the
// production bytes: it drives the real resource methods over a real httptest
// RPC path (tea-rpc -> RpcPost -> json UseNumber), not a constructed model or
// a static self-assertion. It covers the gates the root reported as real
// FAIL/exit1 before the fixes and the lifecycle methods not exercised by the
// main unit test: the json.Number max_partition_number mapping, an empty-array
// tags response mapping to a KNOWN (not null) map, and the Delete lifecycle
// (success + not-found idempotency).

// TestUnitRootStorageReadActualJSONNumber verifies that a real HTTP
// DescribeStorageSets response carrying StorageSetPartitionNumber as a JSON
// number (decoded to json.Number by the SDK's UseNumber, not float64) is
// mapped to the correct int64, not dropped to null. Before the fix this read
// returned Int64Null with empty diagnostics.
func TestUnitRootStorageReadActualJSONNumber(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	stateModel := nullStorageSetModel()
	stateModel.Id = types.StringValue(testStorageSetID)
	stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
	var resp resource.ReadResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Read(context.Background(), resource.ReadRequest{State: storageSetStateFromModel(t, sch, stateModel)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", resp.Diagnostics)
	}
	var model storageSetResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Get: %v", diags)
	}
	if model.MaxPartitionNumber.ValueInt64() != 3 {
		t.Fatalf("expected max_partition_number=3 (json.Number path), got %d (null=%v)", model.MaxPartitionNumber.ValueInt64(), model.MaxPartitionNumber.IsNull())
	}
}

// TestUnitRootStorageReadEmptyTagsIsKnownMap verifies that a real HTTP
// DescribeStorageSets with Tags.Tag=[] maps to an empty KNOWN map, not null.
// A null state value would differ from a config that sets tags={} and force a
// perpetual non-no-op plan (the root ProtocolEmptyTagsNoop FAIL).
func TestUnitRootStorageReadEmptyTagsIsKnownMap(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	stateModel := nullStorageSetModel()
	stateModel.Id = types.StringValue(testStorageSetID)
	stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
	var resp resource.ReadResponse
	resp.State = tfsdk.State{Schema: sch}
	r.Read(context.Background(), resource.ReadRequest{State: storageSetStateFromModel(t, sch, stateModel)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected read diagnostics: %v", resp.Diagnostics)
	}
	var model storageSetResourceModel
	if diags := resp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("state.Get: %v", diags)
	}
	if model.Tags.IsNull() {
		t.Fatal("expected empty KNOWN tags map (matches tags={}), got null — would force a perpetual non-no-op plan")
	}
	if model.Tags.IsUnknown() {
		t.Fatal("expected empty KNOWN tags map, got unknown")
	}
	if len(model.Tags.Elements()) != 0 {
		t.Fatalf("expected empty tags map, got %d elements: %v", len(model.Tags.Elements()), model.Tags.Elements())
	}
}

// TestUnitRootStorageDelete_CallsDeleteStorageSet verifies Delete issues
// DeleteStorageSet with the resource id and succeeds.
func TestUnitRootStorageDelete_CallsDeleteStorageSet(t *testing.T) {
	srv, sent := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DeleteStorageSet" {
			return http.StatusOK, `{"RequestId":"r"}`
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	stateModel := nullStorageSetModel()
	stateModel.Id = types.StringValue(testStorageSetID)
	stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: storageSetStateFromModel(t, sch, stateModel)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected delete diagnostics: %v", resp.Diagnostics)
	}
	loadVal, ok := sent.Load("DeleteStorageSet")
	if !ok {
		t.Fatal("expected DeleteStorageSet to be called")
	}
	params, _ := loadVal.(map[string]string)
	if params["StorageSetId"] != testStorageSetID {
		t.Fatalf("expected StorageSetId=%q, got %q", testStorageSetID, params["StorageSetId"])
	}
}

// TestUnitRootStorageDelete_NotFoundIsNoOp verifies Delete on an already-removed
// storage set (InvalidStorageSetId.NotFound) is idempotent: no error.
func TestUnitRootStorageDelete_NotFoundIsNoOp(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DeleteStorageSet" {
			return http.StatusBadRequest, `{"Code":"InvalidStorageSetId.NotFound","Message":"not found","RequestId":"r"}`
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	r, sch := newStorageSetResource(t, client)
	stateModel := nullStorageSetModel()
	stateModel.Id = types.StringValue(testStorageSetID)
	stateModel.ZoneId = types.StringValue("cn-hangzhou-i")
	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: storageSetStateFromModel(t, sch, stateModel)}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected not-found delete to be an idempotent no-op, got diagnostics: %v", resp.Diagnostics)
	}
}
