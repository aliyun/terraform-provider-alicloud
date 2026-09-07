package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRdsCrossRegionsDataSource(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_rds_cross_regions.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceRdsCrossRegionsConfigDependence)

	rdsCrossRegionsConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existRdsCrossRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        CHECKSET,
			"ids.0":        CHECKSET,
			"regions.#":    CHECKSET,
			"regions.0.id": CHECKSET,
		}
	}

	var fakeRdsCrossRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"regions.#": "0",
		}
	}

	var RdsCrossRegionsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRdsCrossRegionsMapFunc,
		fakeMapFunc:  fakeRdsCrossRegionsMapFunc,
	}

	RdsCrossRegionsCheckInfo.dataSourceTestCheck(t, rand, rdsCrossRegionsConfig)
}

func dataSourceRdsCrossRegionsConfigDependence(name string) string {
	return ""
}
