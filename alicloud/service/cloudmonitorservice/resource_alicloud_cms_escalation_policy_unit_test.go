// Package cloudmonitorservice. Internal TestUnit for the
// suppressCmsEscalationEnabledSubChannelsDiff DiffSuppressFunc, kept in the
// internal test package so it can call the unexported suppressor directly.
package cloudmonitorservice

import (
	"testing"
)

// TestUnitCmsEscalationEnabledSubChannelsDiff verifies the DiffSuppressFunc
// for the enabled_sub_channels TypeSet count key. It must suppress the
// spurious empty/zero count-key drift surfaced by IDRefreshName when the API
// returns a nil/empty set, while never masking real element changes.
func TestUnitCmsEscalationEnabledSubChannelsDiff(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		old      string
		new      string
		expected bool
	}{
		{"EmptyToZero", "enabled_sub_channels.#", "", "0", true},
		{"ZeroToEmpty", "enabled_sub_channels.#", "0", "", true},
		{"ZeroToZero", "enabled_sub_channels.#", "0", "0", false},
		{"EmptyToEmpty", "enabled_sub_channels.#", "", "", false},
		{"NonCountKeyUnchanged", "enabled_sub_channels.0", "ding", "ding", false},
		{"RealElementChange", "enabled_sub_channels.0", "ding", "mail", false},
		{"CountToOne", "enabled_sub_channels.#", "0", "1", false},
		{"CountToTwo", "enabled_sub_channels.#", "1", "2", false},
		{"NestedCountKey", "escalation_stage_list.0.notify_channels.0.enabled_sub_channels.#", "", "0", true},
		{"NestedCountKeyZeroToEmpty", "escalation_stage_list.0.notify_channels.0.enabled_sub_channels.#", "0", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := suppressCmsEscalationEnabledSubChannelsDiff(tc.key, tc.old, tc.new, nil)
			if result != tc.expected {
				t.Errorf("key=%q old=%q new=%q: expected %v got %v", tc.key, tc.old, tc.new, tc.expected, result)
			}
		})
	}
}
