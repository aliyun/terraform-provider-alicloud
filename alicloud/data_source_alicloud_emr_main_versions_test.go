package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAlicloudEmrMainVersionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceID := "data.alicloud_emr_main_versions.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-main-versions%v.abc", defaultRegionToTest, rand),
		dataSourceEmrMainVersionsConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"emr_version": "EMR-3.22.0",
		}),
	}

	var (
		existEmrMainVersionsMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                           CHECKSET,
				"ids.0":                           CHECKSET,
				"main_versions.#":                 CHECKSET,
				"main_versions.0.emr_version":     CHECKSET,
				"main_versions.0.image_id":        CHECKSET,
				"main_versions.0.cluster_types.#": CHECKSET,
				"main_versions.0.cluster_types.0": CHECKSET,
			}
		}

		fakeEmrMainVersionsMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           CHECKSET,
				"main_versions.#": CHECKSET,
			}
		}

		emrMainVersionsCheckInfo = dataSourceAttr{
			resourceId:   resourceID,
			existMapFunc: existEmrMainVersionsMapFunc,
			fakeMapFunc:  fakeEmrMainVersionsMapFunc,
		}
	)

	emrMainVersionsCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrMainVersionsConfigDependence(name string) string {
	return ``
}
