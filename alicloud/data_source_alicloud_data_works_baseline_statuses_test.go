package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAliCloudDataWorksBaselineStatusesDataSource(t *testing.T) {
	resourceId := "data.alicloud_data_works_baseline_statuses.default"
	name := "tf-testacc-dataworks-baseline-statuses"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceDataWorksBaselineStatusesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"bizdate": "2026-09-16T00:00:00Z",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"bizdate": "2026-09-16T00:00:00Z",
			"ids":     []string{"999999999"},
		}),
	}

	var existDataWorksBaselineStatusesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"statuses.#": CHECKSET,
		}
	}

	var fakeDataWorksBaselineStatusesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"statuses.#": "0",
			"ids.#":      "0",
		}
	}

	var DataWorksBaselineStatusesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDataWorksBaselineStatusesMapFunc,
		fakeMapFunc:  fakeDataWorksBaselineStatusesMapFunc,
	}

	DataWorksBaselineStatusesCheckInfo.dataSourceTestCheck(t, acctest.RandInt(), idsConf)
}

func dataSourceDataWorksBaselineStatusesConfigDependence(name string) string {
	return ""
}
