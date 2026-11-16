// lintignore: S013
package oss

import (
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

func TestUnitCommonValidateOssBucketDateTimestamp(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Date",
			value:       "2024-01-15",
			expectError: false,
			description: "Standard YYYY-MM-DD should be valid",
		},
		{
			name:        "Valid_Boundary_LeapYear",
			value:       "2024-02-29",
			expectError: false,
			description: "Leap year date should be valid",
		},
		{
			name:        "Valid_Start_Of_Year",
			value:       "2000-01-01",
			expectError: false,
			description: "Start of year should be valid",
		},
		{
			name:        "Invalid_No_Padding",
			value:       "2024-1-5",
			expectError: true,
			description: "Date without zero-padding should be invalid",
		},
		{
			name:        "Invalid_Slash_Separator",
			value:       "2024/01/15",
			expectError: true,
			description: "Date with slash separator should be invalid",
		},
		{
			name:        "Invalid_String",
			value:       "not-a-date",
			expectError: true,
			description: "Non-date string should be invalid",
		},
		{
			name:        "Invalid_Datetime",
			value:       "2024-01-15T00:00:00Z",
			expectError: true,
			description: "Datetime string should be invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateOssBucketDateTimestamp(tc.value, "date")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}
