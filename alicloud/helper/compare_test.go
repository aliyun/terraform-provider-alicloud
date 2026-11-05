package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for the compare helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonCompareMapWithIgnoreEquivalent(t *testing.T) {
	tests := []struct {
		name   string
		map1   map[string]interface{}
		map2   map[string]interface{}
		ignore []string
		expect bool
	}{
		{
			name:   "EqualMapsNoIgnore",
			map1:   map[string]interface{}{"a": 1, "b": "test"},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{},
			expect: true,
		},
		{
			name:   "EqualMapsWithIgnore",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test", "c": false},
			ignore: []string{"c"},
			expect: true,
		},
		{
			name:   "DifferentMapsIgnoreNotApplied",
			map1:   map[string]interface{}{"a": 1, "b": "test"},
			map2:   map[string]interface{}{"a": 1, "b": "different"},
			ignore: []string{},
			expect: false,
		},
		{
			name:   "DifferentLengthMaps",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{"c"},
			expect: false,
		},
		{
			name:   "IgnoreKeyMissingInOneMap",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 1, "b": "test"},
			ignore: []string{"c"},
			expect: false, // 长度不同
		},
		{
			name:   "IgnoreMultipleKeys",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true, "d": 3.14},
			map2:   map[string]interface{}{"a": 1, "b": "test", "c": false, "d": 6.28},
			ignore: []string{"c", "d"},
			expect: true,
		},
		{
			name:   "DifferentValueNonIgnoredKey",
			map1:   map[string]interface{}{"a": 1, "b": "test", "c": true},
			map2:   map[string]interface{}{"a": 2, "b": "test", "c": true},
			ignore: []string{"c"},
			expect: false,
		},
		{
			name:   "NestedMapsIgnored",
			map1:   map[string]interface{}{"a": 1, "config": map[string]interface{}{"key": "value"}},
			map2:   map[string]interface{}{"a": 1, "config": map[string]interface{}{"key": "different"}},
			ignore: []string{"config"},
			expect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareMapWithIgnoreEquivalent(tt.map1, tt.map2, tt.ignore)
			if result != tt.expect {
				t.Errorf(
					"CompareMapWithIgnoreEquivalent() = %v, want %v\nMap1: %#v\nMap2: %#v\nIgnore: %v",
					result, tt.expect, tt.map1, tt.map2, tt.ignore,
				)
			}
		})
	}
}

func TestUnitCommonCompareJsonTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal templates",
			`{"key":"value","number":123}`,
			`{"number":123,"key":"value"}`,
			true,
			false,
		},
		{
			"different templates",
			`{"key":"value1"}`,
			`{"key":"value2"}`,
			false,
			false,
		},
		{
			"invalid json1",
			`{invalid}`,
			`{"key":"value"}`,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareJsonTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonCompareArrayJsonTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal arrays",
			`[{"key":"value"},{"number":123}]`,
			`[{"key":"value"},{"number":123}]`,
			true,
			false,
		},
		{
			"different arrays",
			`[{"key":"value1"}]`,
			`[{"key":"value2"}]`,
			false,
			false,
		},
		{
			"invalid json",
			`[{invalid}]`,
			`[{"key":"value"}]`,
			false,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareArrayJsonTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonCompareYamlTemplateAreEquivalent(t *testing.T) {
	tests := []struct {
		name     string
		tem1     string
		tem2     string
		expected bool
		hasError bool
	}{
		{
			"equal yaml",
			"key: value\nnumber: 123",
			"number: 123\nkey: value",
			true,
			false,
		},
		{
			"different yaml",
			"key: value1",
			"key: value2",
			false,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CompareYamlTemplateAreEquivalent(tt.tem1, tt.tem2)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonIsSlicesValuesEqual(t *testing.T) {
	tests := []struct {
		name     string
		left     []interface{}
		right    []interface{}
		expected bool
	}{
		{
			name:     "OrderStringSlice",
			left:     []interface{}{"a", "b", "c"},
			right:    []interface{}{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "DisorderStringSlice",
			left:     []interface{}{"b", "c", "a"},
			right:    []interface{}{"a", "b", "c"},
			expected: true,
		},
		{
			name:     "NotEqualStringSlice",
			left:     []interface{}{"a", "b", "d"},
			right:    []interface{}{"a", "b", "c"},
			expected: false,
		},
		{
			name:     "OrderIntSlice",
			left:     []interface{}{1, 2, 3, 4, 5},
			right:    []interface{}{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "DisorderIntSlice",
			left:     []interface{}{5, 3, 1, 4, 2},
			right:    []interface{}{1, 2, 3, 4, 5},
			expected: true,
		},
		{
			name:     "NotEqualIntSlice",
			left:     []interface{}{1, 2, 3, 6, 5},
			right:    []interface{}{1, 2, 3, 4, 5},
			expected: false,
		},
		{
			name:     "OrderFloatSlice",
			left:     []interface{}{1.1, 2.2, 3.3},
			right:    []interface{}{1.1, 2.2, 3.3},
			expected: true,
		},
		{
			name:     "DisorderFloatSlice",
			left:     []interface{}{3.3, 1.1, 2.2},
			right:    []interface{}{1.1, 2.2, 3.3},
			expected: true,
		},
		{
			name:     "NotEqualFloatSlice",
			left:     []interface{}{1.1, 2.2, 6.6},
			right:    []interface{}{1.1, 2.2, 3.3},
			expected: false,
		},
		{
			name:     "OrderBoolSlice",
			left:     []interface{}{true, false},
			right:    []interface{}{true, false},
			expected: true,
		},
		{
			name:     "DisorderBoolSlice",
			left:     []interface{}{false, true},
			right:    []interface{}{true, false},
			expected: true,
		},
		{
			name:     "OrderMixedTypeSlice",
			left:     []interface{}{1, "a", true},
			right:    []interface{}{1, "a", true},
			expected: true,
		},
		{
			name:     "DisorderMixedTypeSlice",
			left:     []interface{}{"a", true, 1},
			right:    []interface{}{1, "a", true},
			expected: true,
		},
		{
			name:     "NotEqualMixedTypeSlice",
			left:     []interface{}{1, "b", true},
			right:    []interface{}{1, "a", true},
			expected: false,
		},
		{
			name:     "EmptySlice",
			left:     []interface{}{},
			right:    []interface{}{},
			expected: true,
		},
		{
			name:     "NilSlice",
			left:     nil,
			right:    nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSlicesValuesEqual(tt.left, tt.right)
			assert.Equal(t, tt.expected, result)
		})
	}
}
