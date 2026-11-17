// lintignore: S013
package adb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonAdbPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		payType     string
		paymentType string
		expected    bool
		description string
	}{
		{
			name:        "PrePaid_PayType",
			payType:     "PrePaid",
			expected:    false,
			description: "PrePaid pay type should not suppress diff",
		},
		{
			name:        "Subscription_PaymentType",
			paymentType: "Subscription",
			expected:    false,
			description: "Subscription payment type should not suppress diff",
		},
		{
			name:        "PostPaid_PayType",
			payType:     "PostPaid",
			expected:    true,
			description: "PostPaid pay type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.payType != "" {
				data["pay_type"] = tc.payType
			}
			if tc.paymentType != "" {
				data["payment_type"] = tc.paymentType
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"pay_type":     {Type: schema.TypeString},
				"payment_type": {Type: schema.TypeString},
			}, data)

			result := adbPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
func TestUnitCommonAdbPostPaidAndRenewDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name          string
		payType       string
		paymentType   string
		renewalStatus string
		expected      bool
		description   string
	}{
		{
			name:          "PrePaid_PayType_NonNotRenewal",
			payType:       "PrePaid",
			renewalStatus: "AutoRenewal",
			expected:      false,
			description:   "PrePaid pay_type with non-NotRenewal status should not suppress diff",
		},
		{
			name:          "Subscription_PaymentType_NonNotRenewal",
			paymentType:   "Subscription",
			renewalStatus: "AutoRenewal",
			expected:      false,
			description:   "Subscription payment_type with non-NotRenewal status should not suppress diff",
		},
		{
			name:          "PrePaid_PayType_NotRenewal",
			payType:       "PrePaid",
			renewalStatus: "NotRenewal",
			expected:      true,
			description:   "PrePaid pay_type with NotRenewal should suppress diff",
		},
		{
			name:          "PostPaid_PayType",
			payType:       "PostPaid",
			renewalStatus: "AutoRenewal",
			expected:      true,
			description:   "PostPaid pay_type should suppress diff",
		},
		{
			name:          "No_PayType",
			renewalStatus: "AutoRenewal",
			expected:      true,
			description:   "No pay_type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"renewal_status": tc.renewalStatus,
			}
			if tc.payType != "" {
				data["pay_type"] = tc.payType
			}
			if tc.paymentType != "" {
				data["payment_type"] = tc.paymentType
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"pay_type":       {Type: schema.TypeString},
				"payment_type":   {Type: schema.TypeString},
				"renewal_status": {Type: schema.TypeString},
			}, data)
			result := adbPostPaidAndRenewDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
