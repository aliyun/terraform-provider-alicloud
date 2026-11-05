package alicloud

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
)

// Tests in this file cover what remains in the root package plus functions
// already sunk to leaf packages and re-exported here as aliases (writeToFile,
// buildClientToken, debugOn, addDebug, GetUserHomeDir). Tests of the utilities
// moved to alicloud/helper live in that package now; test names are unchanged.

// Migrated to alicloud/helper (output_test.go).
func TestUnitCommonWriteToFile(t *testing.T) {
	tmpFile := fmt.Sprintf("%s/testfile.txt", t.TempDir())

	t.Run("write content", func(t *testing.T) {
		err := writeToFile(tmpFile, "test content")
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		assert.Equal(t, "test content", string(data))
	})

	t.Run("write json", func(t *testing.T) {
		content := map[string]interface{}{"key": "value"}
		err := writeToFile(tmpFile, content)
		assert.NoError(t, err)
		data, _ := os.ReadFile(tmpFile)
		var result map[string]interface{}
		json.Unmarshal(data, &result)
		assert.Equal(t, content["key"], result["key"])
	})
}

func TestUnitCommonNewInstanceDiff(t *testing.T) {

	tests := []struct {
		name           string
		state          *terraform.InstanceState
		attributes     map[string]interface{}
		attributesDiff map[string]interface{}
		expectChanges  map[string]terraform.ResourceAttrDiff
	}{
		{
			name: "simple attribute change",
			state: &terraform.InstanceState{
				ID: "vpc-123456",
				Attributes: map[string]string{
					"name":       "old_vpc_name",
					"cidr_block": "192.168.0.0/16",
				},
			},
			attributes: map[string]interface{}{
				"name":       "old_vpc_name",
				"cidr_block": "192.168.0.0/16",
			},
			attributesDiff: map[string]interface{}{
				"name": "new_vpc_name",
			},
			expectChanges: map[string]terraform.ResourceAttrDiff{
				"name": {Old: "old_vpc_name", New: "new_vpc_name"},
			},
		},
		{
			name: "add new attribute",
			state: &terraform.InstanceState{
				ID: "vpc-234567",
				Attributes: map[string]string{
					"cidr_block": "10.0.0.0/16",
				},
			},
			attributes: map[string]interface{}{
				"cidr_block": "10.0.0.0/16",
			},
			attributesDiff: map[string]interface{}{
				"description": "new description",
			},
			expectChanges: map[string]terraform.ResourceAttrDiff{
				"description": {Old: "", New: "new description"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diff, err := newInstanceDiff("alicloud_vpc", tt.attributes, tt.attributesDiff, tt.state)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			for attr, expectedDiff := range tt.expectChanges {
				actualDiff, ok := diff.Attributes[attr]
				if !ok {
					t.Errorf("missing expected diff for attribute: %s", attr)
					continue
				}

				if actualDiff.Old != expectedDiff.Old || actualDiff.New != expectedDiff.New {
					t.Errorf("attribute %s mismatch:\n  expected: Old=%q, New=%q\n  actual:   Old=%q, New=%q",
						attr, expectedDiff.Old, expectedDiff.New, actualDiff.Old, actualDiff.New)
				}
			}

			for attr := range diff.Attributes {
				found := false
				for expectedAttr := range tt.expectChanges {
					if attr == expectedAttr {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("unexpected diff attribute: %s (Old=%q, New=%q)",
						attr, diff.Attributes[attr].Old, diff.Attributes[attr].New)
				}
			}

		})
	}
}

// Migrated to alicloud/helper (client_token_test.go).
func TestUnitCommonBuildClientToken(t *testing.T) {
	token1 := buildClientToken("CreateInstance")
	token2 := buildClientToken("CreateInstance")

	if token1 == "" {
		t.Error("buildClientToken returned empty string")
	}

	if len(token1) > 64 {
		t.Error("buildClientToken returned token longer than 64 characters")
	}

	if token1 == token2 {
		t.Error("buildClientToken should return different tokens")
	}

	if !strings.Contains(token1, "TF-CreateInstance") {
		t.Error("buildClientToken should contain action name")
	}
}

func TestUnitCommonParseSlbListenerId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{"valid 2 parts", "slb-123:http", 2, false},
		{"valid 3 parts", "slb-123:http:80", 3, false},
		{"invalid", "slb-123", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseSlbListenerId(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != tt.expected {
					t.Errorf("Expected %d parts, got %d", tt.expected, len(result))
				}
			}
		})
	}
}

func TestUnitCommonGetCenChildInstanceType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{"vpc instance", "vpc-12345", "VPC", false},
		{"vbr instance", "vbr-12345", "VBR", false},
		{"ccn instance", "ccn-12345", "CCN", false},
		{"invalid", "ecs-12345", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetCenChildInstanceType(tt.input)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

// Migrated to alicloud/logs (logs_test.go).
func TestUnitCommonDebugOn(t *testing.T) {
	// Save original value
	originalDebug := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", originalDebug)

	// Test when DEBUG is not set
	os.Setenv("DEBUG", "")
	if debugOn() {
		t.Error("debugOn should return false when DEBUG is not set")
	}

	// Test when DEBUG contains terraform
	os.Setenv("DEBUG", "terraform")
	if !debugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG contains terraform with other values
	os.Setenv("DEBUG", "something,terraform,other")
	if !debugOn() {
		t.Error("debugOn should return true when DEBUG contains terraform")
	}

	// Test when DEBUG does not contain terraform
	os.Setenv("DEBUG", "other")
	if debugOn() {
		t.Error("debugOn should return false when DEBUG does not contain terraform")
	}
}

// Migrated to alicloud/helper (output_test.go).
func TestUnitCommonWriteToFileWithHomeDir(t *testing.T) {
	// Test with ~ in path
	tmpDir := t.TempDir()
	// We can't directly test ~ expansion without mocking, but we can test the logic
	err := writeToFile(fmt.Sprintf("%s/test.txt", tmpDir), "content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Test with nil data
	err = writeToFile(fmt.Sprintf("%s/test2.txt", tmpDir), nil)
	if err != nil {
		t.Errorf("Unexpected error for nil: %v", err)
	}

	// Test overwriting existing file
	err = writeToFile(fmt.Sprintf("%s/test.txt", tmpDir), "new content")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUnitCommonNewInstanceDiffComplex(t *testing.T) {
	// Test with simple map attribute changes
	state := &terraform.InstanceState{
		ID: "vpc-123",
		Attributes: map[string]string{
			"cidr_block": "192.168.0.0/16",
			"name":       "old-name",
		},
	}

	attributes := map[string]interface{}{
		"cidr_block": "192.168.0.0/16",
		"name":       "old-name",
	}

	attributesDiff := map[string]interface{}{
		"name": "new-name",
	}

	diff, err := newInstanceDiff("alicloud_vpc", attributes, attributesDiff, state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if diff == nil {
		t.Fatal("Expected non-nil diff")
	}

	// Verify the diff contains the change
	if diff.Attributes["name"].Old != "old-name" || diff.Attributes["name"].New != "new-name" {
		t.Errorf("Expected name change from old-name to new-name")
	}
}

// Migrated to alicloud/helper (output_test.go).
func TestUnitCommonWriteToFileEdgeCases(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with map data
	mapData := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	err := writeToFile(fmt.Sprintf("%s/map.json", tmpDir), mapData)
	if err != nil {
		t.Errorf("Unexpected error writing map: %v", err)
	}

	// Test with struct data
	structData := struct {
		Name string
		Age  int
	}{
		Name: "test",
		Age:  30,
	}
	err = writeToFile(fmt.Sprintf("%s/struct.json", tmpDir), structData)
	if err != nil {
		t.Errorf("Unexpected error writing struct: %v", err)
	}

	// Verify file was created and can be read
	content, err := os.ReadFile(fmt.Sprintf("%s/struct.json", tmpDir))
	if err != nil {
		t.Errorf("Failed to read written file: %v", err)
	}
	if len(content) == 0 {
		t.Error("Written file is empty")
	}
}

func TestUnitCommonNewInstanceDiffRemoveAttributes(t *testing.T) {
	// Test removing attributes (setting to empty)
	state := &terraform.InstanceState{
		ID: "vpc-123",
		Attributes: map[string]string{
			"description": "old description",
			"name":        "test-vpc",
		},
	}

	attributes := map[string]interface{}{
		"description": "old description",
		"name":        "test-vpc",
	}

	attributesDiff := map[string]interface{}{
		"description": "",
		"name":        "test-vpc",
	}

	diff, err := newInstanceDiff("alicloud_vpc", attributes, attributesDiff, state)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if diff == nil {
		t.Fatal("Expected non-nil diff")
	}

	// Verify description was removed
	if descDiff, ok := diff.Attributes["description"]; ok {
		if descDiff.New != "" {
			t.Errorf("Expected empty description, got %s", descDiff.New)
		}
	}
}

// Tests for previously uncovered functions

// Migrated to alicloud/helper (output_test.go).
func TestUnitCommonGetUserHomeDir(t *testing.T) {
	// Test getting user home directory
	homeDir, err := GetUserHomeDir()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if homeDir == "" {
		t.Error("Expected non-empty home directory")
	}
	// Verify it's an absolute path
	if !strings.HasPrefix(homeDir, "/") && !strings.Contains(homeDir, ":\\") {
		t.Errorf("Expected absolute path, got %s", homeDir)
	}
}

// Migrated to alicloud/helper (output_test.go).
func TestUnitCommonWriteToFileComprehensive(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		filePath string
		data     interface{}
		wantErr  bool
	}{
		{
			name:     "write string",
			filePath: fmt.Sprintf("%s/string.txt", tmpDir),
			data:     "test string content",
			wantErr:  false,
		},
		{
			name:     "write nil",
			filePath: fmt.Sprintf("%s/nil.txt", tmpDir),
			data:     nil,
			wantErr:  false,
		},
		{
			name:     "write map",
			filePath: fmt.Sprintf("%s/map.json", tmpDir),
			data:     map[string]interface{}{"key": "value", "num": 123},
			wantErr:  false,
		},
		{
			name:     "write struct",
			filePath: fmt.Sprintf("%s/struct.json", tmpDir),
			data: struct {
				Name string
				Age  int
			}{"test", 30},
			wantErr: false,
		},
		{
			name:     "write slice",
			filePath: fmt.Sprintf("%s/slice.json", tmpDir),
			data:     []string{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "unmarshalable data",
			filePath: fmt.Sprintf("%s/channel.json", tmpDir),
			data:     make(chan int),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := writeToFile(tt.filePath, tt.data)

			if tt.wantErr && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify file was created for non-error cases
			if !tt.wantErr && tt.data != nil {
				if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
					t.Error("File was not created")
				}
			}
		})
	}

	// Test overwriting existing file
	t.Run("overwrite existing file", func(t *testing.T) {
		filePath := fmt.Sprintf("%s/overwrite.txt", tmpDir)

		// Write first time
		err := writeToFile(filePath, "first content")
		if err != nil {
			t.Fatalf("First write failed: %v", err)
		}

		// Write second time (overwrite)
		err = writeToFile(filePath, "second content")
		if err != nil {
			t.Fatalf("Overwrite failed: %v", err)
		}

		// Verify content was updated
		content, _ := os.ReadFile(filePath)
		if string(content) != "second content" {
			t.Errorf("Expected 'second content', got %s", string(content))
		}
	})
}

// Migrated to alicloud/logs (logs_test.go).
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

			addDebug(tt.action, tt.content, tt.requestInfo...)
		})
	}
}
