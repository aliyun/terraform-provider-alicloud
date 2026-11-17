// lintignore: S013
package rkvstore

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitCommonRedisSecurityGroupIdDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"SameIds", "sg1,sg2", "sg2,sg1", true},
		{"DifferentOrder", "sg1,sg2", "sg2,sg1", true},
		{"DifferentIds", "sg1,sg2", "sg3,sg4", false},
		{"DifferentCount", "sg1", "sg1,sg2", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := initTestData(t, nil)
			result := redisSecurityGroupIdDiffSuppressFunc("", tc.old, tc.new, d)
			if result != tc.expected {
				t.Errorf("Expected %v got %v", tc.expected, result)
			}
		})
	}
}

func TestUnitCommonRedisPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		chargeType  string
		paymentType string
		expected    bool
		description string
	}{
		{
			name:        "PrePaid_PaymentType",
			paymentType: "PrePaid",
			expected:    false,
			description: "PrePaid payment type should not suppress diff",
		},
		{
			name:        "PrePaid_ChargeType",
			chargeType:  "PrePaid",
			expected:    false,
			description: "PrePaid charge type should not suppress diff",
		},
		{
			name:        "PostPaid_PaymentType",
			paymentType: "PostPaid",
			expected:    true,
			description: "PostPaid payment type should suppress diff",
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

			result := redisPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonRedisPostPaidAndRenewDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		paymentType string
		chargeType  string
		autoRenew   bool
		expected    bool
		description string
	}{
		{
			name:        "PrePaid_PaymentType_AutoRenew",
			paymentType: "PrePaid",
			autoRenew:   true,
			expected:    false,
			description: "PrePaid payment_type with auto_renew should not suppress diff",
		},
		{
			name:        "PrePaid_ChargeType_AutoRenew",
			chargeType:  "PrePaid",
			autoRenew:   true,
			expected:    false,
			description: "PrePaid instance_charge_type with auto_renew should not suppress diff",
		},
		{
			name:        "PrePaid_PaymentType_NoAutoRenew",
			paymentType: "PrePaid",
			autoRenew:   false,
			expected:    true,
			description: "PrePaid without auto_renew should suppress diff",
		},
		{
			name:        "PostPaid_AutoRenew",
			paymentType: "PostPaid",
			autoRenew:   true,
			expected:    true,
			description: "PostPaid with auto_renew should suppress diff",
		},
		{
			name:        "NoType_AutoRenew",
			autoRenew:   true,
			expected:    true,
			description: "No charge type with auto_renew should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"auto_renew": tc.autoRenew,
			}
			if tc.paymentType != "" {
				data["payment_type"] = tc.paymentType
			}
			if tc.chargeType != "" {
				data["instance_charge_type"] = tc.chargeType
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"payment_type":         {Type: schema.TypeString},
				"instance_charge_type": {Type: schema.TypeString},
				"auto_renew":           {Type: schema.TypeBool},
			}, data)
			result := redisPostPaidAndRenewDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func initTestData(t *testing.T, data map[string]interface{}) *schema.ResourceData {
	resourceSchema := map[string]*schema.Schema{
		"listener_forward":        {Type: schema.TypeString},
		"protocol":                {Type: schema.TypeString},
		"sticky_session":          {Type: schema.TypeString},
		"sticky_session_type":     {Type: schema.TypeString},
		"health_check":            {Type: schema.TypeString},
		"instance_charge_type":    {Type: schema.TypeString},
		"payment_type":            {Type: schema.TypeString},
		"enable_backup_log":       {Type: schema.TypeBool},
		"log_backup":              {Type: schema.TypeBool},
		"backup_retention_period": {Type: schema.TypeInt},
	}

	d := schema.TestResourceDataRaw(t, resourceSchema, data)
	return d
}
