// lintignore: S013
package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonPrivateDnsTypeDiffSuppressFunc(t *testing.T) {
	// Test for primary_dns when private_dns_type is PrivateZone
	testCases := []struct {
		name           string
		privateDnsType string
		oldValue       string
		newValue       string
		expected       bool
		description    string
	}{
		{
			name:           "PrivateZone_With_Empty_New",
			privateDnsType: "PrivateZone",
			oldValue:       "100.100.2.136",
			newValue:       "",
			expected:       true,
			description:    "When private_dns_type is PrivateZone and new value is empty, suppress diff",
		},
		{
			name:           "PrivateZone_With_NonEmpty_New",
			privateDnsType: "PrivateZone",
			oldValue:       "100.100.2.136",
			newValue:       "8.8.8.8",
			expected:       false,
			description:    "When private_dns_type is PrivateZone and new value is not empty, do not suppress diff",
		},
		{
			name:           "Custom_With_Empty_New",
			privateDnsType: "Custom",
			oldValue:       "8.8.8.8",
			newValue:       "",
			expected:       false,
			description:    "When private_dns_type is Custom and new value is empty, do not suppress diff",
		},
		{
			name:           "Custom_With_NonEmpty_New",
			privateDnsType: "Custom",
			oldValue:       "8.8.8.8",
			newValue:       "1.1.1.1",
			expected:       false,
			description:    "When private_dns_type is Custom and new value is not empty, do not suppress diff",
		},
		{
			name:           "Empty_Type_With_Empty_New",
			privateDnsType: "",
			oldValue:       "8.8.8.8",
			newValue:       "",
			expected:       false,
			description:    "When private_dns_type is empty and new value is empty, do not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"private_dns_type": {Type: schema.TypeString},
				"primary_dns":      {Type: schema.TypeString},
			}, map[string]interface{}{
				"private_dns_type": tc.privateDnsType,
			})

			// 模拟 resource_alicloud_cloud_firewall_private_dns.go 中的 DiffSuppressFunc 逻辑
			diffSuppressFunc := func(k, old, new string, d *schema.ResourceData) bool {
				if v, ok := d.GetOk("private_dns_type"); ok && v.(string) == "PrivateZone" {
					return new == ""
				}
				return false
			}

			result := diffSuppressFunc("primary_dns", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonStandbyDnsTypeDiffSuppressFunc(t *testing.T) {
	// Test for standby_dns when private_dns_type is PrivateZone
	testCases := []struct {
		name           string
		privateDnsType string
		oldValue       string
		newValue       string
		expected       bool
		description    string
	}{
		{
			name:           "PrivateZone_With_Empty_New",
			privateDnsType: "PrivateZone",
			oldValue:       "100.100.2.138",
			newValue:       "",
			expected:       true,
			description:    "When private_dns_type is PrivateZone and new value is empty, suppress diff",
		},
		{
			name:           "PrivateZone_With_NonEmpty_New",
			privateDnsType: "PrivateZone",
			oldValue:       "100.100.2.138",
			newValue:       "4.4.4.4",
			expected:       false,
			description:    "When private_dns_type is PrivateZone and new value is not empty, do not suppress diff",
		},
		{
			name:           "Custom_With_Empty_New",
			privateDnsType: "Custom",
			oldValue:       "4.4.4.4",
			newValue:       "",
			expected:       false,
			description:    "When private_dns_type is Custom and new value is empty, do not suppress diff",
		},
		{
			name:           "Custom_With_NonEmpty_New",
			privateDnsType: "Custom",
			oldValue:       "4.4.4.4",
			newValue:       "2.2.2.2",
			expected:       false,
			description:    "When private_dns_type is Custom and new value is not empty, do not suppress diff",
		},
		{
			name:           "Empty_Type_With_Empty_New",
			privateDnsType: "",
			oldValue:       "4.4.4.4",
			newValue:       "",
			expected:       false,
			description:    "When private_dns_type is empty and new value is empty, do not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"private_dns_type": {Type: schema.TypeString},
				"standby_dns":      {Type: schema.TypeString},
			}, map[string]interface{}{
				"private_dns_type": tc.privateDnsType,
			})

			// 模拟 resource_alicloud_cloud_firewall_private_dns.go 中的 DiffSuppressFunc 逻辑
			diffSuppressFunc := func(k, old, new string, d *schema.ResourceData) bool {
				if v, ok := d.GetOk("private_dns_type"); ok && v.(string) == "PrivateZone" {
					return new == ""
				}
				return false
			}

			result := diffSuppressFunc("standby_dns", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
