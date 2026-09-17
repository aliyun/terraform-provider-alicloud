package alicloud

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

// TestAccAlicloudDataWorksColumnLineageDataSource verifies the
// data.alicloud_data_works_column_lineage data source against a real
// DataWorks column. A column lineage fixture cannot be created by the
// provider, so the test queries the lineage of a column identified by the
// ALICLOUD_DATA_WORKS_COLUMN_GUID environment variable. When the variable is
// unset the test is skipped; when it is set the data source lists the upstream
// lineage of that column and the fake step queries a non-existent column and
// expects an empty result.
func TestAccAlicloudDataWorksColumnLineageDataSource(t *testing.T) {
	columnGuid := os.Getenv("ALICLOUD_DATA_WORKS_COLUMN_GUID")
	if columnGuid == "" {
		t.Skip("Skipping because ALICLOUD_DATA_WORKS_COLUMN_GUID is not set. Provide a real DataWorks column GUID to run this acceptance test.")
	}
	rand := acctest.RandInt()
	resourceId := "data.alicloud_data_works_column_lineage.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, columnGuid, dataSourceDataWorksColumnLineageConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"column_guid": columnGuid,
			"direction":   "up",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"column_guid": columnGuid + "_fake",
			"direction":   "up",
		}),
	}

	var existDataWorksColumnLineageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         CHECKSET,
			"column_lineages.#":             CHECKSET,
			"column_lineages.0.column_guid": CHECKSET,
			"column_lineages.0.column_name": CHECKSET,
			"column_lineages.0.region_id":   CHECKSET,
		}
	}
	var fakeDataWorksColumnLineageMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":             "0",
			"column_lineages.#": "0",
		}
	}

	var DataWorksColumnLineageCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksColumnLineageMapFunc,
		fakeMapFunc:  fakeDataWorksColumnLineageMapFunc,
	}

	DataWorksColumnLineageCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceDataWorksColumnLineageConfigDependence(name string) string {
	return ""
}
