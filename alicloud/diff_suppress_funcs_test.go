package alicloud

import (
	"github.com/stretchr/testify/assert"
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

func TestUnitCommonHttpDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		protocol string
		expected bool
	}{
		{"HTTP_Protocol", "http", false},
		{"HTTPS_Protocol", "https", true},
		{"TCP_Protocol", "tcp", true},
		{"UDP_Protocol", "udp", true},
		{"Empty_Protocol", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})

			result := httpDiffSuppressFunc("any_key", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result, "Unexpected result for protocol: "+tc.protocol)
		})
	}
}

func TestUnitCommonForwardPortDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		protocol        string
		listenerForward string
		expected        bool
	}{
		{"HTTP_With_Forward", "http", "on", false},
		{"HTTP_Without_Forward", "http", "off", true},
		{"HTTPS_With_Forward", "https", "on", true},
		{"TCP_With_Forward", "tcp", "on", true},
		{"Empty_Protocol_With_Forward", "", "on", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol":         {Type: schema.TypeString},
				"listener_forward": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol":         tc.protocol,
				"listener_forward": tc.listenerForward,
			})

			result := forwardPortDiffSuppressFunc("any_key", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUnitCommonHttpsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name     string
		protocol string
		expected bool
	}{
		{"HTTPS_Protocol", "https", false},
		{"HTTP_Protocol", "http", true},
		{"TCP_Protocol", "tcp", true},
		{"UDP_Protocol", "udp", true},
		{"Empty_Protocol", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})

			result := httpsDiffSuppressFunc("any_key", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUnitCommonStickySessionTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name          string
		protocol      string
		stickySession string
		expected      bool
	}{
		{"HTTP_With_StickySession", "http", "on", false},
		{"HTTP_Without_StickySession", "http", "off", true},
		{"HTTPS_With_StickySession", "https", "on", false},
		{"HTTPS_Without_StickySession", "https", "off", true},
		{"TCP_With_StickySession", "tcp", "on", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol":       {Type: schema.TypeString},
				"sticky_session": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol":       tc.protocol,
				"sticky_session": tc.stickySession,
			})

			result := stickySessionTypeDiffSuppressFunc("any_key", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUnitCommonCookieTimeoutDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name              string
		protocol          string
		stickySession     string
		stickySessionType string
		expected          bool
	}{
		{"HTTP_Insert_StickySession", "http", "on", "insert", false},
		{"HTTP_Server_StickySession", "http", "on", "server", true},
		{"HTTPS_Insert_StickySession", "https", "on", "insert", false},
		{"HTTPS_Server_StickySession", "https", "on", "server", true},
		{"TCP_Insert_StickySession", "tcp", "on", "insert", true},
		{"HTTP_No_StickySession", "http", "off", "insert", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol":            {Type: schema.TypeString},
				"sticky_session":      {Type: schema.TypeString},
				"sticky_session_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol":            tc.protocol,
				"sticky_session":      tc.stickySession,
				"sticky_session_type": tc.stickySessionType,
			})

			result := cookieTimeoutDiffSuppressFunc("any_key", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result)
		})
	}
}
