package elasticsearch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

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
				"enable_public": {Type: schema.TypeBool, Optional: true},
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
				"enable_kibana_public_network": {Type: schema.TypeBool, Optional: true},
			}, map[string]interface{}{
				"enable_kibana_public_network": tc.enableKibanaPublicNetwork,
			})

			result := elasticsearchEnableKibanaPublicDiffSuppressFunc("key", "old", "new", d)
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
				"enable_kibana_private_network": {Type: schema.TypeBool, Optional: true},
			}, map[string]interface{}{
				"enable_kibana_private_network": tc.enableKibanaPrivateNetwork,
			})

			result := elasticsearchEnableKibanaPrivateDiffSuppressFunc("key", "old", "new", d)
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

func TestUnitCommonEsDataNodeDiskPerformanceLevelDiffSuppressFunc(t *testing.T) {
	testCases := []struct {
		name        string
		diskType    string
		expected    bool
		description string
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
				"data_node_disk_type": {Type: schema.TypeString, Optional: true},
			}, map[string]interface{}{
				"data_node_disk_type": tc.diskType,
			})
			result := esDataNodeDiskPerformanceLevelDiffSuppressFunc("key", "old", "new", d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
