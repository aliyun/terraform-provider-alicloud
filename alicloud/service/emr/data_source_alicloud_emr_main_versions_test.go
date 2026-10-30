package emr_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudEmrMainVersionsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceID := "data.alicloud_emr_main_versions.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceID,
		fmt.Sprintf("tf-testacc%semr-main-versions%v.abc", acctest.DefaultRegionToTest, rand),
		dataSourceEmrMainVersionsConfigDependence)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{}),
		FakeConfig: testAccConfig(map[string]interface{}{
			"emr_version": "EMR-3.22.0",
		}),
	}

	var (
		existEmrMainVersionsMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                           acctest.CHECKSET,
				"ids.0":                           acctest.CHECKSET,
				"main_versions.#":                 acctest.CHECKSET,
				"main_versions.0.emr_version":     acctest.CHECKSET,
				"main_versions.0.image_id":        acctest.CHECKSET,
				"main_versions.0.cluster_types.#": acctest.CHECKSET,
				"main_versions.0.cluster_types.0": acctest.CHECKSET,
			}
		}

		fakeEmrMainVersionsMapFunc = func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           acctest.CHECKSET,
				"main_versions.#": acctest.CHECKSET,
			}
		}

		emrMainVersionsCheckInfo = acctest.DataSourceAttr{
			ResourceId:   resourceID,
			ExistMapFunc: existEmrMainVersionsMapFunc,
			FakeMapFunc:  fakeEmrMainVersionsMapFunc,
		}
	)

	emrMainVersionsCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

func dataSourceEmrMainVersionsConfigDependence(name string) string {
	return ``
}
