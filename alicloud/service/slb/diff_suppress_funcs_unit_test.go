// lintignore: S013
package slb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				"listener_sync":  {Type: schema.TypeString},
				"sticky_session": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync":  tc.listenerSync,
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
				"listener_sync":       {Type: schema.TypeString},
				"sticky_session":      {Type: schema.TypeString},
				"sticky_session_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync":       tc.listenerSync,
				"sticky_session":      tc.stickySession,
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
				"listener_sync":       {Type: schema.TypeString},
				"sticky_session":      {Type: schema.TypeString},
				"sticky_session_type": {Type: schema.TypeString},
			}, map[string]interface{}{
				"listener_sync":       tc.listenerSync,
				"sticky_session":      tc.stickySession,
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
		name        string
		protocol    string
		healthCheck string
		expected    bool
		description string
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
				"protocol":         {Type: schema.TypeString},
				"health_check":     {Type: schema.TypeString},
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
