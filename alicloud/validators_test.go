package alicloud

import (
	"bytes"
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

func TestValidateCIDRNetworkAddress(t *testing.T) {
	validCIDRNetworkAddress := []string{"192.168.10.0/24", "0.0.0.0/0", "10.121.10.0/24"}
	for _, v := range validCIDRNetworkAddress {
		_, errors := validateCIDRNetworkAddress(v, "cidr_network_address")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid cidr network address: %q", v, errors)
		}
	}

	invalidCIDRNetworkAddress := []string{"1.2.3.4", "0x38732/21"}
	for _, v := range invalidCIDRNetworkAddress {
		_, errors := validateCIDRNetworkAddress(v, "cidr_network_address")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid cidr network address", v)
		}
	}
}

func TestValidateSwitchCIDRNetworkAddress(t *testing.T) {
	validSwitchCIDRNetworkAddress := []string{"192.168.10.0/24", "0.0.0.0/16", "127.0.0.0/29", "10.121.10.0/24"}
	for _, v := range validSwitchCIDRNetworkAddress {
		_, errors := validateSwitchCIDRNetworkAddress(v, "switch_cidr_network_address")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid switch cidr network address: %q", v, errors)
		}
	}

	invalidSwitchCIDRNetworkAddress := []string{"1.2.3.4", "0x38732/21", "10.121.10.0/15", "10.121.10.0/30", "256.121.10.0/22"}
	for _, v := range invalidSwitchCIDRNetworkAddress {
		_, errors := validateSwitchCIDRNetworkAddress(v, "switch_cidr_network_address")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid switch cidr network address", v)
		}
	}
}

func TestValidateAllowedStringSplitValue(t *testing.T) {
	exceptValues := []string{"aliyun", "alicloud", "alibaba"}
	validValues := "aliyun,alicloud"
	_, errors := validateAllowedSplitStringValue(exceptValues, ",")(validValues, "allowvalue")
	if len(errors) != 0 {
		t.Fatalf("%q should be a valid value in %#v: %q", validValues, exceptValues, errors)
	}

	invalidValues := "ali,alidata"
	_, invalidErr := validateAllowedSplitStringValue(exceptValues, ",")(invalidValues, "allowvalue")
	if len(invalidErr) == 0 {
		t.Fatalf("%q should be an invalid value", invalidValues)
	}
}

func TestValidateStringConvertInt64(t *testing.T) {
	_, errors := validateStringConvertInt64()(math.MaxInt64, "name")
	if len(errors) == 0 {
		t.Fatalf("%q cannot convert to int64", 110)
	}

	_, errors = validateStringConvertInt64()("abcd", "name")
	if len(errors) == 0 {
		t.Fatalf("%q cannot convert to int64", "abcd")
	}

	_, errors = validateStringConvertInt64()("6666666", "name")
	if len(errors) != 0 {
		t.Fatalf("%q can convert to int64", "666666")
	}

	_, errors = validateStringConvertInt64()(strconv.FormatInt(math.MaxInt64, 10), "name")
	if len(errors) != 0 {
		t.Fatalf("%q can convert to int64", math.MaxInt64)
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

func TestValidateIntBetween(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = IntBetween(min, max)(check, attribute)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be in the range (%d - %d), got %d", attribute, min, max, check)))
}

func TestValidateIntAtLeast(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = IntAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at least (%d), got %d", k, limit, v)))
}

func TestValidateIntAtMost(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = IntAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at most (%d), got %d", k, limit, v)))
}

func TestValidateIntInSlice(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = IntInSlice(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be one of %v, got %d", k, limit, v)))
}

func TestValidateStringInSlice(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = StringInSlice(limit, false)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be one of %v, got %s", k, limit, v)))
}

func TestValidateStringLenBetween(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = StringLenBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected length of %s to be in the range (%d - %d), got %s", k, min, max, v)))
}

func TestValidateStringMatch(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = StringMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"invalid value for %s (%s)", k, message)))

	resourceSchemaValidationSkipped = false
	v = "tf-#testAcc-update"
	s, es = StringMatch(limit, "")(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected value of %s to match regular expression %q %s", k, limit, skipResourceSchemaValidationWarning)))

	resourceSchemaValidationSkipped = true
	s, es = StringMatch(limit, "")(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to match regular expression %q", k, limit)))
}

func TestValidateStringDoesNotMatch(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = StringDoesNotMatch(limit, message)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"invalid value for %s (%s)", k, message)))

	resourceSchemaValidationSkipped = false
	v = "http://tf-testAcc-update.com"
	s, es = StringDoesNotMatch(limit, "")(v, k)
	assert.NotNil(t, es)
	assert.Nil(t, s)
	assert.True(t, strings.Contains(es[0].Error(), fmt.Sprintf(
		"expected value of %s to not match regular expression %q %s", k, limit, skipResourceSchemaValidationWarning)))

	resourceSchemaValidationSkipped = true
	s, es = StringDoesNotMatch(limit, "")(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to not match regular expression %q", k, limit)))
}

func TestValidateFloatBetween(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = FloatBetween(min, max)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be in the range (%f - %f), got %f", k, min, max, v)))
}

func TestValidateFloatAtLeast(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = FloatAtLeast(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at least (%f), got %f", k, limit, v)))
}

func TestValidateFloatAtMost(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = FloatAtMost(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected %s to be at most (%f), got %f", k, limit, v)))
}

func TestValidateStringDoesNotContainAny(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = StringDoesNotContainAny(limit)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"expected value of %s to not contain any of %q", k, limit)))
}

func TestValidateRFC3339TimeString(t *testing.T) {
	resourceSchemaValidationSkipped = false

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

	resourceSchemaValidationSkipped = true
	s, es = ValidateRFC3339TimeString(false)(v, k)
	assert.Nil(t, es)
	assert.NotNil(t, s)
	assert.True(t, strings.Contains(s[0], fmt.Sprintf(
		"%q: invalid RFC3339 timestamp", k)))
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

func TestUnitCommonValidateOnsGroupId(t *testing.T) {
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
			result, err := normalizeYamlString(tc.input)
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
			result, err := normalizeJsonString(tc.input)
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
			ws, errors := validateYamlString(tc.value, "yaml_field")
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
			ws, errors := validateNormalName(tc.value, "name")
			assert.Empty(t, ws)
			if tc.expectError {
				assertHasErrors(t, ws, errors, tc.description)
			} else {
				assertNoErrors(t, ws, errors, tc.description)
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

func TestUnitCommonStringLenAtLeast(t *testing.T) {
	// Reset to default (no-skip) before each subtest.
	t.Run("Error_Mode", func(t *testing.T) {
		resourceSchemaValidationSkipped = false
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
		resourceSchemaValidationSkipped = true
		defer func() { resourceSchemaValidationSkipped = false }()
		validate := StringLenAtLeast(5)

		ws, es := validate("hi", "field")
		assert.Empty(t, es, "In skip mode, errors should be downgraded to warnings")
		assert.NotEmpty(t, ws, "In skip mode, violation should produce a warning")
		assert.Contains(t, ws[0], "expected length of")
	})
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
