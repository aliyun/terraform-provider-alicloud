package alicloud

import (
	"bytes"
	"math"
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
