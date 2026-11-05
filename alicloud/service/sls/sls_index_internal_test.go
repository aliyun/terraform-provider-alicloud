package sls

import (
	"testing"
)

func TestAreKeysJsonEquivalent(t *testing.T) {
	tests := []struct {
		name        string
		old         string
		new         string
		expected    bool
		expectError bool
	}{
		{
			name:     "API adds alias and caseSensitive defaults",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"","caseSensitive":false}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only alias default",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only caseSensitive default",
			old:      `{"cost_ms":{"type":"text","doc_value":true,"caseSensitive":false}}`,
			new:      `{"cost_ms":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds chn default for text type field",
			old:      `{"message":{"type":"text","doc_value":true,"chn":false,"caseSensitive":false,"alias":""}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "API adds only chn default",
			old:      `{"message":{"type":"text","doc_value":true,"chn":false}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "user changes chn from false to true should surface diff",
			old:      `{"message":{"type":"text","doc_value":true,"chn":true}}`,
			new:      `{"message":{"type":"text","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "both configs identical with no API defaults",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "identical compound field values",
			old:      `{"message":{"type":"text","token":[",",";"]}}`,
			new:      `{"message":{"type":"text","token":[",",";"]}}`,
			expected: true,
		},
		{
			name:     "user changes type from long to text",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"text","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user changes doc_value from true to false",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":false}}`,
			expected: false,
		},
		{
			name:     "user changes alias to non-empty value",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"old"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"new"}}`,
			expected: false,
		},
		{
			name:     "user deletes alias with non-empty value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":"legacy"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes caseSensitive with non-default value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"caseSensitive":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes doc_value should surface diff",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long"}}`,
			expected: false,
		},
		{
			name:     "user deletes alias with empty value is suppressed",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "user adds a new key",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "user deletes an existing key",
			old:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "multiple keys with API defaults on one key",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""},"duration_ms":{"type":"long","doc_value":true,"caseSensitive":false,"alias":""}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: true,
		},
		{
			name:     "multiple keys one has non-default old-only field",
			old:      `{"cost_ms":{"type":"long","doc_value":true,"alias":""},"duration_ms":{"type":"long","doc_value":true,"alias":"custom"}}`,
			new:      `{"cost_ms":{"type":"long","doc_value":true},"duration_ms":{"type":"long","doc_value":true}}`,
			expected: false,
		},
		{
			name:     "both empty strings",
			old:      ``,
			new:      ``,
			expected: true,
		},
		{
			name:     "one empty one not",
			old:      `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:      ``,
			expected: false,
		},
		{
			name:        "invalid JSON in old falls back to original comparison",
			old:         `{invalid}`,
			new:         `{"cost_ms":{"type":"long","doc_value":true}}`,
			expected:    false,
			expectError: true,
		},
		{
			name:        "invalid JSON in new falls back to original comparison",
			old:         `{"cost_ms":{"type":"long","doc_value":true}}`,
			new:         `{invalid}`,
			expected:    false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := areKeysJsonEquivalent(tt.old, tt.new)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("areKeysJsonEquivalent(%q, %q) = %v, want %v", tt.old, tt.new, result, tt.expected)
			}
		})
	}
}
