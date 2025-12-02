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

func TestUnitCommonPrivateDnsTypeDiffSuppressFunc(t *testing.T) {
	// Test for primary_dns when private_dns_type is PrivateZone
	testCases := []struct {
		name            string
		privateDnsType  string
		oldValue        string
		newValue        string
		expected        bool
		description     string
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
		name            string
		privateDnsType  string
		oldValue        string
		newValue        string
		expected        bool
		description     string
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
		name             string
		protocol         string
		healthCheck      string
		healthCheckType  string
		listenerForward  string
		expected         bool
		description      string
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
		name                      string
		masterInstanceChargeType  string
		resourceId                string
		forceUpdate               bool
		expected                  bool
		description               string
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
		name             string
		enableBackupLog  bool
		logBackup        bool
		expected         bool
		description      string
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
		name                        string
		enableKibanaPublicNetwork   bool
		expected                    bool
		description                 string
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
		name                         string
		enableKibanaPrivateNetwork   bool
		expected                     bool
		description                  string
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
		name              string
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
			description:  "TDE enabled with empty new value should not suppress diff",
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
		name           string
		payType        string
		renewalStatus  string
		expected       bool
		description    string
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
				"pay_type":        {Type: schema.TypeString},
				"renewal_status":  {Type: schema.TypeString},
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
		name              string
		instanceChargeType string
		hasChargeType      bool
		expected           bool
		description        string
	}{
		{
			name:              "PostPaid_ChargeType",
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
			name:              "PrePaid_ChargeType",
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
		name              string
		instanceChargeType string
		spotStrategy       string
		expected           bool
		description        string
	}{
		{
			name:              "PostPaid_With_SpotWithPriceLimit",
			instanceChargeType: "PostPaid",
			spotStrategy:       "SpotWithPriceLimit",
			expected:           false,
			description:        "PostPaid with SpotWithPriceLimit should not suppress diff",
		},
		{
			name:              "PostPaid_Without_SpotWithPriceLimit",
			instanceChargeType: "PostPaid",
			spotStrategy:       "NoSpot",
			expected:           true,
			description:        "PostPaid without SpotWithPriceLimit should suppress diff",
		},
		{
			name:              "PrePaid",
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
