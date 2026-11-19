// lintignore: S013
package cs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

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
		name         string
		diskCategory string
		expected     bool
		description  string
	}{
		{
			name:         "CloudESSD",
			diskCategory: "cloud_essd",
			expected:     false,
			description:  "cloud_essd category should not suppress diff",
		},
		{
			name:         "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:     true,
			description:  "Non-cloud_essd category should suppress diff",
		},
		{
			name:         "Cloud",
			diskCategory: "cloud",
			expected:     true,
			description:  "Cloud category should suppress diff",
		},
		{
			name:         "Empty",
			diskCategory: "",
			expected:     false,
			description:  "Empty disk category (GetOk returns false) should not suppress diff",
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
		name         string
		diskCategory string
		expected     bool
		description  string
	}{
		{
			name:         "CloudESSD",
			diskCategory: "cloud_essd",
			expected:     false,
			description:  "cloud_essd category should not suppress diff",
		},
		{
			name:         "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:     true,
			description:  "Non-cloud_essd category should suppress diff",
		},
		{
			name:         "Empty",
			diskCategory: "",
			expected:     false,
			description:  "Empty disk category (GetOk returns false) should not suppress diff",
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
			name:         "CloudSSD",
			diskCategory: "cloud_ssd",
			expected:     true,
			description:  "cloud_ssd should suppress diff",
		},
		{
			name:         "Empty",
			diskCategory: "",
			expected:     false,
			description:  "Empty system_disk_category (GetOk returns false) should not suppress diff",
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
