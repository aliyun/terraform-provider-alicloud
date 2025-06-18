package alicloud

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestUnitCommonHttpHttpsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name           string
		forwardEnabled bool
		protocol       string
		expected       bool
	}{
		{"ListenerForwardOn", true, "http", true},
		{"ProtocolHTTP", false, "http", false},
		{"ProtocolHTTPS", false, "https", false},
		{"OtherProtocol", false, "udp", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := initTestData(t, map[string]interface{}{
				"listener_forward": mapToStr(tc.forwardEnabled, OnFlag, OffFlag),
				"protocol":         tc.protocol,
			})
			result := httpHttpsDiffSuppressFunc("", "", "", d)
			if result != tc.expected {
				t.Errorf("Expected %v got %v", tc.expected, result)
			}
		})
	}
}

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

func TestUnitCommonLogRetentionPeriodDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		backupLogEnabled bool
		logBackupEnabled bool
		newPeriod        int
		retentionPeriod  int
		expected         bool
	}{
		{"BackupLogEnabled", true, false, 30, 0, false},
		{"LogBackupEnabled", false, true, 30, 0, false},
		{"NewGTBackupRetention", false, false, 30, 20, true},
		{"ValidSuppression", false, false, 10, 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := initTestData(t, map[string]interface{}{
				"enable_backup_log":       tc.backupLogEnabled,
				"log_backup":              tc.logBackupEnabled,
				"backup_retention_period": tc.retentionPeriod,
			})
			result := logRetentionPeriodDiffSuppressFunc(
				"", "", strconv.Itoa(tc.newPeriod), d)
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

func mapToStr(condition bool, trueVal, falseVal FlagType) string {
	if condition {
		return string(trueVal)
	}
	return string(falseVal)
}
