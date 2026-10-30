package validators

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
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

// Tests for the generic validators moved from the root package; names are
// unchanged from alicloud/validators_test.go.

func TestValidateCIDRNetworkAddress(t *testing.T) {
	validCIDRNetworkAddress := []string{"192.168.10.0/24", "0.0.0.0/0", "10.121.10.0/24"}
	for _, v := range validCIDRNetworkAddress {
		_, errors := ValidateCIDRNetworkAddress(v, "cidr_network_address")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid cidr network address: %q", v, errors)
		}
	}

	invalidCIDRNetworkAddress := []string{"1.2.3.4", "0x38732/21"}
	for _, v := range invalidCIDRNetworkAddress {
		_, errors := ValidateCIDRNetworkAddress(v, "cidr_network_address")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid cidr network address", v)
		}
	}
}

func TestValidateSwitchCIDRNetworkAddress(t *testing.T) {
	validSwitchCIDRNetworkAddress := []string{"192.168.10.0/24", "0.0.0.0/16", "127.0.0.0/29", "10.121.10.0/24"}
	for _, v := range validSwitchCIDRNetworkAddress {
		_, errors := ValidateSwitchCIDRNetworkAddress(v, "switch_cidr_network_address")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid switch cidr network address: %q", v, errors)
		}
	}

	invalidSwitchCIDRNetworkAddress := []string{"1.2.3.4", "0x38732/21", "10.121.10.0/15", "10.121.10.0/30", "256.121.10.0/22"}
	for _, v := range invalidSwitchCIDRNetworkAddress {
		_, errors := ValidateSwitchCIDRNetworkAddress(v, "switch_cidr_network_address")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid switch cidr network address", v)
		}
	}
}

func TestValidateAllowedStringSplitValue(t *testing.T) {
	exceptValues := []string{"aliyun", "alicloud", "alibaba"}
	validValues := "aliyun,alicloud"
	_, errors := ValidateAllowedSplitStringValue(exceptValues, ",")(validValues, "allowvalue")
	if len(errors) != 0 {
		t.Fatalf("%q should be a valid value in %#v: %q", validValues, exceptValues, errors)
	}

	invalidValues := "ali,alidata"
	_, invalidErr := ValidateAllowedSplitStringValue(exceptValues, ",")(invalidValues, "allowvalue")
	if len(invalidErr) == 0 {
		t.Fatalf("%q should be an invalid value", invalidValues)
	}
}

func TestValidateStringConvertInt64(t *testing.T) {
	_, errors := ValidateStringConvertInt64()(math.MaxInt64, "name")
	if len(errors) == 0 {
		t.Fatalf("%q cannot convert to int64", 110)
	}

	_, errors = ValidateStringConvertInt64()("abcd", "name")
	if len(errors) == 0 {
		t.Fatalf("%q cannot convert to int64", "abcd")
	}

	_, errors = ValidateStringConvertInt64()("6666666", "name")
	if len(errors) != 0 {
		t.Fatalf("%q can convert to int64", "666666")
	}

	_, errors = ValidateStringConvertInt64()(strconv.FormatInt(math.MaxInt64, 10), "name")
	if len(errors) != 0 {
		t.Fatalf("%q can convert to int64", math.MaxInt64)
	}

}

func resetSkipResourceSchemaValidationForTest() {
	skipResourceSchemaValidation = false
}

func TestValidateIntBetween(t *testing.T) {
	SetSkipResourceSchemaValidation(false)
	defer resetSkipResourceSchemaValidationForTest()

	min := 0
	max := 10
	check := 1
	attribute := "port"
	s, es := IntBetween(min, max)(check, attribute)
	assert.Nil(t, es)
	assert.Nil(t, s)

	check = 20
	s, es = IntBetween(min, max)(check, attribute)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be in the range (%d - %d), got %d %s", attribute, min, max, check, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = IntBetween(min, max)(check, attribute)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be in the range (%d - %d), got %d", attribute, min, max, check)))
}

func TestValidateNullableIntBetween(t *testing.T) {
	SetSkipResourceSchemaValidation(false)
	defer resetSkipResourceSchemaValidationForTest()

	attribute := "size"
	largeRange := ValidateNullableIntBetween(0, 9007199254740991)
	smallRange := ValidateNullableIntBetween(0, 50)

	// valid values pass regardless of range
	for _, v := range []string{"", "0", "10", "9007199254740991"} {
		s, es := largeRange(v, attribute)
		assert.Nil(t, es, "value %q should be valid", v)
		assert.Nil(t, s, "value %q should not produce warnings", v)
	}

	// format errors are hard errors regardless of skip flag
	hardErrorValues := []string{"abc", "+10", "010", " 10", "1.5", "-0"}
	for _, skip := range []bool{false, true} {
		SetSkipResourceSchemaValidation(skip)
		for _, v := range hardErrorValues {
			s, es := largeRange(v, attribute)
			assert.NotNil(t, es, "value %q should be a hard error (skip=%v)", v, skip)
			assert.Nil(t, s, "value %q should not produce warnings (skip=%v)", v, skip)
		}
	}

	// range violation with skip disabled: hard error carrying the skip hint
	SetSkipResourceSchemaValidation(false)
	for _, v := range []string{"-1", "51"} {
		s, es := smallRange(v, attribute)
		assert.NotNil(t, es, "value %q should be a range error", v)
		assert.Nil(t, s)
		assert.True(t, strings.Contains(es[0].Error(), skipResourceSchemaValidationWarning))
		assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
			"expected %s to be in the range (%d - %d), got %s", attribute, 0, 50, v)))
	}

	// range violation with skip enabled: downgraded to a warning
	SetSkipResourceSchemaValidation(true)
	for _, v := range []string{"-1", "51"} {
		s, es := smallRange(v, attribute)
		assert.Nil(t, es, "value %q should be downgraded to a warning", v)
		assert.NotNil(t, s)
		assert.True(t, strings.Contains(s[0], fmt.Sprintf(
			"expected %s to be in the range (%d - %d), got %s", attribute, 0, 50, v)))
	}

	// canonical form check is never downgraded even with skip enabled
	s, es := smallRange("+10", attribute)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), "canonical decimal integer string"))

	// non-string input is a hard error
	SetSkipResourceSchemaValidation(false)
	s, es = smallRange(10, attribute)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf("expected type of %s to be string", attribute)))
}

func TestValidateNullableBool(t *testing.T) {
	SetSkipResourceSchemaValidation(false)
	defer resetSkipResourceSchemaValidationForTest()

	attribute := "enabled"

	// valid values
	for _, v := range []string{"", "true", "false"} {
		s, es := ValidateNullableBool(v, attribute)
		assert.Nil(t, es, "value %q should be valid", v)
		assert.Nil(t, s, "value %q should not produce warnings", v)
	}

	// invalid values with skip disabled: error carrying the skip hint
	for _, v := range []string{"True", "FALSE", "1", "0", "yes"} {
		s, es := ValidateNullableBool(v, attribute)
		assert.NotNil(t, es, "value %q should be an error", v)
		assert.Nil(t, s)
		assert.True(t, strings.Contains(es[0].Error(), skipResourceSchemaValidationWarning))
		assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
			"expected %s to be an empty string or lowercase \"true\"/\"false\", got %s", attribute, v)))
	}

	// invalid value with skip enabled: downgraded to a warning
	SetSkipResourceSchemaValidation(true)
	s, es := ValidateNullableBool("True", attribute)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be an empty string or lowercase \"true\"/\"false\", got %s", attribute, "True")))

	// non-string input is a hard error even with skip enabled
	s, es = ValidateNullableBool(1, attribute)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf("expected type of %s to be string", attribute)))
}

func TestValidateIntAtLeast(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := 10
	v := 10
	k := "port"
	s, es := IntAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 2
	s, es = IntAtLeast(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be at least (%d), got %d %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = IntAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at least (%d), got %d", k, limit, v)))
}

func TestValidateIntAtMost(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := 10
	v := 10
	k := "port"
	s, es := IntAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 20
	s, es = IntAtMost(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be at most (%d), got %d %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = IntAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at most (%d), got %d", k, limit, v)))
}

func TestValidateIntInSlice(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := []int{1, 3, 6, 9, 12}
	v := 6
	k := "period"
	s, es := IntInSlice(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 10
	s, es = IntInSlice(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be one of %v, got %d %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = IntInSlice(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be one of %v, got %d", k, limit, v)))
}

func TestValidateStringInSlice(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := []string{"IPV4", "IPV6"}
	v := "IPV4"
	k := "ip_type"
	s, es := StringInSlice(limit, false)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "IPV4"
	s, es = StringInSlice(limit, true)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "IPV5"
	s, es = StringInSlice(limit, false)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be one of %v, got %s %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringInSlice(limit, false)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be one of %v, got %s", k, limit, v)))
}

func TestValidateStringLenBetween(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	min := 2
	max := 10
	v := "tf-testAcc"
	k := "name"
	s, es := StringLenBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "tf-testAcc-update"
	s, es = StringLenBetween(min, max)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected length of %s to be in the range (%d - %d), got %s %s", k, min, max, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringLenBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected length of %s to be in the range (%d - %d), got %s", k, min, max, v)))
}

func TestValidateStringMatch(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.-]{1,127}$`)
	message := "The name must be 2 to 128 characters in length, and must start with a letter. It can contain digits, periods (.), underscores (_), and hyphens (-)."
	v := "tf-testAcc"
	k := "name"
	s, es := StringMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "tf-#testAcc-update"
	s, es = StringMatch(limit, message)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"invalid value for %s (%s)", k, message)))

	SetSkipResourceSchemaValidation(false)
	v = "tf-#testAcc-update"
	s, es = StringMatch(limit, "")(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected value of %s to match regular expression %q %s", k, limit, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringMatch(limit, "")(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to match regular expression %q", k, limit)))
}

func TestValidateStringDoesNotMatch(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := regexp.MustCompile(`(^http://.*)|(^https://.*)`)
	message := "It must cannot start with `https://`, `https://`"
	v := "tf-testAcc"
	k := "name"
	s, es := StringDoesNotMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "http://tf-testAcc-update.com"
	s, es = StringDoesNotMatch(limit, message)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"invalid value for %s (%s) %s", k, message, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringDoesNotMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"invalid value for %s (%s)", k, message)))

	SetSkipResourceSchemaValidation(false)
	v = "http://tf-testAcc-update.com"
	s, es = StringDoesNotMatch(limit, "")(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected value of %s to not match regular expression %q %s", k, limit, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringDoesNotMatch(limit, "")(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to not match regular expression %q", k, limit)))
}

func TestValidateFloatBetween(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	min := 2.1
	max := 10.9
	v := 5.3
	k := "weight"
	s, es := FloatBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 1.1
	s, es = FloatBetween(min, max)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be in the range (%f - %f), got %f %s", k, min, max, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = FloatBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be in the range (%f - %f), got %f", k, min, max, v)))
}

func TestValidateFloatAtLeast(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := 2.1
	v := 5.3
	k := "weight"
	s, es := FloatAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 1
	s, es = FloatAtLeast(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be at least (%f), got %f %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = FloatAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at least (%f), got %f", k, limit, v)))
}

func TestValidateFloatAtMost(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := 10.9
	v := 5.3
	k := "weight"
	s, es := FloatAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = 20
	s, es = FloatAtMost(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected %s to be at most (%f), got %f %s", k, limit, v, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = FloatAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at most (%f), got %f", k, limit, v)))
}

func TestValidateStringDoesNotContainAny(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	limit := ",."
	v := "tf-testAcc"
	k := "name"
	s, es := StringDoesNotContainAny(limit)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "tf-test,Acce"
	s, es = StringDoesNotContainAny(limit)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected value of %s to not contain any of %q %s", k, limit, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = StringDoesNotContainAny(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to not contain any of %q", k, limit)))
}

func TestValidateRFC3339TimeString(t *testing.T) {
	SetSkipResourceSchemaValidation(false)

	v := "2023-07-26T23:50:22Z"
	k := "auto_release_time"
	s, es := ValidateRFC3339TimeString(false)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "2023-07-26T23:50:22Z07:09"
	s, es = ValidateRFC3339TimeString(false)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"%q: invalid RFC3339 timestamp %s", k, skipResourceSchemaValidationWarning)))

	v = ""
	s, es = ValidateRFC3339TimeString(false)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"%q: invalid RFC3339 timestamp %s", k, skipResourceSchemaValidationWarning)))

	s, es = ValidateRFC3339TimeString(true)(v, k)
	assert.Nil(t, es)
	assert.Nil(t, s)

	v = "2023-07-26T23:50Z"
	s, es = ValidateRFC3339TimeString(false)(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"%q: invalid RFC3339 timestamp %s", k, skipResourceSchemaValidationWarning)))

	SetSkipResourceSchemaValidation(true)
	s, es = ValidateRFC3339TimeString(false)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"%q: invalid RFC3339 timestamp", k)))
}

// ─── TestUnitCommon* ─────────────────────────────────────────────────────────

func TestUnitCommonNormalizeYamlString(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		expectEmpty bool
		description string
	}{
		{
			name:        "Nil_Input",
			input:       nil,
			expectError: false,
			expectEmpty: true,
			description: "nil input should return empty string without error",
		},
		{
			name:        "Empty_String",
			input:       "",
			expectError: false,
			expectEmpty: true,
			description: "empty string should return empty string without error",
		},
		{
			name:        "Valid_Simple_YAML",
			input:       "key: value",
			expectError: false,
			description: "valid simple YAML should normalize without error",
		},
		{
			name:        "Valid_Multi_Key_YAML",
			input:       "name: alice\nage: 30",
			expectError: false,
			description: "valid multi-key YAML should normalize without error",
		},
		{
			name:        "Invalid_YAML",
			input:       "key: :\n  bad: [unclosed",
			expectError: true,
			description: "invalid YAML should return error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NormalizeYamlString(tc.input)
			if tc.expectError {
				assert.Error(t, err, tc.description)
			} else {
				assert.NoError(t, err, tc.description)
				if tc.expectEmpty {
					assert.Equal(t, "", result)
				}
			}
		})
	}
}

func TestUnitCommonNormalizeJsonString(t *testing.T) {
	testCases := []struct {
		name        string
		input       interface{}
		expectError bool
		expectEmpty bool
		description string
	}{
		{
			name:        "Nil_Input",
			input:       nil,
			expectError: false,
			expectEmpty: true,
			description: "nil input should return empty string without error",
		},
		{
			name:        "Empty_String",
			input:       "",
			expectError: false,
			expectEmpty: true,
			description: "empty string should return empty string without error",
		},
		{
			name:        "Valid_Object",
			input:       `{"key":"value","num":42}`,
			expectError: false,
			description: "valid JSON object should normalize without error",
		},
		{
			name:        "Valid_Array",
			input:       `["a","b","c"]`,
			expectError: false,
			description: "valid JSON array should normalize without error",
		},
		{
			name:        "Valid_Unordered_Keys",
			input:       `{"b":2,"a":1}`,
			expectError: false,
			description: "JSON with unordered keys should normalize without error",
		},
		{
			name:        "Invalid_JSON",
			input:       `{"key": value_without_quotes}`,
			expectError: true,
			description: "invalid JSON should return error",
		},
		{
			name:        "Invalid_Truncated",
			input:       `{"key":`,
			expectError: true,
			description: "truncated JSON should return error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NormalizeJsonString(tc.input)
			if tc.expectError {
				assert.Error(t, err, tc.description)
				assert.NotEmpty(t, result, "on error the original string should be returned")
			} else {
				assert.NoError(t, err, tc.description)
				if tc.expectEmpty {
					assert.Equal(t, "", result)
				} else {
					assert.NotEmpty(t, result)
				}
			}
		})
	}
}

func TestUnitCommonValidateYamlString(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_YAML",
			value:       "key: value",
			expectError: false,
			description: "valid YAML should produce no error",
		},
		{
			name:        "Valid_Empty_String",
			value:       "",
			expectError: false,
			description: "empty string is allowed (normalizes to empty)",
		},
		{
			name:        "Valid_Complex_YAML",
			value:       "server:\n  host: localhost\n  port: 8080",
			expectError: false,
			description: "complex YAML should be valid",
		},
		{
			name:        "Invalid_YAML",
			value:       "key: :\n  bad: [unclosed",
			expectError: true,
			description: "invalid YAML should produce an error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := ValidateYamlString(tc.value, "yaml_field")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
				assert.Contains(t, errors[0].Error(), "invalid YAML")
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}

func TestUnitCommonValidateNormalName(t *testing.T) {
	testCases := []struct {
		name        string
		value       string
		expectError bool
		description string
	}{
		{
			name:        "Valid_Min_Length",
			value:       "ab",
			expectError: false,
			description: "name of exactly 2 chars is valid",
		},
		{
			name:        "Valid_Max_Length",
			value:       strings.Repeat("a", 128),
			expectError: false,
			description: "name of exactly 128 chars is valid",
		},
		{
			name:        "Valid_Normal",
			value:       "my-resource-name",
			expectError: false,
			description: "normal name should be valid",
		},
		{
			name:        "Invalid_Too_Short",
			value:       "a",
			expectError: true,
			description: "name of 1 char is below minimum (2)",
		},
		{
			name:        "Invalid_Too_Long",
			value:       strings.Repeat("a", 129),
			expectError: true,
			description: "name of 129 chars exceeds maximum (128)",
		},
		{
			name:        "Invalid_HTTP_Prefix",
			value:       "http://example.com",
			expectError: true,
			description: "name starting with http:// is invalid",
		},
		{
			name:        "Invalid_HTTPS_Prefix",
			value:       "https://example.com",
			expectError: true,
			description: "name starting with https:// is invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ws, errors := ValidateNormalName(tc.value, "name")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
			}
		})
	}
}

func TestUnitCommonStringLenAtLeast(t *testing.T) {
	// Reset to default (no-skip) before each subtest.
	t.Run("Error_Mode", func(t *testing.T) {
		SetSkipResourceSchemaValidation(false)
		validate := StringLenAtLeast(5)

		testCases := []struct {
			name        string
			value       interface{}
			expectError bool
			expectWarn  bool
			description string
		}{
			{
				name:        "Valid_Exact_Min",
				value:       "hello",
				expectError: false,
				description: "string of exactly min length should be valid",
			},
			{
				name:        "Valid_Above_Min",
				value:       "hello world",
				expectError: false,
				description: "string longer than min should be valid",
			},
			{
				name:        "Invalid_Too_Short",
				value:       "hi",
				expectError: true,
				description: "string shorter than min should produce an error",
			},
			{
				name:        "Invalid_Empty",
				value:       "",
				expectError: true,
				description: "empty string is shorter than min",
			},
			{
				name:        "Invalid_Only_Whitespace",
				value:       "   ",
				expectError: true,
				description: "whitespace-only string trims to empty → shorter than min",
			},
			{
				name:        "Invalid_Wrong_Type",
				value:       12345,
				expectError: true,
				description: "non-string type should produce an error",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				ws, es := validate(tc.value, "field")
				if tc.expectError {
					assert.NotEmpty(t, es, tc.description)
				} else {
					assert.Empty(t, es, tc.description)
					assert.Empty(t, ws)
				}
			})
		}
	})

	t.Run("Warn_Mode", func(t *testing.T) {
		SetSkipResourceSchemaValidation(true)
		defer resetSkipResourceSchemaValidationForTest()
		validate := StringLenAtLeast(5)

		ws, es := validate("hi", "field")
		assert.Empty(t, es, "In skip mode, errors should be downgraded to warnings")
		assert.NotEmpty(t, ws, "In skip mode, violation should produce a warning")
		assert.Contains(t, ws[0], "expected length of")
	})
}

// TestSkipResourceSchemaValidationSwitch pins the contract of the switch behind
// skip_resource_schema_validation: Provider() seeds it from
// TF_SKIP_RESOURCE_SCHEMA_VALIDATION, and the provider argument's validator
// overwrites it during PrepareProviderConfig. The value is live, not a latch.
func TestSkipResourceSchemaValidationSwitch(t *testing.T) {
	t.Cleanup(resetSkipResourceSchemaValidationForTest)

	t.Run("not a latch: publishing again takes effect", func(t *testing.T) {
		resetSkipResourceSchemaValidationForTest()
		SetSkipResourceSchemaValidation(true)
		assert.True(t, SkipResourceSchemaValidation())

		SetSkipResourceSchemaValidation(false)
		assert.False(t, SkipResourceSchemaValidation(), "an explicit false must re-enable the errors")
	})

	t.Run("switch controls the validators", func(t *testing.T) {
		resetSkipResourceSchemaValidationForTest()
		SetSkipResourceSchemaValidation(false)
		_, es := IntBetween(0, 10)(20, "port")
		assert.NotEmpty(t, es)

		SetSkipResourceSchemaValidation(true)
		ws, es := IntBetween(0, 10)(20, "port")
		assert.Empty(t, es, "the violation must be downgraded to a warning")
		assert.NotEmpty(t, ws)
	})
}
