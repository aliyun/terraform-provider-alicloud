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
