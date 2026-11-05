// Node group matching is not reachable through a public entry point: the resource
// exposes no lookup helper and DescribeNodeGroups needs a live instance, so these
// cases live in an internal unit test file.
package starrocks

import (
	"strings"
	"testing"
)

func unitTestBackendNodeGroup(id, zone string) map[string]interface{} {
	return map[string]interface{}{
		"NodeGroupId":              id,
		"ComponentType":            "CN",
		"ZoneId":                   zone,
		"Cu":                       float64(8),
		"DiskNumber":               float64(1),
		"ResidentNodeNumber":       float64(3),
		"SpecType":                 "standard",
		"StoragePerformanceLevel":  "pl1",
		"LocalStorageInstanceType": "",
		"StorageSize":              float64(200),
		"Status":                   "RUNNING",
	}
}

func TestUnitMatchBackendNodeGroup(t *testing.T) {
	config := map[string]interface{}{
		"cu":                          8,
		"disk_number":                 1,
		"resident_node_number":        3,
		"zone_id":                     "cn-hangzhou-i",
		"spec_type":                   "standard",
		"storage_performance_level":   "pl1",
		"local_storage_instance_type": "null",
		"storage_size":                300,
	}

	t.Run("resolves the single matching group", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-unrelated-zone", "cn-hangzhou-h"),
			unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i"),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id != "ng-target" {
			t.Fatalf("expected ng-target, got %q", id)
		}
	})

	t.Run("fails when no group matches", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-other-zone", "cn-hangzhou-h"),
		}
		_, err := matchBackendNodeGroup(groups, config)
		if err == nil || !strings.Contains(err.Error(), "no backend node group matches") {
			t.Fatalf("expected a no-match error, got %v", err)
		}
	})

	t.Run("fails instead of guessing between identical groups", func(t *testing.T) {
		groups := []map[string]interface{}{
			unitTestBackendNodeGroup("ng-1", "cn-hangzhou-i"),
			unitTestBackendNodeGroup("ng-2", "cn-hangzhou-i"),
		}
		id, err := matchBackendNodeGroup(groups, config)
		if err == nil {
			t.Fatalf("expected an ambiguity error, got id %q", id)
		}
		if !strings.Contains(err.Error(), "refusing to guess") {
			t.Fatalf("expected an ambiguity error, got %v", err)
		}
	})
}

func TestUnitBackendNodeGroupAttributesMatch(t *testing.T) {
	config := map[string]interface{}{
		"cu":                          8,
		"disk_number":                 1,
		"resident_node_number":        3,
		"zone_id":                     "cn-hangzhou-i",
		"spec_type":                   "standard",
		"storage_performance_level":   "pl1",
		"local_storage_instance_type": "null",
		"storage_size":                300,
	}

	group := unitTestBackendNodeGroup("ng-target", "cn-hangzhou-i")
	if !backendNodeGroupAttributesMatch(group, config) {
		t.Fatal("expected the group to match the configured attributes")
	}

	group["Cu"] = float64(16)
	if backendNodeGroupAttributesMatch(group, config) {
		t.Fatal("expected a different cu to make the group not match")
	}
	group["Cu"] = float64(8)

	group["StorageSize"] = float64(400)
	if !backendNodeGroupAttributesMatch(group, config) {
		t.Fatal("expected storage_size to be ignored when matching")
	}
}

func TestUnitBackendNodeGroupStorageSizeChangeForcesReplacement(t *testing.T) {
	cases := []struct {
		name  string
		old   int
		oldOk bool
		new   int
		newOk bool
		want  bool
	}{
		{name: "growth is in place", old: 200, oldOk: true, new: 300, newOk: true, want: false},
		{name: "same size is in place", old: 200, oldOk: true, new: 200, newOk: true, want: false},
		{name: "shrink replaces", old: 300, oldOk: true, new: 200, newOk: true, want: true},
		{name: "removing the value replaces", old: 300, oldOk: true, new: 0, newOk: false, want: true},
		{name: "no previous value never replaces", old: 0, oldOk: false, new: 300, newOk: true, want: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := backendNodeGroupStorageSizeChangeForcesReplacement(tc.old, tc.oldOk, tc.new, tc.newOk)
			if got != tc.want {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
