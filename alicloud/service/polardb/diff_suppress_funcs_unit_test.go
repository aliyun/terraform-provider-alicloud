// lintignore: S013
package polardb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

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
		name                   string
		dbType                 string
		dbVersion              string
		serverlessType         string
		serverlessSteadySwitch string
		expected               bool
		description            string
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
			name:                   "MySQL_57_SteadyServerless_ON",
			dbType:                 "MySQL",
			dbVersion:              "5.7",
			serverlessType:         "SteadyServerless",
			serverlessSteadySwitch: "ON",
			expected:               false,
			description:            "MySQL 5.7 with SteadyServerless and ON switch should not suppress diff",
		},
		{
			name:                   "MySQL_57_SteadyServerless_OFF",
			dbType:                 "MySQL",
			dbVersion:              "5.7",
			serverlessType:         "SteadyServerless",
			serverlessSteadySwitch: "OFF",
			expected:               true,
			description:            "MySQL 5.7 with SteadyServerless and OFF switch should suppress diff",
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
				"db_type":                  {Type: schema.TypeString},
				"db_version":               {Type: schema.TypeString},
				"serverless_type":          {Type: schema.TypeString},
				"serverless_steady_switch": {Type: schema.TypeString},
			}, map[string]interface{}{
				"db_type":                  tc.dbType,
				"db_version":               tc.dbVersion,
				"serverless_type":          tc.serverlessType,
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
			expected:         false,
			description:      "Normal category must preserve the diff independently of database type",
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
		name           string
		creationOption string
		hotStandby     string
		old            string
		new            string
		expected       bool
		description    string
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
				"creation_option":     {Type: schema.TypeString},
				"hot_standby_cluster": {Type: schema.TypeString},
			}, map[string]interface{}{
				"creation_option":     tc.creationOption,
				"hot_standby_cluster": tc.hotStandby,
			})
			result := polardbStandbyAzDiffSuppressFunc("key", tc.old, tc.new, d)
			assert.Equal(t, tc.expected, result, tc.description)
		})
	}
}
