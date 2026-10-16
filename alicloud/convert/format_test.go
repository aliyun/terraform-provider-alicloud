package convert

import (
	"encoding/json"
	"testing"
)

// Tests for the format helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonFormatInt(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int
		panic    bool
	}{
		{"nil", nil, 0, false},
		{"float64", 3.14, 3, false},
		{"float32", float32(3.0), 3, false},
		{"int64", int64(123), 123, false},
		{"int32", int32(456), 456, false},
		{"int", 789, 789, false},
		{"string_number", "123", 123, false},
		{"string_empty", "", 0, false},
		{"json_number", json.Number("123"), 123, false},
		{"string_invalid", "abc", 0, true},
		{"unsupported_type", true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := FormatInt(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonFormatBool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
		panic    bool
	}{
		{"nil", nil, false, false},
		{"bool_true", true, true, false},
		{"bool_false", false, false, false},
		{"string_true", "true", true, false},
		{"string_false", "false", false, false},
		{"string_invalid", "abc", false, true},
		{"unsupported_type", 123, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := FormatBool(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}

func TestUnitCommonFormatFloat64(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
		panic    bool
	}{
		{"nil", nil, 0.0, false},
		{"float64", 3.14, 3.14, false},
		{"float32", float32(3.0), 3.0, false},
		{"int64", int64(123), 123.0, false},
		{"int32", int32(456), 456.0, false},
		{"int", 789, 789.0, false},
		{"string_number", "3.14", 3.14, false},
		{"string_empty", "", 0.0, false},
		{"json_number", json.Number("123.45"), 123.45, false},
		{"string_invalid", "abc", 0.0, true},
		{"unsupported_type", true, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.panic {
						t.Errorf("Unexpected panic: %v", r)
					}
				} else if tt.panic {
					t.Error("Expected panic but no panic occurred")
				}
			}()

			result := FormatFloat64(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}
