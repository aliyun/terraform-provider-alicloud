package connectivity

import (
	"testing"
)

func TestUnitCommonConvertKebabToSnake(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Normal case", "test-case", "test_case"},
		{"Empty string", "", ""},
		{"No hyphens", "helloWorld", "helloWorld"},
		{"Multiple hyphens", "this-is-a-test", "this_is_a_test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertKebabToSnake(tt.input); got != tt.expected {
				t.Errorf("ConvertKebabToSnake() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnitCommonIsInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Int", 42, true},
		{"Int64", int64(100), true},
		{"Float", 3.14, false},
		{"Valid string", "123", true},
		{"Invalid string", "12a3", false},
		{"Bool", true, false},
		{"Struct", struct{}{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInteger(tt.input); got != tt.expected {
				t.Errorf("isInteger() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUnitCommonIsString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"String type", "hello", true},
		{"Int type", 42, false},
		{"Struct type", struct{}{}, false},
		{"Byte slice", []byte{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isString(tt.input); got != tt.expected {
				t.Errorf("isString() = %v, want %v", got, tt.expected)
			}
		})
	}
}
