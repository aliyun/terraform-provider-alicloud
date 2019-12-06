package alicloud

import (
	"math"
	"strconv"
	"testing"
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
