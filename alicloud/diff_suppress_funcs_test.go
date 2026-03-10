// lintignore: S013
package alicloud

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
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

func TestUnitCommonCookieDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name              string
		protocol          string
		stickySession     string
		stickySessionType string
		expected          bool
		description       string
	}{
		{
			name:              "Server_StickySession_Type",
			protocol:          "http",
			stickySession:     "on",
			stickySessionType: "server",
			expected:          false,
			description:       "Server sticky session type should not suppress diff",
		},
		{
			name:              "Insert_StickySession_Type",
			protocol:          "http",
			stickySession:     "on",
			stickySessionType: "insert",
			expected:          true,
			description:       "Insert sticky session type should suppress diff",
		},
		{
			name:              "No_StickySession",
			protocol:          "http",
			stickySession:     "off",
			stickySessionType: "server",
			expected:          true,
			description:       "No sticky session should suppress diff",
		},
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

			result := cookieDiffSuppressFunc("cookie", "old_value", "new_value", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEstablishedTimeoutDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		protocol    string
		expected    bool
		description string
	}{
		{
			name:        "TCP_Protocol",
			protocol:    "tcp",
			expected:    false,
			description: "TCP protocol should not suppress diff",
		},
		{
			name:        "HTTP_Protocol",
			protocol:    "http",
			expected:    true,
			description: "HTTP protocol should suppress diff",
		},
		{
			name:        "UDP_Protocol",
			protocol:    "udp",
			expected:    true,
			description: "UDP protocol should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})

			result := establishedTimeoutDiffSuppressFunc("established_timeout", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonHttpHttpsTcpDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		protocol        string
		healthCheck     string
		healthCheckType string
		listenerForward string
		expected        bool
		description     string
	}{
		{
			name:            "HTTP_With_HealthCheck_On",
			protocol:        "http",
			healthCheck:     "on",
			listenerForward: "off",
			expected:        false,
			description:     "HTTP with health check on should not suppress diff",
		},
		{
			name:            "TCP_With_HTTP_HealthCheck",
			protocol:        "tcp",
			healthCheckType: "http",
			expected:        false,
			description:     "TCP with HTTP health check should not suppress diff",
		},
		{
			name:        "Other_Cases",
			protocol:    "udp",
			healthCheck: "off",
			expected:    true,
			description: "Other cases should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"protocol": tc.protocol,
			}
			if tc.healthCheck != "" {
				data["health_check"] = tc.healthCheck
			}
			if tc.healthCheckType != "" {
				data["health_check_type"] = tc.healthCheckType
			}
			if tc.listenerForward != "" {
				data["listener_forward"] = tc.listenerForward
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol":          {Type: schema.TypeString},
				"health_check":      {Type: schema.TypeString},
				"health_check_type": {Type: schema.TypeString},
				"listener_forward":  {Type: schema.TypeString},
			}, data)

			result := httpHttpsTcpDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

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

func TestUnitCommonCsKubernetesMasterPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                     string
		masterInstanceChargeType string
		resourceId               string
		forceUpdate              bool
		expected                 bool
		description              string
	}{
		{
			name:                     "PostPaid_ChargeType",
			masterInstanceChargeType: "PostPaid",
			expected:                 true,
			description:              "PostPaid charge type should suppress diff",
		},
		{
			name:                     "PrePaid_No_ForceUpdate_Existing_Resource",
			masterInstanceChargeType: "PrePaid",
			resourceId:               "existing-id",
			forceUpdate:              false,
			expected:                 true,
			description:              "PrePaid with no force update on existing resource should suppress diff",
		},
		{
			name:                     "PrePaid_With_ForceUpdate",
			masterInstanceChargeType: "PrePaid",
			resourceId:               "existing-id",
			forceUpdate:              true,
			expected:                 false,
			description:              "PrePaid with force update should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"master_instance_charge_type": {Type: schema.TypeString},
				"force_update":                {Type: schema.TypeBool},
			}, map[string]interface{}{
				"master_instance_charge_type": tc.masterInstanceChargeType,
				"force_update":                tc.forceUpdate,
			})
			d.SetId(tc.resourceId)

			result := csKubernetesMasterPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEnableBackupLogDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		enableBackupLog bool
		logBackup       bool
		expected        bool
		description     string
	}{
		{
			name:            "BackupLog_Enabled",
			enableBackupLog: true,
			expected:        false,
			description:     "Enable backup log should not suppress diff",
		},
		{
			name:        "LogBackup_Enabled",
			logBackup:   true,
			expected:    false,
			description: "Log backup enabled should not suppress diff",
		},
		{
			name:        "Both_Disabled",
			expected:    true,
			description: "Both disabled should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_backup_log": {Type: schema.TypeBool},
				"log_backup":        {Type: schema.TypeBool},
			}, map[string]interface{}{
				"enable_backup_log": tc.enableBackupLog,
				"log_backup":        tc.logBackup,
			})

			result := enableBackupLogDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonArchiveBackupPeriodDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                  string
		enableBackupLog       bool
		logBackup             bool
		backupRetentionPeriod int
		retentionPeriod       int
		newValue              string
		expected              bool
		description           string
	}{
		{
			name:            "BackupLog_Enabled",
			enableBackupLog: true,
			expected:        false,
			newValue:        "100",
			description:     "Enable backup log should not suppress diff",
		},
		{
			name:        "LogBackup_Enabled",
			logBackup:   true,
			expected:    false,
			newValue:    "100",
			description: "Log backup enabled should not suppress diff",
		},
		{
			name:                  "Both_Disabled_Invalid_New_High_Retention",
			backupRetentionPeriod: 1000,
			retentionPeriod:       1000,
			newValue:              "invalid",
			expected:              true,
			description:           "Both disabled with invalid new and high retention should suppress diff",
		},
		{
			name:                  "Both_Disabled_Invalid_New_Low_Retention",
			backupRetentionPeriod: 500,
			newValue:              "invalid",
			expected:              false,
			description:           "Both disabled with invalid new and low retention (<730) should not suppress diff",
		},
		{
			name:                  "Both_Disabled_Valid_New",
			backupRetentionPeriod: 1000,
			expected:              true,
			newValue:              "100",
			description:           "Both disabled with valid new value should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_backup_log":       {Type: schema.TypeBool},
				"log_backup":              {Type: schema.TypeBool},
				"backup_retention_period": {Type: schema.TypeInt},
				"retention_period":        {Type: schema.TypeInt},
			}, map[string]interface{}{
				"enable_backup_log":       tc.enableBackupLog,
				"log_backup":              tc.logBackup,
				"backup_retention_period": tc.backupRetentionPeriod,
				"retention_period":        tc.retentionPeriod,
			})

			result := archiveBackupPeriodDiffSuppressFunc("key", "old", tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

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

func TestUnitCommonPostPaidAndRenewDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		chargeType  string
		paymentType string
		autoRenew   bool
		expected    bool
		description string
	}{
		{
			name:        "Prepaid_With_AutoRenew",
			chargeType:  "Prepaid",
			autoRenew:   true,
			expected:    false,
			description: "Prepaid with auto renew should not suppress diff",
		},
		{
			name:        "Subscription_With_AutoRenew",
			paymentType: "Subscription",
			autoRenew:   true,
			expected:    false,
			description: "Subscription with auto renew should not suppress diff",
		},
		{
			name:        "PostPaid_No_AutoRenew",
			chargeType:  "PostPaid",
			autoRenew:   false,
			expected:    true,
			description: "PostPaid without auto renew should suppress diff",
		},
		{
			name:        "Prepaid_No_AutoRenew",
			chargeType:  "Prepaid",
			autoRenew:   false,
			expected:    true,
			description: "Prepaid without auto renew should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"auto_renew": tc.autoRenew,
			}
			if tc.chargeType != "" {
				data["instance_charge_type"] = tc.chargeType
			}
			if tc.paymentType != "" {
				data["payment_type"] = tc.paymentType
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
				"payment_type":         {Type: schema.TypeString},
				"auto_renew":           {Type: schema.TypeBool},
			}, data)

			result := PostPaidAndRenewDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
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

func TestUnitCommonElasticsearchEnablePublicDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		enablePublic bool
		expected     bool
		description  string
	}{
		{
			name:         "Enable_Public_True",
			enablePublic: true,
			expected:     false,
			description:  "Enable public true should not suppress diff",
		},
		{
			name:         "Enable_Public_False",
			enablePublic: false,
			expected:     true,
			description:  "Enable public false should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_public": {Type: schema.TypeBool},
			}, map[string]interface{}{
				"enable_public": tc.enablePublic,
			})

			result := elasticsearchEnablePublicDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonElasticsearchEnableKibanaPublicDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                      string
		enableKibanaPublicNetwork bool
		expected                  bool
		description               string
	}{
		{
			name:                      "Enable_Kibana_Public_True",
			enableKibanaPublicNetwork: true,
			expected:                  false,
			description:               "Enable Kibana public true should not suppress diff",
		},
		{
			name:                      "Enable_Kibana_Public_False",
			enableKibanaPublicNetwork: false,
			expected:                  true,
			description:               "Enable Kibana public false should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_kibana_public_network": {Type: schema.TypeBool},
			}, map[string]interface{}{
				"enable_kibana_public_network": tc.enableKibanaPublicNetwork,
			})

			result := elasticsearchEnableKibanaPublicDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		payType     string
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
			name:        "PostPaid_PayType",
			payType:     "PostPaid",
			expected:    true,
			description: "PostPaid pay type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"pay_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"pay_type": tc.payType,
			})

			result := polardbPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonElasticsearchEnableKibanaPrivateDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                       string
		enableKibanaPrivateNetwork bool
		expected                   bool
		description                string
	}{
		{
			name:                       "Enable_Kibana_Private_True",
			enableKibanaPrivateNetwork: true,
			expected:                   false,
			description:                "Enable Kibana private true should not suppress diff",
		},
		{
			name:                       "Enable_Kibana_Private_False",
			enableKibanaPrivateNetwork: false,
			expected:                   true,
			description:                "Enable Kibana private false should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_kibana_private_network": {Type: schema.TypeBool},
			}, map[string]interface{}{
				"enable_kibana_private_network": tc.enableKibanaPrivateNetwork,
			})

			result := elasticsearchEnableKibanaPrivateDiffSuppressFunc("key", "old", "new", d)
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

func TestUnitCommonPolardbDBClusterVersionDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                 string
		oldValue             string
		clusterLatestVersion string
		hasLatestVersion     bool
		expected             bool
		description          string
	}{
		{
			name:             "No_Latest_Version",
			oldValue:         "5.7",
			hasLatestVersion: false,
			expected:         true,
			description:      "No latest version should suppress diff",
		},
		{
			name:                 "Old_Equals_Latest",
			oldValue:             "8.0",
			clusterLatestVersion: "8.0",
			hasLatestVersion:     true,
			expected:             true,
			description:          "Old equals latest version should suppress diff",
		},
		{
			name:                 "Old_Not_Equals_Latest",
			oldValue:             "5.7",
			clusterLatestVersion: "8.0",
			hasLatestVersion:     true,
			expected:             false,
			description:          "Old not equals latest version should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.hasLatestVersion {
				data["cluster_latest_version"] = tc.clusterLatestVersion
			}

			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"cluster_latest_version": {Type: schema.TypeString},
			}, data)

			result := polardbDBClusterVersionDiffSuppressFunc("key", tc.oldValue, "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbTDEAndEnabledDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		tdeStatus   string
		oldValue    string
		newValue    string
		expected    bool
		description string
	}{
		{
			name:        "TDE_Enabled_With_Different_Values",
			tdeStatus:   "Enabled",
			oldValue:    "key1",
			newValue:    "key2",
			expected:    true,
			description: "TDE enabled with different non-empty values should suppress diff",
		},
		{
			name:        "TDE_Disabled",
			tdeStatus:   "Disabled",
			oldValue:    "key1",
			newValue:    "key2",
			expected:    false,
			description: "TDE disabled should not suppress diff",
		},
		{
			name:        "TDE_Enabled_With_Empty_Old",
			tdeStatus:   "Enabled",
			oldValue:    "",
			newValue:    "key2",
			expected:    false,
			description: "TDE enabled with empty old value should not suppress diff",
		},
		{
			name:        "TDE_Enabled_With_Empty_New",
			tdeStatus:   "Enabled",
			oldValue:    "key1",
			newValue:    "",
			expected:    false,
			description: "TDE enabled with empty new value should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"tde_status": {Type: schema.TypeString},
			}, map[string]interface{}{
				"tde_status": tc.tdeStatus,
			})

			result := polardbTDEAndEnabledDiffSuppressFunc("key", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbPostPaidAndRenewDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name          string
		payType       string
		renewalStatus string
		expected      bool
		description   string
	}{
		{
			name:          "PrePaid_With_AutoRenewal",
			payType:       "PrePaid",
			renewalStatus: "AutoRenewal",
			expected:      false,
			description:   "PrePaid with auto renewal should not suppress diff",
		},
		{
			name:          "PrePaid_Without_Renewal",
			payType:       "PrePaid",
			renewalStatus: "NotRenewal",
			expected:      true,
			description:   "PrePaid without renewal should suppress diff",
		},
		{
			name:          "PostPaid",
			payType:       "PostPaid",
			renewalStatus: "NotRenewal",
			expected:      true,
			description:   "PostPaid should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"pay_type":       {Type: schema.TypeString},
				"renewal_status": {Type: schema.TypeString},
			}, map[string]interface{}{
				"pay_type":       tc.payType,
				"renewal_status": tc.renewalStatus,
			})

			result := polardbPostPaidAndRenewDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

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

func TestUnitCommonVpcTypeResourceDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		vswitchId   string
		expected    bool
		description string
	}{
		{
			name:        "With_VSwitch_ID",
			vswitchId:   "vsw-123456",
			expected:    false,
			description: "With vswitch_id should not suppress diff",
		},
		{
			name:        "Empty_VSwitch_ID",
			vswitchId:   "",
			expected:    true,
			description: "Empty vswitch_id should suppress diff",
		},
		{
			name:        "Whitespace_VSwitch_ID",
			vswitchId:   "   ",
			expected:    true,
			description: "Whitespace vswitch_id should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"vswitch_id": {Type: schema.TypeString},
			}, map[string]interface{}{
				"vswitch_id": tc.vswitchId,
			})

			result := vpcTypeResourceDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonWhiteIpListDiffSuppressFunc(t *testing.T) {
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

func TestUnitCommonUpperLowerCaseDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		oldValue    string
		newValue    string
		expected    bool
		description string
	}{
		{
			name:        "Same_Case_Insensitive",
			oldValue:    "MySQL",
			newValue:    "mysql",
			expected:    true,
			description: "Same string different case should suppress diff",
		},
		{
			name:        "Different_String",
			oldValue:    "MySQL",
			newValue:    "PostgreSQL",
			expected:    false,
			description: "Different strings should not suppress diff",
		},
		{
			name:        "Same_String_Same_Case",
			oldValue:    "MySQL",
			newValue:    "MySQL",
			expected:    true,
			description: "Same string same case should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})

			result := UpperLowerCaseDiffSuppressFunc("key", tc.oldValue, tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonTcpUdpDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		protocol    string
		expected    bool
		description string
	}{
		{
			name:        "TCP_Protocol",
			protocol:    "tcp",
			expected:    false,
			description: "TCP protocol should not suppress diff",
		},
		{
			name:        "UDP_Protocol",
			protocol:    "udp",
			expected:    false,
			description: "UDP protocol should not suppress diff",
		},
		{
			name:        "HTTP_Protocol",
			protocol:    "http",
			expected:    true,
			description: "HTTP protocol should suppress diff",
		},
		{
			name:        "HTTPS_Protocol",
			protocol:    "https",
			expected:    true,
			description: "HTTPS protocol should suppress diff",
		},
		{
			name:        "Empty_Protocol",
			protocol:    "",
			expected:    true,
			description: "Empty protocol should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})
			result := tcpUdpDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSslCertificateIdDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		protocol    string
		expected    bool
		description string
	}{
		{
			name:        "HTTPS_Protocol",
			protocol:    "https",
			expected:    false,
			description: "HTTPS protocol should not suppress diff",
		},
		{
			name:        "HTTP_Protocol",
			protocol:    "http",
			expected:    true,
			description: "HTTP protocol should suppress diff",
		},
		{
			name:        "TCP_Protocol",
			protocol:    "tcp",
			expected:    true,
			description: "TCP protocol should suppress diff",
		},
		{
			name:        "Empty_Protocol",
			protocol:    "",
			expected:    true,
			description: "Empty protocol should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})
			result := sslCertificateIdDiffSuppressFunc("key", "old", "new", d)
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

func TestUnitCommonCsKubernetesWorkerPostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                     string
		workerInstanceChargeType string
		resourceId               string
		forceUpdate              bool
		expected                 bool
		description              string
	}{
		{
			name:                     "PostPaid_ChargeType",
			workerInstanceChargeType: "PostPaid",
			expected:                 true,
			description:              "PostPaid charge type should suppress diff",
		},
		{
			name:                     "PrePaid_No_ForceUpdate_Existing",
			workerInstanceChargeType: "PrePaid",
			resourceId:               "existing-id",
			forceUpdate:              false,
			expected:                 true,
			description:              "PrePaid with no force update on existing resource should suppress diff",
		},
		{
			name:                     "PrePaid_ForceUpdate",
			workerInstanceChargeType: "PrePaid",
			resourceId:               "existing-id",
			forceUpdate:              true,
			expected:                 false,
			description:              "PrePaid with force update should not suppress diff",
		},
		{
			name:                     "PrePaid_New_Resource",
			workerInstanceChargeType: "PrePaid",
			resourceId:               "",
			forceUpdate:              false,
			expected:                 false,
			description:              "PrePaid on new resource (empty ID) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"worker_instance_charge_type": {Type: schema.TypeString},
				"force_update":                {Type: schema.TypeBool},
			}, map[string]interface{}{
				"worker_instance_charge_type": tc.workerInstanceChargeType,
				"force_update":                tc.forceUpdate,
			})
			d.SetId(tc.resourceId)
			result := csKubernetesWorkerPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCsNodepoolInstancePostPaidDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		chargeType  string
		expected    bool
		description string
	}{
		{
			name:        "PostPaid",
			chargeType:  "PostPaid",
			expected:    true,
			description: "PostPaid should suppress diff",
		},
		{
			name:        "PrePaid",
			chargeType:  "PrePaid",
			expected:    false,
			description: "PrePaid should not suppress diff",
		},
		{
			name:        "Empty",
			chargeType:  "",
			expected:    false,
			description: "Empty charge type should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"instance_charge_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"instance_charge_type": tc.chargeType,
			})
			result := csNodepoolInstancePostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonMasterDiskPerformanceLevelDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		diskCategory string
		expected    bool
		description string
	}{
		{
			name:        "CloudESSD",
			diskCategory: "cloud_essd",
			expected:    false,
			description: "cloud_essd category should not suppress diff",
		},
		{
			name:        "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:    true,
			description: "Non-cloud_essd category should suppress diff",
		},
		{
			name:        "Cloud",
			diskCategory: "cloud",
			expected:    true,
			description: "Cloud category should suppress diff",
		},
		{
			name:        "Empty",
			diskCategory: "",
			expected:    false,
			description: "Empty disk category (GetOk returns false) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"master_disk_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"master_disk_category": tc.diskCategory,
			})
			result := masterDiskPerformanceLevelDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonWorkerDiskPerformanceLevelDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		diskCategory string
		expected    bool
		description string
	}{
		{
			name:        "CloudESSD",
			diskCategory: "cloud_essd",
			expected:    false,
			description: "cloud_essd category should not suppress diff",
		},
		{
			name:        "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:    true,
			description: "Non-cloud_essd category should suppress diff",
		},
		{
			name:        "Empty",
			diskCategory: "",
			expected:    false,
			description: "Empty disk category (GetOk returns false) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"worker_disk_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"worker_disk_category": tc.diskCategory,
			})
			result := workerDiskPerformanceLevelDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCsNodepoolDiskPerformanceLevelDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		diskCategory string
		expected    bool
		description string
	}{
		{
			name:        "CloudESSD",
			diskCategory: "cloud_essd",
			expected:    false,
			description: "cloud_essd should not suppress diff",
		},
		{
			name:        "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:    true,
			description: "cloud_ssd should suppress diff",
		},
		{
			name:        "Empty",
			diskCategory: "",
			expected:    false,
			description: "Empty system_disk_category (GetOk returns false) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"system_disk_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"system_disk_category": tc.diskCategory,
			})
			result := csNodepoolDiskPerformanceLevelDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCsNodepoolSpotInstanceSettingDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		spotStrategy string
		expected    bool
		description string
	}{
		{
			name:        "SpotWithPriceLimit",
			spotStrategy: "SpotWithPriceLimit",
			expected:    false,
			description: "SpotWithPriceLimit should not suppress diff",
		},
		{
			name:        "NoSpot",
			spotStrategy: "NoSpot",
			expected:    true,
			description: "NoSpot should suppress diff",
		},
		{
			name:        "SpotAsPriceGo",
			spotStrategy: "SpotAsPriceGo",
			expected:    true,
			description: "SpotAsPriceGo should suppress diff",
		},
		{
			name:        "Empty",
			spotStrategy: "",
			expected:    true,
			description: "Empty spot strategy should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"spot_strategy": {Type: schema.TypeString},
			}, map[string]interface{}{
				"spot_strategy": tc.spotStrategy,
			})
			result := csNodepoolSpotInstanceSettingDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCsNodepoolScalingPolicyDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name          string
		scalingConfig string
		expected      bool
		description   string
	}{
		{
			name:          "WithScalingConfig",
			scalingConfig: "some-config",
			expected:      false,
			description:   "With scaling_config should not suppress diff",
		},
		{
			name:          "NoScalingConfig",
			scalingConfig: "",
			expected:      true,
			description:   "Without scaling_config should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"scaling_config": {Type: schema.TypeString},
			}, map[string]interface{}{
				"scaling_config": tc.scalingConfig,
			})
			result := csNodepoolScalingPolicyDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonChargeTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		paymentType string
		expected    bool
		description string
	}{
		{
			name:        "WithPaymentType",
			paymentType: "Subscription",
			expected:    true,
			description: "When payment_type is set should suppress diff",
		},
		{
			name:        "PayAsYouGo",
			paymentType: "PayAsYouGo",
			expected:    true,
			description: "PayAsYouGo payment_type should suppress diff",
		},
		{
			name:        "NoPaymentType",
			paymentType: "",
			expected:    false,
			description: "Without payment_type should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"payment_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"payment_type": tc.paymentType,
			})
			result := ChargeTypeDiffSuppressFunc("key", "old", "new", d)
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

func TestUnitCommonRamSAMLProviderDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "Same_Content",
			old:         "<saml>content</saml>",
			new:         "<saml>content</saml>",
			expected:    true,
			description: "Identical strings should suppress diff",
		},
		{
			name:        "Different_Content",
			old:         "<saml>old</saml>",
			new:         "<saml>new</saml>",
			expected:    false,
			description: "Different strings should not suppress diff",
		},
		{
			name:        "Same_With_Newlines",
			old:         "<saml>\ncontent\n</saml>",
			new:         "<saml>content</saml>",
			expected:    true,
			description: "Strings differing only by newlines should suppress diff",
		},
		{
			name:        "Same_With_Whitespace",
			old:         "  content  ",
			new:         "content",
			expected:    true,
			description: "Strings differing only by surrounding whitespace should suppress diff",
		},
		{
			name:        "Empty_Strings",
			old:         "",
			new:         "",
			expected:    true,
			description: "Empty strings should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ramSAMLProviderDiffSuppressFunc(tc.old, tc.new)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbStorageTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		creationCategory string
		storageType      string
		expected         bool
		description      string
	}{
		{
			name:             "SENormal_ESSDAUTOPL",
			creationCategory: "SENormal",
			storageType:      "ESSDAUTOPL",
			expected:         false,
			description:      "SENormal with ESSDAUTOPL should not suppress diff",
		},
		{
			name:             "SENormal_Other_Storage",
			creationCategory: "SENormal",
			storageType:      "PSL4",
			expected:         true,
			description:      "SENormal with non-ESSDAUTOPL should suppress diff",
		},
		{
			name:             "Non_SENormal",
			creationCategory: "Normal",
			storageType:      "ESSDAUTOPL",
			expected:         true,
			description:      "Non-SENormal should suppress diff",
		},
		{
			name:             "Empty_Category",
			creationCategory: "",
			storageType:      "ESSDAUTOPL",
			expected:         true,
			description:      "Empty creation category should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"creation_category": {Type: schema.TypeString},
				"storage_type":      {Type: schema.TypeString},
			}, map[string]interface{}{
				"creation_category": tc.creationCategory,
				"storage_type":      tc.storageType,
			})
			result := polardbStorageTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbServrelessTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                  string
		dbType                string
		dbVersion             string
		serverlessType        string
		serverlessSteadySwitch string
		expected              bool
		description           string
	}{
		{
			name:           "MySQL_80_AgileServerless",
			dbType:         "MySQL",
			dbVersion:      "8.0",
			serverlessType: "AgileServerless",
			expected:       false,
			description:    "MySQL 8.0 with AgileServerless should not suppress diff",
		},
		{
			name:                  "MySQL_57_SteadyServerless_ON",
			dbType:                "MySQL",
			dbVersion:             "5.7",
			serverlessType:        "SteadyServerless",
			serverlessSteadySwitch: "ON",
			expected:              false,
			description:           "MySQL 5.7 with SteadyServerless and ON switch should not suppress diff",
		},
		{
			name:                  "MySQL_57_SteadyServerless_OFF",
			dbType:                "MySQL",
			dbVersion:             "5.7",
			serverlessType:        "SteadyServerless",
			serverlessSteadySwitch: "OFF",
			expected:              true,
			description:           "MySQL 5.7 with SteadyServerless and OFF switch should suppress diff",
		},
		{
			name:           "PostgreSQL_14_AgileServerless",
			dbType:         "PostgreSQL",
			dbVersion:      "14",
			serverlessType: "AgileServerless",
			expected:       false,
			description:    "PostgreSQL 14 with AgileServerless should not suppress diff",
		},
		{
			name:           "Oracle_14_AgileServerless",
			dbType:         "Oracle",
			dbVersion:      "14",
			serverlessType: "AgileServerless",
			expected:       false,
			description:    "Oracle 14 with AgileServerless should not suppress diff",
		},
		{
			name:           "MySQL_56_AgileServerless",
			dbType:         "MySQL",
			dbVersion:      "5.6",
			serverlessType: "AgileServerless",
			expected:       true,
			description:    "MySQL 5.6 (unsupported version) should suppress diff",
		},
		{
			name:           "MySQL_80_NoServerless",
			dbType:         "MySQL",
			dbVersion:      "8.0",
			serverlessType: "None",
			expected:       true,
			description:    "MySQL 8.0 with no serverless type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"db_type":                 {Type: schema.TypeString},
				"db_version":              {Type: schema.TypeString},
				"serverless_type":         {Type: schema.TypeString},
				"serverless_steady_switch": {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":                 tc.dbType,
				"db_version":              tc.dbVersion,
				"serverless_type":         tc.serverlessType,
				"serverless_steady_switch": tc.serverlessSteadySwitch,
			})
			result := polardbServrelessTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbProxyClassDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		creationCategory string
		expected         bool
		description      string
	}{
		{
			name:             "SENormal",
			creationCategory: "SENormal",
			expected:         false,
			description:      "SENormal should not suppress diff",
		},
		{
			name:             "Normal",
			creationCategory: "Normal",
			expected:         true,
			description:      "Normal should suppress diff",
		},
		{
			name:             "Empty",
			creationCategory: "",
			expected:         true,
			description:      "Empty category should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"creation_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"creation_category": tc.creationCategory,
			})
			result := polardbProxyClassDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbProxyTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		creationCategory string
		expected         bool
		description      string
	}{
		{
			name:             "SENormal",
			creationCategory: "SENormal",
			expected:         false,
			description:      "SENormal should not suppress diff",
		},
		{
			name:             "Normal",
			creationCategory: "Normal",
			expected:         false,
			description:      "Normal should not suppress diff",
		},
		{
			name:             "Other",
			creationCategory: "Archive",
			expected:         true,
			description:      "Other category should suppress diff",
		},
		{
			name:             "Empty",
			creationCategory: "",
			expected:         true,
			description:      "Empty category should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"creation_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"creation_category": tc.creationCategory,
			})
			result := polardbProxyTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name           string
		dbType         string
		creationOption string
		expected       bool
		description    string
	}{
		{
			name:           "MySQL_Normal",
			dbType:         "MySQL",
			creationOption: "Normal",
			expected:       false,
			description:    "MySQL with Normal creation_option should not suppress diff",
		},
		{
			name:           "MySQL_No_Option",
			dbType:         "MySQL",
			creationOption: "",
			expected:       false,
			description:    "MySQL with no creation_option (!optionOk) should not suppress diff",
		},
		{
			name:           "MySQL_CreateGdnStandby",
			dbType:         "MySQL",
			creationOption: "CreateGdnStandby",
			expected:       true,
			description:    "MySQL with CreateGdnStandby should suppress diff",
		},
		{
			name:           "PostgreSQL_Normal",
			dbType:         "PostgreSQL",
			creationOption: "Normal",
			expected:       true,
			description:    "Non-MySQL with Normal should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"db_type":         {Type: schema.TypeString},
				"creation_option": {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":         tc.dbType,
				"creation_option": tc.creationOption,
			})
			result := polardbDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbXengineDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name           string
		dbType         string
		dbVersion      string
		creationOption string
		expected       bool
		description    string
	}{
		{
			name:           "MySQL_80_Normal",
			dbType:         "MySQL",
			dbVersion:      "8.0",
			creationOption: "Normal",
			expected:       false,
			description:    "MySQL 8.0 with Normal option should not suppress diff",
		},
		{
			name:           "MySQL_80_No_Option",
			dbType:         "MySQL",
			dbVersion:      "8.0",
			creationOption: "",
			expected:       false,
			description:    "MySQL 8.0 with no creation_option should not suppress diff",
		},
		{
			name:           "MySQL_80_CreateGdnStandby",
			dbType:         "MySQL",
			dbVersion:      "8.0",
			creationOption: "CreateGdnStandby",
			expected:       true,
			description:    "MySQL 8.0 with CreateGdnStandby should suppress diff",
		},
		{
			name:           "MySQL_57_Normal",
			dbType:         "MySQL",
			dbVersion:      "5.7",
			creationOption: "Normal",
			expected:       true,
			description:    "MySQL 5.7 (not 8.0) should suppress diff",
		},
		{
			name:           "PostgreSQL_80_Normal",
			dbType:         "PostgreSQL",
			dbVersion:      "8.0",
			creationOption: "Normal",
			expected:       true,
			description:    "Non-MySQL should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"db_type":         {Type: schema.TypeString},
				"db_version":      {Type: schema.TypeString},
				"creation_option": {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":         tc.dbType,
				"db_version":      tc.dbVersion,
				"creation_option": tc.creationOption,
			})
			result := polardbXengineDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbAndCreationDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		dbType           string
		creationCategory string
		expected         bool
		description      string
	}{
		{
			name:             "MySQL_Normal",
			dbType:           "MySQL",
			creationCategory: "Normal",
			expected:         false,
			description:      "MySQL with Normal category should not suppress diff",
		},
		{
			name:             "MySQL_NormalMultimaster",
			dbType:           "MySQL",
			creationCategory: "NormalMultimaster",
			expected:         false,
			description:      "MySQL with NormalMultimaster category should not suppress diff",
		},
		{
			name:             "MySQL_No_Category",
			dbType:           "MySQL",
			creationCategory: "",
			expected:         false,
			description:      "MySQL with no creation_category should not suppress diff",
		},
		{
			name:             "MySQL_SENormal",
			dbType:           "MySQL",
			creationCategory: "SENormal",
			expected:         true,
			description:      "MySQL with SENormal should suppress diff",
		},
		{
			name:             "PostgreSQL_Normal",
			dbType:           "PostgreSQL",
			creationCategory: "Normal",
			expected:         true,
			description:      "Non-MySQL should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"db_type":           {Type: schema.TypeString},
				"creation_category": {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":           tc.dbType,
				"creation_category": tc.creationCategory,
			})
			result := polardbAndCreationDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbCompressStorageDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name           string
		dbType         string
		creationOption string
		storageType    string
		expected       bool
		description    string
	}{
		{
			name:           "MySQL_Normal_PSL4",
			dbType:         "MySQL",
			creationOption: "Normal",
			storageType:    "PSL4",
			expected:       false,
			description:    "MySQL Normal with PSL4 storage should not suppress diff",
		},
		{
			name:           "MySQL_No_Option_PSL4",
			dbType:         "MySQL",
			creationOption: "",
			storageType:    "PSL4",
			expected:       false,
			description:    "MySQL with no option and PSL4 should not suppress diff",
		},
		{
			name:           "MySQL_Normal_ESSD",
			dbType:         "MySQL",
			creationOption: "Normal",
			storageType:    "PSL5",
			expected:       true,
			description:    "MySQL Normal with non-PSL4 storage should suppress diff",
		},
		{
			name:           "MySQL_CreateGdnStandby_PSL4",
			dbType:         "MySQL",
			creationOption: "CreateGdnStandby",
			storageType:    "PSL4",
			expected:       true,
			description:    "MySQL with CreateGdnStandby should suppress diff",
		},
		{
			name:           "PostgreSQL_Normal_PSL4",
			dbType:         "PostgreSQL",
			creationOption: "Normal",
			storageType:    "PSL4",
			expected:       true,
			description:    "Non-MySQL should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"db_type":         {Type: schema.TypeString},
				"creation_option": {Type: schema.TypeString},
				"storage_type":    {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":         tc.dbType,
				"creation_option": tc.creationOption,
				"storage_type":    tc.storageType,
			})
			result := polardbCompressStorageDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPolardbStandbyAzDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		creationOption  string
		hotStandby      string
		old             string
		new             string
		expected        bool
		description     string
	}{
		{
			name:           "Normal_HotStandby_NotAuto",
			creationOption: "Normal",
			hotStandby:     "ON",
			old:            "cn-hangzhou-h",
			new:            "cn-hangzhou-i",
			expected:       false,
			description:    "Normal with hot standby and non-auto values should not suppress diff",
		},
		{
			name:           "No_Option_HotStandby_NotAuto",
			creationOption: "",
			hotStandby:     "ON",
			old:            "cn-hangzhou-h",
			new:            "cn-hangzhou-i",
			expected:       false,
			description:    "No creation_option with hot standby and non-auto values should not suppress diff",
		},
		{
			name:           "Normal_HotStandby_Old_Auto",
			creationOption: "Normal",
			hotStandby:     "ON",
			old:            "auto",
			new:            "cn-hangzhou-i",
			expected:       true,
			description:    "Old is auto should suppress diff",
		},
		{
			name:           "Normal_HotStandby_New_Auto",
			creationOption: "Normal",
			hotStandby:     "ON",
			old:            "cn-hangzhou-h",
			new:            "auto",
			expected:       true,
			description:    "New is auto should suppress diff",
		},
		{
			name:           "Normal_HotStandby_OFF",
			creationOption: "Normal",
			hotStandby:     "OFF",
			old:            "cn-hangzhou-h",
			new:            "cn-hangzhou-i",
			expected:       true,
			description:    "hot_standby_cluster OFF should suppress diff",
		},
		{
			name:           "CreateGdnStandby_HotStandby",
			creationOption: "CreateGdnStandby",
			hotStandby:     "ON",
			old:            "cn-hangzhou-h",
			new:            "cn-hangzhou-i",
			expected:       true,
			description:    "Non-Normal option should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"creation_option":    {Type: schema.TypeString},
				"hot_standby_cluster": {Type: schema.TypeString},
			}, map[string]interface{}{
				"creation_option":    tc.creationOption,
				"hot_standby_cluster": tc.hotStandby,
			})
			result := polardbStandbyAzDiffSuppressFunc("key", tc.old, tc.new, d)
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
		name       string
		ipProtocol string
		new        string
		expected   bool
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
			name:       "With_Ipv6CidrIp",
			ipv6CidrIp: "2001:db8::/32",
			expected:   true,
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

func TestUnitCommonSlbAclDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		aclStatus   string
		expected    bool
		description string
	}{
		{
			name:        "ACL_On",
			aclStatus:   "on",
			expected:    false,
			description: "ACL on should not suppress diff",
		},
		{
			name:        "ACL_Off",
			aclStatus:   "off",
			expected:    true,
			description: "ACL off should suppress diff",
		},
		{
			name:        "Empty",
			aclStatus:   "",
			expected:    true,
			description: "Empty acl_status should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"acl_status": {Type: schema.TypeString},
			}, map[string]interface{}{
				"acl_status": tc.aclStatus,
			})
			result := slbAclDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbServerCertificateDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                  string
		alicloudCertificateId string
		expected              bool
		description           string
	}{
		{
			name:                  "With_AlicloudCertificateId",
			alicloudCertificateId: "cert-12345",
			expected:              true,
			description:           "With alicloud_certificate_id should suppress diff",
		},
		{
			name:                  "No_AlicloudCertificateId",
			alicloudCertificateId: "",
			expected:              false,
			description:           "Without alicloud_certificate_id should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"alicloud_certificate_id": {Type: schema.TypeString},
			}, map[string]interface{}{
				"alicloud_certificate_id": tc.alicloudCertificateId,
			})
			result := slbServerCertificateDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbRuleListenerSyncDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		listenerSync string
		expected     bool
		description  string
	}{
		{
			name:         "ListenerSync_Off",
			listenerSync: "off",
			expected:     false,
			description:  "listener_sync off should not suppress diff",
		},
		{
			name:         "ListenerSync_On",
			listenerSync: "on",
			expected:     true,
			description:  "listener_sync on should suppress diff",
		},
		{
			name:         "Empty",
			listenerSync: "",
			expected:     true,
			description:  "Empty listener_sync should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"listener_sync": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync": tc.listenerSync,
			})
			result := slbRuleListenerSyncDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbRuleStickySessionTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name          string
		listenerSync  string
		stickySession string
		expected      bool
		description   string
	}{
		{
			name:          "ListenerSync_Off_StickySession_On",
			listenerSync:  "off",
			stickySession: "on",
			expected:      false,
			description:   "listener_sync off with sticky_session on should not suppress diff",
		},
		{
			name:          "ListenerSync_Off_StickySession_Off",
			listenerSync:  "off",
			stickySession: "off",
			expected:      true,
			description:   "listener_sync off with sticky_session off should suppress diff",
		},
		{
			name:          "ListenerSync_On_StickySession_On",
			listenerSync:  "on",
			stickySession: "on",
			expected:      true,
			description:   "listener_sync on should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"listener_sync": {Type: schema.TypeString},
				"sticky_session": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync": tc.listenerSync,
				"sticky_session": tc.stickySession,
			})
			result := slbRuleStickySessionTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbRuleCookieTimeoutDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name              string
		listenerSync      string
		stickySession     string
		stickySessionType string
		expected          bool
		description       string
	}{
		{
			name:              "ListenerSync_Off_Sticky_On_Insert",
			listenerSync:      "off",
			stickySession:     "on",
			stickySessionType: "insert",
			expected:          false,
			description:       "listener_sync off, sticky on, insert type should not suppress diff",
		},
		{
			name:              "ListenerSync_Off_Sticky_On_Server",
			listenerSync:      "off",
			stickySession:     "on",
			stickySessionType: "server",
			expected:          true,
			description:       "listener_sync off, sticky on, server type should suppress diff",
		},
		{
			name:              "ListenerSync_On",
			listenerSync:      "on",
			stickySession:     "on",
			stickySessionType: "insert",
			expected:          true,
			description:       "listener_sync on should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"listener_sync":      {Type: schema.TypeString},
				"sticky_session":     {Type: schema.TypeString},
				"sticky_session_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync":      tc.listenerSync,
				"sticky_session":     tc.stickySession,
				"sticky_session_type": tc.stickySessionType,
			})
			result := slbRuleCookieTimeoutDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbRuleCookieDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name              string
		listenerSync      string
		stickySession     string
		stickySessionType string
		expected          bool
		description       string
	}{
		{
			name:              "ListenerSync_Off_Sticky_On_Server",
			listenerSync:      "off",
			stickySession:     "on",
			stickySessionType: "server",
			expected:          false,
			description:       "listener_sync off, sticky on, server type should not suppress diff",
		},
		{
			name:              "ListenerSync_Off_Sticky_On_Insert",
			listenerSync:      "off",
			stickySession:     "on",
			stickySessionType: "insert",
			expected:          true,
			description:       "listener_sync off, sticky on, insert type should suppress diff",
		},
		{
			name:              "ListenerSync_On",
			listenerSync:      "on",
			stickySession:     "on",
			stickySessionType: "server",
			expected:          true,
			description:       "listener_sync on should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"listener_sync":      {Type: schema.TypeString},
				"sticky_session":     {Type: schema.TypeString},
				"sticky_session_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync":      tc.listenerSync,
				"sticky_session":     tc.stickySession,
				"sticky_session_type": tc.stickySessionType,
			})
			result := slbRuleCookieDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSlbRuleHealthCheckDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		listenerSync string
		healthCheck  string
		expected     bool
		description  string
	}{
		{
			name:         "ListenerSync_Off_HealthCheck_On",
			listenerSync: "off",
			healthCheck:  "on",
			expected:     false,
			description:  "listener_sync off with health_check on should not suppress diff",
		},
		{
			name:         "ListenerSync_Off_HealthCheck_Off",
			listenerSync: "off",
			healthCheck:  "off",
			expected:     true,
			description:  "listener_sync off with health_check off should suppress diff",
		},
		{
			name:         "ListenerSync_On_HealthCheck_On",
			listenerSync: "on",
			healthCheck:  "on",
			expected:     true,
			description:  "listener_sync on should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"listener_sync": {Type: schema.TypeString},
				"health_check":  {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync": tc.listenerSync,
				"health_check":  tc.healthCheck,
			})
			result := slbRuleHealthCheckDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonHealthCheckDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		protocol     string
		healthCheck  string
		expected     bool
		description  string
	}{
		{
			name:        "TCP_NoHealthCheck",
			protocol:    "tcp",
			expected:    false,
			description: "TCP protocol (httpDiff=true) should not suppress diff",
		},
		{
			name:        "TCP_HealthCheck_On",
			protocol:    "tcp",
			healthCheck: "on",
			expected:    false,
			description: "TCP with health_check on should not suppress diff",
		},
		{
			name:        "HTTP_HealthCheck_On",
			protocol:    "http",
			healthCheck: "on",
			expected:    false,
			description: "HTTP with health_check on should not suppress diff",
		},
		{
			name:        "HTTP_NoHealthCheck",
			protocol:    "http",
			expected:    true,
			description: "HTTP without health_check should suppress diff",
		},
		{
			name:        "UDP_NoHealthCheck",
			protocol:    "udp",
			expected:    false,
			description: "UDP (httpDiff=true) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"protocol": tc.protocol,
			}
			if tc.healthCheck != "" {
				data["health_check"] = tc.healthCheck
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol":        {Type: schema.TypeString},
				"health_check":    {Type: schema.TypeString},
				"listener_forward": {Type: schema.TypeString},
			}, data)
			result := healthCheckDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonHealthCheckTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		protocol    string
		expected    bool
		description string
	}{
		{
			name:        "TCP_Protocol",
			protocol:    "tcp",
			expected:    false,
			description: "TCP protocol should not suppress diff",
		},
		{
			name:        "HTTP_Protocol",
			protocol:    "http",
			expected:    true,
			description: "HTTP protocol should suppress diff",
		},
		{
			name:        "UDP_Protocol",
			protocol:    "udp",
			expected:    true,
			description: "UDP protocol should suppress diff",
		},
		{
			name:        "Empty_Protocol",
			protocol:    "",
			expected:    true,
			description: "Empty protocol should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"protocol": {Type: schema.TypeString},
			}, map[string]interface{}{
				"protocol": tc.protocol,
			})
			result := healthCheckTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEcsInternetDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		maxBandwidth int
		expected    bool
		description string
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

func TestUnitCommonRouterInterfaceAcceptsideDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		role        string
		expected    bool
		description string
	}{
		{
			name:        "AcceptingSide",
			role:        "AcceptingSide",
			expected:    true,
			description: "AcceptingSide role should suppress diff",
		},
		{
			name:        "InitiatingSide",
			role:        "InitiatingSide",
			expected:    false,
			description: "InitiatingSide role should not suppress diff",
		},
		{
			name:        "Empty_Role",
			role:        "",
			expected:    false,
			description: "Empty role should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"role": {Type: schema.TypeString},
			}, map[string]interface{}{
				"role": tc.role,
			})
			result := routerInterfaceAcceptsideDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonRouterInterfaceVBRTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		role        string
		routerType  string
		expected    bool
		description string
	}{
		{
			name:        "AcceptingSide",
			role:        "AcceptingSide",
			expected:    true,
			description: "AcceptingSide role should suppress diff",
		},
		{
			name:        "InitiatingSide_VRouter",
			role:        "InitiatingSide",
			routerType:  "VRouter",
			expected:    true,
			description: "VRouter type should suppress diff",
		},
		{
			name:        "InitiatingSide_VBR",
			role:        "InitiatingSide",
			routerType:  "VBR",
			expected:    false,
			description: "InitiatingSide with VBR should not suppress diff",
		},
		{
			name:        "Empty_Role_VBR",
			role:        "",
			routerType:  "VBR",
			expected:    false,
			description: "Empty role with VBR should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"role":        {Type: schema.TypeString},
				"router_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"role":        tc.role,
				"router_type": tc.routerType,
			})
			result := routerInterfaceVBRTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonWorkerDataDiskSizeSuppressFunc(t *testing.T) {
	testCases := []struct {
		name                    string
		workerDataDiskCategory  string
		resourceId              string
		forceUpdate             bool
		expected                bool
		description             string
	}{
		{
			name:        "No_DiskCategory",
			expected:    true,
			description: "No worker_data_disk_category should suppress diff",
		},
		{
			name:                   "DiskCategory_Set_New_Resource",
			workerDataDiskCategory: "cloud_ssd",
			resourceId:             "",
			forceUpdate:            false,
			expected:               false,
			description:            "DiskCategory set on new resource should not suppress diff",
		},
		{
			name:                   "DiskCategory_Set_Existing_No_ForceUpdate",
			workerDataDiskCategory: "cloud_ssd",
			resourceId:             "existing-id",
			forceUpdate:            false,
			expected:               true,
			description:            "DiskCategory set on existing resource without force_update should suppress diff",
		},
		{
			name:                   "DiskCategory_Set_Existing_ForceUpdate",
			workerDataDiskCategory: "cloud_ssd",
			resourceId:             "existing-id",
			forceUpdate:            true,
			expected:               false,
			description:            "DiskCategory set on existing resource with force_update should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{
				"force_update": tc.forceUpdate,
			}
			if tc.workerDataDiskCategory != "" {
				data["worker_data_disk_category"] = tc.workerDataDiskCategory
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"worker_data_disk_category": {Type: schema.TypeString},
				"force_update":              {Type: schema.TypeBool},
			}, data)
			d.SetId(tc.resourceId)
			result := workerDataDiskSizeSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonImageIdSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		newValue    string
		resourceId  string
		forceUpdate bool
		expected    bool
		description string
	}{
		{
			name:        "Empty_New",
			newValue:    "",
			expected:    true,
			description: "Empty new value should suppress diff",
		},
		{
			name:        "NonEmpty_New_New_Resource",
			newValue:    "aliyun-image-123",
			resourceId:  "",
			forceUpdate: false,
			expected:    false,
			description: "Non-empty new on new resource should not suppress diff",
		},
		{
			name:        "NonEmpty_New_Existing_No_ForceUpdate",
			newValue:    "aliyun-image-123",
			resourceId:  "existing-id",
			forceUpdate: false,
			expected:    true,
			description: "Non-empty new on existing resource without force_update should suppress diff",
		},
		{
			name:        "NonEmpty_New_Existing_ForceUpdate",
			newValue:    "aliyun-image-123",
			resourceId:  "existing-id",
			forceUpdate: true,
			expected:    false,
			description: "Non-empty new on existing resource with force_update should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"force_update": {Type: schema.TypeBool},
			}, map[string]interface{}{
				"force_update": tc.forceUpdate,
			})
			d.SetId(tc.resourceId)
			result := imageIdSuppressFunc("key", "old", tc.newValue, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEsVersionDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "Same_Major_Minor",
			old:         "6.7_with_X-Pack",
			new:         "6.7_with_X-Pack",
			expected:    true,
			description: "Same major.minor versions should suppress diff",
		},
		{
			name:        "Same_Major_Minor_Different_Patch",
			old:         "6.7",
			new:         "6.7.3",
			expected:    true,
			description: "Same major.minor but different patch should suppress diff",
		},
		{
			name:        "Different_Minor",
			old:         "6.7_with_X-Pack",
			new:         "6.8_with_X-Pack",
			expected:    false,
			description: "Different minor version should not suppress diff",
		},
		{
			name:        "Different_Major",
			old:         "6.7_with_X-Pack",
			new:         "7.7_with_X-Pack",
			expected:    false,
			description: "Different major version should not suppress diff",
		},
		{
			name:        "Empty_Versions",
			old:         "",
			new:         "",
			expected:    false,
			description: "Empty versions (len < 2) should not suppress diff",
		},
		{
			name:        "Single_Part_Version",
			old:         "6",
			new:         "6",
			expected:    false,
			description: "Single part version (len < 2) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			result := esVersionDiffSuppressFunc("key", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonVpnSslConnectionsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		enableSsl   bool
		hasSsl      bool
		expected    bool
		description string
	}{
		{
			name:        "SSL_Enabled",
			enableSsl:   true,
			hasSsl:      true,
			expected:    false,
			description: "enable_ssl=true should not suppress diff",
		},
		{
			name:        "SSL_Disabled",
			enableSsl:   false,
			hasSsl:      true,
			expected:    true,
			description: "enable_ssl=false should suppress diff",
		},
		{
			name:        "SSL_Not_Set",
			hasSsl:      false,
			expected:    true,
			description: "enable_ssl not set should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.hasSsl {
				data["enable_ssl"] = tc.enableSsl
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"enable_ssl": {Type: schema.TypeBool},
			}, data)
			result := vpnSslConnectionsDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonKmsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name            string
		password        string
		accountPassword string
		expected        bool
		description     string
	}{
		{
			name:        "With_Password",
			password:    "mypassword",
			expected:    true,
			description: "With password should suppress diff",
		},
		{
			name:            "With_AccountPassword",
			accountPassword: "myaccountpassword",
			expected:        true,
			description:     "With account_password should suppress diff",
		},
		{
			name:        "No_Password",
			expected:    false,
			description: "Without any password should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := map[string]interface{}{}
			if tc.password != "" {
				data["password"] = tc.password
			}
			if tc.accountPassword != "" {
				data["account_password"] = tc.accountPassword
			}
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"password":         {Type: schema.TypeString},
				"account_password": {Type: schema.TypeString},
			}, data)
			result := kmsDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSagDnatEntryTypeDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		entryType   string
		expected    bool
		description string
	}{
		{
			name:        "Intranet_Type",
			entryType:   "Intranet",
			expected:    false,
			description: "Intranet type should not suppress diff",
		},
		{
			name:        "Internet_Type",
			entryType:   "Internet",
			expected:    true,
			description: "Internet type should suppress diff",
		},
		{
			name:        "Empty_Type",
			entryType:   "",
			expected:    true,
			description: "Empty type should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"type": tc.entryType,
			})
			result := sagDnatEntryTypeDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSagClientUserPasswordSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		userName    string
		expected    bool
		description string
	}{
		{
			name:        "With_UserName",
			userName:    "alice",
			expected:    false,
			description: "Non-empty user_name should not suppress diff",
		},
		{
			name:        "Empty_UserName",
			userName:    "",
			expected:    true,
			description: "Empty user_name should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"user_name": {Type: schema.TypeString},
			}, map[string]interface{}{
				"user_name": tc.userName,
			})
			result := sagClientUserPasswordSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSelectdbPostPaidDiffSuppressFunc(t *testing.T) {
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
				"payment_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"payment_type": tc.paymentType,
			})
			result := selectdbPostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCmsClientInfoSuppressFunc(t *testing.T) {
	escalationSchema := map[string]*schema.Schema{
		"statistics":          {Type: schema.TypeString},
		"comparison_operator": {Type: schema.TypeString},
		"threshold":           {Type: schema.TypeString},
		"times":               {Type: schema.TypeString},
	}

	testCases := []struct {
		name        string
		items       []interface{}
		expected    bool
		description string
	}{
		{
			name:        "Empty_List",
			items:       []interface{}{},
			expected:    false,
			description: "Empty escalations_info should not suppress diff",
		},
		{
			name: "Item_With_All_Fields",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": ">=",
					"threshold":           "90",
					"times":               "3",
				},
			},
			expected:    false,
			description: "Item with all fields set should not suppress diff",
		},
		{
			name: "Item_With_Empty_Statistics",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "",
					"comparison_operator": ">=",
					"threshold":           "90",
					"times":               "3",
				},
			},
			expected:    true,
			description: "Item with empty statistics should suppress diff",
		},
		{
			name: "Item_With_Empty_Operator",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": "",
					"threshold":           "90",
					"times":               "3",
				},
			},
			expected:    true,
			description: "Item with empty comparison_operator should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"escalations_info": {
					Type: schema.TypeList,
					Elem: &schema.Resource{Schema: escalationSchema},
				},
			}, map[string]interface{}{
				"escalations_info": tc.items,
			})
			result := cmsClientInfoSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCmsClientWarnSuppressFunc(t *testing.T) {
	escalationSchema := map[string]*schema.Schema{
		"statistics":          {Type: schema.TypeString},
		"comparison_operator": {Type: schema.TypeString},
		"threshold":           {Type: schema.TypeString},
		"times":               {Type: schema.TypeString},
	}

	testCases := []struct {
		name        string
		items       []interface{}
		expected    bool
		description string
	}{
		{
			name:        "Empty_List",
			items:       []interface{}{},
			expected:    false,
			description: "Empty escalations_warn should not suppress diff",
		},
		{
			name: "Item_With_All_Fields",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": ">=",
					"threshold":           "80",
					"times":               "2",
				},
			},
			expected:    false,
			description: "Item with all fields set should not suppress diff",
		},
		{
			name: "Item_With_Empty_Threshold",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": ">=",
					"threshold":           "",
					"times":               "2",
				},
			},
			expected:    true,
			description: "Item with empty threshold should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"escalations_warn": {
					Type: schema.TypeList,
					Elem: &schema.Resource{Schema: escalationSchema},
				},
			}, map[string]interface{}{
				"escalations_warn": tc.items,
			})
			result := cmsClientWarnSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCmsClientCriticalSuppressFunc(t *testing.T) {
	escalationSchema := map[string]*schema.Schema{
		"statistics":          {Type: schema.TypeString},
		"comparison_operator": {Type: schema.TypeString},
		"threshold":           {Type: schema.TypeString},
		"times":               {Type: schema.TypeString},
	}

	testCases := []struct {
		name        string
		items       []interface{}
		expected    bool
		description string
	}{
		{
			name:        "Empty_List",
			items:       []interface{}{},
			expected:    false,
			description: "Empty escalations_critical should not suppress diff",
		},
		{
			name: "Item_With_Empty_Times",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": ">=",
					"threshold":           "95",
					"times":               "",
				},
			},
			expected:    true,
			description: "Item with empty times should suppress diff",
		},
		{
			name: "Item_With_All_Fields",
			items: []interface{}{
				map[string]interface{}{
					"statistics":          "Average",
					"comparison_operator": ">=",
					"threshold":           "95",
					"times":               "1",
				},
			},
			expected:    false,
			description: "Item with all fields set should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"escalations_critical": {
					Type: schema.TypeList,
					Elem: &schema.Resource{Schema: escalationSchema},
				},
			}, map[string]interface{}{
				"escalations_critical": tc.items,
			})
			result := cmsClientCriticalSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonAlikafkaInstanceConfigDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "New_Empty",
			old:         `{"key":"value"}`,
			new:         "",
			expected:    true,
			description: "Empty new value should suppress diff",
		},
		{
			name:        "Old_Empty",
			old:         "",
			new:         `{"key":"value"}`,
			expected:    false,
			description: "Empty old value should not suppress diff",
		},
		{
			name:        "Old_Invalid_JSON",
			old:         "invalid json",
			new:         `{"key":"value"}`,
			expected:    false,
			description: "Invalid old JSON should not suppress diff",
		},
		{
			name:        "New_Invalid_JSON",
			old:         `{"key":"value"}`,
			new:         "invalid json",
			expected:    false,
			description: "Invalid new JSON should not suppress diff",
		},
		{
			name:        "Same_Values",
			old:         `{"key":"value","other":"data"}`,
			new:         `{"key":"value"}`,
			expected:    true,
			description: "New is subset of old with same values should suppress diff",
		},
		{
			name:        "Different_Values",
			old:         `{"key":"old_value"}`,
			new:         `{"key":"new_value"}`,
			expected:    false,
			description: "Different values for same key should not suppress diff",
		},
		{
			name:        "New_Key_Not_In_Old",
			old:         `{"other":"value"}`,
			new:         `{"key":"value"}`,
			expected:    true,
			description: "New key not in old map should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			result := alikafkaInstanceConfigDiffSuppressFunc("key", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonPayTypePostPaidDiffSuppressFunc(t *testing.T) {
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
				"pay_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"pay_type": tc.payType,
			})
			result := payTypePostPaidDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEngineDiffSuppressFunc(t *testing.T) {
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
				"engine": {Type: schema.TypeString},
			}, map[string]interface{}{
				"engine": tc.engine,
			})
			result := engineDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSslEnabledDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		sslEnabled  int
		expected    bool
		description string
	}{
		{
			name:        "SSL_Enabled_1",
			sslEnabled:  1,
			expected:    false,
			description: "ssl_enabled=1 should not suppress diff",
		},
		{
			name:        "SSL_Disabled_0",
			sslEnabled:  0,
			expected:    true,
			description: "ssl_enabled=0 (GetOk returns false for int 0) should suppress diff",
		},
		{
			name:        "SSL_Enabled_2",
			sslEnabled:  2,
			expected:    true,
			description: "ssl_enabled=2 (not 1) should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"ssl_enabled": {Type: schema.TypeInt},
			}, map[string]interface{}{
				"ssl_enabled": tc.sslEnabled,
			})
			result := sslEnabledDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSslActionDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		sslAction   string
		expected    bool
		description string
	}{
		{
			name:        "Open",
			sslAction:   "Open",
			expected:    false,
			description: "ssl_action=Open should not suppress diff",
		},
		{
			name:        "Update",
			sslAction:   "Update",
			expected:    false,
			description: "ssl_action=Update should not suppress diff",
		},
		{
			name:        "Close",
			sslAction:   "Close",
			expected:    true,
			description: "ssl_action=Close should suppress diff",
		},
		{
			name:        "Empty",
			sslAction:   "",
			expected:    true,
			description: "Empty ssl_action should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"ssl_action": {Type: schema.TypeString},
			}, map[string]interface{}{
				"ssl_action": tc.sslAction,
			})
			result := sslActionDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonSecurityIpsDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		securityIps []interface{}
		expected    bool
		description string
	}{
		{
			name:        "With_IPs",
			securityIps: []interface{}{"192.168.0.1", "10.0.0.0/8"},
			expected:    false,
			description: "Non-empty security_ips should not suppress diff",
		},
		{
			name:        "Empty_IPs",
			securityIps: []interface{}{},
			expected:    true,
			description: "Empty security_ips should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"security_ips": {
					Type: schema.TypeSet,
					Elem: &schema.Schema{Type: schema.TypeString},
					Set:  schema.HashString,
				},
			}, map[string]interface{}{
				"security_ips": tc.securityIps,
			})
			result := securityIpsDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonKernelVersionDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		upgradeFlag bool
		expected    bool
		description string
	}{
		{
			name:        "Upgrade_True",
			upgradeFlag: true,
			expected:    false,
			description: "upgrade_db_instance_kernel_version=true should not suppress diff",
		},
		{
			name:        "Upgrade_False",
			upgradeFlag: false,
			expected:    true,
			description: "upgrade_db_instance_kernel_version=false (GetOk returns false) should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"upgrade_db_instance_kernel_version": {Type: schema.TypeBool},
			}, map[string]interface{}{
				"upgrade_db_instance_kernel_version": tc.upgradeFlag,
			})
			result := kernelVersionDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonKernelSmallVersionDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name               string
		targetMinorVersion string
		expected           bool
		description        string
	}{
		{
			name:               "Has_Minor_Version",
			targetMinorVersion: "rds_20240101",
			expected:           false,
			description:        "Non-empty target_minor_version: HasChange=true, GetOk=true → should not suppress diff",
		},
		{
			name:               "Empty_Minor_Version",
			targetMinorVersion: "",
			expected:           true,
			description:        "Empty minor version: HasChange=false → should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"target_minor_version": {Type: schema.TypeString},
			}, map[string]interface{}{
				"target_minor_version": tc.targetMinorVersion,
			})
			result := kernelSmallVersionDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonStorageAutoScaleDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name             string
		storageAutoScale string
		expected         bool
		description      string
	}{
		{
			name:             "Enable",
			storageAutoScale: "Enable",
			expected:         false,
			description:      "storage_auto_scale=Enable should not suppress diff",
		},
		{
			name:             "enable_lowercase",
			storageAutoScale: "enable",
			expected:         false,
			description:      "storage_auto_scale=enable (lowercase) should not suppress diff",
		},
		{
			name:             "Disable",
			storageAutoScale: "Disable",
			expected:         true,
			description:      "storage_auto_scale=Disable should suppress diff",
		},
		{
			name:             "Empty",
			storageAutoScale: "",
			expected:         true,
			description:      "Empty storage_auto_scale should suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"storage_auto_scale": {Type: schema.TypeString},
			}, map[string]interface{}{
				"storage_auto_scale": tc.storageAutoScale,
			})
			result := StorageAutoScaleDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonCmsAlarmDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		old         string
		new         string
		expected    bool
		description string
	}{
		{
			name:        "New_Empty",
			old:         `[{"key":"value"}]`,
			new:         "",
			expected:    true,
			description: "Empty new value should suppress diff",
		},
		{
			name:        "Old_Empty",
			old:         "",
			new:         `[{"key":"value"}]`,
			expected:    false,
			description: "Empty old with non-empty new should not suppress diff",
		},
		{
			name:        "Same_Order",
			old:         `[{"a":"1"},{"b":"2"}]`,
			new:         `[{"a":"1"},{"b":"2"}]`,
			expected:    false,
			description: "Exactly equal strings (new == old check) should not suppress diff",
		},
		{
			name:        "Different_Order_Same_Content",
			old:         `[{"b":"2"},{"a":"1"}]`,
			new:         `[{"a":"1"},{"b":"2"}]`,
			expected:    true,
			description: "Different order but same content after sorting should suppress diff",
		},
		{
			name:        "Different_Content",
			old:         `[{"a":"1"}]`,
			new:         `[{"a":"2"}]`,
			expected:    false,
			description: "Different content should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{}, map[string]interface{}{})
			result := CmsAlarmDiffSuppressFunc("key", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}

func TestUnitCommonEsDataNodeDiskPerformanceLevelDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name         string
		diskType     string
		expected     bool
		description  string
	}{
		{
			name:        "Cloud_ESSD",
			diskType:    "cloud_essd",
			expected:    false,
			description: "cloud_essd disk type should not suppress diff",
		},
		{
			name:        "Cloud_SSD",
			diskType:    "cloud_ssd",
			expected:    true,
			description: "Non-cloud_essd should suppress diff",
		},
		{
			name:        "Cloud",
			diskType:    "cloud",
			expected:    true,
			description: "Cloud disk type should suppress diff",
		},
		{
			name:        "Empty",
			diskType:    "",
			expected:    false,
			description: "Empty disk type (GetOk returns false) should not suppress diff",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
				"data_node_disk_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"data_node_disk_type": tc.diskType,
			})
			result := esDataNodeDiskPerformanceLevelDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
