package rds_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"

	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRdsCrossRegionsDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_rds_cross_regions.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceRdsCrossRegionsConfigDependence)

	rdsCrossRegionsConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{}),
	}

	var existRdsCrossRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        acctest.CHECKSET,
			"ids.0":        acctest.CHECKSET,
			"regions.#":    acctest.CHECKSET,
			"regions.0.id": acctest.CHECKSET,
		}
	}

	var fakeRdsCrossRegionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"regions.#": "0",
		}
	}

	var RdsCrossRegionsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existRdsCrossRegionsMapFunc,
		FakeMapFunc:  fakeRdsCrossRegionsMapFunc,
	}

	RdsCrossRegionsCheckInfo.DataSourceTestCheck(t, rand, rdsCrossRegionsConfig)
}

func dataSourceRdsCrossRegionsConfigDependence(name string) string {
	return ""
}
