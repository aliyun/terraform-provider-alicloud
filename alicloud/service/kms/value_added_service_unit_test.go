package kms

import (
	"testing"
	"time"
)

func TestUnitKmsValueAddedServiceInEffect(t *testing.T) {
	now := time.Date(2026, 8, 4, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		name     string
		instance map[string]interface{}
		expected bool
	}{
		{
			name:     "effective with a future end",
			instance: map[string]interface{}{"InstanceID": "dkr-1", "Status": "Normal", "EndTime": "2027-08-04T16:00:00Z"},
			expected: true,
		},
		{
			name:     "end exactly one second away still counts",
			instance: map[string]interface{}{"InstanceID": "dkr-2", "Status": "Normal", "EndTime": "2026-08-04T12:00:01Z"},
			expected: true,
		},
		{
			name:     "end exactly at now has lapsed",
			instance: map[string]interface{}{"InstanceID": "dkr-8", "Status": "Normal", "EndTime": "2026-08-04T12:00:00Z"},
			expected: false,
		},
		{
			name:     "expired subscription still reported as Normal",
			instance: map[string]interface{}{"InstanceID": "dkr-3", "Status": "Normal", "EndTime": "2025-08-04T16:00:00Z"},
			expected: false,
		},
		{
			name:     "refunded order lingering in Creating",
			instance: map[string]interface{}{"InstanceID": "dkr-4", "Status": "Creating", "EndTime": "2999-09-08T16:00:00Z"},
			expected: false,
		},
		{
			name:     "no EndTime reported",
			instance: map[string]interface{}{"InstanceID": "dkr-5", "Status": "Normal"},
			expected: false,
		},
		{
			name:     "nil EndTime",
			instance: map[string]interface{}{"InstanceID": "dkr-6", "Status": "Normal", "EndTime": nil},
			expected: false,
		},
		{
			name:     "unparseable EndTime",
			instance: map[string]interface{}{"InstanceID": "dkr-7", "Status": "Normal", "EndTime": "not a timestamp"},
			expected: false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := kmsValueAddedServiceInEffect(c.instance, now); got != c.expected {
				t.Errorf("kmsValueAddedServiceInEffect(%v) = %v, want %v", c.instance, got, c.expected)
			}
		})
	}
}
