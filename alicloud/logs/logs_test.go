package logs

import (
	"os"
	"testing"
)

// Tests copied from alicloud/common_test.go; test names are unchanged.
func TestUnitCommonDebugOn(t *testing.T) {
	// Save original value
	originalDebug := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", originalDebug)

	// Test when DEBUG is not set
	os.Setenv("DEBUG", "")
	if DebugOn() {
		t.Error("debugOn should return false when DEBUG is not set")
	}

	// Test when DEBUG contains terraform
	os.Setenv("DEBUG", "terraform")
	if !DebugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG contains terraform with other values
	os.Setenv("DEBUG", "something,terraform,other")
	if !DebugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG does not contain terraform
	os.Setenv("DEBUG", "other")
	if DebugOn() {
		t.Error("debugOn should return false when DEBUG does not contain terraform")
	}
}

func TestUnitCommonAddDebug(t *testing.T) {
	// Save original DEBUG env
	originalDebug := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", originalDebug)

	tests := []struct {
		name        string
		debugEnv    string
		action      interface{}
		content     interface{}
		requestInfo []interface{}
	}{
		{
			name:     "debug off",
			debugEnv: "",
			action:   "TestAction",
			content:  "test content",
		},
		{
			name:     "debug on - simple",
			debugEnv: "terraform",
			action:   "CreateInstance",
			content:  "creating instance",
		},
		{
			name:     "debug on - with map request",
			debugEnv: "terraform",
			action:   "UpdateResource",
			content:  "updating",
			requestInfo: []interface{}{
				map[string]interface{}{"key": "value"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DEBUG", tt.debugEnv)

			// Should not panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("addDebug panicked: %v", r)
				}
			}()

			AddDebug(tt.action, tt.content, tt.requestInfo...)
		})
	}
}
