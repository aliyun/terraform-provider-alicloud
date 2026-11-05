package hbase

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitWhiteIpListDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		oldValue    string
		newValue    string
		expected    bool
		description string
	}{
		{
			name:        "Same_IPs_Different_Order",
			oldValue:    "192.168.1.1,192.168.1.2,192.168.1.3",
			newValue:    "192.168.1.3,192.168.1.1,192.168.1.2",
			expected:    true,
			description: "Same IPs in different order should suppress diff",
		},
		{
			name:        "Different_IPs",
			oldValue:    "192.168.1.1,192.168.1.2",
			newValue:    "192.168.1.3,192.168.1.4",
			expected:    false,
			description: "Different IPs should not suppress diff",
		},
		{
			name:        "Different_Count",
			oldValue:    "192.168.1.1,192.168.1.2",
			newValue:    "192.168.1.1",
			expected:    false,
			description: "Different count of IPs should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

			result := whiteIpListDiffSuppressFunc("key", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitPayTypePostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		payType     string
		expected    bool
		description string
	}{
		{
			name:        "PostPaid",
			payType:     "PostPaid",
			expected:    true,
			description: "PostPaid pay_type should suppress diff",
		},
		{
			name:        "postpaid_lowercase",
			payType:     "postpaid",
			expected:    true,
			description: "postpaid (lowercase) should suppress diff",
		},
		{
			name:        "PrePaid",
			payType:     "PrePaid",
			expected:    false,
			description: "PrePaid should not suppress diff",
		},
		{
			name:        "Empty",
			payType:     "",
			expected:    false,
			description: "Empty pay_type should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"pay_type": {Type: schema.TypeString, Optional: true},
			}, map[string]interface{}{
				"pay_type": tc.payType,
			})
			result := payTypePostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitEngineDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		engine      string
		expected    bool
		description string
	}{
		{
			name:        "BDS_Engine",
			engine:      "bds",
			expected:    true,
			description: "bds engine should suppress diff",
		},
		{
			name:        "BDS_Engine_Uppercase",
			engine:      "BDS",
			expected:    true,
			description: "BDS engine (uppercase) should suppress diff",
		},
		{
			name:        "Other_Engine",
			engine:      "mysql",
			expected:    false,
			description: "Other engine should not suppress diff",
		},
		{
			name:        "Empty_Engine",
			engine:      "",
			expected:    false,
			description: "Empty engine should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"engine": {Type: schema.TypeString, Optional: true},
			}, map[string]interface{}{
				"engine": tc.engine,
			})
			result := engineDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
