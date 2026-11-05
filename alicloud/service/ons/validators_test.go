package ons

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertNoErrors asserts the validator produced no warnings or errors.
func assertNoErrors(t *testing.T, ws []string, errs []error, label string) {
	t.Helper()
	if len(errs) != 0 {
		t.Errorf("%s: expected no errors but got: %v", label, errs)
	}
}

// assertHasErrors asserts the validator produced at least one error.
func assertHasErrors(t *testing.T, ws []string, errs []error, label string) {
	t.Helper()
	if len(errs) == 0 {
		t.Errorf("%s: expected errors but got none", label)
	}
}

func TestUnitAlicloudOnsGroupIdValidator(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_GID_Hyphen",
			value:       "GID-test123",
			expectError: false,
			description: "GID- prefix with valid chars should be valid",
		},
		{
			name:        "Valid_GID_Underscore",
			value:       "GID_test123",
			expectError: false,
			description: "GID_ prefix with valid chars should be valid",
		},
		{
			name:        "Valid_Exact_Min_Length",
			value:       "GID-abc",
			expectError: false,
			description: "GID- prefix with 7 total chars (min) should be valid",
		},
		{
			name:        "Valid_Exact_Max_Length",
			value:       "GID-" + strings.Repeat("a", 60),
			expectError: false,
			description: "GID- prefix with 64 total chars (max) should be valid",
		},
		{
			name:        "Invalid_Wrong_Prefix",
			value:       "GROUP-test123",
			expectError: true,
			description: "Wrong prefix should be invalid",
		},
		{
			name:        "Invalid_Too_Short",
			value:       "GID-a",
			expectError: true,
			description: "Too short (6 chars) should be invalid",
		},
		{
			name:        "Invalid_Too_Long",
			value:       "GID-" + strings.Repeat("a", 61),
			expectError: true,
			description: "Too long (65 chars) should be invalid",
		},
		{
			name:        "Invalid_Special_Chars",
			value:       "GID-test@123",
			expectError: true,
			description: "Special characters not allowed",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateOnsGroupId(tc.value, "group_id")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}
