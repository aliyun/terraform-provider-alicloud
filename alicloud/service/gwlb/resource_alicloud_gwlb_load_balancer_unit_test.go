package gwlb

import (
	"strings"
	"testing"
)

// TestUnitAliCloudGwlbLoadBalancerUpdateWaitTarget covers the regression that N-01 introduced:
// an attribute update that does not touch traffic_mode must keep the status-only wait. State
// written by provider versions before 2.3.0 has no traffic_mode value, so a combined target
// would become "Active/" and never converge.
func TestUnitAliCloudGwlbLoadBalancerUpdateWaitTarget(t *testing.T) {
	cases := []struct {
		name               string
		trafficModeChanged bool
		trafficMode        interface{}
		wantTarget         string
		wantCombined       bool
	}{
		{
			name:               "legacy state without traffic mode only renames",
			trafficModeChanged: false,
			trafficMode:        "",
			wantTarget:         "Active",
			wantCombined:       false,
		},
		{
			name:               "traffic mode change targets the new value",
			trafficModeChanged: true,
			trafficMode:        "ByPass",
			wantTarget:         "Active/ByPass",
			wantCombined:       true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			target, combined := loadBalancerUpdateWaitTarget(c.trafficModeChanged, c.trafficMode)
			if target != c.wantTarget {
				t.Fatalf("unexpected wait target: got %q, want %q", target, c.wantTarget)
			}
			if combined != c.wantCombined {
				t.Fatalf("unexpected combined-wait flag: got %v, want %v", combined, c.wantCombined)
			}
			if !combined && strings.Contains(target, "/") {
				t.Fatalf("status-only wait must not carry a traffic mode, got %q", target)
			}
		})
	}
}
