// lintignore: S013
package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonPrePaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		chargeType  string
		paymentType string
		expected    bool
		description string
	}{
		{
			name:        "Prepaid_ChargeType",
			chargeType:  "Prepaid",
			expected:    true,
			description: "Prepaid charge type should suppress diff",
		},
		{
			name:        "Subscription_PaymentType",
			paymentType: "Subscription",
			expected:    true,
			description: "Subscription payment type should suppress diff",
		},
		{
			name:        "PostPaid_ChargeType",
			chargeType:  "PostPaid",
			expected:    false,
			description: "PostPaid charge type should not suppress diff",
		},
		{
			name:        "PayAsYouGo_PaymentType",
			paymentType: "PayAsYouGo",
			expected:    false,
			description: "PayAsYouGo payment type should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.chargeType != "" {
				data["instance_charge_type"] = tc.chargeType
			}
			if tc.paymentType != "" {
				data["payment_type"] = tc.paymentType
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
				"payment_type":         {Type: schema.TypeString},
			}, data)

			result := PrePaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsNotAutoRenewDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name               string
		instanceChargeType string
		renewalStatus      string
		expected           bool
		description        string
	}{
		{
			name:               "PostPaid_ChargeType",
			instanceChargeType: "PostPaid",
			expected:           true,
			description:        "PostPaid should suppress diff",
		},
		{
			name:               "PrePaid_With_AutoRenewal",
			instanceChargeType: "PrePaid",
			renewalStatus:      "AutoRenewal",
			expected:           false,
			description:        "PrePaid with auto renewal should not suppress diff",
		},
		{
			name:               "PrePaid_Without_AutoRenewal",
			instanceChargeType: "PrePaid",
			renewalStatus:      "Normal",
			expected:           true,
			description:        "PrePaid without auto renewal should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
				"renewal_status":       {Type: schema.TypeString},
			}, map[string]interface{}{
				"instance_charge_type": tc.instanceChargeType,
				"renewal_status":       tc.renewalStatus,
			})

			result := ecsNotAutoRenewDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsSpotStrategyDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name               string
		instanceChargeType string
		hasChargeType      bool
		expected           bool
		description        string
	}{
		{
			name:               "PostPaid_ChargeType",
			instanceChargeType: "PostPaid",
			hasChargeType:      true,
			expected:           false,
			description:        "PostPaid should not suppress diff",
		},
		{
			name:          "No_ChargeType",
			hasChargeType: false,
			expected:      false,
			description:   "No charge type should not suppress diff",
		},
		{
			name:               "PrePaid_ChargeType",
			instanceChargeType: "PrePaid",
			hasChargeType:      true,
			expected:           true,
			description:        "PrePaid should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.hasChargeType {
				data["instance_charge_type"] = tc.instanceChargeType
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
			}, data)

			result := ecsSpotStrategyDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsSpotPriceLimitDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name               string
		instanceChargeType string
		spotStrategy       string
		expected           bool
		description        string
	}{
		{
			name:               "PostPaid_With_SpotWithPriceLimit",
			instanceChargeType: "PostPaid",
			spotStrategy:       "SpotWithPriceLimit",
			expected:           false,
			description:        "PostPaid with SpotWithPriceLimit should not suppress diff",
		},
		{
			name:               "PostPaid_Without_SpotWithPriceLimit",
			instanceChargeType: "PostPaid",
			spotStrategy:       "NoSpot",
			expected:           true,
			description:        "PostPaid without SpotWithPriceLimit should suppress diff",
		},
		{
			name:               "PrePaid",
			instanceChargeType: "PrePaid",
			spotStrategy:       "SpotWithPriceLimit",
			expected:           true,
			description:        "PrePaid should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
				"spot_strategy":        {Type: schema.TypeString},
			}, map[string]interface{}{
				"instance_charge_type": tc.instanceChargeType,
				"spot_strategy":        tc.spotStrategy,
			})

			result := ecsSpotPriceLimitDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsSystemDiskPerformanceLevelSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		diskCategory string
		expected     bool
		description  string
	}{
		{
			name:         "CloudESSD",
			diskCategory: "cloud_essd",
			expected:     false,
			description:  "cloud_essd should not suppress diff",
		},
		{
			name:         "Cloud",
			diskCategory: "cloud",
			expected:     true,
			description:  "Non-cloud_essd should suppress diff",
		},
		{
			name:         "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:     true,
			description:  "cloud_ssd should suppress diff",
		},
		{
			name:         "Empty",
			diskCategory: "",
			expected:     true,
			description:  "Empty disk category should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"system_disk_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"system_disk_category": tc.diskCategory,
			})
			result := ecsSystemDiskPerformanceLevelSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsSecurityGroupRulePortRangeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		ipProtocol  string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "TCP_AllPortRange",
			ipProtocol:  "tcp",
			new:         "-1/-1",
			expected:    true,
			description: "TCP with AllPortRange new value should suppress diff",
		},
		{
			name:        "TCP_SpecificRange",
			ipProtocol:  "tcp",
			new:         "80/80",
			expected:    false,
			description: "TCP with specific port range should not suppress diff",
		},
		{
			name:        "UDP_AllPortRange",
			ipProtocol:  "udp",
			new:         "-1/-1",
			expected:    true,
			description: "UDP with AllPortRange new value should suppress diff",
		},
		{
			name:        "UDP_SpecificRange",
			ipProtocol:  "udp",
			new:         "53/53",
			expected:    false,
			description: "UDP with specific port range should not suppress diff",
		},
		{
			name:        "ICMP_AllPortRange",
			ipProtocol:  "icmp",
			new:         "-1/-1",
			expected:    false,
			description: "Non-TCP/UDP with AllPortRange new value should not suppress diff",
		},
		{
			name:        "ICMP_OtherRange",
			ipProtocol:  "icmp",
			new:         "80/80",
			expected:    true,
			description: "Non-TCP/UDP with non-AllPortRange should suppress diff",
		},
		{
			name:        "ALL_AllPortRange",
			ipProtocol:  "all",
			new:         "-1/-1",
			expected:    false,
			description: "All protocol with AllPortRange should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"ip_protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"ip_protocol": tc.ipProtocol,
			})
			result := ecsSecurityGroupRulePortRangeDiffSuppressFunc("key", "old", tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsSecurityGroupRulePreFixListIdDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                  string
		cidrIp                string
		ipv6CidrIp            string
		sourceSecurityGroupId string
		expected              bool
		description           string
	}{
		{
			name:        "With_CidrIp",
			cidrIp:      "192.168.0.0/24",
			expected:    true,
			description: "With cidr_ip should suppress diff",
		},
		{
			name:        "With_Ipv6CidrIp",
			ipv6CidrIp:  "2001:db8::/32",
			expected:    true,
			description: "With ipv6_cidr_ip should suppress diff",
		},
		{
			name:                  "With_SourceSecurityGroupId",
			sourceSecurityGroupId: "sg-12345",
			expected:              true,
			description:           "With source_security_group_id should suppress diff",
		},
		{
			name:        "No_Source",
			expected:    false,
			description: "Without any source should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.cidrIp != "" {
				data["cidr_ip"] = tc.cidrIp
			}
			if tc.ipv6CidrIp != "" {
				data["ipv6_cidr_ip"] = tc.ipv6CidrIp
			}
			if tc.sourceSecurityGroupId != "" {
				data["source_security_group_id"] = tc.sourceSecurityGroupId
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"cidr_ip":                  {Type: schema.TypeString},
				"ipv6_cidr_ip":             {Type: schema.TypeString},
				"source_security_group_id": {Type: schema.TypeString},
			}, data)
			result := ecsSecurityGroupRulePreFixListIdDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonEcsInternetDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		maxBandwidth int
		expected     bool
		description  string
	}{
		{
			name:         "Bandwidth_Positive",
			maxBandwidth: 100,
			expected:     false,
			description:  "Positive bandwidth should not suppress diff",
		},
		{
			name:         "Bandwidth_One",
			maxBandwidth: 1,
			expected:     false,
			description:  "Bandwidth=1 should not suppress diff",
		},
		{
			name:         "Bandwidth_Zero",
			maxBandwidth: 0,
			expected:     true,
			description:  "Zero bandwidth (GetOk returns false for int 0) should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"internet_max_bandwidth_out": {Type: schema.TypeInt},
			}, map[string]interface{}{
				"internet_max_bandwidth_out": tc.maxBandwidth,
			})
			result := ecsInternetDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
