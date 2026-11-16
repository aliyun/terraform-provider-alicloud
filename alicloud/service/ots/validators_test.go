// lintignore: S013
package ots

import (
	"bytes"
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
