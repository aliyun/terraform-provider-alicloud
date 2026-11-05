package helper

import (
	"reflect"
	"testing"
)

// Tests for the resource-id helpers moved from the root package
// (alicloud/common.go); test names are unchanged from
// alicloud/common_test.go.

func TestUnitCommonParseResourceIds(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"single part", "id1", []string{"id1"}},
		{"two parts", "id1:id2", []string{"id1", "id2"}},
		{"three parts", "vpc:vsw:ecs", []string{"vpc", "vsw", "ecs"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := ParseResourceIds(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseResourceIds(%s) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestUnitCommonParseResourceId(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"valid two parts", "vpc:vsw", 2, []string{"vpc", "vsw"}, false},
		{"valid three parts", "a:b:c", 3, []string{"a", "b", "c"}, false},
		{"invalid length", "vpc:vsw", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceId(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonParseResourceIdN(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"valid split", "a:b:c:d", 2, []string{"a", "b:c:d"}, false},
		{"exact match", "a:b", 2, []string{"a", "b"}, false},
		{"invalid length", "a:b", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceIdN(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonParseResourceIdWithEscaped(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		expected []string
		hasError bool
	}{
		{"no escape", "a:b:c", 3, []string{"a", "b", "c"}, false},
		{"with escape", `a\:b:c`, 2, []string{"a:b", "c"}, false},
		{"multiple escape", `a\:b\:c:d`, 2, []string{"a:b:c", "d"}, false},
		{"invalid length", "a:b", 3, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseResourceIdWithEscaped(tt.input, tt.length)
			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestUnitCommonEscapeColons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no colon", "abc", "abc"},
		{"one colon", "a:b", `a\:b`},
		{"multiple colons", "a:b:c", `a\:b\:c`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeColons(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeColons(%s) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}
