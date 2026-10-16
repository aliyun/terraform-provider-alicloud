package helper

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Tests for the state helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonComputePeriodByUnit(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		createTime  string
		endTime     string
		unit        string
		expect      int
		expectError bool
	}{
		{
			"cal month",
			now.Format(time.RFC3339),
			now.AddDate(0, 2, 0).Format(time.RFC3339),
			"Month",
			2,
			false,
		},
		{
			"cal week",
			now.Format(time.RFC3339),
			now.AddDate(0, 0, 14).Format(time.RFC3339),
			"Week",
			2,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			period, err := ComputePeriodByUnit(
				tt.createTime,
				tt.endTime,
				0,
				tt.unit,
			)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.InDelta(t, tt.expect, period, 1)
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnitCommonComputePeriodByUnitEdgeCases(t *testing.T) {
	now := time.Now()

	// Test with Year unit
	result, err := ComputePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(1, 0, 0).Format(time.RFC3339),
		0,
		"Year",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1 year, got %d", result)
	}

	// Test with nil createTime
	_, err = ComputePeriodByUnit(nil, now.Format(time.RFC3339), 0, "Month")
	if err == nil {
		t.Error("Expected error for nil createTime")
	}

	// Test with nil endTime
	_, err = ComputePeriodByUnit(now.Format(time.RFC3339), nil, 0, "Month")
	if err == nil {
		t.Error("Expected error for nil endTime")
	}

	// Test with int64 time
	_, err = ComputePeriodByUnit(
		now.Unix(),
		now.AddDate(0, 1, 0).Unix(),
		0,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error with int64 time: %v", err)
	}

	// Test with invalid time format
	_, err = ComputePeriodByUnit("invalid", "invalid", 0, "Month")
	if err == nil {
		t.Error("Expected error for invalid time format")
	}

	// Test with unsupported period unit
	_, err = ComputePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(0, 1, 0).Format(time.RFC3339),
		0,
		"Invalid",
	)
	if err == nil {
		t.Error("Expected error for unsupported period unit")
	}

	// Test with currentPeriod set
	result, err = ComputePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(0, 2, 0).Format(time.RFC3339),
		5,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 5 {
		t.Errorf("Expected 5 (currentPeriod), got %d", result)
	}

	// Test UnStandardRFC3339 format
	UnStandardRFC3339 := "2006-01-02T15:04Z07:00"
	create := now.Format(UnStandardRFC3339)
	end := now.AddDate(0, 1, 0).Format(UnStandardRFC3339)
	_, err = ComputePeriodByUnit(create, end, 0, "Month")
	if err != nil {
		t.Errorf("Unexpected error with UnStandardRFC3339: %v", err)
	}

	// Test period > 12
	result, err = ComputePeriodByUnit(
		now.Format(time.RFC3339),
		now.AddDate(2, 0, 0).Format(time.RFC3339),
		0,
		"Month",
	)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 12 {
		t.Errorf("Expected 12 (max period), got %d", result)
	}
}

func TestUnitCommonBuildStateConf(t *testing.T) {
	pending := []string{"Creating", "Updating"}
	target := []string{"Running", "Available"}
	timeout := 5 * time.Minute
	delay := 10 * time.Second

	refreshFunc := func() (interface{}, string, error) {
		return nil, "", nil
	}

	conf := BuildStateConf(pending, target, timeout, delay, refreshFunc)

	if conf == nil {
		t.Fatal("Expected non-nil StateChangeConf")
	}

	if !reflect.DeepEqual(conf.Pending, pending) {
		t.Errorf("Pending mismatch: expected %v, got %v", pending, conf.Pending)
	}

	if !reflect.DeepEqual(conf.Target, target) {
		t.Errorf("Target mismatch: expected %v, got %v", target, conf.Target)
	}

	if conf.Timeout != timeout {
		t.Errorf("Timeout mismatch: expected %v, got %v", timeout, conf.Timeout)
	}

	if conf.Delay != delay {
		t.Errorf("Delay mismatch: expected %v, got %v", delay, conf.Delay)
	}

	if conf.MinTimeout != 3*time.Second {
		t.Errorf("MinTimeout should be 3 seconds, got %v", conf.MinTimeout)
	}

	if conf.Refresh == nil {
		t.Error("Refresh function should not be nil")
	}
}

func TestUnitCommonCheckWaitForReady(t *testing.T) {
	type TestStruct struct {
		Status string
		Count  int
		Ready  bool
	}

	tests := []struct {
		name        string
		object      interface{}
		conditions  map[string]interface{}
		expectReady bool
		expectError bool
	}{
		{
			name: "nil conditions",
			object: TestStruct{
				Status: "Running",
				Count:  5,
			},
			conditions:  nil,
			expectReady: false,
			expectError: false,
		},
		{
			name: "all conditions met",
			object: TestStruct{
				Status: "Running",
				Count:  5,
				Ready:  true,
			},
			conditions: map[string]interface{}{
				"Status": "Running",
				"Count":  5,
				"Ready":  true,
			},
			expectReady: true,
			expectError: false,
		},
		{
			name: "condition not met",
			object: TestStruct{
				Status: "Creating",
				Count:  5,
			},
			conditions: map[string]interface{}{
				"Status": "Running",
			},
			expectReady: false,
			expectError: false,
		},
		{
			name: "missing attribute",
			object: TestStruct{
				Status: "Running",
			},
			conditions: map[string]interface{}{
				"NonExistentField": "value",
			},
			expectReady: false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ready, values, err := CheckWaitForReady(tt.object, tt.conditions)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if ready != tt.expectReady {
				t.Errorf("Expected ready=%v, got %v", tt.expectReady, ready)
			}

			if !tt.expectError && tt.conditions != nil && values == nil {
				t.Error("Expected non-nil values map")
			}
		})
	}
}
