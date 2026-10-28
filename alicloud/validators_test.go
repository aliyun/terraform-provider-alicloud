package alicloud

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── helpers ────────────────────────────────────────────────────────────────

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

// Tests for the product-specific validators that remain in this package;
// the generic validator tests moved to alicloud/helper/validators_test.go.

func TestValidateOTSInstanceName(t *testing.T) {
	_, errors := validateOTSInstanceName("abcd", "")
	assert.Nil(t, errors)
	_, errors = validateOTSInstanceName("abc-dfg", "")
	assert.Nil(t, errors)
	_, errors = validateOTSInstanceName("abc123", "")
	assert.Nil(t, errors)
	_, errors = validateOTSInstanceName("AbC-123456", "")
	assert.Nil(t, errors)
	_, errors = validateOTSInstanceName("abcdefghijklmnj", "")
	assert.Nil(t, errors)

	words := []string{"ali", "ots", "taobao", "admin"}
	for _, w := range words {
		_, errors = validateOTSInstanceName(w, "")
		assert.Nil(t, errors)
	}
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aaaaaaaaaaaaaaaaa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("aa bb", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("$aa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("123", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("abcde-", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
	_, errors = validateOTSInstanceName("abc_defg", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the instance name must consist of"))
}

func TestValidateOTSTableName(t *testing.T) {
	_, errors := validateOTSTableName("_", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("a", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("A", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("abcdefg", "")
	assert.Nil(t, errors)
	var sb bytes.Buffer
	for i := 0; i < 255; i++ {
		sb.WriteByte('a')
	}
	_, errors = validateOTSTableName(sb.String(), "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("abc_defg", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("_abcdefg", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTableName("AbC_defg_", "")
	assert.Nil(t, errors)

	_, errors = validateOTSTableName("", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
	_, errors = validateOTSTableName(" ", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
	_, errors = validateOTSTableName(sb.String()+"a", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
	_, errors = validateOTSTableName("$aaa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
	_, errors = validateOTSTableName("1aaa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
	_, errors = validateOTSTableName("aA$%", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the table name must consist of"))
}

func TestValidateOTSTunnelName(t *testing.T) {
	_, errors := validateOTSTunnelName("_", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("a", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("A", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("abcdefg", "")
	assert.Nil(t, errors)
	var sb bytes.Buffer
	for i := 0; i < 255; i++ {
		sb.WriteByte('a')
	}
	_, errors = validateOTSTunnelName(sb.String(), "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("abc_defg", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("_abcdefg", "")
	assert.Nil(t, errors)
	_, errors = validateOTSTunnelName("AbC_defg_", "")
	assert.Nil(t, errors)

	_, errors = validateOTSTunnelName("", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
	_, errors = validateOTSTunnelName(" ", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
	_, errors = validateOTSTunnelName(sb.String()+"a", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
	_, errors = validateOTSTunnelName("$aaa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
	_, errors = validateOTSTunnelName("1aaa", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
	_, errors = validateOTSTunnelName("aA$%", "")
	assert.NotNil(t, errors)
	assert.True(t, strings.Contains(errors[0].Error(), "the tunnel name must consist of"))
}

// ─── TestUnitCommon* ─────────────────────────────────────────────────────────

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

// TestUnitCommonValidateOnsGroupId moved to alicloud/service/ons/validators_test.go
// when validateOnsGroupId relocated to the ons package.

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

func TestUnitCommonValidateDBConnectionPort(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Empty_String",
			value:       "",
			expectError: false,
			description: "empty string is allowed (no validation)",
		},
		{
			name:        "Valid_Boundary_Low",
			value:       "1000",
			expectError: false,
			description: "port 1000 is the minimum valid value",
		},
		{
			name:        "Valid_Boundary_High",
			value:       "5999",
			expectError: false,
			description: "port 5999 is the maximum valid value",
		},
		{
			name:        "Valid_Middle",
			value:       "3306",
			expectError: false,
			description: "port 3306 (MySQL default) should be valid",
		},
		{
			name:        "Invalid_Below_Min",
			value:       "999",
			expectError: true,
			description: "port 999 is below minimum (1000)",
		},
		{
			name:        "Invalid_Above_Max",
			value:       "6000",
			expectError: true,
			description: "port 6000 is above maximum (5999)",
		},
		{
			name:        "Invalid_Zero",
			value:       "0",
			expectError: true,
			description: "port 0 is below minimum",
		},
		{
			name:        "Invalid_Non_Numeric",
			value:       "abc",
			expectError: true,
			description: "non-numeric string should produce parse error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateDBConnectionPort(tc.value, "port")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}

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

func TestUnitCommonIntBetween(t *testing.T) {
	validate := intBetween(10, 100)

	testCases := []struct {
		name        string
		value       interface{}
		expectError bool
		description string
	}{
		{
			name:        "Valid_Min",
			value:       10,
			expectError: false,
			description: "minimum value should be valid",
		},
		{
			name:        "Valid_Max",
			value:       100,
			expectError: false,
			description: "maximum value should be valid",
		},
		{
			name:        "Valid_Middle",
			value:       55,
			expectError: false,
			description: "value in middle of range should be valid",
		},
		{
			name:        "Invalid_Below_Min",
			value:       9,
			expectError: true,
			description: "value below minimum should be invalid",
		},
		{
			name:        "Invalid_Above_Max",
			value:       101,
			expectError: true,
			description: "value above maximum should be invalid",
		},
		{
			name:        "Invalid_Wrong_Type",
			value:       "50",
			expectError: true,
			description: "string instead of int should be invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validate(tc.value, "count")
			assert.Empty(t, ws)
			if tc.expectError {
				assert.NotEmpty(t, errors, tc.description)
			} else {
				assert.Empty(t, errors, tc.description)
			}
		})
	}
}

func TestUnitCommonValidateOTSIndexName(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Letter_Start",
			value:       "myIndex",
			expectError: false,
			description: "index name starting with letter should be valid",
		},
		{
			name:        "Valid_Underscore_Start",
			value:       "_myIndex",
			expectError: false,
			description: "index name starting with underscore should be valid",
		},
		{
			name:        "Valid_Single_Char",
			value:       "a",
			expectError: false,
			description: "single letter is valid (length 1)",
		},
		{
			name:        "Valid_With_Numbers",
			value:       "index_001",
			expectError: false,
			description: "index name with numbers should be valid",
		},
		{
			name:        "Valid_Max_Length",
			value:       strings.Repeat("a", 255),
			expectError: false,
			description: "index name of exactly 255 chars should be valid",
		},
		{
			name:        "Invalid_Digit_Start",
			value:       "1index",
			expectError: true,
			description: "index name starting with digit should be invalid",
		},
		{
			name:        "Invalid_Hyphen",
			value:       "my-index",
			expectError: true,
			description: "hyphen is not allowed in OTS index name",
		},
		{
			name:        "Invalid_Too_Long",
			value:       strings.Repeat("a", 256),
			expectError: true,
			description: "index name longer than 255 chars should be invalid",
		},
		{
			name:        "Invalid_Empty",
			value:       "",
			expectError: true,
			description: "empty name should be invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := validateOTSIndexName(tc.value, "index_name")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}

func TestUnitCommonValidateRedisConfig(t *testing.T) {
	testCases := []struct {
		name        string
		value       interface{}
		expectError bool
		description string
	}{
		{
			name:        "Valid_Config",
			value:       map[string]interface{}{"maxmemory-policy": "allkeys-lru"},
			expectError: false,
			description: "non-empty map should be valid",
		},
		{
			name:        "Valid_Multiple_Keys",
			value:       map[string]interface{}{"maxmemory-policy": "allkeys-lru", "hz": "15"},
			expectError: false,
			description: "map with multiple keys should be valid",
		},
		{
			name:        "Invalid_Empty_Map",
			value:       map[string]interface{}{},
			expectError: true,
			description: "empty map should be invalid",
		},
		{
			name:        "Invalid_Nil",
			value:       nil,
			expectError: true,
			description: "nil value coerces to empty map, should be invalid",
		},
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
