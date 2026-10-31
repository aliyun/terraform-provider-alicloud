package fc_test

import (
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudFCZonesDataSource_basic(t *testing.T) {
	rand := tfacctest.RandInt()
	resourceId := "data.alicloud_fc_zones.default"

	var fcZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      acctest.CHECKSET,
			"ids.0":      acctest.CHECKSET,
			"zones.#":    acctest.CHECKSET,
			"zones.0.id": acctest.CHECKSET,
		}
	}

	var fakeFCZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var fcZonesCheckInfo = acctest.DataSourceAttr{
		ResourceId:   resourceId,
		ExistMapFunc: fcZonesMapFunc,
		FakeMapFunc:  fakeFCZonesMapFunc,
	}

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudFCZonesDataSourceConfig(),
	}

	fcZonesCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

func testAccCheckAlicloudFCZonesDataSourceConfig() string {
	// alicloud_fc_zones has no filter arguments, so a bare data source block is
	// the only valid config. Returning at least one config makes the test build a
	// real TestStep (previously it generated zero steps, which terraform-plugin-sdk/v2
	// rejects with "TestCase missing Steps").
	return `
data "alicloud_fc_zones" "default" {
}
`
}
