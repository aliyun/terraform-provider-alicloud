package alicloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// TestAccAlicloudDataWorksColumnsDataSource verifies the data.alicloud_data_works_columns
// data source against a real DataWorks table. A table fixture cannot be created by the
// provider, so the test queries a table identified by the ALICLOUD_DATA_WORKS_TABLE_GUID
// environment variable. When the variable is unset the test is skipped; when it is set the
// data source lists the columns of that table and the fake step queries a non-existent
// table and expects an empty result.
func TestAccAlicloudDataWorksColumnsDataSource(t *testing.T) {
	tableGuid := os.Getenv("ALICLOUD_DATA_WORKS_TABLE_GUID")
	if tableGuid == "" {
		t.Skip("Skipping because ALICLOUD_DATA_WORKS_TABLE_GUID is not set. Provide a real DataWorks table GUID to run this acceptance test.")
	}
	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_columns.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, tableGuid, dataSourceDataWorksColumnsConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"table_guid": tableGuid,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"table_guid": tableGuid + "_fake",
		}),
	}

	var existDataWorksColumnsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 CHECKSET,
			"columns.#":             CHECKSET,
			"columns.0.column_guid": CHECKSET,
			"columns.0.column_name": CHECKSET,
			"columns.0.column_type": CHECKSET,
			"columns.0.region_id":   CHECKSET,
		}
	}
	var fakeDataWorksColumnsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"columns.#": "0",
		}
	}

	var DataWorksColumnsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksColumnsMapFunc,
		fakeMapFunc:  fakeDataWorksColumnsMapFunc,
	}

	DataWorksColumnsCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataWorksColumnsConfigDependence(name string) string {
	return ""
}
