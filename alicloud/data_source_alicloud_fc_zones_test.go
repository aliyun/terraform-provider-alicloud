package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudFCZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_fc_zones.default"

	var fcZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"ids.0":      CHECKSET,
			"zones.#":    CHECKSET,
			"zones.0.id": CHECKSET,
		}
	}

	var fakeFCZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var fcZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: fcZonesMapFunc,
		fakeMapFunc:  fakeFCZonesMapFunc,
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudFCZonesDataSourceConfig(),
	}

	fcZonesCheckInfo.dataSourceTestCheck(t, rand, allConf)
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
