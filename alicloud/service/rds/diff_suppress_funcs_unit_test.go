// lintignore: S013
package rds

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

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
