// Package ecs contains the Framework-native resources of the ECS product.
package ecs

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These are no-cloud (httptest) Framework tests for the alicloud_ecs_storage_set
// destroy-check CALLER: testAccCheckEcsStorageSetDestroy delegates to
// CheckEcsStorageSetDestroy, so this exercises the real destroy-check code path
// (not a StorageSetExists helper self-assertion) over the caller gates the root
// reported as a real gap: the check must be target-first and must not treat
// "no usable client" as "destroyed" when a resource is still in state.
//
// The acceptance-test CheckDestroy reads acctest.Provider.Meta(); here the
// client is injected directly from the same httptest harness as the other
// no-cloud unit tests (storageSetTestServer/storageSetTestClient), so the
// DescribeStorageSets responses land on the in-process RPC server with no
// cloud credentials.

// storageSetDestroyState builds a *terraform.State whose root module holds the
// given alicloud_ecs_storage_set instance ids, mirroring the shape the SDK
// CheckDestroy receives after an acceptance test: root-module resources keyed
// by address, each carrying its real cloud id in Primary.ID.
func storageSetDestroyState(ids ...string) *terraform.State {
	resources := map[string]*terraform.ResourceState{}
	for i, id := range ids {
		resources[fmt.Sprintf("alicloud_ecs_storage_set.default%d", i)] = &terraform.ResourceState{
			Type:    "alicloud_ecs_storage_set",
			Primary: &terraform.InstanceState{ID: id},
		}
	}
	return &terraform.State{
		Modules: []*terraform.ModuleState{
			{Path: []string{"root"}, Resources: resources},
		},
	}
}

// CheckEcsStorageSetDestroy is the target-first destroy check shared by the
// acceptance-test CheckDestroy (which reads the provider meta) and these
// no-cloud unit tests (which inject an httptest-backed *connectivity.AliyunClient).
//
// It collects the alicloud_ecs_storage_set instances still in state first, then:
//   - no targets => no-op (the test never reached apply, nothing to verify);
//   - targets remain but meta is nil or not a *connectivity.AliyunClient => an
//     explicit error: "cannot check" must never be reported as "destroyed";
//   - targets remain with a usable client => reuse the strict StorageSetExists
//     not-found logic so a Describe error or a still-existing resource surfaces
//     instead of masking a leaked resource as absent.
//
// This is a test-only seam (it lives in a _test.go file) so the caller guard is
// unit-testable without the SDK provider Configure machinery; the production
// behavior of the acceptance-test CheckDestroy is unchanged.
func CheckEcsStorageSetDestroy(meta interface{}, s *terraform.State) error {
	if s == nil {
		return nil
	}
	var ids []string
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_storage_set" {
			continue
		}
		if rs.Primary == nil || rs.Primary.ID == "" {
			continue
		}
		ids = append(ids, rs.Primary.ID)
	}
	if len(ids) == 0 {
		// Nothing in state to verify: the test never reached apply (e.g. the
		// account lacks Storage Set capability, 403 at create), so the check is
		// a no-op. This is the ONLY legitimate no-op path.
		return nil
	}
	// Targets remain in state: a usable client is mandatory. The previous code
	// returned nil when meta was nil or the wrong type, so a resource left in
	// state was silently treated as destroyed — the gap this guard closes.
	if meta == nil {
		return fmt.Errorf("alicloud_ecs_storage_set destroy check cannot run: provider meta is nil but %d resource(s) remain in state; 'cannot check' must not be treated as 'destroyed'", len(ids))
	}
	client, ok := meta.(*connectivity.AliyunClient)
	if !ok || client == nil {
		return fmt.Errorf("alicloud_ecs_storage_set destroy check cannot run: provider meta is %T (expected *connectivity.AliyunClient) but %d resource(s) remain in state; 'cannot check' must not be treated as 'destroyed'", meta, len(ids))
	}
	for _, id := range ids {
		exists, err := StorageSetExists(client, id)
		if err != nil {
			return fmt.Errorf("checking alicloud_ecs_storage_set %q after destroy: %w", id, err)
		}
		if exists {
			return fmt.Errorf("alicloud_ecs_storage_set %q still exists after destroy", id)
		}
	}
	return nil
}

// TestUnitStorageSetCheckDestroy_NoTargetsIsNoop verifies the no-op path:
// empty state and a state with only non-storage_set resources have no targets,
// so the check returns nil regardless of whether a client exists. This is the
// only legitimate no-op: there must be nothing to verify, not "cannot check".
func TestUnitStorageSetCheckDestroy_NoTargetsIsNoop(t *testing.T) {
	emptyState := &terraform.State{
		Modules: []*terraform.ModuleState{{Path: []string{"root"}}},
	}
	if err := CheckEcsStorageSetDestroy(nil, emptyState); err != nil {
		t.Fatalf("empty state with nil meta: expected nil (no targets), got %v", err)
	}
	nonMatching := &terraform.State{
		Modules: []*terraform.ModuleState{{
			Path: []string{"root"},
			Resources: map[string]*terraform.ResourceState{
				"alicloud_db_instance.default": {Type: "alicloud_db_instance", Primary: &terraform.InstanceState{ID: "db-1"}},
			},
		}},
	}
	if err := CheckEcsStorageSetDestroy(nil, nonMatching); err != nil {
		t.Fatalf("non-storage_set resource state with nil meta: expected nil, got %v", err)
	}
}

// TestUnitStorageSetCheckDestroy_TargetGoneIsNoop verifies that a target whose
// DescribeStorageSets legitimately reports not-found (TotalCount=0 corroborates
// an empty result) returns nil: the resource was destroyed.
func TestUnitStorageSetCheckDestroy_TargetGoneIsNoop(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, `{"TotalCount":0,"RequestId":"r"}`
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	if err := CheckEcsStorageSetDestroy(client, storageSetDestroyState(testStorageSetID)); err != nil {
		t.Fatalf("target gone (TotalCount=0 not-found): expected nil, got %v", err)
	}
}

// TestUnitStorageSetCheckDestroy_TargetStillExistsErrors verifies that a target
// DescribeStorageSets still returns surfaces as an error: the resource leaked
// and was not destroyed.
func TestUnitStorageSetCheckDestroy_TargetStillExistsErrors(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, describeStorageSetObjectBody(testStorageSetID, "d", "n", "cn-hangzhou-i", "cn-hangzhou", "rg-1", 3, `{"Tag":[]}`)
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	if err := CheckEcsStorageSetDestroy(client, storageSetDestroyState(testStorageSetID)); err == nil {
		t.Fatal("target still exists: expected error, got nil (a leaked resource would be masked as destroyed)")
	}
}

// TestUnitStorageSetCheckDestroy_DescribeErrorErrors verifies that a Describe
// inconsistency (positive TotalCount without the StorageSets container — the
// API claims matches but returned none) surfaces as an error rather than being
// swallowed into a false not-found: state must be preserved, not dropped.
func TestUnitStorageSetCheckDestroy_DescribeErrorErrors(t *testing.T) {
	srv, _ := storageSetTestServer(t, func(action string, params map[string]string) (int, string) {
		if action == "DescribeStorageSets" {
			return http.StatusOK, `{"TotalCount":1,"RequestId":"r"}`
		}
		return http.StatusOK, `{}`
	})
	client := storageSetTestClient(t, srv)
	if err := CheckEcsStorageSetDestroy(client, storageSetDestroyState(testStorageSetID)); err == nil {
		t.Fatal("describe inconsistency (positive TotalCount, no container): expected error, got nil (state would be dropped on inconsistent data)")
	}
}

// TestUnitStorageSetCheckDestroy_NilClientWithTargetErrors verifies the closed
// gap: when a target remains in state but the provider meta is nil, the check
// must error — NOT return nil. Returning nil here is what previously treated
// "cannot check" as "destroyed".
func TestUnitStorageSetCheckDestroy_NilClientWithTargetErrors(t *testing.T) {
	if err := CheckEcsStorageSetDestroy(nil, storageSetDestroyState(testStorageSetID)); err == nil {
		t.Fatal("nil meta with target present: expected error, got nil (the gap: 'cannot check' reported as 'destroyed')")
	}
}

// TestUnitStorageSetCheckDestroy_BadTypeClientWithTargetErrors verifies the
// same guard for a wrong-type meta: a target remains in state, so the check
// must surface an error instead of a silent nil pass.
func TestUnitStorageSetCheckDestroy_BadTypeClientWithTargetErrors(t *testing.T) {
	if err := CheckEcsStorageSetDestroy("not-an-aliyun-client", storageSetDestroyState(testStorageSetID)); err == nil {
		t.Fatal("wrong-type meta with target present: expected error, got nil (the gap: 'cannot check' reported as 'destroyed')")
	}
}
