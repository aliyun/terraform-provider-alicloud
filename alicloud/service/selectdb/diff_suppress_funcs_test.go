package selectdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		paymentType string
		expected    bool
		description string
	}{
		{
			name:        "Subscription",
			paymentType: "Subscription",
			expected:    false,
			description: "Subscription should not suppress diff",
		},
		{
			name:        "PayAsYouGo",
			paymentType: "PayAsYouGo",
			expected:    true,
			description: "PayAsYouGo should suppress diff",
		},
		{
			name:        "Empty",
			paymentType: "",
			expected:    true,
			description: "Empty payment_type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"payment_type": {Type: schema.TypeString, Optional: true},
			}, map[string]interface{}{
				"payment_type": tc.paymentType,
			})
			result := postPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
