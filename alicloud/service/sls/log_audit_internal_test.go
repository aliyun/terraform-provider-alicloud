package sls

import (
	"fmt"
	"testing"
)

func TestGetInitParameter(t *testing.T) {
	// Config is a JSON string value in the real API response (not a raw JSON object).
	// getInitParameter does configNew.(string) so Config must be a string type.
	configWithInitParam := `{"initParam":{"aliuid":"1234567890","actiontrail_ttl":"10"}}`

	tests := []struct {
		name            string
		input           string
		expectedDisplay string
		expectError     bool
		expectPanic     bool
		hasInitMap      bool
	}{
		// --- Normal response with DisplayName ---
		{
			name:            "normal response with displayName",
			input:           fmt.Sprintf(`{"AppModel":{"DisplayName":"my-audit-app","Config":%q}}`, configWithInitParam),
			expectedDisplay: "my-audit-app",
			hasInitMap:      true,
		},

		// --- DisplayName is empty string (the bug scenario) ---
		{
			name:            "DisplayName is empty string",
			input:           fmt.Sprintf(`{"AppModel":{"DisplayName":"","Config":%q}}`, configWithInitParam),
			expectedDisplay: "",
			hasInitMap:      true,
		},

		// --- DisplayName field is missing entirely ---
		// Note: getInitParameter uses d["DisplayName"].(string) without comma-ok,
		// so a missing DisplayName field causes a panic. This is a pre-existing
		// limitation in the original code — the displayName == "" fallback in
		// Read() handles the empty-string case, not the missing-field case.
		{
			name:            "DisplayName field missing from AppModel",
			input:           fmt.Sprintf(`{"AppModel":{"Config":%q}}`, configWithInitParam),
			expectedDisplay: "",
			hasInitMap:      false,
			expectPanic:     true,
		},

		// --- AppModel is missing ---
		{
			name:            "AppModel missing entirely",
			input:           `{"someOtherKey":"value"}`,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     false, // no error, just empty return values
		},

		// --- AppModel is not a map ---
		{
			name:            "AppModel is a string not a map",
			input:           `{"AppModel":"not-a-map"}`,
			expectedDisplay: "",
			hasInitMap:      false,
		},

		// --- Invalid JSON ---
		{
			name:            "invalid JSON",
			input:           `{invalid json}`,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     true,
		},

		// --- Empty string input ---
		{
			name:            "empty string input",
			input:           ``,
			expectedDisplay: "",
			hasInitMap:      false,
			expectError:     true, // json.Unmarshal of empty string returns error
		},

		// --- Config is not valid JSON string ---
		// Note: json.Unmarshal of "not-json" fails, m is empty, then
		// m["initParam"].(map[string]interface{}) panics on nil.
		{
			name:            "Config is not a valid JSON string",
			input:           `{"AppModel":{"DisplayName":"test","Config":"not-json"}}`,
			expectedDisplay: "test",
			hasInitMap:      false,
			expectPanic:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.expectPanic {
						t.Fatalf("unexpected panic: %v", r)
					}
					return
				}
				if tt.expectPanic {
					t.Fatalf("expected panic but got none")
				}
			}()

			displayName, initMap, err := getInitParameter(tt.input)

			if tt.expectError && err == nil {
				t.Fatalf("expected an error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if displayName != tt.expectedDisplay {
				t.Errorf("displayName = %q, want %q", displayName, tt.expectedDisplay)
			}
			if tt.hasInitMap && initMap == nil {
				t.Errorf("expected initMap to be non-nil")
			}
			if !tt.hasInitMap && initMap != nil {
				t.Errorf("expected initMap to be nil, got %v", initMap)
			}
		})
	}
}
