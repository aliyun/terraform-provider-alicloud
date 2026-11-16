// lintignore: S013
package alidns

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertNoErrors(t *testing.T, ws []string, errs []error, label string) {
	t.Helper()
	if len(errs) != 0 {
		t.Errorf("%s: expected no errors but got: %v", label, errs)
	}
}

func assertHasErrors(t *testing.T, ws []string, errs []error, label string) {
	t.Helper()
	if len(errs) == 0 {
		t.Errorf("%s: expected errors but got none", label)
	}
}

func TestUnitCommonValidateRR(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Simple",
			value:       "www",
			expectError: false,
			description: "Simple hostname should be valid",
		},
		{
			name:        "Valid_Subdomain",
			value:       "sub.domain",
			expectError: false,
			description: "Subdomain should be valid",
		},
		{
			name:        "Valid_At_Symbol",
			value:       "@",
			expectError: false,
			description: "@ symbol should be valid",
		},
		{
			name:        "Invalid_Starts_With_Hyphen",
			value:       "-invalid",
			expectError: true,
			description: "Starting with hyphen should be invalid",
		},
		{
			name:        "Invalid_Ends_With_Hyphen",
			value:       "invalid-",
			expectError: true,
			description: "Ending with hyphen should be invalid",
		},
		{
			name:        "Invalid_Too_Long",
			value:       strings.Repeat("a", 254),
			expectError: true,
			description: "RR longer than 253 chars should be invalid",
		},
		{
			name:        "Invalid_Part_Too_Long",
			value:       "www." + strings.Repeat("a", 64) + ".com",
			expectError: true,
			description: "Part longer than 63 chars should be invalid",
		},
		{
			name:        "Valid_Part_Exact_Max",
			value:       strings.Repeat("a", 63),
			expectError: false,
			description: "Part of exactly 63 chars should be valid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateRR(tc.value, "rr")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}
