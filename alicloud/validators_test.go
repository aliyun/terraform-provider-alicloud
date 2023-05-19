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
