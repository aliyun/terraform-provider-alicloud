package mongodb_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMongoDBZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_mongodb_zones.default"

	testAccConfig := acctest.DataSourceTestAccConfigFunc(resourceId, "", dataSourceMongoDBZonesConfigDependence)

	multiConfig := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccConfig(map[string]interface{}{
			"multi": "true",
		}),
	}

	var existDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    acctest.CHECKSET,
			"ids.0":                    acctest.CHECKSET,
			"zones.#":                  acctest.CHECKSET,
			"zones.0.id":               acctest.CHECKSET,
			"zones.0.multi_zone_ids.#": acctest.CHECKSET,
		}
	}

	var fakeMongDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var MongDBZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: existDBZonesMapFunc,
		FakeMapFunc:  fakeMongDBZonesMapFunc,
	}

	preCheck := func() {
		acctest.PreCheck(t)
		acctest.PreCheckWithRegions(t, true, []connectivity.Region{"eu-central-1"})
	}
	MongDBZonesCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, multiConfig)
}

func dataSourceMongoDBZonesConfigDependence(name string) string {
	return ""
}
