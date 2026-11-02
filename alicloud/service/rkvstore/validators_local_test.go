package rkvstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertNoErrors/assertHasErrors mirror the helpers in the root package's
// validators_test.go; relocated with the test when validateRedisConfig moved
// into package rkvstore. Behavior preserved verbatim.

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

// TestUnitCommonValidateRedisConfig is the rebuilt unit test for the
// product-specific validateRedisConfig validator that moved from the root
// alicloud package into package rkvstore. Original assertions preserved.

func TestUnitCommonValidateRedisConfig(t *testing.T) {
	testCases := []struct {
		name        string
		value       interface{}
		expectError bool
		description string
	}{
		{name: "Valid_Config", value: map[string]interface{}{"maxmemory-policy": "allkeys-lru"}, expectError: false, description: "non-empty map should be valid"},
		{name: "Valid_Multiple_Keys", value: map[string]interface{}{"maxmemory-policy": "allkeys-lru", "hz": "15"}, expectError: false, description: "map with multiple keys should be valid"},
		{name: "Invalid_Empty_Map", value: map[string]interface{}{}, expectError: true, description: "empty map should be invalid"},
		{name: "Invalid_Nil", value: nil, expectError: true, description: "nil value coerces to empty map, should be invalid"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateRedisConfig(tc.value, "redis_config")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}
