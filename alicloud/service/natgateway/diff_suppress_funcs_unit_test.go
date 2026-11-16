// lintignore: S013
package natgateway

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestUnitCommonPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		chargeType  string
		paymentType string
		expected    bool
	}{
		{"PrepaidCharge", "Prepaid", "", false},
		{"SubscriptionPayment", "", "Subscription", false},
		{"PostPaid", "PostPaid", "", true},
		{"PayAsYouGo", "", "PayAsYouGo", true},
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

			d := initTestData(t, data)
			result := PostPaidDiffSuppressFunc("", "", "", d)
			if result != tc.expected {
				t.Errorf("Expected %v got %v", tc.expected, result)
			}
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
