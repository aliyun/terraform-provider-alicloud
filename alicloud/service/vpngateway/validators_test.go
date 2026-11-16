// lintignore: S013
package vpngateway

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCommonValidateSslVpnPortValue(t *testing.T) {
	excluded := []int{22, 2222, 22222}
	validate := validateSslVpnPortValue(excluded)

	testCases := []struct {
		name        string
		value       int
		expectError bool
		description string
	}{
		{
			name:        "Valid_Port",
			value:       443,
			expectError: false,
			description: "port 443 is valid and not excluded",
		},
		{
			name:        "Valid_Boundary_Low",
			value:       1,
			expectError: false,
			description: "port 1 is the minimum valid",
		},
		{
			name:        "Valid_Boundary_High",
			value:       65535,
			expectError: false,
			description: "port 65535 is the maximum valid",
		},
		{
			name:        "Invalid_Excluded_22",
			value:       22,
			expectError: true,
			description: "port 22 is in the excluded list",
		},
		{
			name:        "Invalid_Excluded_2222",
			value:       2222,
			expectError: true,
			description: "port 2222 is in the excluded list",
		},
		{
			name:        "Invalid_Below_Min",
			value:       0,
			expectError: true,
			description: "port 0 is below minimum (1)",
		},
		{
			name:        "Invalid_Above_Max",
			value:       65536,
			expectError: true,
			description: "port 65536 is above maximum (65535)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validate(tc.value, "port")
			if tc.expectError {
				assert.NotEmpty(t, errors, tc.description)
			} else {
				assert.Empty(t, errors, tc.description)
				assert.Empty(t, ws)
			}
		})
	}
}
