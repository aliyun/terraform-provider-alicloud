// lintignore: S013
package alidns

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonDnsValueDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		dnsType     string
		oldValue    string
		newValue    string
		expected    bool
		description string
	}{
		{
			name:        "NS_Type_With_Trailing_Dot",
			dnsType:     "NS",
			oldValue:    "ns1.example.com",
			newValue:    "ns1.example.com.",
			expected:    true,
			description: "NS type should trim trailing dot",
		},
		{
			name:        "MX_Type_With_Trailing_Dot",
			dnsType:     "MX",
			oldValue:    "mail.example.com",
			newValue:    "mail.example.com.",
			expected:    true,
			description: "MX type should trim trailing dot",
		},
		{
			name:        "CNAME_Type_With_Trailing_Dot",
			dnsType:     "CNAME",
			oldValue:    "www.example.com",
			newValue:    "www.example.com.",
			expected:    true,
			description: "CNAME type should trim trailing dot",
		},
		{
			name:        "A_Type_No_Trim",
			dnsType:     "A",
			oldValue:    "192.168.1.1",
			newValue:    "192.168.1.2",
			expected:    false,
			description: "A type should not trim and values differ",
		},
		{
			name:        "Different_Values",
			dnsType:     "NS",
			oldValue:    "ns1.example.com",
			newValue:    "ns2.example.com",
			expected:    false,
			description: "Different values should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"type": tc.dnsType,
			})

			result := dnsValueDiffSuppressFunc("value", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonDnsPriorityDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		dnsType     string
		expected    bool
		description string
	}{
		{
			name:        "MX_Type",
			dnsType:     "MX",
			expected:    false,
			description: "MX type should not suppress diff",
		},
		{
			name:        "A_Type",
			dnsType:     "A",
			expected:    true,
			description: "A type should suppress diff",
		},
		{
			name:        "CNAME_Type",
			dnsType:     "CNAME",
			expected:    true,
			description: "CNAME type should suppress diff",
		},
		{
			name:        "Empty_Type",
			dnsType:     "",
			expected:    true,
			description: "Empty type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"type": tc.dnsType,
			})
			result := dnsPriorityDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
