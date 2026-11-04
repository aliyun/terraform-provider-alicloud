package rds

import (
	"testing"
)

// TestUnitDBInstanceVswitchIdsOrderDiffSuppressFunc is the bounded-suppression
// unit net for `vswitch_ids`. RDS binds vswitch order to the zone order
// (ZoneId/ZoneIdSlave1/ZoneIdSlave2), so an order-only change is not an
// independently applicable config change and is suppressed; a membership change
// (add/remove/replace a vswitch) is not. The cases are exercised against the
// suppressor's nil-data path with the member lists serialized as comma-joined
// strings, the form the collection-order checker reads.
func TestUnitDBInstanceVswitchIdsOrderDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"ReorderOnly", "vsw-a,vsw-b", "vsw-b,vsw-a", true},
		{"ReorderDuplicates", "vsw-a,vsw-a,vsw-b", "vsw-b,vsw-a,vsw-a", true},
		{"MemberAdded", "vsw-a,vsw-b", "vsw-a,vsw-b,vsw-c", false},
		{"MemberRemoved", "vsw-a,vsw-b,vsw-c", "vsw-a,vsw-b", false},
		{"MemberReplaced", "vsw-a", "vsw-b", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := dbInstanceVswitchIdsOrderDiffSuppressFunc("vswitch_ids.0", tc.old, tc.new, nil); got != tc.expected {
				t.Fatalf("case %s: suppressor = %v, want %v", tc.name, got, tc.expected)
			}
		})
	}
}

// TestUnitDBInstanceVswitchIdDiffSuppressFunc is the regression net for the
// `vswitch_id` DiffSuppressFunc (F-1). `vswitch_ids` is Optional+Computed and
// Read always backfills the API value into state; the prior suppressor used
// d.GetOk("vswitch_ids"), which reads that Computed state value, so a
// single-vswitch `vswitch_id` change (vsw-A→vsw-B) was silently suppressed and
// ForceNew never fired, drifting config from the instance. The decision core
// is exercised directly with explicit inputs so the regression is pinned
// without a live ResourceData; vswitchIdsConfigured must reflect user config,
// never the Computed state value.
func TestUnitDBInstanceVswitchIdDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                 string
		isNew                bool
		vswitchIdsConfigured bool
		newVswitchId         string
		expected             bool
	}{
		{"NewResourceNotSuppressed", true, false, "vsw-b", false},
		// F-1 regression: user only configured vswitch_id (vswitch_ids absent
		// from config); a single-vswitch change must NOT be suppressed so
		// ForceNew fires.
		{"SingleVswitchChangeNotSuppressed", false, false, "vsw-b", false},
		{"MultiZoneVswitchSuppressed", false, false, "vsw-a,vsw-b", true},
		// When vswitch_ids is actually in user config, suppress vswitch_id drift.
		{"VswitchIdsConfiguredSuppressesMulti", false, true, "vsw-a,vsw-b", true},
		{"VswitchIdsConfiguredSuppressesSingle", false, true, "vsw-a", true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := dbInstanceVswitchIdDecide(tc.isNew, tc.vswitchIdsConfigured, tc.newVswitchId); got != tc.expected {
				t.Fatalf("case %s: decide = %v, want %v", tc.name, got, tc.expected)
			}
		})
	}
}
